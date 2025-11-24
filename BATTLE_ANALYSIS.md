# Pokemon Showdown Battle Analysis System

## Overview

VGC Corner now includes a comprehensive battle analysis system that parses Pokemon Showdown replay logs and provides deep insights into battle progression, critical moments, and strategic decisions. This document describes the complete implementation across backend and frontend.

## System Architecture

```
User uploads Showdown .log file
        │
        ▼
POST /api/showdown/analyze
        │
        ▼
┌──────────────────────────┐
│  ParseShowdownLog()      │
│  - Parse log format      │
│  - Extract metadata      │
│  - Extract teams         │
└──────────────────────────┘
        │
        ▼
┌──────────────────────────┐
│  StateTracker            │
│  - Reconstruct battle    │
│  - Track HP/status       │
│  - Calculate scores      │
└──────────────────────────┘
        │
        ▼
┌──────────────────────────┐
│  Analysis Engine         │
│  - Detect turning points │
│  - Calculate statistics  │
│  - Identify key moments  │
└──────────────────────────┘
        │
        ▼
  BattleSummary JSON
        │
        ▼
Frontend Dashboard
- Position score timeline
- Key moments panel
- Team comparison
- Battle statistics
```

## Key Components

### 1. Backend Parser (`internal/analysis/parser.go`)

#### Core Function: `ParseShowdownLog(logContent string) (*BattleSummary, error)`

The parser processes Showdown replay logs in pipe-delimited format and extracts:

- **Metadata**: Format, timestamp, player names, team sizes
- **Team Information**: Pokemon species, levels, gender, abilities, items
- **Battle Events**: Moves, switches, faints, HP changes, status conditions
- **Field Effects**: Tailwind, terrain, weather
- **Stat Changes**: Boosts and unboosts during battle

#### Log Format Examples

```
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|player|p1|Player1|giovanni|1487
|player|p2|Player2|steven|1398
|poke|p1|Ursaluna-Bloodmoon, L50, M|
|teamsize|p1|4
|start
|switch|p1a: Whimsicott|Whimsicott, L50, M|100\/100
|move|p1a: Whimsicott|Tailwind|p1a: Whimsicott
|-damage|p1b: Typhlosion|0 fnt
|faint|p1b: Typhlosion
|turn|1
```

### 2. StateTracker

Maintains complete battle state throughout parsing:

```go
type StateTracker struct {
    playerNames        map[string]string
    teamSizes          map[string]int
    teams              map[string][]Pokémon
    activePokemon      map[string]*Pokémon
    losses             map[string]int
    fieldEffects       map[string][]string
    statBoosts         map[string]map[string]int
}
```

**Key Methods:**
- `SwitchPokemon(playerID, pokeName string, hp int)` - Track active Pokemon changes
- `UpdatePokemonHP(playerID string, currentHP, maxHP int)` - Update HP values
- `FaintPokemon(playerID string)` - Record faints
- `CalculatePositionScore() *PositionScore` - Evaluate battle position

### 3. Position Scoring Algorithm

Calculates a 0-100 position score for each player after every turn:

```
Player Score = (Active Pokemon HP% × 0.6) + (Remaining Team% × 0.4)
```

**Example:**
- Player has 1 active Pokemon at 50% HP
- Player has 3 Pokemon left out of 4 team size
- Score = (50 × 0.6) + ((3/4 × 100) × 0.4) = 30 + 30 = 60

### 4. Turning Point Detection

Identifies critical momentum shifts:

```go
func detectTurningPoints(summary *BattleSummary)
```

- Compares position scores between consecutive turns
- Flags shifts of 15+ points as significant
- Calculates significance rating (1-10 scale)
- Creates TurningPoint structs with momentum data

**Triggering Threshold:** 15+ point momentum shift

## Data Structures

### BattleSummary

Complete battle analysis package:

```go
type BattleSummary struct {
    ID         string
    Format     string
    Timestamp  time.Time
    Duration   int
    Player1    Player
    Player2    Player
    Winner     string
    Turns      []Turn
    Stats      BattleStats
    KeyMoments []KeyMoment
}
```

### Turn

Per-turn battle state and events:

