// Battle Analysis Types

export interface TurnData {
  turnNumber: number;
  events: BattleEvent[];
  boardState: BoardState;
}

export interface BoardState {
  player1Active: ActivePokemon[];
  player2Active: ActivePokemon[];
}

export interface ActivePokemon {
  species: string;
  nickname?: string;
  position: number;
  hp: number;
  maxHp: number;
  status?: string;
  isLead?: boolean;
}

export interface BattleEvent {
  type: 'move' | 'switch' | 'faint' | 'status' | 'weather' | 'terrain' | 'damage' | 'heal' | 'other';
  pokemon: string;
  action: string;
  target?: string;
  result?: EventResult;
  details?: string;
  playerSide: 'player1' | 'player2';
}

export type EventResult =
  | 'critical-hit'
  | 'super-effective'
  | 'not-very-effective'
  | 'miss'
  | 'faint'
  | 'immune'
  | 'success'
  | 'fail';

export interface TurnAnalysis {
  battleId: string;
  turns: TurnData[];
  format: string;
  player1: string;
  player2: string;
  winner?: string;
}
