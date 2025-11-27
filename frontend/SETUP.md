# VGC Corner Refinery Frontend Setup

## Overview

VGC Corner Refinery is the landing page for VGC Corner, allowing users to analyze Pokémon Showdown replays and view their analysis history.

## Environment Configuration

Create a `.env.local` file in the frontend directory:

```bash
# Backend API URL
NEXT_PUBLIC_API_URL=http://localhost:8080
```

For production, change the URL to your deployed backend service.

## Running the Application

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm run dev
```

3. Open [http://localhost:3000](http://localhost:3000) in your browser

## Features

### Landing Page (/)

The VGC Corner Refinery landing page includes:

- **Replay URL Input**: Paste a Pokémon Showdown replay URL to analyze
- **Analyze Button**: Processes the replay and extracts battle data
- **Recent Replays**: Shows recently analyzed battles with timestamps
- **Clickable Replays**: Click any recent replay to view detailed analysis

### Replay Detail Page (/replay/[replayId])

The replay detail page displays:

- **Battle Information**: ID, format, date, duration, and turn count
- **Player Teams**: Detailed Pokémon information for both players
  - Pokémon species and nicknames
  - Abilities and items
  - Move sets
  - Tera types (if applicable)
- **Winner Badge**: Indicates which player won the battle
- **Key Moments** (when available): Important events during the battle

## API Integration

The frontend connects to the following backend endpoints:

- `POST /api/showdown/analyze` - Analyze a new replay
- `GET /api/showdown/replays` - List recent replays
- `GET /api/showdown/replays/{replayId}` - Get specific replay details

### Supported Replay URL Formats

The app accepts various Pokémon Showdown replay URL formats:

- Full URL: `https://replay.pokemonshowdown.com/gen9vgc2024regh-2024567890`
- Without protocol: `replay.pokemonshowdown.com/gen9vgc2024regh-2024567890`
- Just the ID: `gen9vgc2024regh-2024567890`

## Development

### File Structure

```
frontend/src/
├── app/
│   ├── page.tsx                    # VGC Corner Refinery landing page
│   ├── replay/
│   │   └── [replayId]/
│   │       └── page.tsx            # Replay detail page
│   ├── layout.tsx                  # Root layout
│   └── globals.css                 # Global styles
├── lib/
│   ├── api/
│   │   └── client.ts               # API client functions
│   └── types/
│       └── api.ts                  # TypeScript type definitions
└── components/                     # Reusable components (future)
```

### API Client Usage

```typescript
import { analyzeReplayUrl, listReplays, getReplay } from '@/lib/api/client';

// Analyze a replay
const response = await analyzeReplayUrl('https://replay.pokemonshowdown.com/...');

// List recent replays
const replays = await listReplays({ limit: 10 });

// Get specific replay
const replay = await getReplay('gen9vgc2024regh-2024567890');
```

## Testing

Run tests:
```bash
npm test
```

Run tests in watch mode:
```bash
npm run test:watch
```

Generate coverage report:
```bash
npm run test:coverage
```

## Building for Production

```bash
npm run build
npm start
```

## Notes

- The app uses Tailwind CSS for styling
- All pages are client-side rendered (`'use client'`) for dynamic interactions
- Error handling is built into all API calls
- Loading states are displayed during data fetching
