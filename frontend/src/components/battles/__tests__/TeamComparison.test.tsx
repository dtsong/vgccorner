import { render, screen } from '@testing-library/react';
import TeamComparison from '../TeamComparison';
import { Player, Pokémon } from '@/lib/types/showdown';

// Mock Pokemon data
const createMockPokemon = (overrides: Partial<Pokémon> = {}): Pokémon => ({
  id: 'pikachu',
  name: 'Pikachu',
  level: 50,
  gender: 'M',
  ability: 'Static',
  item: 'Light Ball',
  stats: {
    hp: 100,
    attack: 85,
    defense: 60,
    spAtk: 90,
    spDef: 70,
    speed: 110,
  },
  moves: [],
  happiness: 255,
  shiny: false,
  currentHP: 100,
  maxHP: 100,
  status: '',
  teraType: '',
  ...overrides,
});

const mockPlayer1: Player = {
  name: 'Ash',
  team: [
    createMockPokemon({ name: 'Pikachu', currentHP: 80, maxHP: 100 }),
    createMockPokemon({ name: 'Charizard', currentHP: 0, maxHP: 150 }),
    createMockPokemon({ name: 'Blastoise', currentHP: 130, maxHP: 150, status: 'burn' }),
  ],
  active: null,
  losses: 1,
  totalLeft: 2,
  activeIndex: 0,
};

const mockPlayer2: Player = {
  name: 'Gary',
  team: [
    createMockPokemon({ name: 'Gengar', currentHP: 90, maxHP: 120, teraType: 'Ghost' }),
    createMockPokemon({ name: 'Alakazam', currentHP: 30, maxHP: 110 }),
  ],
  active: null,
  losses: 0,
  totalLeft: 2,
  activeIndex: 0,
};

describe('TeamComparison', () => {
  it('renders both player names', () => {
    render(<TeamComparison player1={mockPlayer1} player2={mockPlayer2} />);

    const playerNames = screen.getAllByText('Ash');
    expect(playerNames.length).toBeGreaterThan(0);
    const garyNames = screen.getAllByText('Gary');
    expect(garyNames.length).toBeGreaterThan(0);
  });

  it('displays all Pokemon names', () => {
    render(<TeamComparison player1={mockPlayer1} player2={mockPlayer2} />);

    expect(screen.getByText('Pikachu')).toBeInTheDocument();
    expect(screen.getByText('Charizard')).toBeInTheDocument();
    expect(screen.getByText('Blastoise')).toBeInTheDocument();
    expect(screen.getByText('Gengar')).toBeInTheDocument();
    expect(screen.getByText('Alakazam')).toBeInTheDocument();
  });

  it('shows fainted status for Pokemon with 0 HP', () => {
    render(<TeamComparison player1={mockPlayer1} player2={mockPlayer2} />);

    const faintedBadges = screen.getAllByText(/FAINTED/i);
    expect(faintedBadges).toHaveLength(1);
  });

  it('displays HP bars for active Pokemon', () => {
    render(<TeamComparison player1={mockPlayer1} player2={mockPlayer2} />);

    expect(screen.getByText('80/100 HP')).toBeInTheDocument();
    expect(screen.getByText('130/150 HP')).toBeInTheDocument();
    expect(screen.getByText('90/120 HP')).toBeInTheDocument();
    expect(screen.getByText('30/110 HP')).toBeInTheDocument();
  });

  it('displays HP percentage', () => {
    render(<TeamComparison player1={mockPlayer1} player2={mockPlayer2} />);

    expect(screen.getByText('80%')).toBeInTheDocument();
    expect(screen.getByText('87%')).toBeInTheDocument();
    expect(screen.getByText('75%')).toBeInTheDocument();
    expect(screen.getByText('27%')).toBeInTheDocument();
  });

  it('shows Pokemon abilities', () => {
    render(<TeamComparison player1={mockPlayer1} player2={mockPlayer2} />);

    const abilities = screen.getAllByText('Static');
    expect(abilities.length).toBeGreaterThan(0);
  });

  it('shows Pokemon items', () => {
    render(<TeamComparison player1={mockPlayer1} player2={mockPlayer2} />);

    const items = screen.getAllByText('Light Ball');
    expect(items.length).toBeGreaterThan(0);
  });

  it('shows status conditions', () => {
    render(<TeamComparison player1={mockPlayer1} player2={mockPlayer2} />);

    expect(screen.getByText('burn')).toBeInTheDocument();
  });

  it('shows Tera type when terastallized', () => {
    render(<TeamComparison player1={mockPlayer1} player2={mockPlayer2} />);

    expect(screen.getByText('Tera: Ghost')).toBeInTheDocument();
  });

  it('displays Pokemon remaining count', () => {
    render(<TeamComparison player1={mockPlayer1} player2={mockPlayer2} />);

    expect(screen.getByText('2/3')).toBeInTheDocument();
    expect(screen.getByText('2/2')).toBeInTheDocument();
  });

  it('applies correct styling to fainted Pokemon', () => {
    render(<TeamComparison player1={mockPlayer1} player2={mockPlayer2} />);

    const charizardName = screen.getByText('Charizard');
    expect(charizardName).toHaveClass('line-through');
  });

  it('does not show HP bar for fainted Pokemon', () => {
    render(<TeamComparison player1={mockPlayer1} player2={mockPlayer2} />);

    // Charizard is fainted, should not have HP display
    expect(screen.queryByText('0/150 HP')).not.toBeInTheDocument();
  });
});
