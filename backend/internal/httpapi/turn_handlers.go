package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/dtsong/vgccorner/backend/internal/analysis"
	"github.com/dtsong/vgccorner/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

// TurnAnalysisResponse represents the response for turn-by-turn analysis
type TurnAnalysisResponse struct {
	Status     string        `json:"status"`
	BattleID   string        `json:"battleId"`
	Format     string        `json:"format"`
	Player1    string        `json:"player1"`
	Player2    string        `json:"player2"`
	Winner     string        `json:"winner,omitempty"`
	Turns      []TurnData    `json:"turns"`
	Archetypes ArchetypeInfo `json:"archetypes"`
}

// TurnData represents detailed information about a single turn
type TurnData struct {
	TurnNumber int           `json:"turnNumber"`
	Events     []BattleEvent `json:"events"`
	BoardState BoardState    `json:"boardState"`
}

// BattleEvent represents a single event during a turn
type BattleEvent struct {
	Type       string `json:"type"`    // "move", "switch", "faint", etc.
	Pokemon    string `json:"pokemon"` // Pokemon performing action
	Action     string `json:"action"`  // Description of action
	Target     string `json:"target,omitempty"`
	Result     string `json:"result,omitempty"`
	Details    string `json:"details,omitempty"`
	PlayerSide string `json:"playerSide"` // "player1" or "player2"
}

// BoardState represents the state of the battle at a specific turn
type BoardState struct {
	Player1Active []ActivePokemon `json:"player1Active"`
	Player2Active []ActivePokemon `json:"player2Active"`
}

// ActivePokemon represents a Pokemon currently on the field
type ActivePokemon struct {
	Species  string `json:"species"`
	Nickname string `json:"nickname,omitempty"`
	Position int    `json:"position"`
	HP       int    `json:"hp"`
	MaxHP    int    `json:"maxHp"`
	Status   string `json:"status,omitempty"`
	IsLead   bool   `json:"isLead,omitempty"`
}

// ArchetypeInfo contains team archetype information
type ArchetypeInfo struct {
	Player1 PlayerArchetype `json:"player1"`
	Player2 PlayerArchetype `json:"player2"`
}

// PlayerArchetype contains archetype details for a player
type PlayerArchetype struct {
	Archetype   string   `json:"archetype"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// handleGetTurnAnalysis handles GET /api/showdown/replays/{replayId}/turns requests
func (s *Server) handleGetTurnAnalysis(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	replayID := chi.URLParam(r, "replayId")

	if replayID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "replayId is required",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	s.logger.Infof("Retrieving turn analysis for replay: %s", replayID)

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

	// Retrieve turn data from database
	turnData, err := s.db.GetTurnData(ctx, replayID)
	if err != nil {
		s.logger.Infof("Failed to retrieve turn data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	if turnData == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Battle not found",
			Code:  "NOT_FOUND",
		})
		return
	}

	// Convert to API response format
	response := convertTurnDataToResponse(turnData)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// convertTurnDataToResponse converts database TurnAnalysisData to API response
func convertTurnDataToResponse(data *db.TurnAnalysisData) TurnAnalysisResponse {
	turns := make([]TurnData, 0, len(data.Turns))

	for _, turn := range data.Turns {
		turnData := TurnData{
			TurnNumber: turn.TurnNumber,
			Events:     convertDBActionsToEvents(turn.Actions),
			BoardState: convertDBBoardState(turn.BoardState),
		}
		turns = append(turns, turnData)
	}

	return TurnAnalysisResponse{
		Status:   "success",
		BattleID: data.BattleID,
		Format:   data.Format,
		Player1:  data.Player1,
		Player2:  data.Player2,
		Winner:   data.Winner,
		Turns:    turns,
		Archetypes: ArchetypeInfo{
			Player1: convertArchetype(data.Player1Archetype),
			Player2: convertArchetype(data.Player2Archetype),
		},
	}
}

// convertDBActionsToEvents converts database actions to API events
func convertDBActionsToEvents(actions []*db.ActionData) []BattleEvent {
	events := make([]BattleEvent, 0, len(actions))

	for _, action := range actions {
		event := BattleEvent{
			Type:       action.ActionType,
			Pokemon:    action.Pokemon,
			Action:     formatActionDescription(action),
			Target:     action.Target,
			Result:     action.Result,
			Details:    action.Details,
			PlayerSide: action.Player,
		}

		events = append(events, event)
	}

	return events
}

// formatActionDescription formats an action into a readable description
func formatActionDescription(action *db.ActionData) string {
	switch action.ActionType {
	case "move":
		return "used " + action.Move
	case "switch":
		return "switched in"
	default:
		return action.ActionType
	}
}

// convertDBBoardState converts database board state to API format
func convertDBBoardState(state *db.BoardStateData) BoardState {
	boardState := BoardState{
		Player1Active: []ActivePokemon{},
		Player2Active: []ActivePokemon{},
	}

	for _, poke := range state.Player1Active {
		boardState.Player1Active = append(boardState.Player1Active, ActivePokemon{
			Species:  poke.Species,
			Nickname: poke.Name,
			Position: poke.Position,
			HP:       poke.HP,
			MaxHP:    poke.MaxHP,
			Status:   poke.Status,
			IsLead:   poke.IsLead,
		})
	}

	for _, poke := range state.Player2Active {
		boardState.Player2Active = append(boardState.Player2Active, ActivePokemon{
			Species:  poke.Species,
			Nickname: poke.Name,
			Position: poke.Position,
			HP:       poke.HP,
			MaxHP:    poke.MaxHP,
			Status:   poke.Status,
			IsLead:   poke.IsLead,
		})
	}

	return boardState
}

// convertArchetype converts database archetype to API format
func convertArchetype(archetype *db.TeamArchetypeData) PlayerArchetype {
	if archetype == nil {
		return PlayerArchetype{
			Archetype:   "Unclassified",
			Description: analysis.GetArchetypeDescription("Unclassified"),
			Tags:        []string{},
		}
	}

	return PlayerArchetype{
		Archetype:   archetype.Archetype,
		Description: archetype.Description,
		Tags:        archetype.Tags,
	}
}
