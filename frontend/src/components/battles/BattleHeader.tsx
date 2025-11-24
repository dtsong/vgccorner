'use client';

import React from 'react';
import { BattleSummary } from '@/lib/types/showdown';

interface BattleHeaderProps {
  battle: BattleSummary;
}

/**
 * Displays the battle header with player names, winner, format, and metadata.
 */
export default function BattleHeader({ battle }: BattleHeaderProps) {
  const isPlayer1Winner = battle.winner === 'player1';

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <div className="mb-4">
        <div className="text-sm text-slate-500 font-semibold mb-2">
          {battle.format} • {formatDate(battle.timestamp)}
        </div>
        <h1 className="text-3xl font-bold text-slate-800 mb-2">
          {battle.player1.name} vs {battle.player2.name}
        </h1>
      </div>

      <div className="grid grid-cols-3 gap-4">
        {/* Player 1 */}
        <div className="flex flex-col items-start">
          <div className="text-sm text-slate-600 mb-1">Player 1</div>
          <div className="flex items-center gap-2">
            <div className={`text-lg font-bold ${isPlayer1Winner ? 'text-green-600' : 'text-slate-600'}`}>
              {battle.player1.name}
            </div>
            {isPlayer1Winner && (
              <span className="inline-block px-3 py-1 bg-green-100 text-green-800 rounded-full text-xs font-bold">
                WINNER
              </span>
            )}
          </div>
          <div className="text-sm text-slate-500 mt-1">
            {battle.player1.losses} Pokémon fainted
          </div>
        </div>

        {/* Divider */}
        <div className="flex items-center justify-center">
          <div className="text-2xl font-bold text-slate-400">vs</div>
        </div>

        {/* Player 2 */}
        <div className="flex flex-col items-end">
          <div className="text-sm text-slate-600 mb-1">Player 2</div>
          <div className="flex items-center gap-2 justify-end">
            {!isPlayer1Winner && (
              <span className="inline-block px-3 py-1 bg-green-100 text-green-800 rounded-full text-xs font-bold">
                WINNER
              </span>
            )}
            <div className={`text-lg font-bold ${!isPlayer1Winner ? 'text-green-600' : 'text-slate-600'}`}>
              {battle.player2.name}
            </div>
          </div>
          <div className="text-sm text-slate-500 mt-1">
            {battle.player2.losses} Pokémon fainted
          </div>
        </div>
      </div>

      {/* Additional Info */}
      <div className="mt-6 pt-6 border-t border-slate-200 flex gap-8">
        <div>
          <div className="text-xs text-slate-500 uppercase font-semibold">Total Turns</div>
          <div className="text-2xl font-bold text-slate-800">{battle.stats.totalTurns}</div>
        </div>
        <div>
          <div className="text-xs text-slate-500 uppercase font-semibold">Battle ID</div>
          <div className="text-sm font-mono text-slate-600">{battle.id.slice(0, 12)}...</div>
        </div>
        <div>
          <div className="text-xs text-slate-500 uppercase font-semibold">Critical Hits</div>
          <div className="text-2xl font-bold text-slate-800">{battle.stats.criticalHits}</div>
        </div>
      </div>
    </div>
  );
}
