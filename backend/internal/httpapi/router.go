package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/dtsong/battleforge/backend/internal/observability"
)

type Server struct {
	logger *observability.Logger
}

func NewRouter(logger *observability.Logger) http.Handler {
	s := &Server{logger: logger}

	r := chi.NewRouter()

	// Basic health check
	r.Get("/healthz", s.handleHealth)

	// Showdown analysis endpoints
	r.Post("/api/showdown/analyze", s.handleAnalyzeShowdown)
	r.Get("/api/showdown/battles/{battleID}", s.handleGetShowdownBattle)

	// TCG Live endpoint (planned)
	r.Post("/api/tcglive/analyze", s.handleAnalyzeTCGLive)

	return r
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
