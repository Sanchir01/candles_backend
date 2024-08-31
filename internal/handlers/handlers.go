package httphandlers

import (
	"context"
	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Sanchir01/candles_backend/internal/config"
	pgstoreCategory "github.com/Sanchir01/candles_backend/internal/database/postgres/category"
	"github.com/Sanchir01/candles_backend/internal/gql/directive"
	genGql "github.com/Sanchir01/candles_backend/internal/gql/generated"
	resolver "github.com/Sanchir01/candles_backend/internal/gql/resolvers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
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
	category  *pgstoreCategory.CategoryPostgresStore
}

const (
	maxUploadSize                       = 30 * 1024 * 1024
	queryCacheLRUSize                   = 1000
	complexityLimit                     = 1000
	automaticPersistedQueryCacheLRUSize = 100
)

func New(r *chi.Mux, lg *slog.Logger, cfg *config.Config,
	category *pgstoreCategory.CategoryPostgresStore,
) *HttpRouter {
	return &HttpRouter{chiRouter: r, logger: lg, config: cfg, category: category}
}

func (r *HttpRouter) StartHttpServer() http.Handler {
	r.newChiCors()
	r.chiRouter.Use(middleware.RequestID)
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
	srv.SetQueryCache(lru.New(queryCacheLRUSize))
	srv.AddTransport(transport.MultipartForm{
		MaxUploadSize: maxUploadSize,
		MaxMemory:     maxUploadSize / 10,
	})
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New(automaticPersistedQueryCacheLRUSize)})

	srv.SetRecoverFunc(
		func(ctx context.Context, err interface{}) (userMessage error) {
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
		r.category, r.logger,
	)}
	cfg.Directives.InputUnion = directive.NewInputUnionDirective()
	cfg.Directives.SortRankInput = directive.NewSortRankInputDirective()

	return cfg
}
func (r *HttpRouter) newChiCors() {
	r.chiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:           300,
	}))
}
