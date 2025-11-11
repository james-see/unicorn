# Datasette Deployment for Unicorn Leaderboard

## Quick Fix for Failing Deployments

### The Issue
Deployment fails with: `ValueError: POSTGRES_URL environment variable is required`

### The Solution
Set the `POSTGRES_URL` environment variable in Vercel:

1. **Go to:** https://vercel.com/dashboard
2. **Select:** Your datasette project
3. **Navigate to:** Settings → Environment Variables
4. **Add variable:**
   - Name: `POSTGRES_URL`
   - Value: (copy from your main `unicorn` project)
   - Environments: Production, Preview, Development
5. **Redeploy:** `vercel --prod`

### Get POSTGRES_URL Value

From your main working project:
```bash
vercel env pull
grep POSTGRES_URL .env.local
```

Or from Vercel Dashboard → Main Project → Settings → Environment Variables

## Deployment

```bash
cd /workspace/datasette-deploy
vercel --prod
```

## Test

```bash
curl "https://YOUR-URL.vercel.app/leaderboard/game_scores.json?_size=5"
```

## Files in This Directory

- `index.py` - Datasette app with Postgres connection
- `requirements.txt` - Python dependencies
- `vercel.json` - Vercel configuration
- `runtime.txt` - Python 3.12
- `datasette-metadata.json` - Datasette metadata
- `leaderboard.db` - Local SQLite file (not used in Vercel deployment)

## Architecture

This deployment connects to the same Postgres database as your main API deployment, providing a Datasette web UI and JSON API for data exploration.

Main API handles writes (`/api/submit-score`), Datasette handles reads.

## Do You Need This?

Your main deployment already provides a JSON API at:
- `https://unicorn-green.vercel.app/api/get-leaderboard`
- `https://unicorn-green.vercel.app/leaderboard/game_scores.json`

Deploy Datasette only if you want:
- Rich web UI for browsing data
- Built-in SQL query interface
- Advanced filtering and export features

Otherwise, the main API is sufficient for your game's leaderboard.
