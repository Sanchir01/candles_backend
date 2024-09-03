package resolver

import (
	pgstorecandles "github.com/Sanchir01/candles_backend/internal/database/postgres/candles"
	pgstorecategory "github.com/Sanchir01/candles_backend/internal/database/postgres/category"
	"log/slog"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	lg          *slog.Logger
	categoryStr *pgstorecategory.CategoryPostgresStore
	candlesStr  *pgstorecandles.CandlesPostgresStore
}

func New(category *pgstorecategory.CategoryPostgresStore, candles *pgstorecandles.CandlesPostgresStore, lg *slog.Logger) *Resolver {
	return &Resolver{categoryStr: category, lg: lg, candlesStr: candles}
}
