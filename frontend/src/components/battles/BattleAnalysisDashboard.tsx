'use client';

import React, { useMemo } from 'react';
import { BattleSummary } from '@/lib/types/showdown';
import BattleHeader from './BattleHeader';
import TurnTimeline from './TurnTimeline';
import PositionScoreChart from './PositionScoreChart';
import KeyMomentsPanel from './KeyMomentsPanel';
import TeamComparison from './TeamComparison';
import BattleStatistics from './BattleStatistics';

interface BattleAnalysisDashboardProps {
  battle: BattleSummary;
}

/**
 * Main battle analysis dashboard that displays comprehensive battle insights
 * including position scores, turning points, team states, and key moments.
 */
export default function BattleAnalysisDashboard({ battle }: BattleAnalysisDashboardProps) {
  // Calculate battle insights
  const battleInsights = useMemo(() => {
    const totalTurns = battle.turns.length;
    const player1Wins = battle.winner === 'player1';
    const player2Wins = battle.winner === 'player2';

    // Calculate average position scores
    let totalP1Score = 0;
    let totalP2Score = 0;
    let scoreCount = 0;

    battle.turns.forEach((turn) => {
      if (turn.positionScore) {
        totalP1Score += turn.positionScore.player1Score;
        totalP2Score += turn.positionScore.player2Score;
        scoreCount++;
      }
    });

    const avgP1Score = scoreCount > 0 ? totalP1Score / scoreCount : 50;
    const avgP2Score = scoreCount > 0 ? totalP2Score / scoreCount : 50;

    return {
      totalTurns,
      player1Wins,
      player2Wins,
      avgP1Score,
      avgP2Score,
      turningPoints: battle.stats.turningPoints || [],
      dominantMoments: battle.keyMoments.filter((m) => m.significance >= 7),
    };
  }, [battle]);

  return (
    <div className="w-full bg-gradient-to-br from-slate-50 to-slate-100 min-h-screen p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        {/* Header Section */}
        <BattleHeader battle={battle} />

        {/* Position Score Timeline */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-bold text-slate-800 mb-4">Position Evolution</h2>
          <PositionScoreChart battle={battle} />
        </div>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Column: Turn Timeline & Key Moments */}
          <div className="lg:col-span-2 space-y-6">
            {/* Turn-by-Turn Timeline */}
            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-2xl font-bold text-slate-800 mb-4">Battle Progression</h2>
              <TurnTimeline battle={battle} />
            </div>

            {/* Key Moments */}
            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-2xl font-bold text-slate-800 mb-4">Critical Moments</h2>
              <KeyMomentsPanel keyMoments={battle.keyMoments} turningPoints={battleInsights.turningPoints} />
            </div>
          </div>

          {/* Right Column: Team & Stats */}
          <div className="space-y-6">
            {/* Team Comparison */}
            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-2xl font-bold text-slate-800 mb-4">Team Status</h2>
              <TeamComparison player1={battle.player1} player2={battle.player2} />
            </div>

            {/* Battle Statistics */}
            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-2xl font-bold text-slate-800 mb-4">Statistics</h2>
              <BattleStatistics stats={battle.stats} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
