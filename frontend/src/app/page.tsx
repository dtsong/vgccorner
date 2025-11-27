'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { analyzeReplayUrl, listReplays } from '@/lib/api/client';
import { ReplayListItem } from '@/lib/types/api';

export default function Home() {
  const router = useRouter();
  const [replayUrl, setReplayUrl] = useState('');
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [recentReplays, setRecentReplays] = useState<ReplayListItem[]>([]);
  const [isLoadingReplays, setIsLoadingReplays] = useState(true);

  // Load recent replays on mount
  useEffect(() => {
    loadRecentReplays();
  }, []);

  const loadRecentReplays = async () => {
    try {
      setIsLoadingReplays(true);
      const response = await listReplays({ limit: 10 });
      setRecentReplays(response.data);
    } catch (err) {
      console.error('Failed to load recent replays:', err);
      // Don't show error to user for background loading
    } finally {
      setIsLoadingReplays(false);
    }
  };

  const handleAnalyze = async () => {
    if (!replayUrl.trim()) {
      setError('Please enter a replay URL');
      return;
    }

    setError(null);
    setIsAnalyzing(true);

    try {
      const response = await analyzeReplayUrl(replayUrl, false);
      if (response.battleId) {
        // Navigate to the replay detail page
        router.push(`/replay/${response.battleId}`);
      }
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to analyze replay');
      setIsAnalyzing(false);
    }
  };

  const handleReplayClick = (replayId: string) => {
    router.push(`/replay/${replayId}`);
  };

  const formatTimestamp = (timestamp: string): string => {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

    if (diffDays === 0) {
      return `Today, ${date.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true })}`;
    } else if (diffDays === 1) {
      return `Yesterday, ${date.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true })}`;
    } else {
      return date.toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        hour: 'numeric',
        minute: '2-digit',
        hour12: true,
      });
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="px-8 py-6">
        <h1 className="text-xl font-semibold text-gray-900">VGC Corner Refinery</h1>
      </header>

      {/* Main Content */}
      <main className="flex flex-col items-center justify-start pt-20 px-4">
        {/* Search Input */}
        <div className="w-full max-w-3xl">
          <div className="relative flex items-center bg-white rounded-full shadow-lg overflow-hidden">
            <input
              type="text"
              placeholder="Paste Pokémon Showdown Replay URL..."
              value={replayUrl}
              onChange={(e) => setReplayUrl(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  handleAnalyze();
                }
              }}
              disabled={isAnalyzing}
              className="flex-1 px-8 py-5 text-lg text-gray-700 placeholder-gray-400 bg-transparent border-none outline-none disabled:opacity-50"
            />
            <button
              onClick={handleAnalyze}
              disabled={isAnalyzing}
              className="px-8 py-5 bg-blue-600 text-white font-medium rounded-full hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed m-2"
            >
              {isAnalyzing ? 'Analyzing...' : 'Analyze'}
            </button>
          </div>

          {/* Error Message */}
          {error && (
            <div className="mt-4 px-6 py-3 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-sm text-red-600">{error}</p>
            </div>
          )}
        </div>

        {/* Recent Replays */}
        <div className="w-full max-w-3xl mt-16">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            Recent Replays
          </h2>

          {isLoadingReplays ? (
            <div className="space-y-3">
              {[...Array(5)].map((_, i) => (
                <div
                  key={i}
                  className="flex items-center gap-3 px-4 py-3 bg-white rounded-lg animate-pulse"
                >
                  <div className="w-5 h-5 bg-gray-200 rounded-full"></div>
                  <div className="h-4 bg-gray-200 rounded w-48"></div>
                </div>
              ))}
            </div>
          ) : recentReplays.length === 0 ? (
            <div className="px-4 py-8 text-center text-gray-500">
              <p>No replays analyzed yet. Start by pasting a replay URL above!</p>
            </div>
          ) : (
            <div className="space-y-2">
              {recentReplays.map((replay) => (
                <button
                  key={replay.id}
                  onClick={() => handleReplayClick(replay.id)}
                  className="w-full flex items-center gap-3 px-4 py-3 bg-white rounded-lg hover:bg-gray-50 transition-colors text-left group"
                >
                  {/* Clock Icon */}
                  <svg
                    className="w-5 h-5 text-gray-400 flex-shrink-0"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <circle cx="12" cy="12" r="10" strokeWidth="2" />
                    <path
                      strokeLinecap="round"
                      strokeWidth="2"
                      d="M12 6v6l4 2"
                    />
                  </svg>

                  {/* Replay Info */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="text-sm text-gray-600">
                        {formatTimestamp(replay.timestamp)}
                      </span>
                      {replay.isPrivate && (
                        <span className="px-2 py-0.5 text-xs bg-gray-100 text-gray-600 rounded">
                          Private
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-gray-500 truncate">
                      {replay.player1} vs {replay.player2}
                      {replay.format && (
                        <span className="ml-2 text-gray-400">
                          • {replay.format}
                        </span>
                      )}
                    </p>
                  </div>

                  {/* Arrow Icon */}
                  <svg
                    className="w-5 h-5 text-gray-400 opacity-0 group-hover:opacity-100 transition-opacity"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M9 5l7 7-7 7"
                    />
                  </svg>
                </button>
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
