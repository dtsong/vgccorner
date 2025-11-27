package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dtsong/vgccorner/backend/internal/analysis"
	"github.com/dtsong/vgccorner/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

// It supports three analysis types via discriminator: replayId, username, or rawLog.
type AnalyzeShowdownRequest struct {
	// Discriminator field
	AnalysisType string `json:"analysisType"`

	// For replayId analysis
	ReplayID string `json:"replayId,omitempty"`

	// For username analysis
	Username string `json:"username,omitempty"`
	Format   string `json:"format,omitempty"`
	Limit    int    `json:"limit,omitempty"`

	// For rawLog analysis
	RawLog string `json:"rawLog,omitempty"`

	// Common field
	IsPrivate bool `json:"isPrivate"`
}

// AnalyzeResponse is the response for analyze requests.
type AnalyzeResponse struct {
	Status   string                  `json:"status"`
	BattleID string                  `json:"battleId,omitempty"`
	Data     *analysis.BattleSummary `json:"data,omitempty"`
	Metadata *ResponseMetadata       `json:"metadata,omitempty"`
}

// ResponseMetadata contains metadata about the analysis.
type ResponseMetadata struct {
	ParseTimeMs    int  `json:"parseTimeMs"`
	AnalysisTimeMs int  `json:"analysisTimeMs"`
	Cached         bool `json:"cached"`
}

// ErrorResponse is the response for errors.
type ErrorResponse struct {
	Error   string      `json:"error"`
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

// ListReplaysRequest represents query parameters for listing replays.
type ListReplaysRequest struct {
	Username  string
	Format    string
	IsPrivate *bool
	Limit     int
	Offset    int
}

// handleAnalyzeShowdown handles POST /api/showdown/analyze requests.
func (s *Server) handleAnalyzeShowdown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	start := time.Now()

	var req AnalyzeShowdownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Infof("Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// Validate request based on analysis type
	var battleSummary *analysis.BattleSummary
	var battlelLog string
	var err error

	switch req.AnalysisType {
	case "replayId":
		if req.ReplayID == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Error: "replayId is required for replayId analysis",
				Code:  "INVALID_REQUEST",
			})
			return
		}
		// TODO: Fetch replay from Showdown API or cache
		w.WriteHeader(http.StatusNotImplemented)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Replay ID analysis not yet implemented",
			Code:  "NOT_IMPLEMENTED",
		})
		return

	case "username":
		if req.Username == "" || req.Format == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Error: "username and format are required for username analysis",
				Code:  "INVALID_REQUEST",
			})
			return
		}
		// TODO: Fetch recent battles by username
		w.WriteHeader(http.StatusNotImplemented)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Username analysis not yet implemented",
			Code:  "NOT_IMPLEMENTED",
		})
		return

	case "rawLog":
		if req.RawLog == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Error: "rawLog is required for rawLog analysis",
				Code:  "INVALID_REQUEST",
			})
			return
		}
		battlelLog = req.RawLog

	default:
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "analysisType must be one of: replayId, username, rawLog",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// Parse battle log with enhanced turn tracking
	parseStart := time.Now()
	battleSummary, err = analysis.ParseEnhancedShowdownLog(battlelLog)
	parseTime := time.Since(parseStart).Milliseconds()

	if err != nil {
		s.logger.Infof("Failed to parse battle log: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Failed to parse battle log: " + err.Error(),
			Code:  "PARSE_ERROR",
		})
		return
	}

	// Store battle in database (if database is configured)
	battleID := battleSummary.ID
	if s.db != nil {
		ctx := r.Context()
		battleRecord := &db.Battle{
			ID:          battleSummary.ID,
			Format:      battleSummary.Format,
			Timestamp:   battleSummary.Timestamp,
			DurationSec: battleSummary.Duration,
			Winner:      battleSummary.Winner,
			Player1ID:   battleSummary.Player1.Name,
			Player2ID:   battleSummary.Player2.Name,
			BattleLog:   battlelLog,
			IsPrivate:   req.IsPrivate,
			Analysis:    convertBattleStats(battleSummary),
			KeyMoments:  convertKeyMoments(battleSummary),
		}

		// Store battle and basic analysis
		storedID, err := s.db.StoreBattle(ctx, battleRecord)
		if err != nil {
			s.logger.Infof("Failed to store battle: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Error: "Failed to store battle",
				Code:  "INTERNAL_ERROR",
			})
			return
		}
		battleID = storedID

		// Store detailed turn-by-turn data
		if err := s.db.StoreTurnData(ctx, battleID, battleSummary); err != nil {
			s.logger.Infof("Failed to store turn data: %v", err)
			// Don't fail the request, just log the error
		}
	}

	analysisTime := time.Since(start).Milliseconds()

	s.logger.Infof("Successfully analyzed Showdown battle: %s (Player1: %s, Player2: %s)",
		battleSummary.ID, battleSummary.Player1.TeamArchetype, battleSummary.Player2.TeamArchetype)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(AnalyzeResponse{
		Status:   "success",
		BattleID: battleID,
		Data:     battleSummary,
		Metadata: &ResponseMetadata{
			ParseTimeMs:    int(parseTime),
			AnalysisTimeMs: int(analysisTime),
			Cached:         false,
		},
	})
}

