# Detailed Implementation Guide: Pokemon Showdown Battle Analysis

## System Overview

The battle analysis system transforms raw Pokemon Showdown replay logs into actionable insights through four-stage processing:

```
Raw Log File → Parse & Reconstruct → Analyze & Score → Visualize
```

## Stage 1: Parsing & Reconstruction

### What Happens

The `ParseShowdownLog()` function reads a pipe-delimited battle log and extracts:

```
|tier|[Gen 9] VGC 2025 Reg H (Bo3)       → Format metadata
|player|p1|Player1|giovanni|1487         → Player information
|poke|p1|Ursaluna-Bloodmoon, L50, M|    → Team composition
|teamsize|p1|4                            → Initial team size
|start                                    → Battle begins
|switch|p1a: Whimsicott|...|100\/100    → Pokemon swap
|move|p1a: Whimsicott|Tailwind|...      → Move execution
|-damage|p1b: Typhlosion|0 fnt            → HP/fainting
|faint|p1b: Typhlosion                    → Confirm faint
|turn|1                                   → Turn separator
```

### Implementation Details

```go
func ParseShowdownLog(logContent string) (*BattleSummary, error) {
    // Pass 1: Extract metadata and teams
    for each line {
        if "tier" → Format
        if "player" → Player names
        if "poke" → Team members
        if "teamsize" → Team size tracking
    }

    // Initialize StateTracker with teams

    // Pass 2: Process all battle events
    for each line {
        if "turn" → Save previous turn, start new one
        if "switch" → Track active Pokemon change
        if "move" → Record action
        if "-damage" → Update HP
        if "-heal" → Update HP
        if "faint" → Record fainting
        if "-boost/-unboost" → Track stat changes
        if "-status" → Track conditions
        if "-terastallize" → Track tera type
        if "-sidestart/-sideend" → Track field effects
        if "win" → Record winner
    }

    // Calculate all statistics
    // Detect turning points
    return summary
}
```

### Example Execution

**Input Log (simplified):**
```
|tier|VGC 2025
|player|p1|Alice|1500
|player|p2|Bob|1400
|poke|p1|Ursaluna, L50|
|poke|p1|Whimsicott, L50|
|poke|p2|Gholdengo, L50|
|poke|p2|Dragonite, L50|
|teamsize|p1|2
|teamsize|p2|2
|start
|switch|p1a: Ursaluna|Ursaluna, L50|100\/100
|switch|p2a: Gholdengo|Gholdengo, L50|100\/100
|turn|1
|move|p1a: Ursaluna|Protect|p1a: Ursaluna
|move|p2a: Gholdengo|Nasty Plot|p2a: Gholdengo
|-boost|p2a: Gholdengo|spa|2
|turn|2
|move|p1a: Ursaluna|Hyper Voice|p2a: Gholdengo
|-damage|p2a: Gholdengo|50\/100
|move|p2a: Gholdengo|Shadow Ball|p1a: Ursaluna
|-supereffective|p1a: Ursaluna
|-damage|p1a: Ursaluna|25\/100
|turn|3
|move|p1a: Ursaluna|Hyper Voice|p2a: Gholdengo
|-damage|p2a: Gholdengo|0 fnt
|faint|p2a: Gholdengo
|switch|p2a: Dragonite|Dragonite, L50|100\/100
|move|p2a: Dragonite|Extreme Speed|p1a: Ursaluna
|-damage|p1a: Ursaluna|0 fnt
|faint|p1a: Ursaluna
|win|Bob
```

**Parsing Output:**
```go
BattleSummary{
    Format: "VGC 2025",
    Player1: Player{
        Name: "Alice",
        Team: [Ursaluna(L50), Whimsicott(L50)],
        Losses: 2
    },
    Player2: Player{
        Name: "Bob",
        Team: [Gholdengo(L50), Dragonite(L50)],
        Losses: 1
    },
    Winner: "player2",
    Turns: [
        Turn 1: {
            Actions: [
                Action{player: "player1", type: "move", move: "Protect"},
                Action{player: "player2", type: "move", move: "Nasty Plot"}
            ],
            PositionScore: {Player1: 100, Player2: 102, Momentum: "player2"}
        },
        Turn 2: { ... },
        Turn 3: { ... }
    ]
}
```

## Stage 2: State Tracking

### StateTracker Purpose

Maintains the battle state as events are processed, enabling accurate position scoring.

### Critical State Variables

```go
activePokemon["p1"]  // Pointer to current active Pokemon
activePokemon["p2"]  // Updated on each switch
losses["p1"]         // Count of fainted Pokemon
teamSizes["p1"]      // Original team size
```

### HP Tracking

```go
// When HP change event occurs:
if activePokemon["p1"] != nil {
    activePokemon["p1"].CurrentHP = newHP
    activePokemon["p1"].MaxHP = maxHP  // if not set
}
```

### State Update Sequence

