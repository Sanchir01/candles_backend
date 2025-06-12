package grpcclientauth

import (
	"context"
	authv1 "github.com/Sanchir01/auth-proto/gen/go/auth"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/Sanchir01/candles_backend/pkg/lib/logger/sl"
	"github.com/google/uuid"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"strings"

	"log/slog"
	"time"
)

type ClientAuthGRPC struct {
	api authv1.AuthClient
	log *slog.Logger
}

func NewClientAuthGRPC(l *slog.Logger, addr string, timeout time.Duration, retriesCount int) (*ClientAuthGRPC, error) {
	const op = "grpc.auth.new"
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithBackoff(grpcretry.BackoffLinear(500 * time.Millisecond)),
		grpcretry.WithCodes(codes.Aborted, codes.Unavailable, codes.NotFound, codes.DeadlineExceeded),
		grpcretry.WithPerRetryTimeout(timeout),
		grpcretry.WithMax(uint(retriesCount)),
	}
	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}
	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
			grpclog.UnaryClientInterceptor(sl.InterceptorsLogger(l), logOpts...),
		),
	)
	if err != nil {
		return nil, err
	}

	return &ClientAuthGRPC{
		api: authv1.NewAuthClient(cc),
		log: l,
	}, nil
}
func (gr *ClientAuthGRPC) Login(ctx context.Context, email, password string) (*model.User, error) {
	const op = "grpc.auth.login"
	resp, err := gr.api.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	userrole := strings.ToLower(resp.GetRole().String())
	return &model.User{
		Email:     resp.GetEmail(),
		Phone:     resp.GetPhone(),
		Title:     resp.GetTitle(),
		Role:      userrole,
		Version:   1,
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Password:  "",
	}, err
}
