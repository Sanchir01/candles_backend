package ordergrpc

import (
	"context"
	orderv1 "github.com/makiyarodzy/mahakala_proto/gen/go/order"
	"google.golang.org/grpc"
)

type serverApi struct {
	orderv1.UnimplementedOrderServer
}

func RegisterOrderServerApi(gRPC *grpc.Server) {
	orderv1.RegisterOrderServer(gRPC, &serverApi{})
}

func (s *serverApi) TelegramOrder(
	ctx context.Context, req *orderv1.OrderRequest,
) (*orderv1.OrderResponse, error) {
	panic("implement me")
}
