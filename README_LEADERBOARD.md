# Global Leaderboard Integration

The Unicorn game now supports a **global leaderboard** where players can submit their scores and compete worldwide!

## Features

- ?? **Global Competition** - See top players from around the world
- ?? **Easy Submission** - One-click score submission from the game
- ?? **Multiple Views** - Filter by difficulty, sort by Net Worth or ROI
- ?? **Real-time Updates** - Leaderboard updates automatically
- ?? **Beautiful UI** - Modern, responsive design on GitHub Pages
- ?? **Privacy First** - Only game stats are submitted, no personal data
- ?? **Free Hosting** - Powered by Vercel (free tier) and Datasette

## How It Works

```
???????????????
?  Game Over  ?
???????????????
       ?
       ?? Save to Local DB ???????????
       ?                             ?
       ?? Prompt: Submit to Global? ??
                                     ?
                    ??????????????????
                    ?
                    v
         ????????????????????
         ? POST /api/submit ? (Vercel Serverless)
         ????????????????????
                  ?
                  v
         ???????????????????
         ?  SQLite Database? (Stored in Vercel)
         ???????????????????
                  ?
                  v
         ???????????????????
         ?    Datasette    ? (Read-only JSON API)
         ???????????????????
                  ?
                  v
         ???????????????????
         ?  GitHub Pages   ? (Displays leaderboard)
         ???????????????????
```

## Quick Start

See [QUICKSTART_LEADERBOARD.md](QUICKSTART_LEADERBOARD.md) for a 5-minute setup guide.

## Architecture

### Components

1. **Game Client** (`leaderboard/leaderboard.go`)
   - HTTP client for score submission
   - Validates API availability
   - User-friendly prompts

2. **Submission API** (`api/submit-score.go`)
   - Go-based Vercel serverless function
   - Accepts POST requests with score data
   - Writes to SQLite database with UUID

3. **Datasette** (`leaderboard.db` + `datasette-metadata.json`)
   - Read-only JSON API for leaderboard data
   - Supports filtering, sorting, pagination
   - Hosted on Vercel

4. **GitHub Pages** (`docs/index.html`)
   - JavaScript fetches from Datasette API
   - Displays formatted leaderboard
   - Filter tabs and refresh button

### Data Flow

```
Player finishes game
  ?
Prompt: "Submit to global leaderboard?"
  ? (if yes)
POST to https://YOUR-PROJECT.vercel.app/api/submit-score
  {
    "player_name": "Alice",
    "final_net_worth": 5000000,
    "roi": 1900.0,
    "successful_exits": 5,
    "turns_played": 120,
    "difficulty": "Expert"
  }
  ?
Vercel function validates and saves to SQLite
  ?
Datasette exposes via JSON API
  ?
GitHub Pages fetches and displays
```

## Files Added/Modified

### New Files
- `.github/workflows/datasette-deploy.yml` - GitHub Actions for auto-deployment
- `api/submit-score.go` - Serverless function for score submission
- `api/go.mod` - Dependencies for API
- `leaderboard/leaderboard.go` - HTTP client for game
- `datasette-metadata.json` - Datasette configuration
- `vercel.json` - Vercel deployment config
- `leaderboard.db` - SQLite database with sample data
- `DATASETTE_SETUP.md` - Detailed setup guide
- `QUICKSTART_LEADERBOARD.md` - Quick start guide

### Modified Files
- `main.go` - Added score submission prompt after game ends
- `docs/index.html` - Added global leaderboard display section
- `.gitignore` - Added Vercel and DB files

## Database Schema

```sql
CREATE TABLE game_scores (
  id TEXT PRIMARY KEY,              -- UUID
  player_name TEXT NOT NULL,
  final_net_worth INTEGER NOT NULL,
  roi REAL NOT NULL,
  successful_exits INTEGER NOT NULL,
  turns_played INTEGER NOT NULL,
  difficulty TEXT NOT NULL,
  played_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_net_worth ON game_scores(final_net_worth DESC);
CREATE INDEX idx_roi ON game_scores(roi DESC);
CREATE INDEX idx_player ON game_scores(player_name);
CREATE INDEX idx_difficulty ON game_scores(difficulty);
```

