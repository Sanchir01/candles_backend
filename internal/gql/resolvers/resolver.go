package resolver

import (
	pgstoreauth "github.com/Sanchir01/candles_backend/internal/database/postgres/auth"
	pgstorecandles "github.com/Sanchir01/candles_backend/internal/database/postgres/candles"
	pgstorecategory "github.com/Sanchir01/candles_backend/internal/database/postgres/category"
	pgstorecolor "github.com/Sanchir01/candles_backend/internal/database/postgres/color"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	lg          *slog.Logger
	categoryStr *pgstorecategory.CategoryPostgresStore
	candlesStr  *pgstorecandles.CandlesPostgresStore
	colorStr    *pgstorecolor.ColorPostgresStore
	authStr     *pgstoreauth.AuthPostgresStore
	pgxdb       *pgxpool.Pool
}

func New(category *pgstorecategory.CategoryPostgresStore, candles *pgstorecandles.CandlesPostgresStore, color *pgstorecolor.ColorPostgresStore, auth *pgstoreauth.AuthPostgresStore,
	lg *slog.Logger, pgxdb *pgxpool.Pool) *Resolver {
	return &Resolver{categoryStr: category, lg: lg, candlesStr: candles, pgxdb: pgxdb, colorStr: color, authStr: auth}
}
