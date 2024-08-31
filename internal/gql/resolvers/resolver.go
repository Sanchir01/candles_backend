package resolver

import (
	pgstoreCategory "github.com/Sanchir01/candles_backend/internal/database/postgres/category"
	"log/slog"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	lg       *slog.Logger
	category *pgstoreCategory.CategoryPostgresStore
}

func New(category *pgstoreCategory.CategoryPostgresStore, lg *slog.Logger) *Resolver {
	return &Resolver{category: category, lg: lg}
}
