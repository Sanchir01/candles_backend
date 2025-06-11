package main

import (
	"context"
	"errors"
	"github.com/Sanchir01/candles_backend/internal/app"
	httphandlers "github.com/Sanchir01/candles_backend/internal/handlers"
	httpserver "github.com/Sanchir01/candles_backend/internal/server/http"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	env, err := app.NewEnv()
	if err != nil {
		panic(err)
	}
	serve := httpserver.NewHTTPServer(env.Config.Host, env.Config.Port,
		env.Config.Timeout, env.Config.IdleTimeout)
	rout := chi.NewRouter()
	//runtime.SetMutexProfileFraction(1)
	//runtime.SetBlockProfileRate(1)
	//_, err = pyroscope.Start(pyroscope.Config{
	//	ApplicationName: "candles.backend",
	//	ServerAddress:   "http://host.docker.internal:4040",
	//	Logger:          pyroscope.StandardLogger,
	//	ProfileTypes: []pyroscope.ProfileType{
	//		pyroscope.ProfileCPU,
	//		pyroscope.ProfileAllocObjects,
	//		pyroscope.ProfileAllocSpace,
	//		pyroscope.ProfileInuseObjects,
	//		pyroscope.ProfileInuseSpace,
	//		pyroscope.ProfileGoroutines,
	//		pyroscope.ProfileMutexCount,
	//		pyroscope.ProfileMutexDuration,
	//		pyroscope.ProfileBlockCount,
	//		pyroscope.ProfileBlockDuration,
	//	},
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}

	prometheusserver := httpserver.NewHTTPServer(env.Config.Prometheus.Host, env.Config.Prometheus.Port, env.Config.Prometheus.Timeout,
		env.Config.Prometheus.IdleTimeout)
	handlers := httphandlers.New(rout, env)
	env.Logger.Info("start server", slog.String("port", env.Config.Port))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	env.Services.EventService.StartCreateEvent(ctx, 5*time.Second, 10)
	defer cancel()
	defer env.KafkaProducer.Close()

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
	<-ctx.Done()

	if err := serve.Gracefull(ctx); err != nil {
		env.Logger.Error("Graphql serve gracefull")
	}

}
