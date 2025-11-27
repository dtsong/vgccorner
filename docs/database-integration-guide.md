# Database Integration Guide

## Overview

This guide explains how the database integration works for turn-by-turn analysis in VGC Corner.

## Database Setup

### 1. Run Migrations

Apply the database migrations to create the necessary schema:

```bash
# From the backend directory
cd backend

# Run initial schema migration
psql -U vgccorner -d vgccorner -f migrations/001_initial_schema.sql

# Run turn analysis enhancements migration
psql -U vgccorner -d vgccorner -f migrations/002_turn_analysis_enhancements.sql
```

### 2. Configure Database Connection

Set environment variables for database connection:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=vgccorner
export DB_PASSWORD=your_password
export DB_NAME=vgccorner
export DB_SSL_MODE=disable  # Use 'require' in production
```

Or use docker-compose:

```bash
docker-compose up -d postgres
```

## Database Schema

### Core Tables

#### battles
Stores basic battle information and team archetypes.

- `id` (UUID) - Primary key
- `format` (VARCHAR) - Battle format (e.g., gen9vgc2024regh)
- `timestamp` (TIMESTAMP) - When the battle occurred
- `duration_sec` (INT) - Battle duration in seconds
- `winner` (VARCHAR) - Winner identifier
- `player1_id`, `player2_id` (VARCHAR) - Player usernames
- `player1_archetype`, `player2_archetype` (VARCHAR) - Team archetypes
- `player1_archetype_data`, `player2_archetype_data` (JSONB) - Detailed classification data
- `battle_log` (TEXT) - Raw battle log
- `is_private` (BOOLEAN) - Privacy flag

#### battle_turns
Stores turn information.

- `id` (UUID) - Primary key
- `battle_id` (UUID) - Foreign key to battles
- `turn_number` (INT) - Turn number

#### battle_actions
Stores actions taken during turns.

- `id` (UUID) - Primary key
- `battle_turn_id` (UUID) - Foreign key to battle_turns
- `player_number` (INT) - 1 or 2
- `action_type` (VARCHAR) - "move", "switch", "item"
- `pokemon_name` (VARCHAR) - Pokemon performing the action
- `target_pokemon` (VARCHAR) - Target of the action
- `result` (VARCHAR) - Action result
- `details` (TEXT) - Additional details
- `order_in_turn` (INT) - Order within the turn

#### move_impacts
Stores detailed impact information for moves.

- `id` (UUID) - Primary key
- `action_id` (UUID) - Foreign key to battle_actions
- `damage_dealt` (INT) - Damage dealt
- `healing_done` (INT) - Healing done
- `status_inflicted` (VARCHAR) - Status condition inflicted
- `speed_control` (VARCHAR) - Speed control type
- `weather_set`, `terrain_set` (VARCHAR) - Field effects
- `fake_out`, `protect_used`, `critical`, `missed` (BOOLEAN) - Move flags
- `effectiveness` (VARCHAR) - Type effectiveness

#### stat_changes
Tracks stat modifications.

- `id` (UUID) - Primary key
- `move_impact_id` (UUID) - Foreign key to move_impacts
- `pokemon_name` (VARCHAR) - Pokemon affected
- `stat` (VARCHAR) - Stat modified
- `stages` (INT) - Number of stages (positive or negative)

#### fainted_pokemon
Tracks which Pokemon fainted.

- `id` (UUID) - Primary key
- `move_impact_id` (UUID) - Foreign key to move_impacts
- `pokemon_name` (VARCHAR) - Pokemon that fainted
- `turn_number` (INT) - Turn it fainted on

#### turn_board_states
Stores the board state at each turn.

- `id` (UUID) - Primary key
- `battle_turn_id` (UUID) - Foreign key to battle_turns
- `player_number` (INT) - 1 or 2
- `pokemon_name` (VARCHAR) - Pokemon nickname
- `pokemon_species` (VARCHAR) - Pokemon species
- `position` (INT) - Position on field
- `hp`, `max_hp` (INT) - HP values
- `status` (VARCHAR) - Status condition
- `is_lead` (BOOLEAN) - Whether this is a lead Pokemon

## Code Architecture

### Database Layer (`backend/internal/db/`)

#### `db.go`
Core database connection and basic operations.

- `NewDatabase(connString)` - Initialize database connection
- `StoreBattle()` - Store basic battle data
- `GetBattle()` - Retrieve battle data
- `ListBattles()` - List battles with filtering

#### `turn_storage.go`
Turn-by-turn data storage and retrieval.

- `StoreTurnData()` - Store detailed turn analysis
- `GetTurnData()` - Retrieve turn-by-turn analysis
- Helper functions for storing/retrieving:
  - Team archetypes
  - Turn data
  - Board states
  - Actions
  - Move impacts

### Analysis Layer (`backend/internal/analysis/`)

#### `parser.go`
Basic battle log parsing.

- `ParseShowdownLog()` - Parse battle logs (basic)
- `ParseEnhancedShowdownLog()` - Parse with enhanced turn tracking

#### `turn_parser.go`
Enhanced turn-by-turn parsing.

- `TurnParser` - Stateful parser for detailed turn tracking
- `ProcessTurnEvent()` - Process individual battle events
- `StartNewTurn()` - Begin tracking a new turn
- `FinalizeTurn()` - Complete turn and calculate state

#### `team_classifier.go`
Team archetype classification.

- `ClassifyTeam()` - Analyze team composition and determine archetype
- `GetArchetypeDescription()` - Get human-readable description

#### `move_impact.go`
Move impact tracking.

- `EnhanceActionWithImpact()` - Add detailed impact data to actions
- Tracks: damage, healing, status, speed control, weather, terrain, stat changes, faints

### HTTP API Layer (`backend/internal/httpapi/`)

#### `showdown_handlers.go`
Handlers for battle analysis.

- `handleAnalyzeShowdown()` - Analyze new battles
  - Parses battle log
  - Stores in database
  - Returns analysis
- `handleGetShowdownReplay()` - Get specific replay
- `handleListShowdownReplays()` - List replays

#### `turn_handlers.go`
Handlers for turn-by-turn analysis.

- `handleGetTurnAnalysis()` - Get detailed turn-by-turn data
  - Retrieves from database
  - Converts to API format
  - Returns with team archetypes

## Data Flow

### Analyzing a New Battle

1. **Client sends replay URL** → `POST /api/showdown/analyze`
2. **Handler parses battle log** → `ParseEnhancedShowdownLog()`
3. **Team classification** → `ClassifyTeam()` for both players
4. **Store basic battle data** → `StoreBattle()`
5. **Store turn-by-turn data** → `StoreTurnData()`
6. **Return analysis** → Client receives battle summary

### Retrieving Turn Analysis

1. **Client requests turn data** → `GET /api/showdown/replays/{id}/turns`
2. **Retrieve from database** → `GetTurnData()`
3. **Convert to API format** → `convertTurnDataToResponse()`
4. **Return to client** → Client displays turn-by-turn

## Testing Database Integration

### 1. Test Connections

```bash
# Test database connection
psql -U vgccorner -d vgccorner -c "SELECT 1;"
```

### 2. Run Backend Tests

```bash
cd backend
go test ./internal/db/...
go test ./internal/analysis/...
```

### 3. Test API Endpoints

```bash
# Health check
curl http://localhost:8080/healthz

