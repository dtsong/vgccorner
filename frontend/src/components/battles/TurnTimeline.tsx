'use client';

import React, { useState } from 'react';
import { BattleSummary } from '@/lib/types/showdown';

interface TurnTimelineProps {
  battle: BattleSummary;
}

/**
 * Displays a timeline of turns showing actions and position scores.
 * Users can expand individual turns to see detailed information.
 */
export default function TurnTimeline({ battle }: TurnTimelineProps) {
  const [expandedTurn, setExpandedTurn] = useState<number | null>(null);

  const toggleTurnExpand = (turnNumber: number) => {
    setExpandedTurn(expandedTurn === turnNumber ? null : turnNumber);
  };

  return (
    <div className="space-y-2">
      {battle.turns.map((turn) => {
        const isExpanded = expandedTurn === turn.turnNumber;
        const score = turn.positionScore;
        const p1Advantage = score && score.player1Score > score.player2Score;
        const p2Advantage = score && score.player2Score > score.player1Score;

        return (
          <div key={`turn-${turn.turnNumber}`} className="border border-slate-200 rounded-lg overflow-hidden">
            {/* Turn Header - Clickable */}
            <button
              onClick={() => toggleTurnExpand(turn.turnNumber)}
              className="w-full p-4 bg-slate-50 hover:bg-slate-100 transition-colors text-left flex items-center justify-between"
            >
              <div className="flex items-center gap-4 flex-1">
                <div className="font-bold text-slate-800 min-w-16">Turn {turn.turnNumber}</div>

                {/* Action Summary */}
                <div className="text-sm text-slate-600">
                  {turn.actions.map((action) => {
                    if (action.actionType === 'move' && action.move) {
                      return `${action.player === 'player1' ? 'P1' : 'P2'} used ${action.move.name}`;
                    } else if (action.actionType === 'switch') {
                      return `${action.player === 'player1' ? 'P1' : 'P2'} switched to ${action.switchTo}`;
                    }
                    return '';
                  })
                    .filter(Boolean)
                    .join(' • ')}
                </div>
              </div>

              {/* Position Score Indicator */}
              {score && (
                <div className="flex items-center gap-2 ml-4">
                  <div className="flex items-center gap-1">
                    <div className={`text-sm font-bold ${p1Advantage ? 'text-blue-600' : 'text-slate-600'}`}>
                      {score.player1Score.toFixed(0)}
                    </div>
                    <div className="text-xs text-slate-500">/</div>
                    <div className={`text-sm font-bold ${p2Advantage ? 'text-red-600' : 'text-slate-600'}`}>
                      {score.player2Score.toFixed(0)}
                    </div>
                  </div>
                  <div
                    className={`text-xs font-semibold px-2 py-1 rounded ${
                      score.momentumPlayer === 'player1'
                        ? 'bg-blue-100 text-blue-700'
                        : score.momentumPlayer === 'player2'
                          ? 'bg-red-100 text-red-700'
                          : 'bg-slate-100 text-slate-700'
                    }`}
                  >
                    {score.momentumPlayer === 'neutral' ? 'Neutral' : `${score.momentumPlayer} ahead`}
                  </div>
                </div>
              )}

              {/* Expand Icon */}
              <div className="ml-2 text-slate-500">
                {isExpanded ? '▼' : '▶'}
              </div>
            </button>

            {/* Expanded Details */}
            {isExpanded && (
              <div className="p-4 bg-white border-t border-slate-200 space-y-4">
                {/* Actions Detail */}
                <div>
                  <h4 className="font-semibold text-slate-800 mb-2">Actions</h4>
                  <div className="space-y-1">
                    {turn.actions.length === 0 ? (
                      <p className="text-sm text-slate-500">No actions recorded</p>
                    ) : (
                      turn.actions.map((action, idx) => (
                        <div key={idx} className="text-sm text-slate-700">
                          <span className={`font-semibold ${action.player === 'player1' ? 'text-blue-600' : 'text-red-600'}`}>
                            {action.player === 'player1' ? 'P1' : 'P2'}
                          </span>
                          {' • '}
                          {action.actionType === 'move' && action.move && `Used ${action.move.name}`}
                          {action.actionType === 'switch' && `Switched to ${action.switchTo}`}
                          {action.actionType === 'item' && `Used item: ${action.item}`}
                        </div>
                      ))
                    )}
                  </div>
                </div>

                {/* Damage & Healing */}
                {(Object.keys(turn.damageDealt).length > 0 || Object.keys(turn.healingDone).length > 0) && (
                  <div>
                    <h4 className="font-semibold text-slate-800 mb-2">Damage & Healing</h4>
                    <div className="space-y-1 text-sm">
                      {Object.entries(turn.damageDealt).map(([player, damage]) => (
                        <div key={`damage-${player}`} className="text-red-600">
                          {player}: -{damage} HP
                        </div>
                      ))}
                      {Object.entries(turn.healingDone).map(([player, healing]) => (
                        <div key={`heal-${player}`} className="text-green-600">
                          {player}: +{healing} HP
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Position Score Details */}
                {score && (
                  <div>
                    <h4 className="font-semibold text-slate-800 mb-2">Position Evaluation</h4>
                    <div className="grid grid-cols-2 gap-2 text-sm">
                      <div className="bg-blue-50 p-2 rounded">
                        <div className="text-blue-700 font-semibold">Player 1</div>
                        <div className="text-blue-600">{score.player1Score.toFixed(1)}/100</div>
                      </div>
                      <div className="bg-red-50 p-2 rounded">
                        <div className="text-red-700 font-semibold">Player 2</div>
                        <div className="text-red-600">{score.player2Score.toFixed(1)}/100</div>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
