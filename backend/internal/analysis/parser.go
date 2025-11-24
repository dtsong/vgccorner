package analysis

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

// ParseShowdownLog parses a Pokémon Showdown battle log and returns a comprehensive BattleSummary.
func ParseShowdownLog(logContent string) (*BattleSummary, error) {
	lines := strings.Split(logContent, "\n")

	summary := &BattleSummary{
		ID:         generateUUID(),
		Timestamp:  time.Now(),
		Turns:      []Turn{},
		KeyMoments: []KeyMoment{},
		Stats:      BattleStats{},
	}

	// Create a state tracker to maintain battle state throughout
	tracker := NewStateTracker()

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
		case "tier":
			if len(parts) > 2 {
				summary.Format = strings.Join(parts[2:], "|")
			}

		case "player":
			if len(parts) > 3 {
				playerID := parts[2]
				playerName := parts[3]
				tracker.SetPlayerName(playerID, playerName)
				switch playerID {
				case "p1":
					summary.Player1.Name = playerName
				case "p2":
					summary.Player2.Name = playerName
				}
			}

		case "teamsize":
			if len(parts) > 3 {
				playerID := parts[2]
				teamSize := parseInt(parts[3])
				tracker.SetTeamSize(playerID, teamSize)
			}

		case "poke":
			if len(parts) > 3 {
				playerID := parts[2]
				pokeStr := parts[3]
				poke := parsePokemonFromTeamPreview(pokeStr)
				tracker.AddPokemonToTeam(playerID, poke)
			}
		}
	}

	// Initialize tracker with teams
	summary.Player1.Team = tracker.GetTeam("p1")
	summary.Player2.Team = tracker.GetTeam("p2")
	summary.Player1.TotalLeft = tracker.GetTeamSize("p1")
	summary.Player2.TotalLeft = tracker.GetTeamSize("p2")

	// Second pass: process all battle events
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
			// Save previous turn and start new one
			if currentTurn != nil {
				// Calculate position score for the turn
				currentTurn.PositionScore = tracker.CalculatePositionScore()
				summary.Turns = append(summary.Turns, *currentTurn)
			}
			turnNumber = parseInt(parts[2])
			currentTurn = &Turn{
				TurnNumber:  turnNumber,
				Actions:     []Action{},
				DamageDealt: make(map[string]int),
				HealingDone: make(map[string]int),
			}

		case "switch":
			if len(parts) >= 4 {
				action := parseSwitch(parts)
				if currentTurn != nil {
					currentTurn.Actions = append(currentTurn.Actions, action)
				}
				// Update tracker state
				playerID := extractRawPlayerID(parts[2])
				pokeName := extractPokemonName(parts[3])
				pokehp := extractHPFromSwitch(parts)
				tracker.SwitchPokemon(playerID, pokeName, pokehp)
			}

		case "move":
			if len(parts) >= 4 {
				action := parseMove(parts)
				if currentTurn != nil {
					currentTurn.Actions = append(currentTurn.Actions, action)
				}
			}

		case "-damage":
			if len(parts) >= 4 {
				playerID := extractRawPlayerID(parts[2])
				hpStr := parts[3]
				hp, maxHP := parseHP(hpStr)
				tracker.UpdatePokemonHP(playerID, hp, maxHP)
			}

		case "-heal":
			if len(parts) >= 4 {
				playerID := extractRawPlayerID(parts[2])
				hpStr := parts[3]
				hp, maxHP := parseHP(hpStr)
				tracker.UpdatePokemonHP(playerID, hp, maxHP)
			}

		case "faint":
			if len(parts) > 2 {
				playerID := extractRawPlayerID(parts[2])
				tracker.FaintPokemon(playerID)
				if currentTurn != nil {
					addKeyMoment(summary, turnNumber, "KO", "Pokémon fainted", 8)
				}
			}

		case "-boost", "-unboost":
			// Track stat changes for position scoring
			if len(parts) > 3 {
				tracker.RecordStatChange(parts)
			}

		case "-status":
			// Track status conditions
			if len(parts) > 2 {
				playerID := extractRawPlayerID(parts[2])
				status := parts[3]
				tracker.UpdatePokemonStatus(playerID, status)
			}

		case "-terastallize":
			// Track terastallization
			if len(parts) > 3 {
				playerID := extractRawPlayerID(parts[2])
				teraType := parts[3]
				tracker.TerastallizePokemon(playerID, teraType)
			}

		case "-sidestart", "-sideend":
			// Track field effects like Tailwind
			tracker.RecordFieldEffect(parts)

		case "-crit":
			summary.Stats.CriticalHits++

		case "-supereffective":
			summary.Stats.SuperEffective++

		case "-resisted":
			summary.Stats.NotVeryEffective++

		case "win":
			if len(parts) > 2 {
				winner := parts[2]
				summary.Winner = tracker.PlayerToID(winner)
			}
		}
	}

	// Add the last turn
	if currentTurn != nil {
		currentTurn.PositionScore = tracker.CalculatePositionScore()
		summary.Turns = append(summary.Turns, *currentTurn)
	}

	// Update player losses from tracker
	summary.Player1.Losses = tracker.losses["p1"]
	summary.Player2.Losses = tracker.losses["p2"]
	summary.Player1.TotalLeft = tracker.GetTeamSize("p1") - tracker.losses["p1"]
	summary.Player2.TotalLeft = tracker.GetTeamSize("p2") - tracker.losses["p2"]

	// Calculate statistics and turning points
	calculateStats(summary)
	detectTurningPoints(summary)

	return summary, nil
}