```
Event: |switch|p1a: Ursaluna|...|100\/100
  ↓
Extract: playerID="p1", name="Ursaluna", hp=100
  ↓
StateTracker.SwitchPokemon("p1", "Ursaluna", 100)
  ↓
activePokemon["p1"] = &teams["p1"][0]  // Find by name
  ↓
teams["p1"][0].CurrentHP = 100

Event: |-damage|p1a: Ursaluna|50\/100
  ↓
Extract: playerID="p1", hp=50, maxHP=100
  ↓
StateTracker.UpdatePokemonHP("p1", 50, 100)
  ↓
activePokemon["p1"].CurrentHP = 50
activePokemon["p1"].MaxHP = 100  // if was 0

Event: |faint|p1a: Ursaluna
  ↓
Extract: playerID="p1"
  ↓
StateTracker.FaintPokemon("p1")
  ↓
activePokemon["p1"].CurrentHP = 0
losses["p1"]++  // Now 2 total losses
```

## Stage 3: Position Scoring

### Algorithm

```
For each turn:
    activeP1HP = active Pokemon HP% (or 0 if fainted)
    activeP2HP = active Pokemon HP% (or 0 if fainted)
    teamP1Remaining = (teamSize - losses) / teamSize * 100
    teamP2Remaining = (teamSize - losses) / teamSize * 100

    P1 Score = (activeP1HP × 0.6) + (teamP1Remaining × 0.4)
    P2 Score = (activeP2HP × 0.6) + (teamP2Remaining × 0.4)
```

### Example Calculation

**Scenario 1: Early game (both teams full)**
- P1: Active Pokemon 100% HP, 4 team members remaining
- P2: Active Pokemon 100% HP, 4 team members remaining

```
P1Score = (100 × 0.6) + (100 × 0.4) = 60 + 40 = 100
P2Score = (100 × 0.6) + (100 × 0.4) = 60 + 40 = 100
Result: Neutral (tied at 100)
```

**Scenario 2: Mid-game (one team ahead)**
- P1: Active Pokemon 75% HP, 3 team members remaining (out of 4)
- P2: Active Pokemon 50% HP, 2 team members remaining (out of 4)

```
P1Score = (75 × 0.6) + (75 × 0.4) = 45 + 30 = 75
P2Score = (50 × 0.6) + (50 × 0.4) = 30 + 20 = 50
Result: P1 Advantage (25 points)
```

**Scenario 3: Late game (one Pokemon left)**
- P1: Active Pokemon 100% HP, 1 team member remaining (out of 4)
- P2: Active Pokemon 100% HP, 2 team members remaining (out of 4)

```
P1Score = (100 × 0.6) + (25 × 0.4) = 60 + 10 = 70
P2Score = (100 × 0.6) + (50 × 0.4) = 60 + 20 = 80
Result: P2 Advantage (10 points)
```

### Why This Weighting?

- **60% Active Pokemon HP**: Immediate battle threat
  - A Pokemon at low HP can be knocked out quickly
  - A fresh Pokemon can control the board

- **40% Team Remaining**: Strategic depth
  - More Pokemon means more switches and flexibility
  - Last Pokemon leaves no room for mistakes

## Stage 4: Turning Point Detection

### Algorithm

```go
func detectTurningPoints(summary *BattleSummary) {
    for each turn (starting from turn 2) {
        prevScore := turns[i-1].positionScore
        currScore := turns[i].positionScore

        p1Delta = currScore.p1 - prevScore.p1
        p2Delta = currScore.p2 - prevScore.p2

        momentumShift = p1Delta - p2Delta

        if abs(momentumShift) >= 15 {  // Threshold: 15+ points
            significance = int(abs(momentumShift) / 10)
            if significance > 10 { significance = 10 }

            turningPoint := TurningPoint{
                turn: i,
                score1Before: prevScore.p1,
                score1After: currScore.p1,
                score2Before: prevScore.p2,
                score2After: currScore.p2,
                momentumShift: momentumShift,
                significance: significance,
                description: describe(momentumShift)
            }

            addToTurningPoints(turningPoint)
        }
    }
}
```

### Example Turning Point

**Turn 2 to Turn 3:**

Turn 2 End:
```
P1Score: 85 (Ursaluna at 70% HP, 3/4 team)
P2Score: 65 (Gholdengo at 50% HP, 4/4 team)
Momentum: P1 ahead by 20 points
```

Turn 3 Happens: "P2 KOs P1's Pokemon, P1 switches in weakened Whimsicott"

Turn 3 End:
```
P1Score: 45 (Whimsicott at 30% HP, 2/4 team)  // Lost Pokemon + low HP
P2Score: 75 (Gholdengo at 50% HP, 4/4 team)   // Slight HP damage
Momentum: P2 ahead by 30 points
```

Momentum Shift Calculation:
```
P1 Delta = 45 - 85 = -40
P2 Delta = 75 - 65 = +10
MomentumShift = -40 - 10 = -50 (P2 gained 50 points)

Significance = int(50 / 10) = 5/10 (Major)
Type: Turning Point (critical KO)
```

## Stage 5: Statistics Aggregation

### Collected Statistics

