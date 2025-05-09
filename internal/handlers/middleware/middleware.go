package customMiddleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Sanchir01/candles_backend/internal/app"
	userFeature "github.com/Sanchir01/candles_backend/internal/feature/user"
	"github.com/prometheus/client_golang/prometheus"
)

const responseWriterKey = "responseWriter"

var requestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "mahakala-count",
		Subsystem: "http",
		Name:      "request_total",
		Help:      "Total number of HTTP requests",
	},
	[]string{"path", "method"},
)

var requesDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path"},
)

func GetJWTClaimsFromCtx(ctx context.Context) (*userFeature.Claims, error) {
	claims, ok := ctx.Value(app.AccessTokenContextKey).(*userFeature.Claims)
	if !ok {
		return nil, errors.New("no JWT claims found in context")
	}
	return claims, nil
}

func WithResponseWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), responseWriterKey, w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetResponseWriter(ctx context.Context) http.ResponseWriter {
	return ctx.Value(responseWriterKey).(http.ResponseWriter)
}

func AuthMiddleware(domain string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			access, err := r.Cookie("refreshToken")
			if err != nil {
				refresh, err := r.Cookie("accessToken")
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				accessToken, err := userFeature.NewAccessToken(refresh.Value, 0, w, domain)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				token, err := userFeature.ParseToken(accessToken)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				ctx := context.WithValue(r.Context(), app.AccessTokenContextKey, token)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			validAccessToken, err := userFeature.ParseToken(access.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), app.AccessTokenContextKey, validAccessToken)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func NewDataLoadersMiddleware(env *app.Env) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loaders := app.NewDataLoaders(env.Repositories)

			// Добавляем DataLoaders в контекст запроса
			ctx := context.WithValue(r.Context(), app.DataLoadersContextKey, loaders)

			// Передаем контекст дальше по цепочке
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()
		requestCount.WithLabelValues(r.URL.Path, r.Method).Inc()
		requesDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}
