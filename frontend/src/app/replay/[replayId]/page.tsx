'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { getReplay } from '@/lib/api/client';
import { BattleSummary } from '@/lib/types/api';

export default function ReplayDetailPage() {
  const params = useParams();
  const router = useRouter();
  const replayId = params.replayId as string;

  const [battleData, setBattleData] = useState<BattleSummary | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadReplay = async () => {
      try {
        setIsLoading(true);
        setError(null);
        const response = await getReplay(replayId);
        if (response.data) {
          setBattleData(response.data);
        } else {
          setError('Replay data not found');
        }
      } catch (err: unknown) {
        setError((err as Error).message || 'Failed to load replay');
      } finally {
        setIsLoading(false);
      }
    };
    loadReplay();
  }, [replayId]);


  const formatDuration = (seconds: number): string => {
    const minutes = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${minutes}:${secs.toString().padStart(2, '0')}`;
  };

  const formatTimestamp = (timestamp: string): string => {
    const date = new Date(timestamp);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit',
      hour12: true,
    });
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <header className="px-8 py-6 bg-white border-b border-gray-200">
          <button
            onClick={() => router.push('/')}
            className="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors"
          >
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M15 19l-7-7 7-7"
              />
            </svg>
            Back to VGC Corner Refinery
          </button>
        </header>
        <div className="flex items-center justify-center min-h-[60vh]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading replay...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error || !battleData) {
    return (
      <div className="min-h-screen bg-gray-50">
        <header className="px-8 py-6 bg-white border-b border-gray-200">
          <button
            onClick={() => router.push('/')}
            className="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors"
          >
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M15 19l-7-7 7-7"
              />
            </svg>
            Back to VGC Corner Refinery
          </button>
        </header>
        <div className="flex items-center justify-center min-h-[60vh]">
          <div className="text-center">
            <div className="mb-4 text-red-500">
              <svg
                className="w-16 h-16 mx-auto"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <h2 className="text-xl font-semibold text-gray-900 mb-2">
              Failed to Load Replay
            </h2>
            <p className="text-gray-600 mb-6">{error}</p>
            <button
              onClick={() => router.push('/')}
              className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Return Home
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="px-8 py-6 bg-white border-b border-gray-200">
        <button
          onClick={() => router.push('/')}
          className="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors"
        >
          <svg
            className="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M15 19l-7-7 7-7"
            />
          </svg>
          Back to Battle Forge
        </button>
      </header>

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-8 py-8">
        {/* Battle Header */}
        <div className="bg-white rounded-xl shadow-sm p-8 mb-6">
          <div className="flex items-center justify-between mb-4">
            <h1 className="text-3xl font-bold text-gray-900">Battle Analysis</h1>
            <div className="flex items-center gap-3">
              <button
                onClick={() => router.push(`/replay/${replayId}/analysis`)}
                className="px-4 py-2 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2"
              >
                <svg
                  className="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                  />
                </svg>
                Turn-by-Turn Analysis
              </button>
              <span className="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm font-medium">
                {battleData.format}
              </span>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm text-gray-600">
            <div>
              <span className="font-medium">Battle ID:</span> {battleData.id}
            </div>
            <div>
              <span className="font-medium">Date:</span>{' '}
              {formatTimestamp(battleData.timestamp)}
            </div>
            <div>
              <span className="font-medium">Duration:</span>{' '}
              {formatDuration(battleData.duration)} ({battleData.turns} turns)
            </div>
          </div>
        </div>

        {/* Players */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          {/* Player 1 */}
          <div className="bg-white rounded-xl shadow-sm p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900">
                {battleData.player1.name}
              </h2>
              {battleData.winner === battleData.player1.name && (
                <span className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-medium">
                  Winner
                </span>
              )}
            </div>

            <div className="space-y-3">
              {battleData.player1.team.map((pokemon, idx) => (
                <div
                  key={idx}
                  className="p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <div className="flex items-center justify-between mb-2">
                    <h3 className="font-semibold text-gray-900">
                      {pokemon.name || pokemon.species}
                    </h3>
                    <span className="text-sm text-gray-500">
                      Lv. {pokemon.level}
                    </span>
                  </div>
                  <div className="text-sm text-gray-600 space-y-1">
                    <div>
                      <span className="font-medium">Ability:</span>{' '}
                      {pokemon.ability}
                    </div>
                    <div>
                      <span className="font-medium">Item:</span> {pokemon.item}
                    </div>
                    {pokemon.teraType && (
                      <div>
                        <span className="font-medium">Tera Type:</span>{' '}
                        {pokemon.teraType}
                      </div>
                    )}
                    <div>
                      <span className="font-medium">Moves:</span>{' '}
                      {pokemon.moves.join(', ')}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Player 2 */}
          <div className="bg-white rounded-xl shadow-sm p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900">
                {battleData.player2.name}
              </h2>
              {battleData.winner === battleData.player2.name && (
                <span className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-medium">
                  Winner
                </span>
              )}
            </div>

            <div className="space-y-3">
              {battleData.player2.team.map((pokemon, idx) => (
                <div
                  key={idx}
                  className="p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <div className="flex items-center justify-between mb-2">
                    <h3 className="font-semibold text-gray-900">
                      {pokemon.name || pokemon.species}
                    </h3>
                    <span className="text-sm text-gray-500">
                      Lv. {pokemon.level}
                    </span>
                  </div>
                  <div className="text-sm text-gray-600 space-y-1">
                    <div>
                      <span className="font-medium">Ability:</span>{' '}
                      {pokemon.ability}
                    </div>
                    <div>
                      <span className="font-medium">Item:</span> {pokemon.item}
                    </div>
                    {pokemon.teraType && (
                      <div>
                        <span className="font-medium">Tera Type:</span>{' '}
                        {pokemon.teraType}
                      </div>
                    )}
                    <div>
                      <span className="font-medium">Moves:</span>{' '}
                      {pokemon.moves.join(', ')}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Key Moments */}
        {battleData.keyMoments && battleData.keyMoments.length > 0 && (
          <div className="bg-white rounded-xl shadow-sm p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">
              Key Moments
            </h2>
            <div className="space-y-3">
              {battleData.keyMoments.map((moment, idx) => (
                <div
                  key={idx}
                  className="p-4 border-l-4 border-blue-500 bg-blue-50 rounded-r-lg"
                >
                  <div className="flex items-center gap-2 mb-1">
                    <span className="text-sm font-semibold text-blue-900">
                      Turn {moment.turn}
                    </span>
                    <span className="text-xs px-2 py-0.5 bg-blue-200 text-blue-800 rounded">
                      {moment.type}
                    </span>
                  </div>
                  <p className="text-sm text-gray-700">{moment.description}</p>
                </div>
              ))}
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
