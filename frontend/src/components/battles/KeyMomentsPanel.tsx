'use client';

import React from 'react';
import { KeyMoment, TurningPoint } from '@/lib/types/showdown';

interface KeyMomentsPanelProps {
  keyMoments: KeyMoment[];
  turningPoints: TurningPoint[];
}

/**
 * Displays key moments and turning points from the battle.
 * Highlights critical events that shaped the outcome.
 */
export default function KeyMomentsPanel({ keyMoments, turningPoints }: KeyMomentsPanelProps) {
  // Merge and sort key moments and turning points by turn
  interface MomentDisplay {
    turnNumber: number;
    description: string;
    type: string;
    significance: number;
    isTurningPoint: boolean;
    momentumShift?: number;
  }

  const allMoments: MomentDisplay[] = [
    ...keyMoments.map((km) => ({
      turnNumber: km.turnNumber,
      description: km.description,
      type: km.type,
      significance: km.significance,
      isTurningPoint: false,
    })),
    ...turningPoints.map((tp) => ({
      turnNumber: tp.turnNumber,
      description: tp.description,
      type: 'turning_point' as const,
      significance: tp.significance,
      isTurningPoint: true,
      momentumShift: tp.momentumShift,
    })),
  ].sort((a, b) => a.turnNumber - b.turnNumber);

  const getSignificanceColor = (significance: number): string => {
    if (significance >= 9) return 'bg-red-100 text-red-700 border-red-300';
    if (significance >= 7) return 'bg-orange-100 text-orange-700 border-orange-300';
    if (significance >= 5) return 'bg-yellow-100 text-yellow-700 border-yellow-300';
    return 'bg-slate-100 text-slate-700 border-slate-300';
  };

  const getTypeIcon = (type: string): string => {
    switch (type) {
      case 'ko':
        return '💥';
      case 'switch':
        return '🔄';
      case 'status':
        return '⚠️';
      case 'weather':
        return '🌦️';
      case 'turning_point':
        return '📈';
      default:
        return '📌';
    }
  };

  if (allMoments.length === 0) {
    return (
      <div className="flex items-center justify-center h-40 bg-slate-50 rounded border border-slate-200">
        <p className="text-slate-500">No key moments recorded</p>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {allMoments.map((moment, idx) => (
        <div
          key={idx}
          className={`p-4 rounded-lg border-2 ${getSignificanceColor(
            moment.significance,
          )} transition-shadow hover:shadow-md`}
        >
          <div className="flex items-start justify-between gap-4">
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 mb-1">
                <span className="text-xl">{getTypeIcon(moment.type)}</span>
                <span className="inline-block px-2 py-0.5 bg-slate-700 text-white rounded text-xs font-semibold">
                  Turn {moment.turnNumber}
                </span>
                {moment.isTurningPoint && (
                  <span className="inline-block px-2 py-0.5 bg-purple-200 text-purple-800 rounded text-xs font-semibold">
                    TURNING POINT
                  </span>
                )}
              </div>
              <p className="text-sm font-medium break-words">{moment.description}</p>

              {/* Show momentum shift for turning points */}
              {moment.isTurningPoint && moment.momentumShift !== undefined && (
                <div className="mt-2 pt-2 border-t border-current border-opacity-20">
                  <div className="text-xs font-semibold">
                    Momentum Shift:{' '}
                    {moment.momentumShift > 0 ? (
                      <span className="text-blue-600">+{moment.momentumShift.toFixed(1)} (Player 1)</span>
                    ) : (
                      <span className="text-red-600">{moment.momentumShift.toFixed(1)} (Player 2)</span>
                    )}
                  </div>
                </div>
              )}
            </div>

            {/* Significance Badge */}
            <div className="flex flex-col items-end">
              <div className="font-bold text-lg">{moment.significance}</div>
              <div className="text-xs opacity-75 whitespace-nowrap">
                {moment.significance >= 8 ? 'Critical' : moment.significance >= 6 ? 'Major' : 'Minor'}
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