// StateTracker maintains the game state throughout the battle
type StateTracker struct {
	playerNames        map[string]string
	teamSizes          map[string]int
	teams              map[string][]Pokémon
	activePokemon      map[string]*Pokémon // Current active mon for each player
	activePokemonIndex map[string]int
	losses             map[string]int // Fainted pokemon count
	fieldEffects       map[string][]string // Side effects like Tailwind
	statBoosts         map[string]map[string]int // Player->stat->boost level
}

func NewStateTracker() *StateTracker {
	return &StateTracker{
		playerNames:        make(map[string]string),
		teamSizes:          make(map[string]int),
		teams:              make(map[string][]Pokémon),
		activePokemon:      make(map[string]*Pokémon),
		activePokemonIndex: make(map[string]int),
		losses:             make(map[string]int),
		fieldEffects:       make(map[string][]string),
		statBoosts:         make(map[string]map[string]int),
	}
}

func (st *StateTracker) SetPlayerName(playerID, name string) {
	st.playerNames[playerID] = name
}

func (st *StateTracker) SetTeamSize(playerID string, size int) {
	st.teamSizes[playerID] = size
}

func (st *StateTracker) AddPokemonToTeam(playerID string, poke Pokémon) {
	st.teams[playerID] = append(st.teams[playerID], poke)
}

func (st *StateTracker) GetTeam(playerID string) []Pokémon {
	return st.teams[playerID]
}

func (st *StateTracker) GetTeamSize(playerID string) int {
	if size, ok := st.teamSizes[playerID]; ok {
		return size
	}
	return len(st.teams[playerID])
}

func (st *StateTracker) SwitchPokemon(playerID, pokeName string, hp int) {
	team := st.teams[playerID]
	for i, poke := range team {
		if poke.Name == pokeName {
			st.activePokemon[playerID] = &team[i]
			st.activePokemonIndex[playerID] = i
			team[i].CurrentHP = hp
			if team[i].MaxHP == 0 {
				team[i].MaxHP = 100 // Default to 100 for now
			}
			break
		}
	}
}

func (st *StateTracker) UpdatePokemonHP(playerID string, currentHP, maxHP int) {
	if poke, ok := st.activePokemon[playerID]; ok {
		poke.CurrentHP = currentHP
		if poke.MaxHP == 0 {
			poke.MaxHP = maxHP
		}
	}
}

func (st *StateTracker) FaintPokemon(playerID string) {
	if poke, ok := st.activePokemon[playerID]; ok {
		poke.CurrentHP = 0
	}
	st.losses[playerID]++
}

func (st *StateTracker) UpdatePokemonStatus(playerID, status string) {
	if poke, ok := st.activePokemon[playerID]; ok {
		poke.Status = status
	}
}

func (st *StateTracker) TerastallizePokemon(playerID, teraType string) {
	if poke, ok := st.activePokemon[playerID]; ok {
		poke.TeraType = teraType
	}
}

func (st *StateTracker) RecordFieldEffect(parts []string) {
	if len(parts) < 4 {
		return
	}
	playerID := extractRawPlayerID(parts[2])
	effect := parts[3]
	if !contains(st.fieldEffects[playerID], effect) {
		st.fieldEffects[playerID] = append(st.fieldEffects[playerID], effect)
	}
}