## API Endpoints

### Submit Score
```bash
POST https://YOUR-PROJECT.vercel.app/api/submit-score
Content-Type: application/json

{
  "player_name": "Alice",
  "final_net_worth": 5000000,
  "roi": 1900.0,
  "successful_exits": 5,
  "turns_played": 120,
  "difficulty": "Expert"
}

# Response
{
  "success": true,
  "message": "Score submitted successfully!",
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Get Leaderboard
```bash
GET https://YOUR-DATASETTE.vercel.app/leaderboard/game_scores.json?_size=10&_sort_desc=final_net_worth

# Optional parameters:
?difficulty=Expert          # Filter by difficulty
&_sort_desc=roi            # Sort by ROI instead
&_size=20                  # Get 20 results
&_offset=10                # Pagination
```

## Deployment

### Option 1: Vercel (Recommended)

```bash
# Deploy API and serve datasette
vercel deploy --prod

# Deploy datasette separately for better performance
datasette publish vercel leaderboard.db \
  --project=unicorn-leaderboard \
  --install=datasette-cors \
  --metadata=datasette-metadata.json
```

### Option 2: GitHub Actions

Push to GitHub and the workflow will auto-deploy datasette:
- Runs every 6 hours
- Can be triggered manually
- Requires `VERCEL_TOKEN` in GitHub secrets

### Option 3: Other Platforms

- **Fly.io**: `datasette publish fly leaderboard.db`
- **Cloudflare**: Deploy as Workers
- **Self-hosted**: Run datasette on your own server

## Security Considerations

1. **No Authentication** - Anyone can submit scores
   - Consider adding: Rate limiting, CAPTCHA, token validation

2. **Data Validation** - Basic validation in API
   - Add: Score range validation, profanity filter

3. **Privacy** - Only game stats are stored
   - No IP addresses, emails, or personal data

4. **Database Size** - SQLite has practical limits
   - Free tier supports ~100k scores easily
   - Consider archiving old scores monthly

## Cost Analysis

### Vercel Free Tier Limits
- 100 GB bandwidth/month
- 100 hours serverless execution/month
- 100 deployments/day

### Estimated Usage
- **1000 players/month**: 
  - ~50 MB bandwidth
  - ~2 hours function time
  - Well within free tier ?

- **10,000 players/month**:
  - ~500 MB bandwidth  
  - ~20 hours function time
  - Still within free tier ?

- **100,000+ players/month**:
  - Consider Vercel Pro ($20/mo)
  - Or switch to self-hosted solution

## Monitoring

### View Logs
```bash
vercel logs
vercel logs --follow
```

### Check Deployments
```bash
vercel ls
vercel inspect YOUR-DEPLOYMENT-URL
```

### Database Stats
Visit your Datasette instance and go to `/leaderboard` to see:
- Total scores
- Size on disk
- Table schema

## Troubleshooting

See [DATASETTE_SETUP.md](DATASETTE_SETUP.md#troubleshooting) for detailed troubleshooting guide.

## Future Enhancements

Potential features to add:

- [ ] Player profiles with historical stats
- [ ] Weekly/monthly leaderboards
- [ ] Achievement showcase on leaderboard
- [ ] Score verification (prevent cheating)
- [ ] Rate limiting (1 submission per player per day)
- [ ] Email notifications for top 10
- [ ] Discord bot integration
- [ ] Historical charts and trends
- [ ] Regional leaderboards
- [ ] Tournament mode

## Credits

Built with:
- [Datasette](https://datasette.io) by Simon Willison
- [Vercel](https://vercel.com) for hosting
- [SQLite](https://sqlite.org) for database
- [Go](https://golang.org) for the game and API

## License

Same as main project: MIT License
