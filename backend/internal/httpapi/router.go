package httpapi

import (
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

	// Placeholder for Showdown analysis endpoint
	r.Post("/api/showdown/analyze", s.handleAnalyzeShowdown)

	// Placeholder for TCG Live endpoint
	r.Post("/api/tcglive/analyze", s.handleAnalyzeTCGLive)

	return r
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
