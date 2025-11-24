import { render, screen } from '@testing-library/react';
import BattleHeader from '../BattleHeader';
import { BattleSummary } from '@/lib/types/showdown';

// Mock battle data
const mockBattle: BattleSummary = {
  id: 'battle-12345678-abcd-efgh',
  format: 'Regulation H',
  timestamp: '2024-11-24T10:30:00Z',
  duration: 900,
  winner: 'player1',
  player1: {
    name: 'Ash',
    team: [],
    active: null,
    losses: 2,
    totalLeft: 4,
    activeIndex: 0,
  },
  player2: {
    name: 'Gary',
    team: [],
    active: null,
    losses: 4,
    totalLeft: 2,
    activeIndex: 0,
  },
  turns: [],
  stats: {
    totalTurns: 25,
    moveFrequency: {},
    typeCoverage: {},
    switches: 10,
    criticalHits: 3,
    superEffective: 5,
    notVeryEffective: 2,
    avgDamagePerTurn: 50,
    avgHealPerTurn: 10,
    player1Stats: {
      moveCount: 20,
      switchCount: 5,
      damageDealt: 500,
      damageTaken: 300,
      healingDone: 100,
      healingReceived: 50,
      movesByType: {},
      effectiveness: {
        superEffective: 3,
        notVeryEffective: 1,
        neutral: 16,
      },
    },
    player2Stats: {
      moveCount: 18,
      switchCount: 5,
      damageDealt: 300,
      damageTaken: 500,
      healingDone: 50,
      healingReceived: 100,
      movesByType: {},
      effectiveness: {
        superEffective: 2,
        notVeryEffective: 1,
        neutral: 15,
      },
    },
    turningPoints: [],
  },
  keyMoments: [],
};

describe('BattleHeader', () => {
  it('renders battle format and player names', () => {
    render(<BattleHeader battle={mockBattle} />);

    expect(screen.getByText(/Regulation H/i)).toBeInTheDocument();
    expect(screen.getByText(/Ash vs Gary/i)).toBeInTheDocument();
  });

  it('shows winner badge for player 1', () => {
    render(<BattleHeader battle={mockBattle} />);

    const winnerBadges = screen.getAllByText(/WINNER/i);
    expect(winnerBadges).toHaveLength(1);
  });

  it('shows winner badge for player 2 when they win', () => {
    const player2WinBattle = { ...mockBattle, winner: 'player2' as const };
    render(<BattleHeader battle={player2WinBattle} />);

    const winnerBadges = screen.getAllByText(/WINNER/i);
    expect(winnerBadges).toHaveLength(1);
  });

  it('displays Pokemon fainted counts', () => {
    render(<BattleHeader battle={mockBattle} />);

    expect(screen.getByText('2 Pokémon fainted')).toBeInTheDocument();
    expect(screen.getByText('4 Pokémon fainted')).toBeInTheDocument();
  });

  it('displays total turns', () => {
    render(<BattleHeader battle={mockBattle} />);

    expect(screen.getByText('25')).toBeInTheDocument();
  });

  it('displays battle ID (truncated)', () => {
    render(<BattleHeader battle={mockBattle} />);

    expect(screen.getByText(/battle-12345/i)).toBeInTheDocument();
  });

  it('displays critical hits count', () => {
    render(<BattleHeader battle={mockBattle} />);

    expect(screen.getByText(/Critical Hits/i)).toBeInTheDocument();
    // Check that the value 3 appears in the document
    const critValues = screen.getAllByText('3');
    expect(critValues.length).toBeGreaterThan(0);
  });

  it('formats timestamp correctly', () => {
    render(<BattleHeader battle={mockBattle} />);

    // The exact format depends on locale, but should contain month and year
    expect(screen.getByText(/Nov.*2024/i)).toBeInTheDocument();
  });
});
