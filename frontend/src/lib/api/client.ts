// API Client for VGC Corner Backend

import {
  AnalyzeShowdownRequest,
  AnalyzeResponse,
  ListReplaysResponse,
  ErrorResponse,
} from '@/lib/types/api';

// Use environment variable for API URL, fallback to localhost
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ApiError extends Error {
  constructor(
    message: string,
    public code: string,
    public details?: unknown,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const errorData: ErrorResponse = await response.json().catch(() => ({
      error: 'Unknown error occurred',
      code: 'UNKNOWN_ERROR',
    }));
    throw new ApiError(errorData.error, errorData.code, errorData.details);
  }
  return response.json();
}

/**
 * Analyze a Pokémon Showdown replay
 */
export async function analyzeShowdown(
  request: AnalyzeShowdownRequest,
): Promise<AnalyzeResponse> {
  const response = await fetch(`${API_BASE_URL}/api/showdown/analyze`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });
  return handleResponse<AnalyzeResponse>(response);
}

/**
 * Analyze a replay by its Showdown URL
 */
export async function analyzeReplayUrl(
  url: string,
  isPrivate: boolean = false,
): Promise<AnalyzeResponse> {
  // Extract replay ID from URL
  // Example: https://replay.pokemonshowdown.com/gen9vgc2024regh-2024567890
  const replayId = extractReplayId(url);
  if (!replayId) {
    throw new ApiError('Invalid replay URL', 'INVALID_URL');
  }

  return analyzeShowdown({
    analysisType: 'replayId',
    replayId,
    isPrivate,
  });
}

/**
 * List recent replays
 */
export async function listReplays(params?: {
  username?: string;
  format?: string;
  isPrivate?: boolean;
  limit?: number;
  offset?: number;
}): Promise<ListReplaysResponse> {
  const queryParams = new URLSearchParams();
  if (params?.username) queryParams.append('username', params.username);
  if (params?.format) queryParams.append('format', params.format);
  if (params?.isPrivate !== undefined)
    queryParams.append('isPrivate', String(params.isPrivate));
  if (params?.limit) queryParams.append('limit', String(params.limit));
  if (params?.offset) queryParams.append('offset', String(params.offset));

  const response = await fetch(
    `${API_BASE_URL}/api/showdown/replays?${queryParams}`,
    {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    },
  );
  return handleResponse<ListReplaysResponse>(response);
}

/**
 * Get a specific replay by ID
 */
export async function getReplay(replayId: string): Promise<AnalyzeResponse> {
  const response = await fetch(
    `${API_BASE_URL}/api/showdown/replays/${replayId}`,
    {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    },
  );
  return handleResponse<AnalyzeResponse>(response);
}

/**
 * Extract replay ID from various Pokémon Showdown URL formats
 */
function extractReplayId(url: string): string | null {
  // Handle different URL formats:
  // https://replay.pokemonshowdown.com/gen9vgc2024regh-2024567890
  // replay.pokemonshowdown.com/gen9vgc2024regh-2024567890
  // gen9vgc2024regh-2024567890

  // Remove protocol if present
  const withoutProtocol = url.replace(/^https?:\/\//, '');

  // Remove domain if present
  const withoutDomain = withoutProtocol.replace(
    /^replay\.pokemonshowdown\.com\//,
    '',
  );

  // The remaining should be the replay ID
  const replayId = withoutDomain.trim();

  // Basic validation: should contain at least one hyphen and alphanumeric chars
  if (replayId && /^[a-z0-9]+-[a-z0-9]+$/i.test(replayId)) {
    return replayId;
  }

  return null;
}

export { ApiError };