# Analyze a battle (you'll need a real battle log)
curl -X POST http://localhost:8080/api/showdown/analyze \
  -H "Content-Type: application/json" \
  -d '{"analysisType": "rawLog", "rawLog": "...", "isPrivate": false}'

# List replays
curl http://localhost:8080/api/showdown/replays

# Get turn analysis
curl http://localhost:8080/api/showdown/replays/{battle-id}/turns
```

## Performance Considerations

### Indexes

The schema includes indexes on:
- `battles(format)` - Fast filtering by format
- `battles(timestamp)` - Fast chronological queries
- `battles(player1_id, player2_id)` - Fast player lookups
- `battle_turns(battle_id)` - Fast turn retrieval
- `battle_actions(battle_turn_id)` - Fast action lookups
- `move_impacts(action_id)` - Fast impact lookups

### Optimization Tips

1. **Use transactions** - All data for a battle is stored in one transaction
2. **Batch inserts** - Turns and actions are inserted efficiently
3. **JSONB for archetypes** - Flexible storage for classification data
4. **Connection pooling** - Reuse database connections

## Troubleshooting

### Database Connection Failed

Check environment variables:
```bash
echo $DB_HOST
echo $DB_USER
# etc.
```

### Migrations Failed

Check PostgreSQL logs:
```bash
docker-compose logs postgres
```

### Data Not Storing

Check logs for errors:
```bash
# Backend logs will show database errors
tail -f backend/logs/app.log
```

### Turn Data Not Found

Ensure migrations ran successfully:
```sql
SELECT COUNT(*) FROM battle_turns;
SELECT COUNT(*) FROM battle_actions;
```

## Next Steps

1. **Add caching** - Cache frequently accessed battles
2. **Add pagination** - Implement cursor-based pagination for large turn lists
3. **Add search** - Full-text search on player names and battle IDs
4. **Add analytics** - Aggregate stats across battles
5. **Add archiving** - Archive old battles to separate storage