func (st *StateTracker) RecordStatChange(parts []string) {
	if len(parts) < 4 {
		return
	}
	playerID := extractRawPlayerID(parts[2])
	stat := parts[3]
	if _, ok := st.statBoosts[playerID]; !ok {
		st.statBoosts[playerID] = make(map[string]int)
	}
	boost := parseInt(parts[4])
	st.statBoosts[playerID][stat] = boost
}

func (st *StateTracker) PlayerToID(playerName string) string {
	for id, name := range st.playerNames {
		if name == playerName {
			if id == "p1" {
				return "player1"
			}
			return "player2"
		}
	}
	return "player2"
}

func (st *StateTracker) CalculatePositionScore() *PositionScore {
	score := &PositionScore{}

	// Calculate Player 1 score
	p1Poke := st.activePokemon["p1"]
	p1HP := 0.0
	if p1Poke != nil && p1Poke.MaxHP > 0 {
		p1HP = float64(p1Poke.CurrentHP) / float64(p1Poke.MaxHP) * 100
	}
	p1Team := 0.0
	if st.teamSizes["p1"] > 0 {
		p1Team = float64((st.teamSizes["p1"] - st.losses["p1"]) * 100 / st.teamSizes["p1"])
	}
	score.Player1Score = (p1HP * 0.6) + (p1Team * 0.4)

	// Calculate Player 2 score
	p2Poke := st.activePokemon["p2"]
	p2HP := 0.0
	if p2Poke != nil && p2Poke.MaxHP > 0 {
		p2HP = float64(p2Poke.CurrentHP) / float64(p2Poke.MaxHP) * 100
	}
	p2Team := 0.0
	if st.teamSizes["p2"] > 0 {
		p2Team = float64((st.teamSizes["p2"] - st.losses["p2"]) * 100 / st.teamSizes["p2"])
	}
	score.Player2Score = (p2HP * 0.6) + (p2Team * 0.4)

	// Determine momentum
	if score.Player1Score > score.Player2Score+5 {
		score.MomentumPlayer = "player1"
	} else if score.Player2Score > score.Player1Score+5 {
		score.MomentumPlayer = "player2"
	} else {
		score.MomentumPlayer = "neutral"
	}

	return score
}

// Helper parsing functions

func parsePokemonFromTeamPreview(pokeStr string) Pokémon {
	// Format: "Ursaluna-Bloodmoon, L50, M"
	parts := strings.Split(pokeStr, ",")
	name := strings.TrimSpace(parts[0])

	poke := Pokémon{
		ID:        normalizeID(name),
		Name:      name,
		MaxHP:     100, // Default max HP for level 50
		CurrentHP: 100,
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
	playerID := extractPlayerIDFromRef(parts[2])
	moveName := strings.TrimSpace(parts[3])

	return Action{
		Player:     playerID,
		ActionType: "move",
		Move: &Move{
			ID:   normalizeID(moveName),
			Name: moveName,
		},
	}
}

func parseSwitch(parts []string) Action {
	// |switch|p1b: Typhlosion|Typhlosion-Hisui, L50, M|100\/100
	playerID := extractPlayerIDFromRef(parts[2])
	switchToPoke := extractPokemonName(parts[3])

	return Action{
		Player:     playerID,
		ActionType: "switch",
		SwitchTo:   switchToPoke,
	}
}

func extractPlayerIDFromRef(ref string) string {
	// Convert "p1a: Whimsicott" to "player1" or "p2b: Maushold" to "player2"
	if strings.HasPrefix(ref, "p1") {
		return "player1"
	}
	return "player2"
}

func extractRawPlayerID(ref string) string {
	// Convert "p1a: Whimsicott" to "p1" or "p2b: Maushold" to "p2"
	if strings.HasPrefix(ref, "p1") {
		return "p1"
	}
	return "p2"
}

func extractPokemonName(ref string) string {
	// From "Typhlosion-Hisui, L50, M" extract "Typhlosion-Hisui"
	parts := strings.Split(ref, ",")
	return strings.TrimSpace(parts[0])
}

func extractHPFromSwitch(parts []string) int {
	// From "100\/100" extract 100
	if len(parts) > 4 {
		hpStr := parts[4]
		hpParts := strings.Split(hpStr, "\\/")
		return parseInt(hpParts[0])
	}
	return 100
}

func parseHP(hpStr string) (int, int) {
	// "63\/100" -> (63, 100) or "0 fnt" -> (0, maxHP)
	if strings.Contains(hpStr, "\\/") {
		parts := strings.Split(hpStr, "\\/")
		current := parseInt(parts[0])
		max := parseInt(parts[1])
		return current, max
	}
	// Handle "0 fnt" format
	if strings.Contains(hpStr, "fnt") {
		return 0, 100
	}
	return parseInt(hpStr), 100
}

func normalizeID(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, "-", ""))
}

