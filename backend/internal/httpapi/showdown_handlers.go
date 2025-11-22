package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/dtsong/battleforge/backend/internal/analysis"
)

type AnalyzeShowdownRequest struct {
	BattleLog string `json:"battleLog"`
	Metadata  struct {
		Format     string `json:"format"`
		UploadedAt string `json:"uploadedAt"`
	} `json:"metadata"`
}

type AnalyzeResponse struct {
	BattleID string                `json:"battleId"`
	Status   string                `json:"status"`
	Error    string                `json:"error,omitempty"`
	Data     *analysis.BattleSummary `json:"data,omitempty"`
}

func (s *Server) handleAnalyzeShowdown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req AnalyzeShowdownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Infof("Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(AnalyzeResponse{
			Status: "error",
			Error:  "Invalid request body",
		})
		return
	}

	if req.BattleLog == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(AnalyzeResponse{
			Status: "error",
			Error:  "battleLog is required",
		})
		return
	}

	// Parse the battle log
	battleSummary, err := analysis.ParseShowdownLog(req.BattleLog)
	if err != nil {
		s.logger.Infof("Failed to parse battle log: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(AnalyzeResponse{
			Status: "error",
			Error:  "Failed to parse battle log: " + err.Error(),
		})
		return
	}

	s.logger.Infof("Successfully analyzed Showdown battle: %s", battleSummary.ID)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(AnalyzeResponse{
		BattleID: battleSummary.ID,
		Status:   "success",
		Data:     battleSummary,
	})
}

func (s *Server) handleGetShowdownBattle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	battleID := chi.URLParam(r, "battleID")

	if battleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(AnalyzeResponse{
			Status: "error",
			Error:  "battleID is required",
		})
		return
	}

	s.logger.Infof("Retrieving battle: %s", battleID)

	// TODO: Retrieve battle from database
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(AnalyzeResponse{
		Status: "error",
		Error:  "Battle not found",
	})
}
