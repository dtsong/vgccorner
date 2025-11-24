'use client';

import React, { useMemo } from 'react';
import { BattleSummary } from '@/lib/types/showdown';

interface PositionScoreChartProps {
  battle: BattleSummary;
}

/**
 * Displays a line chart showing position scores for both players across turns.
 * Uses pure SVG for lightweight rendering without external charting dependencies.
 */
export default function PositionScoreChart({ battle }: PositionScoreChartProps) {
  const chartData = useMemo(() => {
    const data = battle.turns
      .filter((turn) => turn.positionScore)
      .map((turn) => ({
        turn: turn.turnNumber,
        p1: turn.positionScore!.player1Score,
        p2: turn.positionScore!.player2Score,
      }));

    return data;
  }, [battle.turns]);

  if (chartData.length === 0) {
    return (
      <div className="flex items-center justify-center h-64 bg-slate-50 rounded border border-slate-200">
        <p className="text-slate-500">No position data available</p>
      </div>
    );
  }

  const width = 800;
  const height = 300;
  const padding = 40;
  const innerWidth = width - padding * 2;
  const innerHeight = height - padding * 2;

  const maxTurn = Math.max(...chartData.map((d) => d.turn));
  const minScore = 0;
  const maxScore = 100;

  const xScale = (turn: number) => padding + (turn / maxTurn) * innerWidth;
  const yScale = (score: number) => padding + innerHeight - ((score - minScore) / (maxScore - minScore)) * innerHeight;

  // Generate path points for player 1
  const p1Path = chartData.map((d, i) => `${i === 0 ? 'M' : 'L'} ${xScale(d.turn)} ${yScale(d.p1)}`).join(' ');

  // Generate path points for player 2
  const p2Path = chartData.map((d, i) => `${i === 0 ? 'M' : 'L'} ${xScale(d.turn)} ${yScale(d.p2)}`).join(' ');

  return (
    <div className="overflow-x-auto">
      <svg width={width} height={height} className="mx-auto">
        {/* Grid lines */}
        {[0, 25, 50, 75, 100].map((score) => (
          <line
            key={`grid-${score}`}
            x1={padding}
            y1={yScale(score)}
            x2={width - padding}
            y2={yScale(score)}
            stroke="#e2e8f0"
            strokeWidth="1"
            strokeDasharray="2,2"
          />
        ))}

        {/* Axes */}
        <line x1={padding} y1={padding} x2={padding} y2={height - padding} stroke="#64748b" strokeWidth="2" />
        <line x1={padding} y1={height - padding} x2={width - padding} y2={height - padding} stroke="#64748b" strokeWidth="2" />

        {/* Player 1 Line (Blue) */}
        <path d={p1Path} stroke="#3b82f6" strokeWidth="3" fill="none" vectorEffect="non-scaling-stroke" />

        {/* Player 2 Line (Red) */}
        <path d={p2Path} stroke="#ef4444" strokeWidth="3" fill="none" vectorEffect="non-scaling-stroke" />

        {/* Data points */}
        {chartData.map((d) => (
          <g key={`point-${d.turn}`}>
            <circle cx={xScale(d.turn)} cy={yScale(d.p1)} r="3" fill="#3b82f6" />
            <circle cx={xScale(d.turn)} cy={yScale(d.p2)} r="3" fill="#ef4444" />
          </g>
        ))}

        {/* Y-axis labels */}
        {[0, 25, 50, 75, 100].map((score) => (
          <text
            key={`label-${score}`}
            x={padding - 10}
            y={yScale(score) + 4}
            textAnchor="end"
            fontSize="12"
            fill="#64748b"
          >
            {score}
          </text>
        ))}

        {/* X-axis labels (every 5 turns or so) */}
        {chartData
          .filter((_, i) => i % Math.ceil(chartData.length / 5) === 0 || i === chartData.length - 1)
          .map((d) => (
            <text
              key={`turn-${d.turn}`}
              x={xScale(d.turn)}
              y={height - padding + 20}
              textAnchor="middle"
              fontSize="12"
              fill="#64748b"
            >
              T{d.turn}
            </text>
          ))}
      </svg>

      {/* Legend */}
      <div className="flex justify-center gap-8 mt-6">
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 rounded-full bg-blue-500"></div>
          <span className="text-sm font-medium text-slate-700">{battle.player1.name}</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 rounded-full bg-red-500"></div>
          <span className="text-sm font-medium text-slate-700">{battle.player2.name}</span>
        </div>
      </div>

      {/* Stats */}
      <div className="mt-6 grid grid-cols-2 gap-4 text-sm">
        <div className="bg-blue-50 p-3 rounded">
          <div className="text-blue-700 font-semibold">{battle.player1.name}</div>
          <div className="text-blue-600">
            Avg Score: {(chartData.reduce((sum, d) => sum + d.p1, 0) / chartData.length).toFixed(1)}
          </div>
        </div>
        <div className="bg-red-50 p-3 rounded">
          <div className="text-red-700 font-semibold">{battle.player2.name}</div>
          <div className="text-red-600">
            Avg Score: {(chartData.reduce((sum, d) => sum + d.p2, 0) / chartData.length).toFixed(1)}
          </div>
        </div>
      </div>
    </div>
  );
}
