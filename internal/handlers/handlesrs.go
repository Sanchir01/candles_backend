package httphandlers

import (
	"github.com/Sanchir01/candles_backend/internal/config"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

type HttpRouter struct {
	chiRouter *chi.Mux
	logger    *slog.Logger
	config    *config.Config
}

func NewChiRouter(r *chi.Mux, lg *slog.Logger, cfg *config.Config) *HttpRouter {
	return &HttpRouter{chiRouter: r, logger: lg, config: cfg}
}