func parseInt(s string) int {
	var result int
	_, _ = fmt.Sscanf(strings.TrimSpace(s), "%d", &result)
	return result
}

func addKeyMoment(summary *BattleSummary, turnNumber int, mType, description string, significance int) {
	summary.KeyMoments = append(summary.KeyMoments, KeyMoment{
		TurnNumber:   turnNumber,
		Type:         mType,
		Description:  description,
		Significance: significance,
	})
}

func detectTurningPoints(summary *BattleSummary) {
	if len(summary.Turns) < 2 {
		return
	}

	var turningPoints []TurningPoint

	for i := 1; i < len(summary.Turns); i++ {
		prev := summary.Turns[i-1]
		curr := summary.Turns[i]

		if prev.PositionScore == nil || curr.PositionScore == nil {
			continue
		}

		// Calculate score change
		p1Delta := curr.PositionScore.Player1Score - prev.PositionScore.Player1Score
		p2Delta := curr.PositionScore.Player2Score - prev.PositionScore.Player2Score
		momentumShift := p1Delta - p2Delta

		// Significant momentum shifts are 15+ points
		if absFloat(momentumShift) >= 15 {
			significance := int(absFloat(momentumShift) / 10)
			if significance > 10 {
				significance = 10
			}

			direction := "Player 1"
			if momentumShift < 0 {
				direction = "Player 2"
			}

			tp := TurningPoint{
				TurnNumber:   curr.TurnNumber,
				Score1Before: prev.PositionScore.Player1Score,
				Score1After:  curr.PositionScore.Player1Score,
				Score2Before: prev.PositionScore.Player2Score,
				Score2After:  curr.PositionScore.Player2Score,
				MomentumShift: momentumShift,
				Significance: significance,
				Description: fmt.Sprintf("%s gained significant momentum this turn", direction),
			}

			turningPoints = append(turningPoints, tp)

			// Add to key moments as well
			addKeyMoment(summary, curr.TurnNumber, "turning_point",
				fmt.Sprintf("Turn %d: %s", curr.TurnNumber, tp.Description), significance)
		}
	}

	summary.Stats.TurningPoints = turningPoints
}

func calculateStats(summary *BattleSummary) {
	summary.Stats.TotalTurns = len(summary.Turns)
	summary.Stats.MoveFrequency = make(map[string]int)
	summary.Stats.Player1Stats = PlayerStats{
		MovesByType: make(map[string]int),
	}
	summary.Stats.Player2Stats = PlayerStats{
		MovesByType: make(map[string]int),
	}

	totalDamageDealt1 := 0
	totalDamageDealt2 := 0
	totalHealing1 := 0
	totalHealing2 := 0

	for _, turn := range summary.Turns {
		for _, action := range turn.Actions {
			if action.ActionType == "move" && action.Move != nil {
				summary.Stats.MoveFrequency[action.Move.ID]++

				if action.Player == "player1" {
					summary.Stats.Player1Stats.MoveCount++
				} else {
					summary.Stats.Player2Stats.MoveCount++
				}
			} else if action.ActionType == "switch" {
				summary.Stats.Switch++
				if action.Player == "player1" {
					summary.Stats.Player1Stats.SwitchCount++
				} else {
					summary.Stats.Player2Stats.SwitchCount++
				}
			}
		}

		// Accumulate damage and healing
		for player, damage := range turn.DamageDealt {
			if player == "player1" {
				totalDamageDealt1 += damage
			} else {
				totalDamageDealt2 += damage
			}
		}
		for player, healing := range turn.HealingDone {
			if player == "player1" {
				totalHealing1 += healing
			} else {
				totalHealing2 += healing
			}
		}
	}

	summary.Stats.Player1Stats.DamageDealt = totalDamageDealt1
	summary.Stats.Player2Stats.DamageDealt = totalDamageDealt2
	summary.Stats.Player1Stats.HealingDone = totalHealing1
	summary.Stats.Player2Stats.HealingDone = totalHealing2

	if summary.Stats.TotalTurns > 0 {
		summary.Stats.AvgDamagePerTurn = float64(totalDamageDealt1+totalDamageDealt2) / float64(summary.Stats.TotalTurns)
		summary.Stats.AvgHealPerTurn = float64(totalHealing1+totalHealing2) / float64(summary.Stats.TotalTurns)
	}
}

func generateUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("battle-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
