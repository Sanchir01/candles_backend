package app

import (
	"github.com/Sanchir01/candles_backend/pkg/lib/logger/handlers/slogpretty"
	"log/slog"
	"os"
)

var (
	development = "development"
	production  = "production"
)

func setupLogger(env string) *slog.Logger {
	var lg *slog.Logger
	switch env {
	case production:
		lg = setupPrettySlog()
	case development:
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
