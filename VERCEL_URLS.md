# Vercel URLs Configuration

## Production Deployment

**Main App URL:** `https://unicorn-green.vercel.app`

## API Endpoints

### 1. Submit Score API
- **URL:** `https://unicorn-green.vercel.app/api/submit-score`
- **Method:** POST
- **Status:** ✅ Working
- **Purpose:** Submit game scores to the leaderboard
- **Request Body:**
```json
{
  "player_name": "YourName",
  "final_net_worth": 5000000,
  "roi": 1500.0,
  "successful_exits": 4,
  "turns_played": 120,
  "difficulty": "Medium"
}
```

### 2. Get Leaderboard API
- **URL:** `https://unicorn-green.vercel.app/api/get-leaderboard`
- **Method:** GET
- **Status:** ⚠️ Deployed (Postgres driver issue)
- **Purpose:** Retrieve leaderboard data
- **Query Parameters:**
  - `_size`: Number of results (default: 10)
  - `_sort_desc`: Sort column (final_net_worth, roi, etc.)
  - `difficulty`: Filter by difficulty (Easy, Medium, Hard, Expert)

### 3. Datasette-compatible Endpoint
- **URL:** `https://unicorn-green.vercel.app/leaderboard/game_scores.json`
- **Method:** GET
- **Status:** ⚠️ Deployed (routes to get-leaderboard)
- **Purpose:** Datasette-compatible JSON API for frontend
- **Same as:** `/api/get-leaderboard` (routes to same handler)

## Frontend Configuration

**GitHub Pages:** `https://james-see.github.io/unicorn`

**Frontend API URL (docs/index.html line 618):**
```javascript
const DATASETTE_URL = 'https://unicorn-green.vercel.app/leaderboard/game_scores.json';
```

## Game Client Configuration

**Game API URL (leaderboard/leaderboard.go line 14):**
```go
DefaultAPIEndpoint = "https://unicorn-green.vercel.app/api/submit-score"
```

## Testing Commands

### Test Submit Score
```bash
curl -X POST https://unicorn-green.vercel.app/api/submit-score \
  -H "Content-Type: application/json" \
  -d '{
    "player_name": "Test",
    "final_net_worth": 5000000,
    "roi": 1500.0,
    "successful_exits": 4,
    "turns_played": 120,
    "difficulty": "Medium"
  }'
```

### Test Leaderboard
```bash
curl "https://unicorn-green.vercel.app/leaderboard/game_scores.json?_size=10&_sort_desc=final_net_worth"
```

## Known Issues

1. **Get Leaderboard Endpoint**: Currently returns Postgres driver error
   - Error: `sql: unknown driver "postgres" (forgotten import?)`
   - Submit score works because it uses the same Postgres connection
   - Build issue with get-leaderboard subdirectory

## Routes Configuration (vercel.json)

```json
{
  "routes": [
    {
      "src": "/api/submit-score",
      "dest": "/api/submit-score"
    },
    {
      "src": "/api/get-leaderboard",
      "dest": "/api/get-leaderboard"
    },
    {
      "src": "/leaderboard/game_scores.json",
      "dest": "/api/get-leaderboard"
    }
  ]
}
```

## Environment Variables Required

- `POSTGRES_URL`: Vercel Postgres connection string (set in Vercel dashboard)

## Last Updated

November 1, 2025