```go
type BattleStats struct {
    TotalTurns        int                      // Turn count
    MoveFrequency     map[string]int           // Move -> Count
    Switches          int                      // Total switches both players
    CriticalHits      int                      // "-crit" events
    SuperEffective    int                      // "-supereffective" events
    NotVeryEffective  int                      // "-resisted" events
    AvgDamagePerTurn  float64                  // Total damage / turns
    Player1Stats      PlayerStats              // Per-player breakdown
    Player2Stats      PlayerStats              // Per-player breakdown
    TurningPoints     []TurningPoint           // Momentum shifts
}
```

### Calculation Process

```go
for each turn {
    for each action {
        if move {
            moveFrequency[moveID]++
            player.moveCount++
        }
        if switch {
            totalSwitches++
            player.switchCount++
        }
    }
}

for each "-crit" event { criticalHits++ }
for each "-supereffective" event { superEffective++ }
for each "-resisted" event { notVeryEffective++ }

totalDamage1, totalDamage2 = sum all damage
avgDamagePerTurn = (totalDamage1 + totalDamage2) / totalTurns
```

## Frontend Integration

### Data Flow

```
BattleAnalysisDashboard receives BattleSummary
    ↓
useMemo calculates derived data:
    - Average position scores
    - Dominant moments (significance >= 7)
    - Turning point count
    ↓
Passes to subcomponents:
    - BattleHeader(battle)
    - PositionScoreChart(battle.turns)
    - TurnTimeline(battle.turns)
    - KeyMomentsPanel(battle.keyMoments, battle.stats.turningPoints)
    - TeamComparison(battle.player1, battle.player2)
    - BattleStatistics(battle.stats)
```

### Component Responsibilities

| Component | Input | Output |
|-----------|-------|--------|
| BattleHeader | BattleSummary | Metadata display |
| PositionScoreChart | Turns[] | SVG chart |
| TurnTimeline | Turns[] | Expandable list |
| KeyMomentsPanel | KeyMoment[], TurningPoint[] | Chronological list |
| TeamComparison | Player, Player | Side-by-side view |
| BattleStatistics | BattleStats | Metrics display |

## Error Handling

### Parser Error Cases

1. **Malformed Log**
   - Missing required sections
   - Invalid HP format
   - → Returns error, doesn't crash

2. **Team Size Mismatch**
   - teamsize says 4 but only 3 Pokemon
   - → Uses actual team size

3. **Faint Without Active Pokemon**
   - Faint event but no active Pokemon set
   - → Skips (defensive programming)

### Frontend Error Cases

1. **Missing Position Scores**
   - Some turns might not have scores
   - → Charts filter out null values

2. **Empty Teams**
   - Player has no Pokemon
   - → Displays empty state

3. **Null Moments**
   - No key moments in battle
   - → Shows "No key moments recorded"

## Performance Optimization

### Current
- 500-line parser optimized for clarity
- Two-pass approach (metadata, then events)
- O(n) complexity where n = log lines

### Future Improvements
- Cache team Pokemon objects (avoid name lookups)
- Lazy evaluate statistics (only compute if requested)
- Stream parsing for very large logs
- Batch state updates to reduce object allocations

## Testing Strategy

### Backend Tests
1. Parse simple 3-turn battle
2. Verify all turns captured
3. Check position scores calculated
4. Confirm turning points detected

### Frontend Tests
1. Component renders with valid data
2. Expandable items toggle correctly
3. Charts display without errors
4. Responsive layout on mobile

## Integration Points

### API Gateway
- POST /api/showdown/analyze
- Takes rawLog parameter
- Returns BattleSummary JSON

### Database (future)
- Store BattleSummary in battles table
- Store TurningPoints separately for analytics
- Index by player name for quick lookup

### Pokedex Integration (future)
- Look up move type effectiveness
- Get actual base stats for Pokemon
- Enhance damage calculations

## Debugging Workflow

### "Parser isn't extracting my Pokemon"
1. Verify log has |poke| lines
2. Check they come before |start|
3. Trace parsePokemonFromTeamPreview()
4. Confirm team is assigned to correct player

### "Position scores are wrong"
1. Check active Pokemon is set correctly
2. Verify HP values are parsed correctly
3. Trace CalculatePositionScore() calculation
4. Print intermediate values

### "Turning points not detected"
1. Verify score delta calculation
2. Check threshold (15 points)
3. Print before/after scores
4. Confirm significance calculation

## Key Takeaways

1. **Two-Pass Parsing**: First extract teams, then process events
2. **State Mutation**: StateTracker maintains mutable state safely
3. **Relative Scoring**: Scores are relative, not absolute values
4. **Thresholded Detection**: Turning points need minimum momentum shift
5. **Component Composition**: Each component has single responsibility
6. **Type Safety**: TypeScript ensures frontend matches backend data

This architecture provides:
- ✅ Correctness through careful state management
- ✅ Clarity through two-phase processing
- ✅ Performance through single O(n) parse pass
- ✅ Extensibility through modular components
- ✅ Maintainability through clear separation of concerns
