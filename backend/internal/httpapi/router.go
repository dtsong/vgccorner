package httpapi

import (
	"net/http"

	"github.com/dtsong/vgccorner/backend/internal/db"
	"github.com/dtsong/vgccorner/backend/internal/observability"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	logger *observability.Logger
	db     *db.Database
}

func NewRouter(logger *observability.Logger, database *db.Database) http.Handler {
	s := &Server{logger: logger, db: database}

	r := chi.NewRouter()

	// Health check endpoint
	r.Get("/healthz", s.handleHealth)

	// Showdown analysis endpoints
	r.Post("/api/showdown/analyze", s.handleAnalyzeShowdown)
	r.Get("/api/showdown/replays", s.handleListShowdownReplays)
	r.Get("/api/showdown/replays/{replayId}", s.handleGetShowdownReplay)
	r.Get("/api/showdown/replays/{replayId}/turns", s.handleGetTurnAnalysis)

	// TCG Live endpoint (planned)
	r.Post("/api/tcglive/analyze", s.handleAnalyzeTCGLive)

	return r
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
