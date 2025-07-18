package httphandlers

import (
	"fmt"
	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Sanchir01/candles_backend/internal/app"
	"github.com/Sanchir01/candles_backend/internal/gql/directive"
	genGql "github.com/Sanchir01/candles_backend/internal/gql/generated"
	resolver "github.com/Sanchir01/candles_backend/internal/gql/resolvers"
	customMiddleware "github.com/Sanchir01/candles_backend/internal/handlers/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vektah/gqlparser/v2/ast"
	"net/http"
)

type HttpRouter struct {
	chiRouter *chi.Mux
	env       *app.Env
}

const (
	maxUploadSize                       = 30 * 1024 * 1024
	queryCacheLRUSize                   = 1000
	complexityLimit                     = 1000
	automaticPersistedQueryCacheLRUSize = 100
)

func New(
	r *chi.Mux, env *app.Env,
) *HttpRouter {
	return &HttpRouter{
		chiRouter: r,
		env:       env,
	}
}

func (r *HttpRouter) StartHttpServer(domain string) http.Handler {
	defer func() {
		if err := recover(); err != nil {
			r.env.Logger.Error("panic error", err)
		}
	}()
	r.newChiCors()
	r.chiRouter.Use(
		middleware.RequestID,
		customMiddleware.PrometheusMiddleware,
		customMiddleware.WithResponseWriter,
		customMiddleware.AuthMiddleware(domain),
		customMiddleware.RecoverMiddleware,
		customMiddleware.UserAgentMiddleware,
	)
	r.chiRouter.Handle("/graphql", playground.ApolloSandboxHandler("Candles", "/"))
	r.chiRouter.Handle("/", r.NewGraphQLHandler())
	r.chiRouter.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, User!"))
		fmt.Println("hello world endpoint finished")
	})
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

	srv.Use(extension.FixedComplexityLimit(complexityLimit))

	return srv
}

func (r *HttpRouter) newSchemaConfig() genGql.Config {
	cfg := genGql.Config{Resolvers: resolver.New(
		r.env,
	)}
	cfg.Directives.InputUnion = directive.NewInputUnionDirective()
	cfg.Directives.SortRankInput = directive.NewSortRankInputDirective()
	cfg.Directives.HasRole = directive.RoleDirective()
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

func (r *HttpRouter) StartPrometheusHandlers() http.Handler {
	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	return router
}
