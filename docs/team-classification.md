# Team Classification System

VGC Corner automatically classifies teams based on their composition, moves, abilities, and items to help players understand their opponent's strategy.

## Classification Criteria

Teams are analyzed and classified into the following archetypes (in priority order):

### 1. **Hard Trick Room**
- **Criteria**: Trick Room on 2+ Pokémon
- **Description**: A team built around Trick Room with multiple setters for reliability
- **Example**: Cresselia + Dusclops with Trick Room

### 2. **TailRoom**
- **Criteria**: Both Tailwind AND Trick Room on the team
- **Description**: A flexible team that can operate under both Tailwind and Trick Room
- **Example**: Talonflame (Tailwind) + Cresselia (Trick Room)

### 3. **Sun Offense**
- **Criteria**: Drought ability OR Sunny Day move
- **Description**: An offensive team utilizing sun weather to power up Fire-type attacks
- **Example**: Torkoal (Drought) + Venusaur + Charizard

### 4. **Rain Offense**
- **Criteria**: Drizzle ability OR Rain Dance move
- **Description**: An offensive team utilizing rain weather to power up Water-type attacks
- **Example**: Pelipper (Drizzle) + Kingdra + Ludicolo

### 5. **Balance Bros**
- **Criteria**: Incineroar + Rillaboom on the same team
- **Description**: A balanced team featuring Incineroar and Rillaboom for defensive synergy
- **Strategy**: Uses Intimidate and Grassy Terrain for control

### 6. **Psy-Spam**
- **Criteria**: Psychic Terrain + Expanding Force user
- **Description**: A team focused on Psychic Terrain with Expanding Force for massive spread damage
- **Example**: Indeedee (Psychic Terrain) + Armarouge (Expanding Force)

### 7. **Tailwind Hyper Offense**
- **Criteria**: Tailwind + Choice Specs/Band/Scarf
- **Description**: An aggressive team using Tailwind and Choice items for overwhelming speed and power
- **Example**: Talonflame (Tailwind) + Choice Specs Dragapult + Choice Band Rillaboom

### 8. **Generic Archetypes**
- **Tailwind**: Team with Tailwind but no Choice items
- **Trick Room**: Team with Trick Room (but only 1 user)
- **Sun/Rain/Sand/Snow**: Weather teams without specific offensive setup

### 9. **Unclassified**
- **Criteria**: Doesn't fit any of the above patterns
- **Description**: A team that doesn't fit standard VGC archetypes
- **Note**: Could be a unique strategy or "goodstuff" team

## How Classification Works

1. **Team Analysis**: Each Pokémon is examined for:
   - Abilities (e.g., Drought, Drizzle, Intimidate)
   - Moves (e.g., Trick Room, Tailwind, Expanding Force)
   - Items (e.g., Choice Specs, Choice Band)
   - Species (e.g., Incineroar, Rillaboom)

2. **Criteria Matching**: The system checks criteria in priority order

3. **First Match Wins**: The first archetype that matches becomes the team's classification

4. **Detailed Metadata**: Additional information is stored:
   - List of Trick Room users
   - List of Tailwind users
   - Weather setters and weather type
   - Choice item users
   - Additional descriptive tags

## API Response

Team classification is included in the battle analysis response:

```json
{
  "player1": {
    "name": "Player1",
    "team": [...],
    "teamArchetype": "Hard Trick Room",
    "classification": {
      "archetype": "Hard Trick Room",
      "hasTrickRoom": true,
      "trickRoomUsers": ["Cresselia", "Dusclops"],
      "hasTailwind": false,
      "tailwindUsers": [],
      "hasWeatherSetter": false,
      "weatherType": "",
      "weatherSetters": [],
      "hasPsyTerrain": false,
      "psyTerrainUsers": [],
      "hasBalanceBros": false,
      "hasChoiceItems": false,
      "choiceUsers": [],
      "tags": []
    }
  }
}
```

## Turn-by-Turn Analysis

The turn analysis endpoint (`GET /api/showdown/replays/{replayId}/turns`) provides:

1. **Board State**: Current Pokémon on the field with HP and status
2. **Events**: Ordered list of actions taken during each turn
3. **Move Impact**: Detailed tracking of:
   - Damage dealt
   - Speed control (Trick Room, Tailwind, paralysis)
   - Weather/terrain changes
   - Status conditions
   - Critical hits
   - Type effectiveness
   - Fake Out usage
   - Stat changes
   - Fainted Pokémon

4. **Team Archetypes**: Both players' team classifications and descriptions

## Example Turn Event

```json
{
  "turnNumber": 1,
  "events": [
    {
      "type": "move",
      "pokemon": "Gengar",
      "action": "used Shadow Ball",
      "target": "on Dusclops",
      "result": "critical-hit",
      "details": "Critical Hit",
      "playerSide": "player1"
    }
  ],
  "boardState": {
    "player1Active": [{
      "species": "Gengar",
      "position": 0,
      "hp": 100,
      "maxHp": 100,
      "isLead": true
    }],
    "player2Active": [{
      "species": "Dusclops",
      "position": 0,
      "hp": 42,
      "maxHp": 100,
      "isLead": true
    }]
  }
}
```

## Future Enhancements

Planned improvements to the classification system:

- [ ] Detect "Goodstuff" vs structured teams
- [ ] Identify specific restricted Pokémon strategies
- [ ] Track team synergies (e.g., Follow Me + Setup sweeper)
- [ ] Recognize common speed tiers
- [ ] Identify protect/fake out patterns
- [ ] Track Tera type usage patterns
- [ ] Detect offensive vs defensive team compositions
