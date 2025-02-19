package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	apiprotov1 "github.com/makiyarodzy/mahakala_proto/gen/go/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	apiv1 apiprotov1.OrderClient
	log   *slog.Logger
}

func NewGRPCClient(ctx context.Context, addr string, logger *slog.Logger, retryCount int, timeout time.Duration) (*GRPCClient, error) {
	const op = "grpc.client.NewGRPCClient"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Aborted, codes.NotFound, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retryCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}
	logOpt := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}
	cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(logger), logOpt...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &GRPCClient{
		apiv1: apiprotov1.NewOrderClient(cc),
	}, nil
}

func (c *GRPCClient) IsOrderCreated(ctx context.Context, products []*apiprotov1.OrderMessage) (string, error) {
	const op = "grpc.order.create"

	resp, err := c.apiv1.TelegramOrder(ctx, &apiprotov1.OrderRequest{
		Orders: products,
	})
	if err != nil {
		return "", err
	}
	return resp.Ok, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
