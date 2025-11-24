package analysis

import "testing"

// Tests for uncovered functions

func TestUpdatePokemonStatus(t *testing.T) {
	tracker := NewStateTracker()

	// Add a pokemon to the team
	poke := Pokémon{Name: "Pikachu", CurrentHP: 100, MaxHP: 100}
	tracker.AddPokemonToTeam("p1", poke)

	// Set as active
	tracker.SwitchPokemon("p1", "Pikachu", 100)

	// Update status
	tracker.UpdatePokemonStatus("p1", "par")

	if tracker.activePokemon["p1"].Status != "par" {
		t.Errorf("expected status 'par', got %q", tracker.activePokemon["p1"].Status)
	}

	// Test with non-existent player
	tracker.UpdatePokemonStatus("p3", "brn")
	// Should not panic
}

func TestTerastallizePokemon(t *testing.T) {
	tracker := NewStateTracker()

	// Add a pokemon to the team
	poke := Pokémon{Name: "Charizard", CurrentHP: 100, MaxHP: 100}
	tracker.AddPokemonToTeam("p2", poke)

	// Set as active
	tracker.SwitchPokemon("p2", "Charizard", 100)

	// Terastallize
	tracker.TerastallizePokemon("p2", "Dragon")

	if tracker.activePokemon["p2"].TeraType != "Dragon" {
		t.Errorf("expected tera type 'Dragon', got %q", tracker.activePokemon["p2"].TeraType)
	}

	// Test with non-existent player
	tracker.TerastallizePokemon("p3", "Fire")
	// Should not panic
}

func TestRecordFieldEffect(t *testing.T) {
	tracker := NewStateTracker()

	// Test with valid parts
	parts := []string{"", "-sidestart", "p1a: Pikachu", "Tailwind"}
	tracker.RecordFieldEffect(parts)

	if len(tracker.fieldEffects["p1"]) == 0 {
		t.Error("expected field effect to be recorded")
	}

	if tracker.fieldEffects["p1"][0] != "Tailwind" {
		t.Errorf("expected 'Tailwind', got %q", tracker.fieldEffects["p1"][0])
	}

	// Test duplicate effect (should not add twice)
	tracker.RecordFieldEffect(parts)
	if len(tracker.fieldEffects["p1"]) != 1 {
		t.Errorf("expected 1 field effect, got %d", len(tracker.fieldEffects["p1"]))
	}

	// Test with insufficient parts
	shortParts := []string{"", "-sidestart", "p1a: Pikachu"}
	tracker.RecordFieldEffect(shortParts)
	// Should not panic
}

func TestRecordStatChange(t *testing.T) {
	tracker := NewStateTracker()

	// Test with valid parts
	parts := []string{"", "-boost", "p1a: Pikachu", "atk", "2"}
	tracker.RecordStatChange(parts)

	if tracker.statBoosts["p1"]["atk"] != 2 {
		t.Errorf("expected atk boost of 2, got %d", tracker.statBoosts["p1"]["atk"])
	}

	// Test with another stat
	parts2 := []string{"", "-unboost", "p2a: Charizard", "def", "-1"}
	tracker.RecordStatChange(parts2)

	if tracker.statBoosts["p2"]["def"] != -1 {
		t.Errorf("expected def boost of -1, got %d", tracker.statBoosts["p2"]["def"])
	}

	// Test with insufficient parts
	shortParts := []string{"", "-boost", "p1a: Pikachu"}
	tracker.RecordStatChange(shortParts)
	// Should not panic
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "item exists",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "banana",
			expected: true,
		},
		{
			name:     "item does not exist",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "orange",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			item:     "test",
			expected: false,
		},
		{
			name:     "empty string",
			slice:    []string{"", "test"},
			item:     "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseHPEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		hpStr       string
		expectedCur int
		expectedMax int
	}{
		{
			name:        "normal HP",
			hpStr:       "50/100",
			expectedCur: 50,
			expectedMax: 100,
		},
		{
			name:        "fainted",
			hpStr:       "0 fnt",
			expectedCur: 0,
			expectedMax: 100,
		},
		{
			name:        "full HP",
			hpStr:       "100/100",
			expectedCur: 100,
			expectedMax: 100,
		},
		{
			name:        "single number (no slash)",
			hpStr:       "75",
			expectedCur: 75,
			expectedMax: 100,
		},
		{
			name:        "zero HP",
			hpStr:       "0/100",
			expectedCur: 0,
			expectedMax: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cur, max := parseHP(tt.hpStr)
			if cur != tt.expectedCur {
				t.Errorf("expected current HP %d, got %d", tt.expectedCur, cur)
			}
			if max != tt.expectedMax {
				t.Errorf("expected max HP %d, got %d", tt.expectedMax, max)
			}
		})
	}
}

func TestExtractHPFromSwitchEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		expected int
	}{
		{
			name:     "with HP",
			parts:    []string{"", "switch", "p1a: Pikachu", "Pikachu, L50", "80/100"},
			expected: 80,
		},
		{
			name:     "without HP (too few parts)",
			parts:    []string{"", "switch", "p1a: Pikachu", "Pikachu, L50"},
			expected: 100,
		},
		{
			name:     "full HP",
			parts:    []string{"", "switch", "p1a: Pikachu", "Pikachu, L50", "100/100"},
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractHPFromSwitch(tt.parts)
			if result != tt.expected {
				t.Errorf("expected HP %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestSwitchPokemonEdgeCases(t *testing.T) {
	tracker := NewStateTracker()

	// Add multiple pokemon to team
	poke1 := Pokémon{Name: "Pikachu", CurrentHP: 80, MaxHP: 100}
	poke2 := Pokémon{Name: "Charizard", CurrentHP: 90, MaxHP: 100}
	tracker.AddPokemonToTeam("p1", poke1)
	tracker.AddPokemonToTeam("p1", poke2)

	// Switch to first pokemon
	tracker.SwitchPokemon("p1", "Pikachu", 80)
	if tracker.activePokemon["p1"].Name != "Pikachu" {
		t.Errorf("expected Pikachu to be active, got %s", tracker.activePokemon["p1"].Name)
	}

	// Switch to second pokemon
	tracker.SwitchPokemon("p1", "Charizard", 90)
	if tracker.activePokemon["p1"].Name != "Charizard" {
		t.Errorf("expected Charizard to be active, got %s", tracker.activePokemon["p1"].Name)
	}

	// Try to switch to non-existent pokemon (should not panic)
	tracker.SwitchPokemon("p1", "Blastoise", 100)
}

func TestUpdatePokemonHPEdgeCases(t *testing.T) {
	tracker := NewStateTracker()

	// Add pokemon
	poke := Pokémon{Name: "Pikachu", CurrentHP: 100, MaxHP: 100}
	tracker.AddPokemonToTeam("p1", poke)
	tracker.SwitchPokemon("p1", "Pikachu", 100)

	// Update HP normally
	tracker.UpdatePokemonHP("p1", 50, 100)
	if tracker.activePokemon["p1"].CurrentHP != 50 {
		t.Errorf("expected HP 50, got %d", tracker.activePokemon["p1"].CurrentHP)
	}

	// Update HP when MaxHP is 0 (should set it)
	tracker.activePokemon["p1"].MaxHP = 0
	tracker.UpdatePokemonHP("p1", 80, 120)
	if tracker.activePokemon["p1"].MaxHP != 120 {
		t.Errorf("expected MaxHP 120, got %d", tracker.activePokemon["p1"].MaxHP)
	}

	// Update non-existent player (should not panic)
	tracker.UpdatePokemonHP("p3", 50, 100)
}

func TestGenerateUUIDFormat(t *testing.T) {
	uuid1 := generateUUID()
	uuid2 := generateUUID()

	// Should generate different UUIDs
	if uuid1 == uuid2 {
		t.Error("expected different UUIDs")
	}

	// Should have correct format (8-4-4-4-12 hex digits with dashes)
	if len(uuid1) != 36 {
		t.Errorf("expected UUID length 36, got %d", len(uuid1))
	}

	// Check for dashes at correct positions
	if uuid1[8] != '-' || uuid1[13] != '-' || uuid1[18] != '-' || uuid1[23] != '-' {
		t.Errorf("UUID has incorrect format: %s", uuid1)
	}
}
