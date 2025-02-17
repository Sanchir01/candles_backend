package grpcserver

import (
	"fmt"
	ordergrpc "github.com/Sanchir01/candles_backend/internal/grpc"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewGrpc(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	ordergrpc.RegisterOrderServerApi(gRPCServer)
	return &App{
		gRPCServer: gRPCServer,
		port:       port,
		log:        log,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op), slog.Int("port", a.port),
	)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {

		panic(err)
	}
}
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(
		slog.String("op", op)).Info("stopping grpc server")

	a.gRPCServer.GracefulStop()
}
