# Datasette Global Leaderboard Setup

This document explains how to set up and deploy the global leaderboard for Unicorn using Datasette and Vercel.

## Architecture

The global leaderboard system consists of three components:

1. **Datasette Instance** - Hosted on Vercel with SQLite database, provides read-only JSON API for leaderboard data
2. **Serverless API** - Go-based Vercel function that accepts score submissions
3. **GitHub Pages** - Displays the leaderboard with live data fetching

**Note:** This setup uses SQLite (not Postgres) which is simpler and works perfectly with Vercel's serverless deployment.

## Prerequisites

- Vercel account (free tier works)
- Vercel CLI: `npm install -g vercel`
- Python 3.11+ (for datasette)
- pip packages: `pip install datasette datasette-publish-vercel`

## Deployment Steps

### Step 1: Initial Database Setup

Create an empty leaderboard database:

```bash
sqlite3 leaderboard.db <<EOF
CREATE TABLE IF NOT EXISTS game_scores (
  id TEXT PRIMARY KEY,
  player_name TEXT NOT NULL,
  final_net_worth INTEGER NOT NULL,
  roi REAL NOT NULL,
  successful_exits INTEGER NOT NULL,
  turns_played INTEGER NOT NULL,
  difficulty TEXT NOT NULL,
  played_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_net_worth ON game_scores(final_net_worth DESC);
CREATE INDEX IF NOT EXISTS idx_roi ON game_scores(roi DESC);
CREATE INDEX IF NOT EXISTS idx_player ON game_scores(player_name);
CREATE INDEX IF NOT EXISTS idx_difficulty ON game_scores(difficulty);
EOF
```

### Step 2: Deploy Datasette to Vercel

```bash
# Login to Vercel
vercel login

# Deploy datasette with the database
datasette publish vercel leaderboard.db \
  --project=unicorn-leaderboard \
  --install=datasette-cors \
  --metadata=datasette-metadata.json \
  --setting sql_time_limit_ms 3500 \
  --setting allow_download off
```

This will give you a URL like: `https://unicorn-leaderboard.vercel.app`

### Step 3: Deploy the Submission API

```bash
# Deploy the Vercel serverless functions
vercel deploy --prod
```

This deploys the Go-based API endpoint at `/api/submit-score`

### Step 4: Update Configuration

Update the following files with your Vercel URLs:

1. **docs/index.html** - Line 617:
   ```javascript
   const DATASETTE_URL = 'https://YOUR-PROJECT.vercel.app/leaderboard/game_scores.json';
   ```

2. **leaderboard/leaderboard.go** - Line 17:
   ```go
   DefaultAPIEndpoint = "https://YOUR-PROJECT.vercel.app/api/submit-score"
   ```

### Step 5: Set Up GitHub Actions (Optional)

The workflow at `.github/workflows/datasette-deploy.yml` can automatically redeploy datasette.

Add your Vercel token as a GitHub secret:
1. Go to https://vercel.com/account/tokens
2. Create a new token
3. Add it to GitHub repo secrets as `VERCEL_TOKEN`

The workflow will run every 6 hours or can be triggered manually.

## Testing the Integration

### Test Score Submission

```bash
curl -X POST https://YOUR-PROJECT.vercel.app/api/submit-score \
  -H "Content-Type: application/json" \
  -d '{
    "player_name": "TestPlayer",
    "final_net_worth": 1000000,
    "roi": 300.5,
    "successful_exits": 3,
    "turns_played": 120,
    "difficulty": "Medium"
  }'
```

Expected response:
```json
{
  "success": true,
  "message": "Score submitted successfully!",
  "id": "uuid-here"
}
```

### Test Datasette API

```bash
curl https://YOUR-PROJECT.vercel.app/leaderboard/game_scores.json?_size=10&_sort_desc=final_net_worth
```

### Test from Game

1. Build and run the game: `go build && ./unicorn`
2. Complete a game
3. Choose "Yes" when asked to submit to global leaderboard
4. Check the GitHub Pages to see your score

## Database Management

### View Current Data

Visit your Datasette instance: `https://YOUR-PROJECT.vercel.app`

### Backup Database

```bash
# Download current database from Vercel
vercel env pull .env.local
# Access via Vercel dashboard or datasette download feature
```

### Add Sample Data

```bash
# Add to leaderboard.db
sqlite3 leaderboard.db <<EOF
INSERT INTO game_scores VALUES 
  ('uuid1', 'Alice', 5000000, 1900.0, 5, 120, 'Expert', '2025-01-01 12:00:00'),
  ('uuid2', 'Bob', 3000000, 1100.0, 3, 120, 'Hard', '2025-01-02 12:00:00'),
  ('uuid3', 'Charlie', 1500000, 500.0, 2, 120, 'Medium', '2025-01-03 12:00:00');
EOF

# Redeploy
datasette publish vercel leaderboard.db --project=unicorn-leaderboard
```

## Troubleshooting

### API Not Working

1. Check Vercel logs: `vercel logs`
2. Verify the API endpoint is accessible: `curl YOUR-API-URL`
3. Check CORS settings in `vercel.json`

### Datasette Not Loading

1. Verify deployment: `vercel ls`
2. Check datasette URL is correct
3. Look at browser console for errors

### Game Can't Submit Scores

1. Ensure API endpoint is correct in `leaderboard/leaderboard.go`
2. Rebuild the game: `go build`
3. Check network connectivity
4. Verify API is accepting POST requests

## Cost Considerations

- **Vercel Free Tier**: 
  - 100 GB bandwidth/month
  - 100 hours serverless function execution/month
  - Sufficient for most indie games

- **Scaling**: If you exceed free tier:
  - Consider rate limiting submissions
  - Cache datasette responses
  - Use Vercel Pro ($20/month)

## Security Notes

1. The API accepts any submissions (no authentication)
2. Consider adding:
   - Rate limiting (e.g., 1 submission per player per day)
   - CAPTCHA for web submissions
   - Validation of score ranges
   - Profanity filter for player names

## Alternative Hosting

Instead of Vercel, you can use:

- **Fly.io**: `datasette publish fly leaderboard.db`
- **Cloudflare Pages**: Deploy static JSON exports
- **GitHub Actions + Pages**: Generate static JSON files periodically

## Features

- ? Real-time global leaderboard
- ? Filter by difficulty
- ? Sort by net worth or ROI
- ? Responsive design
- ? Automatic updates
- ? No backend maintenance for you!

## Support

For issues or questions:
- GitHub Issues: https://github.com/jamesacampbell/unicorn/issues
- Datasette Docs: https://docs.datasette.io
- Vercel Docs: https://vercel.com/docs
