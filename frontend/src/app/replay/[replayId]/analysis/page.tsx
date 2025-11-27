'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import MiniBoardstate from '@/components/battles/MiniBoardstate';
import TurnEvents from '@/components/battles/TurnEvents';
import { TurnData, TurnAnalysis, EventResult } from '@/lib/types/battle';
import { getTurnAnalysis } from '@/lib/api/client';

export default function TurnAnalysisPage() {
  const params = useParams();
  const router = useRouter();
  const replayId = params.replayId as string;

  const [turnAnalysis, setTurnAnalysis] = useState<TurnAnalysis | null>(null);
  const [currentTurn, setCurrentTurn] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadTurnAnalysis();
  }, [replayId]);

  const loadTurnAnalysis = async () => {
    try {
      setIsLoading(true);
      setError(null);

      const response = await getTurnAnalysis(replayId);

      // Convert API response to TurnAnalysis format
      const analysis: TurnAnalysis = {
        battleId: response.battleId,
        format: response.format,
        player1: response.player1,
        player2: response.player2,
        winner: response.winner,
        turns: response.turns.map((turn) => ({
          turnNumber: turn.turnNumber,
          events: turn.events.map((event) => ({
            type: (event.type === 'damage' || event.type === 'heal' ? 'other' : event.type) as 'move' | 'switch' | 'faint' | 'status' | 'weather' | 'terrain' | 'other',
            pokemon: event.pokemon,
            action: event.action,
            target: event.target,
            result: event.result as EventResult | undefined,
            details: event.details,
            playerSide: event.playerSide as 'player1' | 'player2',
          })),
          boardState: {
            player1Active: turn.boardState.player1Active.map((p) => ({
              species: p.species,
              nickname: p.nickname,
              position: p.position,
              hp: p.hp,
              maxHp: p.maxHp,
              status: p.status,
              isLead: p.isLead,
            })),
            player2Active: turn.boardState.player2Active.map((p) => ({
              species: p.species,
              nickname: p.nickname,
              position: p.position,
              hp: p.hp,
              maxHp: p.maxHp,
              status: p.status,
              isLead: p.isLead,
            })),
          },
        })),
      };

      setTurnAnalysis(analysis);
    } catch (err: unknown) {
      setError((err as Error).message || 'Failed to load turn analysis');
    } finally {
      setIsLoading(false);
    }
  };

  const handlePreviousTurn = () => {
    if (currentTurn > 0) {
      setCurrentTurn(currentTurn - 1);
    }
  };

  const handleNextTurn = () => {
    if (turnAnalysis && currentTurn < turnAnalysis.turns.length - 1) {
      setCurrentTurn(currentTurn + 1);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <header className="px-8 py-6 bg-white border-b border-gray-200">
          <button
            onClick={() => router.push(`/replay/${replayId}`)}
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
            Back to Replay Details
          </button>
        </header>
        <div className="flex items-center justify-center min-h-[60vh]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading turn analysis...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error || !turnAnalysis) {
    return (
      <div className="min-h-screen bg-gray-50">
        <header className="px-8 py-6 bg-white border-b border-gray-200">
          <button
            onClick={() => router.push(`/replay/${replayId}`)}
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
            Back to Replay Details
          </button>
        </header>
        <div className="flex items-center justify-center min-h-[60vh]">
          <div className="text-center">
            <p className="text-red-600 mb-4">{error}</p>
            <button
              onClick={() => router.push(`/replay/${replayId}`)}
              className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >
              Return to Replay
            </button>
          </div>
        </div>
      </div>
    );
  }

  const currentTurnData = turnAnalysis.turns[currentTurn];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="px-8 py-6 bg-white border-b border-gray-200">
        <div className="flex items-center justify-between">
          <button
            onClick={() => router.push(`/replay/${replayId}`)}
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
            Back to Replay Details
          </button>

          {turnAnalysis.winner && (
            <span className="px-4 py-2 bg-green-100 text-green-800 font-medium rounded-lg">
              Winner: {turnAnalysis.winner}
            </span>
          )}
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-8 py-8">
        {/* Title */}
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-900">
            Matchup Details - Turn {currentTurnData.turnNumber} Analysis
          </h1>
          <p className="text-gray-600 mt-1">
            {turnAnalysis.player1} vs {turnAnalysis.player2} • {turnAnalysis.format}
          </p>
        </div>

        {/* Mini Boardstate */}
        <MiniBoardstate
          boardState={currentTurnData.boardState}
          player1Name={turnAnalysis.player1}
          player2Name={turnAnalysis.player2}
        />

        {/* Turn Events */}
        <div className="bg-white rounded-xl p-6 mb-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            Turn {currentTurnData.turnNumber} - In Depth
          </h2>
          <TurnEvents events={currentTurnData.events} />
        </div>

        {/* Turn Navigation */}
        <div className="flex items-center justify-between bg-white rounded-xl p-6">
          <button
            onClick={handlePreviousTurn}
            disabled={currentTurn === 0}
            className="flex items-center gap-2 px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
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
            Previous Turn
          </button>

          <div className="text-center">
            <p className="text-sm text-gray-600">
              Turn {currentTurn + 1} of {turnAnalysis.turns.length}
            </p>
          </div>

          <button
            onClick={handleNextTurn}
            disabled={currentTurn === turnAnalysis.turns.length - 1}
            className="flex items-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next Turn
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
                d="M9 5l7 7-7 7"
              />
            </svg>
          </button>
        </div>
      </main>
    </div>
  );
}
