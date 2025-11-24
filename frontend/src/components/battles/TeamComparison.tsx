'use client';

import React from 'react';
import { Player } from '@/lib/types/showdown';

interface TeamComparisonProps {
  player1: Player;
  player2: Player;
}

/**
 * Displays team composition and status for both players.
 * Shows which Pokémon are still in battle and their current HP.
 */
export default function TeamComparison({ player1, player2 }: TeamComparisonProps) {
  const renderTeam = (player: Player, playerLabel: string) => {
    return (
      <div className="space-y-3">
        <div className="font-semibold text-slate-800">{playerLabel}</div>
        <div className="space-y-2">
          {player.team.map((poke, idx) => {
            const isFainted = poke.currentHP === 0;
            const hpPercent = poke.maxHP > 0 ? (poke.currentHP / poke.maxHP) * 100 : 0;

            return (
              <div key={idx} className={`p-2 rounded border ${isFainted ? 'bg-slate-100 border-slate-300' : 'bg-white border-slate-200'}`}>
                <div className="flex justify-between items-start gap-2 mb-1">
                  <div>
                    <div className={`text-sm font-semibold ${isFainted ? 'text-slate-400 line-through' : 'text-slate-800'}`}>
                      {poke.name}
                    </div>
                    <div className="text-xs text-slate-500">Lv. {poke.level}</div>
                  </div>
                  {isFainted && <span className="text-xs font-bold text-slate-500">FAINTED</span>}
                </div>

                {!isFainted && (
                  <>
                    {/* HP Bar */}
                    <div className="mb-1">
                      <div className="flex justify-between mb-1">
                        <span className="text-xs text-slate-600">
                          {poke.currentHP}/{poke.maxHP} HP
                        </span>
                        <span className="text-xs text-slate-500">{hpPercent.toFixed(0)}%</span>
                      </div>
                      <div className="w-full h-2 bg-slate-200 rounded overflow-hidden">
                        <div
                          className={`h-full transition-all ${
                            hpPercent > 50 ? 'bg-green-500' : hpPercent > 25 ? 'bg-yellow-500' : 'bg-red-500'
                          }`}
                          style={{ width: `${hpPercent}%` }}
                        ></div>
                      </div>
                    </div>

                    {/* Status & Item */}
                    <div className="flex gap-2 flex-wrap">
                      {poke.ability && (
                        <span className="text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded">
                          {poke.ability}
                        </span>
                      )}
                      {poke.item && (
                        <span className="text-xs bg-amber-100 text-amber-700 px-2 py-0.5 rounded">
                          {poke.item}
                        </span>
                      )}
                      {poke.status && (
                        <span className="text-xs bg-red-100 text-red-700 px-2 py-0.5 rounded">
                          {poke.status}
                        </span>
                      )}
                      {poke.teraType && (
                        <span className="text-xs bg-purple-100 text-purple-700 px-2 py-0.5 rounded">
                          Tera: {poke.teraType}
                        </span>
                      )}
                    </div>
                  </>
                )}
              </div>
            );
          })}
        </div>
      </div>
    );
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div className="space-y-3">
        {renderTeam(player1, player1.name)}
      </div>
      <div className="space-y-3">
        {renderTeam(player2, player2.name)}
      </div>

      {/* Summary Stats */}
      <div className="md:col-span-2 pt-4 border-t border-slate-200">
        <div className="grid grid-cols-2 gap-2 text-sm">
          <div className="text-center">
            <div className="font-semibold text-slate-700">{player1.name}</div>
            <div className="text-2xl font-bold text-slate-800">
              {player1.team.length - player1.losses}/{player1.team.length}
            </div>
            <div className="text-xs text-slate-500">Pokémon Remaining</div>
          </div>
          <div className="text-center">
            <div className="font-semibold text-slate-700">{player2.name}</div>
            <div className="text-2xl font-bold text-slate-800">
              {player2.team.length - player2.losses}/{player2.team.length}
            </div>
            <div className="text-xs text-slate-500">Pokémon Remaining</div>
          </div>
        </div>
      </div>
    </div>
  );
}
