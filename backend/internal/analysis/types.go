package analysis

import "time"

// BattleSummary represents the complete analysis of a Pokémon battle.
type BattleSummary struct {
	// Metadata about the battle
	ID        string    `json:"id"`
	Format    string    `json:"format"`        // e.g., "Regulation H"
	Timestamp time.Time `json:"timestamp"`
	Duration  int       `json:"duration"`      // in seconds

	// Player information
	Player1 Player `json:"player1"`
	Player2 Player `json:"player2"`
	Winner  string `json:"winner"` // "player1", "player2", or "draw"

	// Battle progression
	Turns []Turn `json:"turns"`

	// Overall statistics
	Stats BattleStats `json:"stats"`

	// Key moments and highlights
	KeyMoments []KeyMoment `json:"keyMoments"`
}

// Player represents a single player in the battle.
type Player struct {
	Name      string     `json:"name"`
	Team      []Pokémon  `json:"team"`
	Active    *Pokémon   `json:"active"` // Currently active Pokémon
	Losses    int        `json:"losses"` // Number of fainted Pokémon
	TotalLeft int        `json:"totalLeft"` // Total Pokémon still in battle
}

// Pokémon represents a single Pokémon with its stats and moves.
type Pokémon struct {
	ID       string `json:"id"`       // e.g., "pikachu"
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Gender   string `json:"gender"`   // "M", "F", or ""
	Ability  string `json:"ability"`
	Item     string `json:"item"`
	Stats    Stats  `json:"stats"`    // Base stats
	Moves    []Move `json:"moves"`
	Happiness int   `json:"happiness"` // 0-255
	Shiny    bool   `json:"shiny"`
}

// Move represents a move a Pokémon knows.
type Move struct {
	ID     string `json:"id"`     // e.g., "thunderbolt"
	Name   string `json:"name"`
	Type   string `json:"type"`   // e.g., "Electric"
	Power  int    `json:"power"`  // 0 if N/A
	Accuracy int   `json:"accuracy"` // 0-100, 0 if N/A
	PP     int    `json:"pp"`     // Power Points
}

// Stats represents Pokémon base stats.
type Stats struct {
	HP     int `json:"hp"`
	Attack int `json:"attack"`
	Defense int `json:"defense"`
	SpAtk   int `json:"spAtk"`
	SpDef   int `json:"spDef"`
	Speed   int `json:"speed"`
}

// Turn represents a single turn in the battle.
type Turn struct {
	TurnNumber    int              `json:"turnNumber"`
	Actions       []Action         `json:"actions"`
	StateAfter    BattleState      `json:"stateAfter"`
	DamageDealt   map[string]int   `json:"damageDealt"`   // Player name -> damage dealt
	HealingDone   map[string]int   `json:"healingDone"`   // Player name -> healing done
}

// Action represents an action taken by a player during a turn.
type Action struct {
	Player     string `json:"player"`    // "player1" or "player2"
	ActionType string `json:"actionType"` // "move", "switch", "item"
	Move       *Move  `json:"move,omitempty"`
	SwitchTo   string `json:"switchTo,omitempty"` // Pokémon name if switch
	Item       string `json:"item,omitempty"`     // Item used if item action
}

// BattleState represents the state of the battle at a point in time.
type BattleState struct {
	Player1Active *Pokémon `json:"player1Active"`
	Player2Active *Pokémon `json:"player2Active"`
	Player1Team   []string `json:"player1Team"`   // List of alive Pokémon names
	Player2Team   []string `json:"player2Team"`
}

// BattleStats represents aggregate statistics about the battle.
type BattleStats struct {
	TotalTurns        int                `json:"totalTurns"`
	MoveFrequency     map[string]int     `json:"moveFrequency"`     // Move ID -> count
	TypeCoverage      map[string]int     `json:"typeCoverage"`      // Type -> count
	Switch            int                `json:"switches"`          // Total switches by both players
	CriticalHits      int                `json:"criticalHits"`
	SuperEffective    int                `json:"superEffective"`
	NotVeryEffective  int                `json:"notVeryEffective"`
	AvgDamagePerTurn  float64            `json:"avgDamagePerTurn"`
	AvgHealPerTurn    float64            `json:"avgHealPerTurn"`
	Player1Stats      PlayerStats        `json:"player1Stats"`
	Player2Stats      PlayerStats        `json:"player2Stats"`
}

// PlayerStats represents stats for an individual player.
type PlayerStats struct {
	MoveCount       int            `json:"moveCount"`
	SwitchCount     int            `json:"switchCount"`
	DamageDealt     int            `json:"damageDealt"`
	DamageTaken     int            `json:"damageTaken"`
	HealingDone     int            `json:"healingDone"`
	HealingReceived int            `json:"healingReceived"`
	MovesByType     map[string]int `json:"movesByType"`     // Type -> count
	Effectiveness   EffectivenessStats `json:"effectiveness"`
}

// EffectivenessStats tracks type effectiveness in the battle.
type EffectivenessStats struct {
	SuperEffective   int `json:"superEffective"`
	NotVeryEffective int `json:"notVeryEffective"`
	Neutral          int `json:"neutral"`
}

// KeyMoment represents a significant moment in the battle.
type KeyMoment struct {
	TurnNumber  int    `json:"turnNumber"`
	Description string `json:"description"` // e.g., "Player 2 switched to Charizard"
	Type        string `json:"type"`        // "switch", "kO", "status", "weather", etc.
	Significance int    `json:"significance"` // 1-10 scale
}
