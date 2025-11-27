package analysis

import "strings"

// ClassifyTeam analyzes a team and determines its archetype based on VGC criteria
func ClassifyTeam(team []Pokémon) TeamClassification {
	classification := TeamClassification{
		TrickRoomUsers:  []string{},
		TailwindUsers:   []string{},
		WeatherSetters:  []string{},
		PsyTerrainUsers: []string{},
		ChoiceUsers:     []string{},
		Tags:            []string{},
	}

	// Track team composition
	hasIncineroar := false
	hasRillaboom := false
	trickRoomCount := 0
	expandingForceUsers := []string{}

	// Analyze each Pokémon
	for _, poke := range team {
		pokeID := strings.ToLower(poke.ID)

		// Check for specific Pokémon
		if strings.Contains(pokeID, "incineroar") {
			hasIncineroar = true
		}
		if strings.Contains(pokeID, "rillaboom") {
			hasRillaboom = true
		}

		// Check ability-based weather setters
		ability := strings.ToLower(poke.Ability)
		switch ability {
		case "drought":
			classification.HasWeatherSetter = true
			classification.WeatherType = "sun"
			classification.WeatherSetters = append(classification.WeatherSetters, poke.Name)
		case "drizzle":
			classification.HasWeatherSetter = true
			classification.WeatherType = "rain"
			classification.WeatherSetters = append(classification.WeatherSetters, poke.Name)
		case "sand stream":
			classification.HasWeatherSetter = true
			classification.WeatherType = "sand"
			classification.WeatherSetters = append(classification.WeatherSetters, poke.Name)
		case "snow warning":
			classification.HasWeatherSetter = true
			classification.WeatherType = "snow"
			classification.WeatherSetters = append(classification.WeatherSetters, poke.Name)
		}

		// Check item
		item := strings.ToLower(poke.Item)
		if item == "choice specs" || item == "choice band" || item == "choice scarf" {
			classification.HasChoiceItems = true
			classification.ChoiceUsers = append(classification.ChoiceUsers, poke.Name)
		}

		// Check moves
		for _, move := range poke.Moves {
			moveName := strings.ToLower(move.Name)
			moveID := strings.ToLower(move.ID)

			switch {
			case moveName == "trick room" || moveID == "trickroom":
				classification.HasTrickRoom = true
				classification.TrickRoomUsers = append(classification.TrickRoomUsers, poke.Name)
				trickRoomCount++

			case moveName == "tailwind" || moveID == "tailwind":
				classification.HasTailwind = true
				classification.TailwindUsers = append(classification.TailwindUsers, poke.Name)

			case moveName == "sunny day" || moveID == "sunnyday":
				classification.HasWeatherSetter = true
				if classification.WeatherType == "" {
					classification.WeatherType = "sun"
				}
				classification.WeatherSetters = append(classification.WeatherSetters, poke.Name)

			case moveName == "rain dance" || moveID == "raindance":
				classification.HasWeatherSetter = true
				if classification.WeatherType == "" {
					classification.WeatherType = "rain"
				}
				classification.WeatherSetters = append(classification.WeatherSetters, poke.Name)

			case moveName == "psychic terrain" || moveID == "psychicterrain":
				classification.HasPsyTerrain = true
				classification.PsyTerrainUsers = append(classification.PsyTerrainUsers, poke.Name)

			case moveName == "expanding force" || moveID == "expandingforce":
				expandingForceUsers = append(expandingForceUsers, poke.Name)
			}
		}
	}

	// Check for Balance Bros (Incineroar + Rillaboom)
	classification.HasBalanceBros = hasIncineroar && hasRillaboom

	// Determine primary archetype based on criteria
	classification.Archetype = determineArchetype(classification, trickRoomCount, expandingForceUsers)

	return classification
}

// determineArchetype determines the primary team archetype based on classification data
func determineArchetype(c TeamClassification, trickRoomCount int, expandingForceUsers []string) string {
	// Priority order based on VGC criteria

	// Hard Trick Room: Trick Room on 2+ Pokémon
	if trickRoomCount >= 2 {
		return "Hard Trick Room"
	}

	// TailRoom: Both Tailwind and Trick Room
	if c.HasTailwind && c.HasTrickRoom {
		return "TailRoom"
	}

	// Sun Offense: Drought or Sunny Day
	if c.HasWeatherSetter && c.WeatherType == "sun" {
		return "Sun Offense"
	}

	// Rain Offense: Drizzle or Rain Dance
	if c.HasWeatherSetter && c.WeatherType == "rain" {
		return "Rain Offense"
	}

	// Balance Bros: Incineroar + Rillaboom
	if c.HasBalanceBros {
		return "Balance Bros"
	}

	// Psy-Spam: Psychic Terrain + Expanding Force user
	if c.HasPsyTerrain && len(expandingForceUsers) > 0 {
		return "Psy-Spam"
	}

	// Tailwind Hyper Offense: Tailwind + Choice items
	if c.HasTailwind && c.HasChoiceItems {
		return "Tailwind Hyper Offense"
	}

	// Tailwind (generic)
	if c.HasTailwind {
		return "Tailwind"
	}

	// Trick Room (generic)
	if c.HasTrickRoom {
		return "Trick Room"
	}

	// Weather teams (generic)
	if c.HasWeatherSetter {
		switch c.WeatherType {
		case "sun":
			return "Sun"
		case "rain":
			return "Rain"
		case "sand":
			return "Sand"
		case "snow":
			return "Snow"
		}
	}

	// Default: Unclassified
	return "Unclassified"
}

// GetArchetypeDescription returns a human-readable description of the team archetype
func GetArchetypeDescription(archetype string) string {
	descriptions := map[string]string{
		"Hard Trick Room":        "A team built around Trick Room with multiple setters for reliability",
		"TailRoom":               "A flexible team that can operate under both Tailwind and Trick Room",
		"Sun Offense":            "An offensive team utilizing sun weather to power up Fire-type attacks",
		"Rain Offense":           "An offensive team utilizing rain weather to power up Water-type attacks",
		"Balance Bros":           "A balanced team featuring Incineroar and Rillaboom for defensive synergy",
		"Psy-Spam":               "A team focused on Psychic Terrain with Expanding Force for massive spread damage",
		"Tailwind Hyper Offense": "An aggressive team using Tailwind and Choice items for overwhelming speed and power",
		"Tailwind":               "A speed-based team utilizing Tailwind for speed control",
		"Trick Room":             "A team utilizing Trick Room for speed control",
		"Sun":                    "A team utilizing sun weather",
		"Rain":                   "A team utilizing rain weather",
		"Sand":                   "A team utilizing sandstorm weather",
		"Snow":                   "A team utilizing snow weather",
		"Unclassified":           "A team that doesn't fit standard VGC archetypes",
	}

	if desc, ok := descriptions[archetype]; ok {
		return desc
	}
	return "A unique team composition"
}
