# Quick Start: Global Leaderboard

Get your global leaderboard up and running in 5 minutes!

## Prerequisites

```bash
# Install Vercel CLI
npm install -g vercel

# Install Python and datasette
pip install datasette datasette-publish-vercel
```

## Deploy in 3 Steps

### 1. Deploy the Submission API

```bash
cd unicorn
vercel login
vercel deploy --prod
```

Note the deployment URL (e.g., `https://unicorn-abc123.vercel.app`)

### 2. Deploy Datasette

```bash
datasette publish vercel leaderboard.db \
  --project=unicorn-leaderboard \
  --install=datasette-cors \
  --metadata=datasette-metadata.json
```

Note the Datasette URL (e.g., `https://unicorn-leaderboard.vercel.app`)

### 3. Update URLs

**In `docs/index.html` (line 617):**
```javascript
const DATASETTE_URL = 'https://YOUR-DATASETTE-URL/leaderboard/game_scores.json';
```

**In `leaderboard/leaderboard.go` (line 17):**
```go
DefaultAPIEndpoint = "https://YOUR-VERCEL-URL/api/submit-score"
```

**Rebuild the game:**
```bash
go build
```

## Test It

```bash
# Test the API
curl -X POST https://YOUR-VERCEL-URL/api/submit-score \
  -H "Content-Type: application/json" \
  -d '{"player_name":"Test","final_net_worth":1000000,"roi":300,"successful_exits":2,"turns_played":120,"difficulty":"Medium"}'

# Play the game
./unicorn

# View the leaderboard
# Open https://jamesacampbell.github.io/unicorn
```

## GitHub Pages

Push your changes to GitHub:

```bash
git add .
git commit -m "Add global leaderboard with datasette"
git push
```

Your leaderboard will be live at: `https://USERNAME.github.io/unicorn`

## Troubleshooting

**API returns 404:**
- Make sure you deployed with `vercel deploy --prod`
- Check the URL matches what's in `leaderboard.go`

**Leaderboard shows error:**
- Verify Datasette is deployed: Visit `https://YOUR-DATASETTE-URL` 
- Check browser console for errors
- Ensure CORS is enabled (datasette-cors plugin)

**Game can't submit scores:**
- Rebuild after updating `leaderboard.go`: `go build`
- Check you have internet connection
- Try the curl test command above

## What's Included

? **Real-time global leaderboard** on your GitHub Pages  
? **Score submission** from the game with one click  
? **Filter by difficulty** - Easy, Medium, Hard, Expert  
? **Sort by Net Worth or ROI**  
? **Free hosting** on Vercel (100GB/month bandwidth)  
? **No database management** - datasette handles everything  

## Next Steps

- **Customize**: Edit `datasette-metadata.json` for branding
- **Add analytics**: Track submissions with Vercel Analytics
- **Set up automation**: Enable GitHub Actions for auto-redeployment
- **Add features**: Rate limiting, player profiles, historical charts

For detailed setup, see [DATASETTE_SETUP.md](DATASETTE_SETUP.md)
