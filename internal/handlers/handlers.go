package httphandlers

import (
	"context"
	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Sanchir01/candles_backend/internal/config"
	pgstoreauth "github.com/Sanchir01/candles_backend/internal/database/postgres/auth"
	pgstorecandles "github.com/Sanchir01/candles_backend/internal/database/postgres/candles"
	pgstoreCategory "github.com/Sanchir01/candles_backend/internal/database/postgres/category"
	pgstorecolor "github.com/Sanchir01/candles_backend/internal/database/postgres/color"
	pgstoreuser "github.com/Sanchir01/candles_backend/internal/database/postgres/user"
	s3store "github.com/Sanchir01/candles_backend/internal/database/s3"
	"github.com/Sanchir01/candles_backend/internal/gql/directive"
	genGql "github.com/Sanchir01/candles_backend/internal/gql/generated"
	resolver "github.com/Sanchir01/candles_backend/internal/gql/resolvers"
	customMiddleware "github.com/Sanchir01/candles_backend/internal/handlers/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"log"
	"log/slog"
	"net/http"
	"runtime"
)

type HttpRouter struct {
	chiRouter *chi.Mux
	logger    *slog.Logger
	config    *config.Config
	db        *sqlx.DB
	s3store   *s3store.S3Store
	pgxdb     *pgxpool.Pool
	authStr   *pgstoreauth.AuthStorePosrtgres
	category  *pgstoreCategory.CategoryPostgresStore
	candles   *pgstorecandles.CandlesPostgresStore
	color     *pgstorecolor.ColorPostgresStore
	userStr   *pgstoreuser.UserPostgresStore
}

const (
	maxUploadSize                       = 30 * 1024 * 1024
	queryCacheLRUSize                   = 1000
	complexityLimit                     = 1000
	automaticPersistedQueryCacheLRUSize = 100
)

func New(r *chi.Mux, lg *slog.Logger, cfg *config.Config,
	s3store *s3store.S3Store, pgxdb *pgxpool.Pool,
	category *pgstoreCategory.CategoryPostgresStore, candlesStr *pgstorecandles.CandlesPostgresStore, colorStr *pgstorecolor.ColorPostgresStore,
	userStr *pgstoreuser.UserPostgresStore, authStr *pgstoreauth.AuthStorePosrtgres,
) *HttpRouter {
	return &HttpRouter{
		chiRouter: r, logger: lg, config: cfg, category: category, color: colorStr,
		candles: candlesStr, userStr: userStr,
		s3store: s3store, pgxdb: pgxdb, authStr: authStr,
	}
}

func (r *HttpRouter) StartHttpServer() http.Handler {
	r.newChiCors()
	r.chiRouter.Use(middleware.RequestID)
	r.chiRouter.Use(customMiddleware.WithResponseWriter)
	r.chiRouter.Use(customMiddleware.AuthMiddleware())
	r.chiRouter.Handle("/graphql", playground.ApolloSandboxHandler("Candles", "/"))
	r.chiRouter.Handle("/", r.NewGraphQLHandler())
	return r.chiRouter
}

func (r *HttpRouter) NewGraphQLHandler() *gqlhandler.Server {
	srv := gqlhandler.New(
		genGql.NewExecutableSchema(r.newSchemaConfig()),
	)
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Options{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](queryCacheLRUSize))
	srv.AddTransport(transport.MultipartForm{
		MaxUploadSize: maxUploadSize,
		MaxMemory:     maxUploadSize / 10,
	})
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](automaticPersistedQueryCacheLRUSize)})

	srv.SetRecoverFunc(
		func(ctx context.Context, err interface{}) error {
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			log.Printf("Panic: %v\nStack: %s\n", err, buf[:n])

			return gqlerror.Errorf("internal server error graphql обработка паники")
		})
	srv.Use(extension.FixedComplexityLimit(complexityLimit))

	return srv
}

func (r *HttpRouter) newSchemaConfig() genGql.Config {
	cfg := genGql.Config{Resolvers: resolver.New(
		r.category, r.candles, r.color, r.logger,
		r.userStr, r.config, r.s3store, r.pgxdb,
		r.authStr,
	)}
	cfg.Directives.InputUnion = directive.NewInputUnionDirective()
	cfg.Directives.SortRankInput = directive.NewSortRankInputDirective()
	cfg.Directives.HasRole = directive.RoleDirective()
	return cfg
}
func (r *HttpRouter) newChiCors() {
	switch r.config.Env {
	case "development":
		r.chiRouter.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowCredentials: true,
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			MaxAge:           300,
		}))
	case "production":
		r.chiRouter.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://mahakala.ru", "http://mahakala.ru"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowCredentials: true,
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			MaxAge:           300,
		}))
	}
}
