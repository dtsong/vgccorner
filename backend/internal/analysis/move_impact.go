package analysis

import "strings"

// EnhanceActionWithImpact adds detailed impact information to an action based on subsequent battle events
func EnhanceActionWithImpact(action *Action, moveName string, events []string) {
	if action.Impact == nil {
		action.Impact = &MoveImpact{
			Fainted:     []string{},
			StatChanges: []StatChange{},
		}
	}

	moveID := strings.ToLower(moveName)

	// Track special move types
	switch {
	case strings.Contains(moveID, "fake out") || strings.Contains(moveID, "fakeout"):
		action.Impact.FakeOut = true
		action.Impact.SpeedControl = "flinch"
	case strings.Contains(moveID, "protect") || strings.Contains(moveID, "detect") ||
		strings.Contains(moveID, "baneful bunker") || strings.Contains(moveID, "king's shield"):
		action.Impact.Protect = true
	case strings.Contains(moveID, "trick room") || strings.Contains(moveID, "trickroom"):
		action.Impact.SpeedControl = "trick-room"
	case strings.Contains(moveID, "tailwind"):
		action.Impact.SpeedControl = "tailwind"
	case strings.Contains(moveID, "icy wind"):
		action.Impact.SpeedControl = "speed-drop"
	case strings.Contains(moveID, "thunder wave"):
		action.Impact.SpeedControl = "paralysis"
	}

	// Parse events to extract impact
	for _, event := range events {
		if !strings.HasPrefix(event, "|") {
			continue
		}
		parts := strings.Split(event, "|")
		if len(parts) < 2 {
			continue
		}

		eventType := parts[1]

		switch eventType {
		case "-damage":
			// Damage dealt - parse HP change
			if len(parts) >= 4 {
				hpBefore, hpAfter := parseHPChange(parts)
				if hpBefore > hpAfter {
					action.Impact.DamageDealt += (hpBefore - hpAfter)
				}
			}

		case "-heal":
			// Healing done
			if len(parts) >= 4 {
				hpBefore, hpAfter := parseHPChange(parts)
				if hpAfter > hpBefore {
					action.Impact.HealingDone += (hpAfter - hpBefore)
				}
			}

		case "-status":
			// Status inflicted
			if len(parts) >= 4 {
				action.Impact.StatusInflicted = parts[3]
			}

		case "faint":
			// Pokémon fainted
			if len(parts) >= 3 {
				faintedPoke := extractPokemonName(parts[2])
				action.Impact.Fainted = append(action.Impact.Fainted, faintedPoke)
			}

		case "-crit":
			// Critical hit
			action.Impact.Critical = true
			action.Result = "critical-hit"

		case "-supereffective":
			// Super effective
			action.Impact.Effectiveness = "super-effective"
			action.Result = "super-effective"

		case "-resisted":
			// Not very effective
			action.Impact.Effectiveness = "not-very-effective"
			action.Result = "not-very-effective"

		case "-immune":
			// Immune
			action.Impact.Effectiveness = "immune"
			action.Result = "immune"

		case "-miss":
			// Move missed
			action.Impact.Missed = true
			action.Result = "miss"

		case "-weather":
			// Weather set
			if len(parts) >= 3 {
				action.Impact.WeatherSet = parts[2]
			}

		case "-fieldstart":
			// Terrain or room set
			if len(parts) >= 3 {
				field := strings.ToLower(parts[2])
				if strings.Contains(field, "terrain") {
					action.Impact.TerrainSet = parts[2]
				} else if strings.Contains(field, "trick room") {
					action.Impact.SpeedControl = "trick-room"
				} else if strings.Contains(field, "tailwind") {
					action.Impact.SpeedControl = "tailwind"
				}
			}

		case "-boost", "-unboost":
			// Stat changes
			if len(parts) >= 5 {
				pokemon := extractPokemonName(parts[2])
				stat := parts[3]
				stages := parseInt(parts[4])
				if eventType == "-unboost" {
					stages = -stages
				}
				action.Impact.StatChanges = append(action.Impact.StatChanges, StatChange{
					Pokemon: pokemon,
					Stat:    stat,
					Stages:  stages,
				})
			}
		}
	}

	// Set result if not already set
	if action.Result == "" {
		if action.Impact.Missed {
			action.Result = "miss"
		} else if len(action.Impact.Fainted) > 0 {
			action.Result = "faint"
		} else if action.Impact.DamageDealt > 0 {
			action.Result = "success"
		}
	}

	// Generate details string
	action.Details = generateActionDetails(action)
}

// parseHPChange extracts HP before and after from an HP string
func parseHPChange(parts []string) (int, int) {
	if len(parts) < 4 {
		return 0, 0
	}
	// HP format: "100/100" or "50/100"
	hpStr := parts[3]
	hp, _ := parseHP(hpStr)
	return 100, hp // Simplified - in real impl, track actual HP
}

// generateActionDetails creates a human-readable description of the action's impact
func generateActionDetails(action *Action) string {
	impact := action.Impact
	if impact == nil {
		return ""
	}

	var details []string

	if impact.Critical {
		details = append(details, "Critical Hit")
	}

	switch impact.Effectiveness {
	case "super-effective":
		details = append(details, "It's super effective")
	case "not-very-effective":
		details = append(details, "It's not very effective")
	case "immune":
		details = append(details, "It doesn't affect the target")
	}

	if impact.FakeOut {
		details = append(details, "Target flinched")
	}

	switch impact.SpeedControl {
	case "trick-room":
		details = append(details, "Dimensions twisted")
	case "tailwind":
		details = append(details, "Tailwind blew")
	case "paralysis":
		details = append(details, "Target was paralyzed")
	}

	if len(impact.Fainted) > 0 {
		for _, poke := range impact.Fainted {
			details = append(details, poke+" flinched")
		}
	}

	if impact.Missed {
		details = append(details, "But it missed")
	}

	if len(details) > 0 {
		return strings.Join(details, ", ")
	}

	return ""
}