```go
type Turn struct {
    TurnNumber    int
    Actions       []Action
    StateAfter    BattleState
    DamageDealt   map[string]int
    HealingDone   map[string]int
    PositionScore *PositionScore
}
```

### PositionScore

Post-turn evaluation:

```go
type PositionScore struct {
    Player1Score  float64 // 0-100
    Player2Score  float64 // 0-100
    MomentumPlayer string  // "player1", "player2", or "neutral"
}
```

### TurningPoint

Critical momentum shifts:

```go
type TurningPoint struct {
    TurnNumber    int
    Score1Before  float64
    Score1After   float64
    Score2Before  float64
    Score2After   float64
    MomentumShift float64 // Negative = P2 gained, Positive = P1 gained
    Significance  int     // 1-10 scale
    Description   string
}
```

### KeyMoment

Important events in battle:

```go
type KeyMoment struct {
    TurnNumber   int
    Description  string
    Type         string // "switch", "ko", "status", "weather", "turning_point"
    Significance int    // 1-10 scale
}
```

## API Endpoint

### POST /api/showdown/analyze

**Request:**
```json
{
  "analysisType": "rawLog",
  "rawLog": "[pipe-delimited battle log content]",
  "isPrivate": false
}
```

**Response:**
```json
{
  "status": "success",
  "battleId": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "format": "[Gen 9] VGC 2025 Reg H (Bo3)",
    "timestamp": "2025-11-24T10:00:00Z",
    "duration": 300,
    "player1": { ... },
    "player2": { ... },
    "winner": "player1",
    "turns": [ ... ],
    "stats": { ... },
    "keyMoments": [ ... ]
  },
  "metadata": {
    "parseTimeMs": 45,
    "analysisTimeMs": 120,
    "cached": false
  }
}
```

## Frontend Components

### BattleAnalysisDashboard

Main container component orchestrating all subcomponents:

- **BattleHeader** - Player names, winner, format, metadata
- **PositionScoreChart** - SVG line chart of position scores
- **TurnTimeline** - Expandable list of turns with actions and scores
- **KeyMomentsPanel** - Highlighted critical moments and turning points
- **TeamComparison** - Side-by-side team status and Pokemon HP
- **BattleStatistics** - Aggregated stats, move frequency, effectiveness

### PositionScoreChart

Lightweight SVG chart showing position score evolution:

- X-axis: Turn numbers
- Y-axis: Position score (0-100)
- Two lines: Player 1 (blue) and Player 2 (red)
- Grid lines at 25-point intervals
- Data points marked with circles
- Legend and average score annotations

### TurnTimeline

Interactive expandable turns:

- **Collapsed View**: Turn number, action summary, position scores
- **Expanded View**: Detailed actions, damage/healing, position scores
- **Color Coding**: Momentum indicators (blue/red/neutral)
- **Significance Badges**: Showing who has advantage

### KeyMomentsPanel

Sorted chronological list of critical events:

- Merges KeyMoments and TurningPoints
- Color-coded by significance (red/orange/yellow/gray)
- Icons for event types (💥 KO, 🔄 Switch, ⚠️ Status, 🌦️ Weather, 📈 Turning Point)
- Shows momentum shift values for turning points
- Significance rating (1-10 scale)

### TeamComparison

Side-by-side team viewer:

- Pokemon name, level, status
- HP bar with color gradient (green/yellow/red)
- Ability, item, status condition badges
- Tera type indicator
- Faint status marking
- Team remainder counter

### BattleStatistics

Aggregated battle metrics:

- Total turns, switches, critical hits
- Type effectiveness counts
- Damage per turn
- Most used moves (top 5)
- Player-by-player comparison

## TypeScript Type System

Complete type definitions in `frontend/src/lib/types/showdown.ts`:

- `BattleSummary` - Root battle data
- `Player` - Player information
- `Pokémon` - Pokemon with HP, status, teratype
- `Turn` - Turn events and position score
- `PositionScore` - Position evaluation
- `TurningPoint` - Momentum shifts
- `KeyMoment` - Important events
- `BattleStats` - Aggregated statistics
- `Move`, `Stats`, `Action`, `BattleState` - Supporting types

