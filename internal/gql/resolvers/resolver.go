package resolver

import (
	"github.com/Sanchir01/candles_backend/internal/app"
	pgstoreauth "github.com/Sanchir01/candles_backend/internal/database/postgres/auth"
	pgstorecandles "github.com/Sanchir01/candles_backend/internal/database/postgres/candles"
	pgstorecategory "github.com/Sanchir01/candles_backend/internal/database/postgres/category"
	pgstorecolor "github.com/Sanchir01/candles_backend/internal/database/postgres/color"
	pgstoreuser "github.com/Sanchir01/candles_backend/internal/database/postgres/user"
	s3store "github.com/Sanchir01/candles_backend/internal/database/s3"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	categoryStr *pgstorecategory.CategoryPostgresStore
	candlesStr  *pgstorecandles.CandlesPostgresStore
	colorStr    *pgstorecolor.ColorPostgresStore
	s3store     *s3store.S3Store
	authStr     *pgstoreauth.AuthStorePosrtgres
	userStr     *pgstoreuser.UserPostgresStore
	env         *app.Env
}

func New(
	category *pgstorecategory.CategoryPostgresStore, candles *pgstorecandles.CandlesPostgresStore,
	color *pgstorecolor.ColorPostgresStore, user *pgstoreuser.UserPostgresStore, s3store *s3store.S3Store,
	authStr *pgstoreauth.AuthStorePosrtgres, env *app.Env,
) *Resolver {
	return &Resolver{
		categoryStr: category, candlesStr: candles, userStr: user, colorStr: color,
		s3store: s3store, authStr: authStr, env: env,
	}
}
