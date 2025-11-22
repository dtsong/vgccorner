package httpapi

import (
	"net/http"
)

func (s *Server) handleAnalyzeTCGLive(w http.ResponseWriter, r *http.Request) {
	// TODO: implement TCG Live analysis
	http.Error(w, "TCG Live analysis not implemented yet", http.StatusNotImplemented)
}
