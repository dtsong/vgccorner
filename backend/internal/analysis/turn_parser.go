package analysis

import "strings"

// TurnParser handles parsing detailed turn-by-turn information from battle logs
type TurnParser struct {
	currentTurn      *Turn
	pendingEvents    []string
	actionOrder      int
	lastMovedPokemon map[string]string // tracks which Pokemon just moved for impact attribution
}

// NewTurnParser creates a new turn parser
func NewTurnParser() *TurnParser {
	return &TurnParser{
		pendingEvents:    []string{},
		lastMovedPokemon: make(map[string]string),
		actionOrder:      0,
	}
}

// ProcessTurnEvent processes a single line from the battle log during a turn
func (tp *TurnParser) ProcessTurnEvent(line string, tracker *StateTracker) {
	if line == "" || !strings.HasPrefix(line, "|") {
		return
	}

	parts := strings.Split(line, "|")
	if len(parts) < 2 {
		return
	}

	command := parts[1]

	switch command {
	case "move":
		// Process any pending events for the previous action
		tp.flushPendingEvents()

		// Parse the move
		if len(parts) >= 4 {
			action := tp.parseMove(parts)
			action.OrderInTurn = tp.actionOrder
			tp.actionOrder++

			if tp.currentTurn != nil {
				tp.currentTurn.Actions = append(tp.currentTurn.Actions, action)
			}

			// Track which Pokemon just moved
			tp.lastMovedPokemon["last"] = action.Pokemon
		}

	case "switch":
		tp.flushPendingEvents()

		if len(parts) >= 4 {
			action := tp.parseSwitch(parts)
			action.OrderInTurn = tp.actionOrder
			tp.actionOrder++

			if tp.currentTurn != nil {
				tp.currentTurn.Actions = append(tp.currentTurn.Actions, action)
			}
		}

	case "-damage", "-heal", "-status", "faint", "-crit", "-supereffective", "-resisted",
		"-immune", "-miss", "-weather", "-fieldstart", "-boost", "-unboost":
		// Collect events that relate to the last action
		tp.pendingEvents = append(tp.pendingEvents, line)

	default:
		// Other events - might want to track these too
	}
}

// StartNewTurn starts tracking a new turn
func (tp *TurnParser) StartNewTurn(turnNumber int) *Turn {
	// Flush any pending events from previous turn
	tp.flushPendingEvents()

	tp.currentTurn = &Turn{
		TurnNumber:  turnNumber,
		Actions:     []Action{},
		DamageDealt: make(map[string]int),
		HealingDone: make(map[string]int),
	}
	tp.actionOrder = 0
	tp.lastMovedPokemon = make(map[string]string)

	return tp.currentTurn
}

// FinalizeTurn completes the current turn and returns it
func (tp *TurnParser) FinalizeTurn(tracker *StateTracker) *Turn {
	// Flush any remaining events
	tp.flushPendingEvents()

	if tp.currentTurn != nil {
		tp.currentTurn.PositionScore = tracker.CalculatePositionScore()
	}

	turn := tp.currentTurn
	tp.currentTurn = nil
	return turn
}

// flushPendingEvents applies pending events to the last action
func (tp *TurnParser) flushPendingEvents() {
	if len(tp.pendingEvents) == 0 {
		return
	}

	if tp.currentTurn != nil && len(tp.currentTurn.Actions) > 0 {
		// Get the last action
		lastAction := &tp.currentTurn.Actions[len(tp.currentTurn.Actions)-1]

		// Enhance it with impact information
		if lastAction.Move != nil {
			EnhanceActionWithImpact(lastAction, lastAction.Move.Name, tp.pendingEvents)
		}
	}

	// Clear pending events
	tp.pendingEvents = []string{}
}

// parseMove parses a move command with enhanced details
func (tp *TurnParser) parseMove(parts []string) Action {
	// |move|p1a: Gengar|Shadow Ball|p2a: Dusclops
	playerID := extractPlayerIDFromRef(parts[2])
	pokemonName := extractPokemonName(parts[2])
	moveName := strings.TrimSpace(parts[3])

	action := Action{
		Player:     playerID,
		ActionType: "move",
		Pokemon:    pokemonName,
		Move: &Move{
			ID:   normalizeID(moveName),
			Name: moveName,
		},
	}

	// Parse target if present
	if len(parts) >= 5 {
		targetName := extractPokemonName(parts[4])
		action.Target = targetName
	}

	return action
}

// parseSwitch parses a switch command
func (tp *TurnParser) parseSwitch(parts []string) Action {
	// |switch|p1b: Typhlosion|Typhlosion-Hisui, L50, M|100/100
	playerID := extractPlayerIDFromRef(parts[2])
	pokemonName := extractPokemonName(parts[2])
	switchToPoke := extractPokemonName(parts[3])

	return Action{
		Player:     playerID,
		ActionType: "switch",
		Pokemon:    pokemonName,
		SwitchTo:   switchToPoke,
	}
}

// ParseEnhancedShowdownLog is an enhanced version of ParseShowdownLog with better turn tracking
func ParseEnhancedShowdownLog(logContent string) (*BattleSummary, error) {
	// First do the basic parsing
	summary, err := ParseShowdownLog(logContent)
	if err != nil {
		return nil, err
	}

	// Now do enhanced turn parsing for more detailed action tracking
	lines := strings.Split(logContent, "\n")
	tracker := NewStateTracker()
	turnParser := NewTurnParser()

	// First pass: set up tracker
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
		case "player":
			if len(parts) > 3 {
				playerID := parts[2]
				playerName := parts[3]
				tracker.SetPlayerName(playerID, playerName)
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

	// Second pass: detailed turn parsing
	var enhancedTurns []Turn
	var currentTurnNumber int

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
			// Finalize previous turn
			if currentTurnNumber > 0 {
				turn := turnParser.FinalizeTurn(tracker)
				if turn != nil {
					enhancedTurns = append(enhancedTurns, *turn)
				}
			}

			// Start new turn
			currentTurnNumber = parseInt(parts[2])
			turnParser.StartNewTurn(currentTurnNumber)

		case "switch":
			turnParser.ProcessTurnEvent(line, tracker)
			// Update tracker
			if len(parts) >= 4 {
				playerID := extractRawPlayerID(parts[2])
				pokeName := extractPokemonName(parts[3])
				pokehp := extractHPFromSwitch(parts)
				tracker.SwitchPokemon(playerID, pokeName, pokehp)
			}

		case "move", "-damage", "-heal", "-status", "faint", "-crit",
			"-supereffective", "-resisted", "-immune", "-miss", "-weather",
			"-fieldstart", "-boost", "-unboost":
			turnParser.ProcessTurnEvent(line, tracker)

			// Update tracker for damage/healing
			if command == "-damage" && len(parts) >= 4 {
				playerID := extractRawPlayerID(parts[2])
				hpStr := parts[3]
				hp, maxHP := parseHP(hpStr)
				tracker.UpdatePokemonHP(playerID, hp, maxHP)
			}
		}
	}

	// Finalize last turn
	if currentTurnNumber > 0 {
		turn := turnParser.FinalizeTurn(tracker)
		if turn != nil {
			enhancedTurns = append(enhancedTurns, *turn)
		}
	}

	// Replace turns in summary with enhanced turns if we got more detailed data
	if len(enhancedTurns) > 0 {
		summary.Turns = enhancedTurns
	}

	return summary, nil
}
