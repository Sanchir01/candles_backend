package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/app"
	telegrambot "github.com/Sanchir01/candles_backend/internal/bot"
	pgstoreauth "github.com/Sanchir01/candles_backend/internal/database/postgres/auth"
	pgstorecandles "github.com/Sanchir01/candles_backend/internal/database/postgres/candles"
	pgstorecategory "github.com/Sanchir01/candles_backend/internal/database/postgres/category"
	pgstorecolor "github.com/Sanchir01/candles_backend/internal/database/postgres/color"
	pgstoreuser "github.com/Sanchir01/candles_backend/internal/database/postgres/user"
	s3store "github.com/Sanchir01/candles_backend/internal/database/s3"
	httphandlers "github.com/Sanchir01/candles_backend/internal/handlers"
	httpserver "github.com/Sanchir01/candles_backend/internal/server/http"
	"github.com/go-chi/chi/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vikstrous/dataloadgen"
	"log"
	"time"

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
		categoryStr = pgstorecategory.New(env.DataBase.PrimaryDB)
		candlesStr  = pgstorecandles.New(env.DataBase.PrimaryDB)
		colorStr    = pgstorecolor.New(env.DataBase.PrimaryDB)
		userStr     = pgstoreuser.New(env.DataBase.PrimaryDB)
		authStr     = pgstoreauth.New(env.DataBase.PrimaryDB)
		s3str       = s3store.New(env.Storages.ImageStorage, context.Background(), env.Config)
		handlers    = httphandlers.New(rout, env, s3str, categoryStr, candlesStr, colorStr, userStr, authStr)
	)
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
	loader := dataloadgen.NewLoader(func(ctx context.Context, keys []string) ([]string, []error) {
		items := make([]string, len(keys))
		errs := make([]error, len(keys))

		for i, key := range keys {
			items[i] = key
			if key == "errorKey" {
				errs[i] = fmt.Errorf("произошла ошибка с ключом: %s", key)
			} else {
				errs[i] = nil
			}
		}
		return items, errs
	},
		dataloadgen.WithBatchCapacity(10),
		dataloadgen.WithWait(2*time.Second),
	)
	myStrings := []string{"123", "213123"}
	one, err := loader.LoadAll(ctx, myStrings)
	if err != nil {
		panic(err)
	}
	env.Logger.Warn("loader user", one)
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
