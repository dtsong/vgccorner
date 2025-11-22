/**
 * Type definitions for BattleSummary
 * Aligned with backend Go structs in internal/analysis/types.go
 */

export interface BattleSummary {
  // Metadata about the battle
  id: string;
  format: string; // e.g., "Regulation H"
  timestamp: string; // ISO 8601 date string
  duration: number; // in seconds

  // Player information
  player1: Player;
  player2: Player;
  winner: "player1" | "player2" | "draw";

  // Battle progression
  turns: Turn[];

  // Overall statistics
  stats: BattleStats;

  // Key moments and highlights
  keyMoments: KeyMoment[];
}

export interface Player {
  name: string;
  team: Pokémon[];
  active: Pokémon | null; // Currently active Pokémon
  losses: number; // Number of fainted Pokémon
  totalLeft: number; // Total Pokémon still in battle
}

export interface Pokémon {
  id: string; // e.g., "pikachu"
  name: string;
  level: number;
  gender: "M" | "F" | ""; // Empty string if unknown
  ability: string;
  item: string;
  stats: Stats;
  moves: Move[];
  happiness: number; // 0-255
  shiny: boolean;
}

export interface Move {
  id: string; // e.g., "thunderbolt"
  name: string;
  type: string; // e.g., "Electric"
  power: number; // 0 if N/A
  accuracy: number; // 0-100, 0 if N/A
  pp: number; // Power Points
}

export interface Stats {
  hp: number;
  attack: number;
  defense: number;
  spAtk: number;
  spDef: number;
  speed: number;
}

export interface Turn {
  turnNumber: number;
  actions: Action[];
  stateAfter: BattleState;
  damageDealt: Record<string, number>; // Player name -> damage dealt
  healingDone: Record<string, number>; // Player name -> healing done
}

export interface Action {
  player: "player1" | "player2";
  actionType: "move" | "switch" | "item";
  move?: Move;
  switchTo?: string; // Pokémon name if switch
  item?: string; // Item used if item action
}

export interface BattleState {
  player1Active: Pokémon | null;
  player2Active: Pokémon | null;
  player1Team: string[]; // List of alive Pokémon names
  player2Team: string[];
}

export interface BattleStats {
  totalTurns: number;
  moveFrequency: Record<string, number>; // Move ID -> count
  typeCoverage: Record<string, number>; // Type -> count
  switches: number; // Total switches by both players
  criticalHits: number;
  superEffective: number;
  notVeryEffective: number;
  avgDamagePerTurn: number;
  avgHealPerTurn: number;
  player1Stats: PlayerStats;
  player2Stats: PlayerStats;
}

export interface PlayerStats {
  moveCount: number;
  switchCount: number;
  damageDealt: number;
  damageTaken: number;
  healingDone: number;
  healingReceived: number;
  movesByType: Record<string, number>; // Type -> count
  effectiveness: EffectivenessStats;
}

export interface EffectivenessStats {
  superEffective: number;
  notVeryEffective: number;
  neutral: number;
}

export interface KeyMoment {
  turnNumber: number;
  description: string; // e.g., "Player 2 switched to Charizard"
  type: "switch" | "ko" | "status" | "weather" | "other";
  significance: number; // 1-10 scale
}
