package resolver

import (
	"github.com/Sanchir01/candles_backend/internal/config"
	pgstoreauth "github.com/Sanchir01/candles_backend/internal/database/postgres/auth"
	pgstorecandles "github.com/Sanchir01/candles_backend/internal/database/postgres/candles"
	pgstorecategory "github.com/Sanchir01/candles_backend/internal/database/postgres/category"
	pgstorecolor "github.com/Sanchir01/candles_backend/internal/database/postgres/color"
	pgstoreuser "github.com/Sanchir01/candles_backend/internal/database/postgres/user"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	s3store     *s3.Client
	cfg         *config.Config
	userStr     *pgstoreuser.UserPostgresStore
}

func New(
	category *pgstorecategory.CategoryPostgresStore, candles *pgstorecandles.CandlesPostgresStore,
	color *pgstorecolor.ColorPostgresStore, auth *pgstoreauth.AuthPostgresStore,
	lg *slog.Logger, user *pgstoreuser.UserPostgresStore, config *config.Config, s3store *s3.Client,
) *Resolver {
	return &Resolver{categoryStr: category, lg: lg, candlesStr: candles, userStr: user, colorStr: color, authStr: auth, cfg: config, s3store: s3store}
}
