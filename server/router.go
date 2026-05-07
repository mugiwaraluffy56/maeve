package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RouterConfig struct {
	Version string
}

func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Get("/api/v1/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"` + cfg.Version + `"}`))
	})
	return r
}
