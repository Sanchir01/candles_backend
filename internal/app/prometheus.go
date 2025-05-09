package app

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ListenPrometheus(env *Env) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(env.Config.Prometheus.Host+":"+env.Config.Prometheus.Port, mux)
}
