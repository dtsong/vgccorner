# Turn-by-Turn Analysis Implementation

## Overview

VGC Corner now includes comprehensive turn-by-turn battle analysis with team classification, move impact tracking, and detailed battle state visualization.

## Backend Implementation (Golang)

### 1. Team Classification System

**File**: `backend/internal/analysis/team_classifier.go`

Automatically classifies teams into VGC archetypes:
- Hard Trick Room
- TailRoom
- Sun Offense
- Rain Offense
- Balance Bros (Incineroar + Rillaboom)
- Psy-Spam (Psychic Terrain + Expanding Force)
- Tailwind Hyper Offense
- Generic archetypes (Tailwind, Trick Room, Weather)
- Unclassified

**Key Functions**:
- `ClassifyTeam(team []Pokémon)` - Analyzes team composition
- `GetArchetypeDescription(archetype string)` - Returns human-readable descriptions

### 2. Enhanced Type System

**File**: `backend/internal/analysis/types.go`

**New Types**:
```go
type TeamClassification struct {
    Archetype       string
    HasTrickRoom    bool
    TrickRoomUsers  []string
    HasTailwind     bool
    TailwindUsers   []string
    HasWeatherSetter bool
    WeatherType     string
    // ... more fields
}

type MoveImpact struct {
    DamageDealt     int
    HealingDone     int
    StatusInflicted string
    SpeedControl    string
    WeatherSet      string
    TerrainSet      string
    FakeOut         bool
    Protect         bool
    StatChanges     []StatChange
    Fainted         []string
    Critical        bool
    Effectiveness   string
    Missed          bool
}

type StatChange struct {
    Pokemon string
    Stat    string
    Stages  int
}
```

**Enhanced Action Type**:
```go
type Action struct {
    Player       string
    ActionType   string
    Pokemon      string
    Move         *Move
    SwitchTo     string
    Item         string
    Target       string
    Result       string
    Details      string
    Impact       *MoveImpact
    OrderInTurn  int
}
```

### 3. Move Impact Tracking

**File**: `backend/internal/analysis/move_impact.go`

Tracks detailed impact of moves:
- Damage dealt and healing done
- Speed control (Trick Room, Tailwind, paralysis, flinch)
- Weather and terrain changes
- Status conditions inflicted
- Critical hits and type effectiveness
- Fake Out and Protect usage
- Stat boosts/drops
- Pokémon faints

**Key Function**:
```go
EnhanceActionWithImpact(action *Action, moveName string, events []string)
```

### 4. Turn Analysis API Endpoint

**File**: `backend/internal/httpapi/turn_handlers.go`

**New Endpoint**: `GET /api/showdown/replays/{replayId}/turns`

**Response Structure**:
```json
{
  "status": "success",
  "battleId": "gen9vgc2024regh-123",
  "format": "gen9vgc2024regh",
  "player1": "Player1",
  "player2": "Player2",
  "winner": "Player1",
  "turns": [
    {
      "turnNumber": 1,
      "events": [
        {
          "type": "move",
          "pokemon": "Gengar",
          "action": "used Shadow Ball",
          "target": "on Dusclops",
          "result": "critical-hit",
          "details": "Critical Hit",
          "playerSide": "player1"
        }
      ],
      "boardState": {
        "player1Active": [...],
        "player2Active": [...]
      }
    }
  ],
  "archetypes": {
    "player1": {
      "archetype": "Hard Trick Room",
      "description": "A team built around Trick Room...",
      "tags": []
    },
    "player2": {
      "archetype": "Tailwind Hyper Offense",
      "description": "An aggressive team using...",
      "tags": []
    }
  }
}
```

### 5. Parser Enhancements

**File**: `backend/internal/analysis/parser.go`

Updated to:
- Call team classification after parsing
- Track detailed turn order
- Populate team archetype fields

### 6. Tests

**File**: `backend/internal/analysis/team_classifier_test.go`

Comprehensive test suite covering all team archetypes:
- `TestClassifyTeam_HardTrickRoom`
- `TestClassifyTeam_TailRoom`
- `TestClassifyTeam_SunOffense`
- `TestClassifyTeam_RainOffense`
- `TestClassifyTeam_BalanceBros`
- `TestClassifyTeam_PsySpam`
- `TestClassifyTeam_TailwindHyperOffense`
- `TestClassifyTeam_Unclassified`

## Frontend Implementation (Next.js/React)

### 1. Turn Analysis Page

**File**: `frontend/src/app/replay/[replayId]/analysis/page.tsx`

Features:
- Turn-by-turn navigation (Previous/Next buttons)
- Turn counter showing progress
- Mini-boardstate visualization
- Event list with color coding
- Winner indicator
- Loading and error states

**Route**: `/replay/[replayId]/analysis`

### 2. Mini-Boardstate Component

**File**: `frontend/src/components/battles/MiniBoardstate.tsx`

