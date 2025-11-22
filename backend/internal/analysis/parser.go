package analysis

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

// ParseShowdownLog parses a Pokémon Showdown battle log and returns a BattleSummary.
// The log format is pipe-delimited with commands like:
// |move|p1a: Whimsicott|Tailwind|p1a: Whimsicott
// |switch|p1b: Typhlosion|Typhlosion-Hisui, L50, M|100\/100
func ParseShowdownLog(logContent string) (*BattleSummary, error) {
	lines := strings.Split(logContent, "\n")

	summary := &BattleSummary{
		ID:         generateUUID(),
		Timestamp:  time.Now(),
		Turns:      []Turn{},
		KeyMoments: []KeyMoment{},
		Stats:      BattleStats{},
	}

	// Parse metadata and teams
	teams := make(map[string][]Pokémon)
	playerNames := make(map[string]string)
	activeTeamSizes := make(map[string]int)

	// First pass: extract metadata and team information
	for _, line := range lines {
		if line == "" || !strings.HasPrefix(line, "|") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}

		command := parts[1]

		switch command {
		case "t:":
			// Timestamp
			if len(parts) > 2 {
				ts, err := time.Parse("2006-01-02", time.Unix(0, 0).String())
				_ = err // Handle parsing if needed
				summary.Timestamp = ts
			}

		case "gen":
			// Generation
			if len(parts) > 2 {
				summary.Format = fmt.Sprintf("Gen %s", parts[2])
			}

		case "tier":
			// Tier/Format name
			if len(parts) > 2 {
				summary.Format = strings.Join(parts[2:], "|")
			}

		case "player":
			// Player information: |player|p1|liamvgc1|giovanni|1487
			if len(parts) > 3 {
			playerID := parts[2]
			playerName := parts[3]
			playerNames[playerID] = playerName
			switch playerID {
			case "p1":
				summary.Player1.Name = playerName
			case "p2":
				summary.Player2.Name = playerName
			}
			}

		case "teamsize":
			// Team size: |teamsize|p1|4
			if len(parts) > 3 {
				playerID := parts[2]
				teamSize := parseTeamSize(parts[3])
				activeTeamSizes[playerID] = teamSize
			}

		case "poke":
			// Pokemon in team: |poke|p1|Ursaluna-Bloodmoon, L50, M|
			if len(parts) > 3 {
				playerID := parts[2]
				pokeStr := parts[3]
				poke := parsePokemon(pokeStr)
				teams[playerID] = append(teams[playerID], poke)
			}

		case "start":
			// Battle has started
		}
	}

	summary.Player1.Team = teams["p1"]
	summary.Player2.Team = teams["p2"]
	summary.Player1.TotalLeft = activeTeamSizes["p1"]
	summary.Player2.TotalLeft = activeTeamSizes["p2"]

	// Second pass: process battle actions and turns
	var currentTurn *Turn
	var turnNumber int

	for _, line := range lines {
		if line == "" || !strings.HasPrefix(line, "|") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}

		command := parts[1]

		switch command {
		case "turn":
			// New turn: |turn|N
			if currentTurn != nil {
				summary.Turns = append(summary.Turns, *currentTurn)
			}
			turnNumber = parseInt(parts[2])
			currentTurn = &Turn{
				TurnNumber:  turnNumber,
				Actions:     []Action{},
				DamageDealt: make(map[string]int),
				HealingDone: make(map[string]int),
			}

		case "move":
			// Move action: |move|p1a: Whimsicott|Tailwind|p1a: Whimsicott
			if currentTurn != nil && len(parts) >= 4 {
				action := parseMove(parts)
				currentTurn.Actions = append(currentTurn.Actions, action)
			}

		case "switch":
			// Switch action: |switch|p1b: Typhlosion|Typhlosion-Hisui, L50, M|100\/100
			if currentTurn != nil && len(parts) >= 4 {
				action := parseSwitch(parts)
				currentTurn.Actions = append(currentTurn.Actions, action)
			}

		case "faint":
			// Pokemon fainted: |faint|p1b: Typhlosion
			if currentTurn != nil {
				playerID := extractPlayerID(parts[2])
				switch playerID {
				case "player1":
					summary.Player1.Losses++
				case "player2":
					summary.Player2.Losses++
				}
			}
			addKeyMoment(summary, turnNumber, "KO", "Pokémon fainted", 8)

		case "win":
			// Battle won: |win|Heliosan
			if len(parts) > 2 {
				winner := parts[2]
				summary.Winner = playerToID(winner, playerNames)
			}

		case "damage", "heal":
			// Damage/Healing: |-damage|p1b: Dragapult|0 fnt
			// Extract HP changes for statistics
		}
	}

	// Add the last turn
	if currentTurn != nil {
		summary.Turns = append(summary.Turns, *currentTurn)
	}

	// Calculate statistics
	calculateStats(summary)

	return summary, nil
}