Full type safety ensures frontend components correctly consume backend data.

## Key Features

### 1. Comprehensive Battle Reconstruction
- Parses complete battle state from log file
- Tracks all Pokemon HP, status, items, abilities
- Maintains accurate team composition changes

### 2. Position Scoring
- Evaluates each player's position after every turn
- Weighs active Pokemon health (60%) and team remaining (40%)
- Enables momentum analysis

### 3. Automatic Turning Point Detection
- Identifies turns where battle momentum shifted significantly
- Flags 15+ point shifts as critical moments
- Assigns significance ratings for prioritization

### 4. Rich Event Tracking
- Captures all battle events (moves, switches, faints)
- Records damage, healing, stat changes
- Tracks field effects and status conditions

### 5. Interactive Visualization
- Expandable turn timeline for detailed exploration
- SVG position score chart
- Side-by-side team comparison
- Chronological key moments panel
- Comprehensive statistics dashboard

## Usage Example

### Backend

```go
battleLog := readFromFile("replay.log")
summary, err := analysis.ParseShowdownLog(battleLog)
if err != nil {
    log.Fatal(err)
}

// Access battle data
fmt.Printf("Winner: %s\n", summary.Winner)
fmt.Printf("Total Turns: %d\n", summary.Stats.TotalTurns)
fmt.Printf("Turning Points: %d\n", len(summary.Stats.TurningPoints))
```

### Frontend

```typescript
import BattleAnalysisDashboard from '@/components/battles/BattleAnalysisDashboard';

// Fetch analysis from API
const response = await fetch('/api/showdown/analyze', {
  method: 'POST',
  body: JSON.stringify({ analysisType: 'rawLog', rawLog: logContent }),
});
const { data: battle } = await response.json();

// Render dashboard
<BattleAnalysisDashboard battle={battle} />
```

## Future Enhancements

### Short-term
- [ ] Store analyses in database for historical tracking
- [ ] Implement replay ID and username analysis modes
- [ ] Add move damage calculations using Pokedex data
- [ ] Track weather and terrain effects in scoring
- [ ] Export battle analysis as PDF report

### Medium-term
- [ ] Type advantage weighting in position scores
- [ ] AI-generated insights and strategic recommendations
- [ ] Comparison with player's historical performance
- [ ] Team synergy analysis
- [ ] Move effectiveness heatmaps

### Long-term
- [ ] Machine learning-based turning point prediction
- [ ] Automated coaching suggestions
- [ ] Integration with Smogon Dex for comprehensive team analysis
- [ ] Community battle sharing and analysis browsing
- [ ] Real-time battle analysis streaming

## Technical Debt & Notes

1. **HP Calculations**: Currently uses approximate HP values. Should integrate actual EV/IV calculations when available.

2. **Damage Tracking**: Position score implementation could be enhanced with:
   - Type advantage calculations
   - Stat boost integration
   - Ability effects (Speed reduction, Sp. Atk drop, etc.)

3. **Performance**: For very long battles (50+ turns), consider:
   - Pagination of turn timeline
   - Lazy loading of turn details
   - Optimized SVG rendering

4. **Error Handling**: Currently lenient with invalid log formats. Could add strict mode for validation.

## Testing

Run backend tests:
```bash
cd backend
go test ./internal/analysis -v
```

Run frontend components:
```bash
cd frontend
npm test -- components/battles
```

## Related Files

- Backend: `backend/internal/analysis/`
  - `types.go` - Data structures
  - `parser.go` - Parsing logic
  - `parser_test.go` - Unit tests

- Frontend: `frontend/src/`
  - `lib/types/showdown.ts` - TypeScript types
  - `components/battles/` - React components
    - `BattleAnalysisDashboard.tsx` - Main container
    - `BattleHeader.tsx` - Header display
    - `PositionScoreChart.tsx` - SVG chart
    - `TurnTimeline.tsx` - Turn list
    - `KeyMomentsPanel.tsx` - Event list
    - `TeamComparison.tsx` - Team viewer
    - `BattleStatistics.tsx` - Stats display

- API: `backend/internal/httpapi/showdown_handlers.go`
  - `handleAnalyzeShowdown()` - Analysis endpoint
