-- Migration: Add turn-by-turn analysis enhancements
-- Version: 002_turn_analysis_enhancements.sql

-- Add team archetype fields to battles table
ALTER TABLE battles
ADD COLUMN IF NOT EXISTS player1_archetype VARCHAR(100),
ADD COLUMN IF NOT EXISTS player1_archetype_data JSONB,
ADD COLUMN IF NOT EXISTS player2_archetype VARCHAR(100),
ADD COLUMN IF NOT EXISTS player2_archetype_data JSONB;

-- Enhance battle_actions table with detailed impact tracking
ALTER TABLE battle_actions
ADD COLUMN IF NOT EXISTS pokemon_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS target_pokemon VARCHAR(100),
ADD COLUMN IF NOT EXISTS result VARCHAR(50),
ADD COLUMN IF NOT EXISTS details TEXT,
ADD COLUMN IF NOT EXISTS order_in_turn INT DEFAULT 0;

-- Create move_impacts table for detailed move impact tracking
CREATE TABLE IF NOT EXISTS move_impacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action_id UUID NOT NULL REFERENCES battle_actions(id) ON DELETE CASCADE,
    damage_dealt INT DEFAULT 0,
    healing_done INT DEFAULT 0,
    status_inflicted VARCHAR(50),
    speed_control VARCHAR(50),
    weather_set VARCHAR(50),
    terrain_set VARCHAR(50),
    fake_out BOOLEAN DEFAULT FALSE,
    protect_used BOOLEAN DEFAULT FALSE,
    critical BOOLEAN DEFAULT FALSE,
    effectiveness VARCHAR(50),
    missed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create stat_changes table for tracking stat modifications
CREATE TABLE IF NOT EXISTS stat_changes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    move_impact_id UUID NOT NULL REFERENCES move_impacts(id) ON DELETE CASCADE,
    pokemon_name VARCHAR(100) NOT NULL,
    stat VARCHAR(20) NOT NULL,
    stages INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create fainted_pokemon table for tracking KOs
CREATE TABLE IF NOT EXISTS fainted_pokemon (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    move_impact_id UUID NOT NULL REFERENCES move_impacts(id) ON DELETE CASCADE,
    pokemon_name VARCHAR(100) NOT NULL,
    turn_number INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create turn_board_states table for storing board state at each turn
CREATE TABLE IF NOT EXISTS turn_board_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    battle_turn_id UUID NOT NULL REFERENCES battle_turns(id) ON DELETE CASCADE,
    player_number INT NOT NULL,
    pokemon_name VARCHAR(100) NOT NULL,
    pokemon_species VARCHAR(100) NOT NULL,
    position INT NOT NULL,
    hp INT NOT NULL,
    max_hp INT NOT NULL,
    status VARCHAR(50),
    is_lead BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_move_impacts_action ON move_impacts(action_id);
CREATE INDEX IF NOT EXISTS idx_stat_changes_impact ON stat_changes(move_impact_id);
CREATE INDEX IF NOT EXISTS idx_fainted_pokemon_impact ON fainted_pokemon(move_impact_id);
CREATE INDEX IF NOT EXISTS idx_turn_board_states_turn ON turn_board_states(battle_turn_id);
CREATE INDEX IF NOT EXISTS idx_turn_board_states_player ON turn_board_states(player_number);
CREATE INDEX IF NOT EXISTS idx_battles_player1_archetype ON battles(player1_archetype);
CREATE INDEX IF NOT EXISTS idx_battles_player2_archetype ON battles(player2_archetype);

-- Add comment documenting the schema
COMMENT ON TABLE move_impacts IS 'Detailed impact information for each move action in a battle';
COMMENT ON TABLE stat_changes IS 'Stat modifications (boosts/drops) caused by moves';
COMMENT ON TABLE fainted_pokemon IS 'Track which Pokemon fainted as a result of moves';
COMMENT ON TABLE turn_board_states IS 'Board state (active Pokemon) at each turn';
COMMENT ON COLUMN battles.player1_archetype IS 'Team archetype classification (e.g., Hard Trick Room, Psy-Spam)';
COMMENT ON COLUMN battles.player1_archetype_data IS 'Detailed JSON data about team classification';
