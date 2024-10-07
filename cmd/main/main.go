package main

import (
	"context"
	"errors"
	"github.com/Sanchir01/candles_backend/internal/app"
	telegrambot "github.com/Sanchir01/candles_backend/internal/bot"
	httphandlers "github.com/Sanchir01/candles_backend/internal/handlers"
	httpserver "github.com/Sanchir01/candles_backend/internal/server/http"
	"github.com/go-chi/chi/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var (
	development = "development"
	production  = "production"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
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
			env.Logger.Error("Listen server error", slog.String("error", err.Error()))
		}
	}(ctx)

	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	tgbot := telegrambot.New(bot, env.Logger)

	if err := tgbot.Start(env.Config); err != nil {
		env.Logger.Error("error for get updates bot")
	}
	if err := serve.Gracefull(ctx); err != nil {
		env.Logger.Error("Graphql serve gracefull")
	}
}