Visual representation of battle field:
- Green background for player 1's side
- Red/pink background for player 2's side
- Pokémon sprites from PokémonDB
- "Lead" badges for starting Pokémon
- HP percentages
- Arrow between sides

### 3. Turn Events Component

**File**: `frontend/src/components/battles/TurnEvents.tsx`

Color-coded event display:
- **Green** = Favorable for player 1 (critical hits, super effective)
- **Red/Pink** = Favorable for player 2
- **Gray** = Neutral events
- Pokémon icons next to each action
- Result badges (Critical Hit, Super Effective, etc.)

### 4. Type Definitions

**File**: `frontend/src/lib/types/battle.ts`

TypeScript types for battle analysis:
```typescript
interface TurnData {
  turnNumber: number;
  events: BattleEvent[];
  boardState: BoardState;
}

interface BattleEvent {
  type: 'move' | 'switch' | 'faint' | ...;
  pokemon: string;
  action: string;
  target?: string;
  result?: EventResult;
  details?: string;
  playerSide: 'player1' | 'player2';
}

interface MoveImpact {
  damageDealt: number;
  healingDone: number;
  speedControl: string;
  fakeOut: boolean;
  // ... more fields
}
```

### 5. API Client Updates

**File**: `frontend/src/lib/api/client.ts`

New function:
```typescript
async function getTurnAnalysis(replayId: string): Promise<TurnAnalysisResponse>
```

Calls: `GET /api/showdown/replays/{replayId}/turns`

### 6. Navigation Button

**File**: `frontend/src/app/replay/[replayId]/page.tsx`

Added "Turn-by-Turn Analysis" button to replay detail page:
- Blue button with chart icon
- Prominently placed in header
- Navigates to `/replay/[replayId]/analysis`

## User Flow

1. **Landing Page**: User pastes Pokémon Showdown replay URL
2. **Analyze**: Click "Analyze" button
3. **Replay Detail**: View team composition, battle stats, winner
4. **Turn Analysis**: Click "Turn-by-Turn Analysis" button
5. **Navigate Turns**: Use Previous/Next buttons to step through battle
6. **View Details**: See board state and events for each turn

## Key Features

### Turn Analysis Display

- **Turn order tracking**: Actions shown in the order they occurred
- **Move impact**: Displays damage, speed control, weather, status effects
- **Visual board state**: Shows active Pokémon with HP
- **Color-coded effectiveness**: Visual indicators for critical hits, super effective moves
- **Fake Out detection**: Highlights flinch-inducing moves
- **Speed control tracking**: Identifies Trick Room, Tailwind, paralysis
- **Weather/terrain tracking**: Shows when weather/terrain is set
- **Stat changes**: Tracks boosts and drops

### Team Classification

Teams are automatically classified with:
- Primary archetype (Hard Trick Room, Psy-Spam, etc.)
- Detailed breakdown of criteria met
- Human-readable descriptions
- Lists of key Pokémon for each strategy component

## Testing

### Backend Tests

Run backend tests:
```bash
cd backend
go test ./internal/analysis/...
```

Specifically test team classification:
```bash
go test ./internal/analysis/ -run TestClassifyTeam
```

### Frontend Tests

Run frontend tests:
```bash
cd frontend
npm test
```

## Next Steps

### Backend TODO:
- [ ] Integrate with database to store/retrieve turn analysis
- [ ] Implement actual replay fetching from Pokémon Showdown
- [ ] Add caching for analyzed battles
- [ ] Track Tera type usage in move impact
- [ ] Detect Follow Me + Setup patterns
- [ ] Calculate damage ranges

### Frontend TODO:
- [ ] Add keyboard shortcuts for turn navigation (arrow keys)
- [ ] Add turn slider for quick navigation
- [ ] Display team archetypes on analysis page
- [ ] Add tooltips for move effects
- [ ] Show damage calculations
- [ ] Add replay playback speed control
- [ ] Export turn analysis to PDF/image

## API Endpoints Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/showdown/analyze` | Analyze a new replay |
| GET | `/api/showdown/replays` | List recent replays |
| GET | `/api/showdown/replays/{id}` | Get replay details |
| **GET** | **`/api/showdown/replays/{id}/turns`** | **Get turn-by-turn analysis** |
| POST | `/api/tcglive/analyze` | Analyze TCG Live battle (planned) |

## Documentation

- **Team Classification**: See `TEAM_CLASSIFICATION.md`
- **Backend Setup**: See `backend/README.md`
- **Frontend Setup**: See `frontend/SETUP.md`

## Dependencies

### Backend
- Go 1.21+
- github.com/go-chi/chi/v5 (routing)

### Frontend
- Next.js 16.0.3
- React 19.2.0
- TypeScript 5+
- Tailwind CSS 4

## Conclusion

The turn-by-turn analysis system provides comprehensive battle insights including:
- Automated team archetype detection
- Detailed move impact tracking
- Visual board state representation
- Turn-by-turn event replay
- Color-coded effectiveness indicators

This foundation enables coaches and players to deeply analyze their battles and understand team strategies.
