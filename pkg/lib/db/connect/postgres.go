package connect

import (
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/config"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"os"
)

func PostgresCon(cfg *config.Config, lg *slog.Logger) *sqlx.DB {
	postgresString := fmt.Sprintf(
		"user=%s dbname=%s sslmode=%s password=%s port=%s host=%s",
		cfg.DB.User, cfg.DB.Database, cfg.DB.SSL, os.Getenv("PASSWORD_POSTGRES"),
		cfg.DB.Port, cfg.DB.Host,
	)

	db, err := sqlx.Open("postgres", postgresString)
	if err != nil {
		lg.Error("sqlx.Connect error", slog.String("error", err.Error()))
	}
	return db
}