// Helper functions

func parsePokemon(pokeStr string) Pokémon {
	// Format: "Ursaluna-Bloodmoon, L50, M"
	parts := strings.Split(pokeStr, ",")
	name := strings.TrimSpace(parts[0])

	poke := Pokémon{
		ID:   strings.ToLower(strings.ReplaceAll(name, "-", "")),
		Name: name,
	}

	if len(parts) > 1 {
		levelStr := strings.TrimSpace(parts[1])
		poke.Level = parseInt(strings.TrimPrefix(levelStr, "L"))
	}

	if len(parts) > 2 {
		genderStr := strings.TrimSpace(parts[2])
		poke.Gender = genderStr
	}

	return poke
}

func parseMove(parts []string) Action {
	// |move|p1a: Whimsicott|Tailwind|p1a: Whimsicott
	playerID := extractPlayerID(parts[2])
	moveName := strings.TrimSpace(parts[3])

	return Action{
		Player:     playerID,
		ActionType: "move",
		Move: &Move{
			ID:   strings.ToLower(moveName),
			Name: moveName,
		},
	}
}

func parseSwitch(parts []string) Action {
	// |switch|p1b: Typhlosion|Typhlosion-Hisui, L50, M|100\/100
	playerID := extractPlayerID(parts[2])
	switchToPoke := strings.TrimSpace(parts[3])

	return Action{
		Player:     playerID,
		ActionType: "switch",
		SwitchTo:   switchToPoke,
	}
}

func extractPlayerID(fullID string) string {
	// Convert "p1a: Whimsicott" to "player1" or "p2b: Maushold" to "player2"
	if strings.HasPrefix(fullID, "p1") {
		return "player1"
	}
	return "player2"
}

func playerToID(playerName string, playerNames map[string]string) string {
	for id, name := range playerNames {
		if name == playerName {
			if id == "p1" {
				return "player1"
			}
			return "player2"
		}
	}
	return "player2" // Default to player2
}

func parseInt(s string) int {
	var result int
	_, _ = fmt.Sscanf(s, "%d", &result)
	return result
}

func parseTeamSize(s string) int {
	return parseInt(s)
}

func addKeyMoment(summary *BattleSummary, turnNumber int, mType, description string, significance int) {
	summary.KeyMoments = append(summary.KeyMoments, KeyMoment{
		TurnNumber:  turnNumber,
		Type:        mType,
		Description: description,
		Significance: significance,
	})
}

func calculateStats(summary *BattleSummary) {
	summary.Stats.TotalTurns = len(summary.Turns)
	summary.Stats.MoveFrequency = make(map[string]int)
	summary.Stats.TypeCoverage = make(map[string]int)
	summary.Stats.Player1Stats = PlayerStats{
		MovesByType:  make(map[string]int),
		Effectiveness: EffectivenessStats{},
	}
	summary.Stats.Player2Stats = PlayerStats{
		MovesByType:  make(map[string]int),
		Effectiveness: EffectivenessStats{},
	}

	// Count moves by frequency
	for _, turn := range summary.Turns {
		for _, action := range turn.Actions {
			if action.ActionType == "move" && action.Move != nil {
				summary.Stats.MoveFrequency[action.Move.ID]++
			}
		}
	}

	// TODO: Calculate damage/heal stats, type coverage, etc.
	// This will require parsing damage output from battle log
}

// generateUUID generates a simple UUID-like string using random bytes.
// For production, consider using google.golang.org/uuid or similar.
func generateUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("battle-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
