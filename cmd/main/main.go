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
	_ "gopkg.in/gomail.v2"
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

	// m := gomail.NewMessage()
	// m.SetHeader("From", "emgushovs.ru")
	// m.SetHeader("To", "emgushovs@mail.ru")
	// m.SetHeader("Subject", "My letter subject")
	// m.SetBody("text/html", "<html>This is the HTML message body</html>")

	// d := gomail.NewDialer(
	// 	"connect.smtp.bz",
	// 	465,
	// 	"emgushovs@mail.ru",
	// 	"tWYoeqsXUB1A",
	// )

	// d.SSL = true

	// if err := d.DialAndSend(m); err != nil {
	// 	panic(err)
	// }

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
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
	go func() { env.GRPCSrv.MustRun() }()
	<-ctx.Done()

	if err := serve.Gracefull(ctx); err != nil {
		env.Logger.Error("Graphql serve gracefull")
	}
	env.GRPCSrv.Stop()
}
