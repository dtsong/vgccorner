import React from 'react';
import { BoardState } from '@/lib/types/battle';

interface MiniBoardstateProps {
  boardState: BoardState;
  player1Name: string;
  player2Name: string;
}

export default function MiniBoardstate({
  boardState,
  player1Name,
  player2Name,
}: MiniBoardstateProps) {
  const getPokemonSprite = (species: string): string => {
    // Use PokéAPI sprites or a placeholder
    const cleanSpecies = species.toLowerCase().replace(/[^a-z0-9]/g, '');
    return `https://img.pokemondb.net/sprites/home/normal/${cleanSpecies}.png`;
  };

  return (
    <div className="bg-white rounded-xl p-6 mb-6">
      <h2 className="text-lg font-semibold text-gray-700 mb-4">
        Mini-Boardstate
      </h2>
      <div className="grid grid-cols-2 gap-8">
        {/* Player 1 Side */}
        <div>
          <h3 className="text-sm font-medium text-gray-600 mb-3">
            {player1Name}
          </h3>
          <div className="relative bg-gradient-to-br from-green-100 to-green-200 rounded-2xl p-6 min-h-[200px]">
            <div className="flex justify-center gap-4 items-end">
              {boardState.player1Active.map((pokemon, idx) => (
                <div key={idx} className="relative">
                  <img
                    src={getPokemonSprite(pokemon.species)}
                    alt={pokemon.nickname || pokemon.species}
                    className="w-24 h-24 object-contain drop-shadow-lg"
                    onError={(e) => {
                      const target = e.target as HTMLImageElement;
                      target.src = '/placeholder-pokemon.png';
                    }}
                  />
                  {pokemon.isLead && (
                    <span className="absolute bottom-0 left-1/2 -translate-x-1/2 px-2 py-0.5 bg-gray-700 text-white text-xs rounded-full">
                      Lead
                    </span>
                  )}
                  <div className="mt-1 text-center">
                    <div className="text-xs font-medium text-gray-700">
                      {pokemon.nickname || pokemon.species}
                    </div>
                    <div className="text-xs text-gray-500">
                      {Math.round((pokemon.hp / pokemon.maxHp) * 100)}% HP
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Arrow */}
        <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-10">
          <svg
            className="w-12 h-12 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M14 5l7 7m0 0l-7 7m7-7H3"
            />
          </svg>
        </div>

        {/* Player 2 Side */}
        <div>
          <h3 className="text-sm font-medium text-gray-600 mb-3 text-right">
            Opponent&apos;s Team
          </h3>
          <div className="relative bg-gradient-to-br from-red-100 to-red-200 rounded-2xl p-6 min-h-[200px]">
            <div className="flex justify-center gap-4 items-end">
              {boardState.player2Active.map((pokemon, idx) => (
                <div key={idx} className="relative">
                  <img
                    src={getPokemonSprite(pokemon.species)}
                    alt={pokemon.nickname || pokemon.species}
                    className="w-24 h-24 object-contain drop-shadow-lg"
                    onError={(e) => {
                      const target = e.target as HTMLImageElement;
                      target.src = '/placeholder-pokemon.png';
                    }}
                  />
                  {pokemon.isLead && (
                    <span className="absolute bottom-0 left-1/2 -translate-x-1/2 px-2 py-0.5 bg-gray-700 text-white text-xs rounded-full">
                      Lead
                    </span>
                  )}
                  <div className="mt-1 text-center">
                    <div className="text-xs font-medium text-gray-700">
                      {pokemon.nickname || pokemon.species}
                    </div>
                    <div className="text-xs text-gray-500">
                      {Math.round((pokemon.hp / pokemon.maxHp) * 100)}% HP
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