// convertBattleStats converts analysis stats to database format
func convertBattleStats(summary *analysis.BattleSummary) *db.BattleAnalysis {
	return &db.BattleAnalysis{
		TotalTurns:            summary.Stats.TotalTurns,
		AvgDamagePerTurn:      summary.Stats.AvgDamagePerTurn,
		AvgHealPerTurn:        summary.Stats.AvgHealPerTurn,
		MovesUsedCount:        summary.Stats.Player1Stats.MoveCount + summary.Stats.Player2Stats.MoveCount,
		SwitchesCount:         summary.Stats.Player1Stats.SwitchCount + summary.Stats.Player2Stats.SwitchCount,
		SuperEffectiveMoves:   summary.Stats.SuperEffective,
		NotVeryEffectiveMoves: summary.Stats.NotVeryEffective,
		CriticalHits:          summary.Stats.CriticalHits,
		Player1DamageDealt:    summary.Stats.Player1Stats.DamageDealt,
		Player1DamageTaken:    summary.Stats.Player1Stats.DamageTaken,
		Player1HealingDone:    summary.Stats.Player1Stats.HealingDone,
		Player2DamageDealt:    summary.Stats.Player2Stats.DamageDealt,
		Player2DamageTaken:    summary.Stats.Player2Stats.DamageTaken,
		Player2HealingDone:    summary.Stats.Player2Stats.HealingDone,
	}
}

// convertKeyMoments converts key moments to database format
func convertKeyMoments(summary *analysis.BattleSummary) []*db.KeyMoment {
	moments := make([]*db.KeyMoment, 0, len(summary.KeyMoments))
	for _, km := range summary.KeyMoments {
		moments = append(moments, &db.KeyMoment{
			TurnNumber:   km.TurnNumber,
			MomentType:   km.Type,
			Description:  km.Description,
			Significance: km.Significance,
		})
	}
	return moments
}

// handleGetShowdownReplay handles GET /api/showdown/replays/{replayId} requests.
func (s *Server) handleGetShowdownReplay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	battleID := chi.URLParam(r, "replayId")

	if battleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "replayId is required",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	s.logger.Infof("Retrieving replay: %s", battleID)

	// Database required for this endpoint
	if s.db == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Database not configured",
			Code:  "SERVICE_UNAVAILABLE",
		})
		return
	}

	ctx := r.Context()
	battle, err := s.db.GetBattle(ctx, battleID)
	if err != nil {
		s.logger.Infof("Failed to retrieve battle: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	if battle == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Replay not found",
			Code:  "NOT_FOUND",
		})
		return
	}

	// Parse the battle log to get full summary
	summary, err := analysis.ParseEnhancedShowdownLog(battle.BattleLog)
	if err != nil {
		s.logger.Infof("Failed to parse battle log: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Failed to parse battle log",
			Code:  "PARSE_ERROR",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(AnalyzeResponse{
		Status:   "success",
		BattleID: battle.ID,
		Data:     summary,
	})
}

// handleListShowdownReplays handles GET /api/showdown/replays requests.
func (s *Server) handleListShowdownReplays(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	username := r.URL.Query().Get("username")
	format := r.URL.Query().Get("format")
	isPrivateStr := r.URL.Query().Get("isPrivate")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var isPrivate *bool
	if isPrivateStr != "" {
		val := isPrivateStr == "true"
		isPrivate = &val
	}

	limit := 10
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	offset := 0
	if offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
			offset = v
		}
	}

	s.logger.Infof("Listing replays: username=%s format=%s isPrivate=%v limit=%d offset=%d", username, format, isPrivate, limit, offset)

	// Database required for this endpoint
	if s.db == nil {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   []interface{}{},
			"pagination": map[string]int{
				"limit":  limit,
				"offset": offset,
				"total":  0,
			},
		})
		return
	}

	ctx := r.Context()
	filter := &db.BattleFilter{
		Format:    format,
		IsPrivate: isPrivate,
	}
	battles, total, err := s.db.ListBattles(ctx, filter, limit, offset)
	if err != nil {
		s.logger.Infof("Failed to list battles: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   battles,
		"pagination": map[string]int{
			"limit":  limit,
			"offset": offset,
			"total":  total,
		},
	})
}
