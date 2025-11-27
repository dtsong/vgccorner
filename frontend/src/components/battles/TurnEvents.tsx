import React from 'react';
import { BattleEvent } from '@/lib/types/battle';

interface TurnEventsProps {
  events: BattleEvent[];
}

export default function TurnEvents({ events }: TurnEventsProps) {
  const getPokemonIcon = (species: string): string => {
    const cleanSpecies = species.toLowerCase().replace(/[^a-z0-9]/g, '');
    return `https://img.pokemondb.net/sprites/home/normal/${cleanSpecies}.png`;
  };

  const getEventColor = (event: BattleEvent): string => {
    if (!event.result) return 'bg-gray-50 border-l-gray-400';

    switch (event.result) {
      case 'critical-hit':
      case 'super-effective':
        return event.playerSide === 'player1'
          ? 'bg-green-50 border-l-green-500'
          : 'bg-red-50 border-l-red-500';
      case 'not-very-effective':
      case 'miss':
        return event.playerSide === 'player1'
          ? 'bg-red-50 border-l-red-300'
          : 'bg-green-50 border-l-green-300';
      case 'faint':
        return event.playerSide === 'player1'
          ? 'bg-red-100 border-l-red-600'
          : 'bg-green-100 border-l-green-600';
      default:
        return 'bg-gray-50 border-l-gray-400';
    }
  };

  const formatEventText = (event: BattleEvent): string => {
    let text = `${event.pokemon} ${event.action}`;

    if (event.target) {
      text += ` ${event.target}`;
    }

    if (event.details) {
      text += `! (${event.details})`;
    } else {
      text += '!';
    }

    return text;
  };

  return (
    <div className="space-y-3">
      {events.map((event, idx) => (
        <div
          key={idx}
          className={`flex items-start gap-3 p-4 border-l-4 rounded-r-lg ${getEventColor(event)}`}
        >
          <img
            src={getPokemonIcon(event.pokemon)}
            alt={event.pokemon}
            className="w-10 h-10 object-contain flex-shrink-0"
            onError={(e) => {
              const target = e.target as HTMLImageElement;
              target.style.display = 'none';
            }}
          />
          <div className="flex-1">
            <p className="text-sm text-gray-800 font-medium">
              {formatEventText(event)}
            </p>
            {event.result && (
              <span className="inline-block mt-1 text-xs px-2 py-0.5 bg-white/50 rounded">
                {event.result.replace(/-/g, ' ')}
              </span>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
