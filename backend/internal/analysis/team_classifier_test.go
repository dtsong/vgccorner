package analysis

import (
	"testing"
)

func TestClassifyTeam_HardTrickRoom(t *testing.T) {
	// Team with 2 Trick Room users
	team := []Pokémon{
		{
			Name: "Cresselia",
			Moves: []Move{
				{Name: "Trick Room"},
			},
		},
		{
			Name: "Dusclops",
			Moves: []Move{
				{Name: "Trick Room"},
			},
		},
		{Name: "Torkoal"},
		{Name: "Rhyperior"},
	}

	classification := ClassifyTeam(team)

	if classification.Archetype != "Hard Trick Room" {
		t.Errorf("Expected 'Hard Trick Room', got '%s'", classification.Archetype)
	}

	if !classification.HasTrickRoom {
		t.Error("Expected HasTrickRoom to be true")
	}

	if len(classification.TrickRoomUsers) != 2 {
		t.Errorf("Expected 2 Trick Room users, got %d", len(classification.TrickRoomUsers))
	}
}

func TestClassifyTeam_TailRoom(t *testing.T) {
	// Team with both Tailwind and Trick Room
	team := []Pokémon{
		{
			Name: "Talonflame",
			Moves: []Move{
				{Name: "Tailwind"},
			},
		},
		{
			Name: "Cresselia",
			Moves: []Move{
				{Name: "Trick Room"},
			},
		},
	}

	classification := ClassifyTeam(team)

	if classification.Archetype != "TailRoom" {
		t.Errorf("Expected 'TailRoom', got '%s'", classification.Archetype)
	}

	if !classification.HasTrickRoom || !classification.HasTailwind {
		t.Error("Expected both HasTrickRoom and HasTailwind to be true")
	}
}

func TestClassifyTeam_SunOffense(t *testing.T) {
	// Team with Drought ability
	team := []Pokémon{
		{
			Name:    "Torkoal",
			Ability: "Drought",
		},
		{Name: "Venusaur"},
		{Name: "Charizard"},
	}

	classification := ClassifyTeam(team)

	if classification.Archetype != "Sun Offense" {
		t.Errorf("Expected 'Sun Offense', got '%s'", classification.Archetype)
	}

	if !classification.HasWeatherSetter {
		t.Error("Expected HasWeatherSetter to be true")
	}

	if classification.WeatherType != "sun" {
		t.Errorf("Expected WeatherType 'sun', got '%s'", classification.WeatherType)
	}
}

func TestClassifyTeam_RainOffense(t *testing.T) {
	// Team with Drizzle ability
	team := []Pokémon{
		{
			Name:    "Pelipper",
			Ability: "Drizzle",
		},
		{Name: "Kingdra"},
		{Name: "Ludicolo"},
	}

	classification := ClassifyTeam(team)

	if classification.Archetype != "Rain Offense" {
		t.Errorf("Expected 'Rain Offense', got '%s'", classification.Archetype)
	}

	if classification.WeatherType != "rain" {
		t.Errorf("Expected WeatherType 'rain', got '%s'", classification.WeatherType)
	}
}

func TestClassifyTeam_BalanceBros(t *testing.T) {
	// Team with Incineroar and Rillaboom
	team := []Pokémon{
		{ID: "incineroar", Name: "Incineroar"},
		{ID: "rillaboom", Name: "Rillaboom"},
		{Name: "Landorus"},
		{Name: "Flutter Mane"},
	}

	classification := ClassifyTeam(team)

	if classification.Archetype != "Balance Bros" {
		t.Errorf("Expected 'Balance Bros', got '%s'", classification.Archetype)
	}

	if !classification.HasBalanceBros {
		t.Error("Expected HasBalanceBros to be true")
	}
}

func TestClassifyTeam_PsySpam(t *testing.T) {
	// Team with Psychic Terrain and Expanding Force
	team := []Pokémon{
		{
			Name: "Indeedee",
			Moves: []Move{
				{Name: "Psychic Terrain"},
			},
		},
		{
			Name: "Armarouge",
			Moves: []Move{
				{Name: "Expanding Force"},
			},
		},
	}

	classification := ClassifyTeam(team)

	if classification.Archetype != "Psy-Spam" {
		t.Errorf("Expected 'Psy-Spam', got '%s'", classification.Archetype)
	}

	if !classification.HasPsyTerrain {
		t.Error("Expected HasPsyTerrain to be true")
	}
}

func TestClassifyTeam_TailwindHyperOffense(t *testing.T) {
	// Team with Tailwind and Choice items
	team := []Pokémon{
		{
			Name: "Talonflame",
			Moves: []Move{
				{Name: "Tailwind"},
			},
		},
		{
			Name: "Dragapult",
			Item: "Choice Specs",
		},
		{
			Name: "Rillaboom",
			Item: "Choice Band",
		},
	}

	classification := ClassifyTeam(team)

	if classification.Archetype != "Tailwind Hyper Offense" {
		t.Errorf("Expected 'Tailwind Hyper Offense', got '%s'", classification.Archetype)
	}

	if !classification.HasTailwind {
		t.Error("Expected HasTailwind to be true")
	}

	if !classification.HasChoiceItems {
		t.Error("Expected HasChoiceItems to be true")
	}

	if len(classification.ChoiceUsers) != 2 {
		t.Errorf("Expected 2 Choice users, got %d", len(classification.ChoiceUsers))
	}
}

func TestClassifyTeam_Unclassified(t *testing.T) {
	// Team that doesn't fit any archetype
	team := []Pokémon{
		{Name: "Garchomp"},
		{Name: "Metagross"},
		{Name: "Salamence"},
		{Name: "Tyranitar"},
	}

	classification := ClassifyTeam(team)

	if classification.Archetype != "Unclassified" {
		t.Errorf("Expected 'Unclassified', got '%s'", classification.Archetype)
	}
}

func TestGetArchetypeDescription(t *testing.T) {
	tests := []struct {
		archetype string
		wantDesc  bool // Should have a description
	}{
		{"Hard Trick Room", true},
		{"TailRoom", true},
		{"Sun Offense", true},
		{"Psy-Spam", true},
		{"Unclassified", true},
		{"Unknown Archetype", true}, // Should return default
	}

	for _, tt := range tests {
		t.Run(tt.archetype, func(t *testing.T) {
			desc := GetArchetypeDescription(tt.archetype)
			if tt.wantDesc && desc == "" {
				t.Errorf("Expected description for '%s', got empty string", tt.archetype)
			}
		})
	}
}
