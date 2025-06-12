package grpcclientorder

import (
	"context"
	orderv1 "github.com/Sanchir01/auth-proto/gen/go/order"
	"github.com/Sanchir01/candles_backend/pkg/lib/logger/sl"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type ClientOrderGRPC struct {
	api orderv1.OrderClient
	log *slog.Logger
}

func NewClientOrderGRPC(l *slog.Logger, addr string, timeout time.Duration, retriesCount int) (*ClientOrderGRPC, error) {
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

	return &ClientOrderGRPC{
		api: orderv1.NewOrderClient(cc),
		log: l,
	}, nil
}

func (gr *ClientOrderGRPC) OrderPush(ctx context.Context, order *orderv1.OrderRequest) (*orderv1.OrderResponse, error) {
	const op = "grpc.order.push"
	resp, err := gr.api.TelegramOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
