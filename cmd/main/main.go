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
	serve := httpserver.NewHTTPServer(env.Config.Host, env.Config.Port,
		env.Config.Timeout, env.Config.IdleTimeout)
	rout := chi.NewRouter()

	prometheusserver := httpserver.NewHTTPServer(env.Config.Prometheus.Host, env.Config.Prometheus.Port, env.Config.Prometheus.Timeout,
		env.Config.Prometheus.IdleTimeout)
	handlers := httphandlers.New(rout, env)
	env.Logger.Info("start server", slog.String("port", env.Config.Port))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()

	go func() {
		if err := serve.Run(handlers.StartHttpServer()); err != nil {
			if !errors.Is(err, context.Canceled) {
				env.Logger.Error("Listen server error", slog.String("error", err.Error()))
				return
			}
		}
	}()
	go func() {
		if err := prometheusserver.Run(handlers.StartPrometheusHandlers()); err != nil {
			if !errors.Is(err, context.Canceled) {
				env.Logger.Error("Listen prometheus server error", slog.String("error", err.Error()))
				return
			}
		}
	}()
	go func() { env.GRPCSrv.MustRun() }()
	<-ctx.Done()
	if err := serve.Gracefull(ctx); err != nil {
		env.Logger.Error("Graphql serve gracefull")
	}
	env.GRPCSrv.Stop()
}
