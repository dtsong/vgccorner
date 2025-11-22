package httpapi

import (
	"encoding/json"
	"net/http"
)

type AnalyzeShowdownRequest struct {
	ReplayID  string `json:"replayId,omitempty"`
	ReplayURL string `json:"replayUrl,omitempty"`
}

type AnalyzeShowdownResponse struct {
	Message string `json:"message"`
	// later: add BattleSummary fields here
}

func (s *Server) handleAnalyzeShowdown(w http.ResponseWriter, r *http.Request) {
	var req AnalyzeShowdownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	s.logger.Infof("received showdown analyze request: replayId=%s replayUrl=%s", req.ReplayID, req.ReplayURL)

	resp := AnalyzeShowdownResponse{
		Message: "Showdown analysis not implemented yet",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
