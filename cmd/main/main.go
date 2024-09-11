package main

import (
	"context"
	"errors"
	telegrambot "github.com/Sanchir01/candles_backend/internal/bot"
	"github.com/Sanchir01/candles_backend/internal/config"
	pgstoreauth "github.com/Sanchir01/candles_backend/internal/database/postgres/auth"
	pgstorecandles "github.com/Sanchir01/candles_backend/internal/database/postgres/candles"
	pgstorecategory "github.com/Sanchir01/candles_backend/internal/database/postgres/category"
	pgstorecolor "github.com/Sanchir01/candles_backend/internal/database/postgres/color"
	httphandlers "github.com/Sanchir01/candles_backend/internal/handlers"
	httpserver "github.com/Sanchir01/candles_backend/internal/server/http"
	"github.com/Sanchir01/candles_backend/pkg/lib/db/connect"
	"github.com/Sanchir01/candles_backend/pkg/lib/logger/handlers/slogpretty"
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

func main() {
	cfg := config.InitConfig()
	lg := setupLogger(cfg.Env)
	lg.Info("Graphql server starting up...", slog.String("port", cfg.HttpServer.Port))

	db := connect.PostgresCon(cfg, lg)
	defer db.Close()
	pgxdb, err := connect.PGXNew(cfg, lg, context.Background())
	if err != nil {
		lg.Error("pgx error connect", err.Error())
	}

	serve := httpserver.NewHttpServer(cfg)
	rout := chi.NewRouter()
	var (
		categoryStr = pgstorecategory.New(pgxdb)
		candlesStr  = pgstorecandles.New(pgxdb)
		colorStr    = pgstorecolor.New(pgxdb)
		authStr     = pgstoreauth.New(pgxdb)
		handlers    = httphandlers.New(rout, lg, cfg, categoryStr, candlesStr, colorStr, authStr, pgxdb)
	)
	callcat, err := categoryStr.AllCategories(context.Background())
	lg.Warn("categoru", callcat)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()

	go func(ctx context.Context) {
		if err := serve.Run(handlers.StartHttpServer()); err != nil {
			if !errors.Is(err, context.Canceled) {
				lg.Error("Listen server error", slog.String("error", err.Error()))
				return
			}
			lg.Error("Listen server error", slog.String("error", err.Error()))
		}
	}(ctx)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	tgbot := telegrambot.New(bot, lg)
	if err := tgbot.Start(cfg); err != nil {
		lg.Error("error for get updates bot")
	}

	if err := serve.Gracefull(ctx); err != nil {
		log.Fatalf("Graphql serve gracefull")
	}
}
func setupLogger(env string) *slog.Logger {
	var lg *slog.Logger
	switch env {
	case development:
		lg = setupPrettySlog()
	case production:
		lg = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return lg
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
