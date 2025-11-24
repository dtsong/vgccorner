'use client';

import React from 'react';
import { BattleStats } from '@/lib/types/showdown';

interface BattleStatisticsProps {
  stats: BattleStats;
}

/**
 * Displays aggregated battle statistics comparing both players.
 */
export default function BattleStatistics({ stats }: BattleStatisticsProps) {
  return (
    <div className="space-y-4">
      {/* Key Stats */}
      <div className="space-y-2">
        <div className="p-3 bg-slate-50 rounded-lg">
          <div className="text-xs text-slate-600 font-semibold mb-1">TOTAL TURNS</div>
          <div className="text-2xl font-bold text-slate-800">{stats.totalTurns}</div>
        </div>

        <div className="p-3 bg-slate-50 rounded-lg">
          <div className="text-xs text-slate-600 font-semibold mb-1">TOTAL SWITCHES</div>
          <div className="text-2xl font-bold text-slate-800">{stats.switches}</div>
        </div>

        <div className="p-3 bg-slate-50 rounded-lg">
          <div className="text-xs text-slate-600 font-semibold mb-1">CRITICAL HITS</div>
          <div className="text-2xl font-bold text-slate-800">{stats.criticalHits}</div>
        </div>
      </div>

      {/* Type Effectiveness */}
      <div className="border-t border-slate-200 pt-4">
        <h4 className="font-semibold text-slate-800 mb-2">Type Effectiveness</h4>
        <div className="space-y-2">
          <div className="flex justify-between items-center text-sm">
            <span className="text-slate-700">Super Effective:</span>
            <span className="font-bold text-green-600">{stats.superEffective}</span>
          </div>
          <div className="flex justify-between items-center text-sm">
            <span className="text-slate-700">Not Very Effective:</span>
            <span className="font-bold text-orange-600">{stats.notVeryEffective}</span>
          </div>
        </div>
      </div>

      {/* Damage Stats */}
      <div className="border-t border-slate-200 pt-4">
        <h4 className="font-semibold text-slate-800 mb-2">Damage Per Turn</h4>
        <div className="text-sm">
          <div className="text-slate-700">
            Average: <span className="font-bold text-slate-800">{stats.avgDamagePerTurn.toFixed(1)}</span>
          </div>
        </div>
      </div>

      {/* Most Used Moves */}
      {Object.keys(stats.moveFrequency).length > 0 && (
        <div className="border-t border-slate-200 pt-4">
          <h4 className="font-semibold text-slate-800 mb-2">Most Used Moves</h4>
          <div className="space-y-1">
            {Object.entries(stats.moveFrequency)
              .sort((a, b) => b[1] - a[1])
              .slice(0, 5)
              .map(([moveId, count]) => (
                <div key={moveId} className="flex justify-between items-center text-sm">
                  <span className="text-slate-700 capitalize">{moveId}</span>
                  <span className="bg-slate-200 text-slate-800 px-2 py-1 rounded text-xs font-semibold">
                    {count}x
                  </span>
                </div>
              ))}
          </div>
        </div>
      )}

      {/* Player Comparison */}
      <div className="border-t border-slate-200 pt-4">
        <h4 className="font-semibold text-slate-800 mb-2">Player Comparison</h4>
        <div className="space-y-2 text-sm">
          <div className="flex justify-between">
            <span className="text-slate-700">Moves Used:</span>
            <span className="font-semibold">
              {stats.player1Stats.moveCount} vs {stats.player2Stats.moveCount}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-700">Switches:</span>
            <span className="font-semibold">
              {stats.player1Stats.switchCount} vs {stats.player2Stats.switchCount}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-700">Total Damage:</span>
            <span className="font-semibold">
              {stats.player1Stats.damageDealt} vs {stats.player2Stats.damageDealt}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
