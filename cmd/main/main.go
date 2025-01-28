package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sanchir01/candles_backend/internal/app"
	httphandlers "github.com/Sanchir01/candles_backend/internal/handlers"
	httpserver "github.com/Sanchir01/candles_backend/internal/server/http"
	"github.com/go-chi/chi/v5"
)

func main() {
	env, err := app.NewEnv()
	if err != nil {
		panic(err)
	}
	serve := httpserver.NewHttpServer(env.Config)
	rout := chi.NewRouter()

	var (
		handlers = httphandlers.New(rout, env)
	)
	env.Logger.Info("start server", slog.String("port", env.Config.Port))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()

	go func(ctx context.Context) {
		if err := serve.Run(handlers.StartHttpServer()); err != nil {
			if !errors.Is(err, context.Canceled) {
				env.Logger.Error("Listen server error", slog.String("error", err.Error()))
				return
			}

		}
	}(ctx)
	go func(ctx context.Context) { env.GRPCSrv.MustRun() }(ctx)
	if err := env.Bot.Start(ctx); err != nil {
		env.Logger.Error("error for get updates bot")
	}
	if err := serve.Gracefull(ctx); err != nil {
		env.Logger.Error("Graphql serve gracefull")
	}
	env.GRPCSrv.Stop()
}
