// API Types for VGC Corner Backend

export interface AnalyzeShowdownRequest {
  analysisType: 'replayId' | 'username' | 'rawLog';
  replayId?: string;
  username?: string;
  format?: string;
  limit?: number;
  rawLog?: string;
  isPrivate: boolean;
}

export interface BattleSummary {
  id: string;
  format: string;
  timestamp: string;
  duration: number;
  winner: string;
  player1: PlayerInfo;
  player2: PlayerInfo;
  turns: number;
  keyMoments?: KeyMoment[];
}

export interface PlayerInfo {
  name: string;
  team: Pokemon[];
  rating?: number;
}

export interface Pokemon {
  name: string;
  species: string;
  ability: string;
  item: string;
  moves: string[];
  teraType?: string;
  level: number;
}

export interface KeyMoment {
  turn: number;
  type: string;
  description: string;
  importance: number;
}

export interface ResponseMetadata {
  parseTimeMs: number;
  analysisTimeMs: number;
  cached: boolean;
}

export interface AnalyzeResponse {
  status: string;
  battleId?: string;
  data?: BattleSummary;
  metadata?: ResponseMetadata;
}

export interface ErrorResponse {
  error: string;
  code: string;
  details?: unknown;
}

export interface ReplayListItem {
  id: string;
  format: string;
  timestamp: string;
  player1: string;
  player2: string;
  winner?: string;
  isPrivate: boolean;
}

export interface ListReplaysResponse {
  status: string;
  data: ReplayListItem[];
  pagination: {
    limit: number;
    offset: number;
    total: number;
  };
}
