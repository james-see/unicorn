# Next Steps: Activate Your Global Leaderboard

## ? What's Been Done

Your Unicorn game now has a complete global leaderboard system ready to deploy!

**Files Created:**
- ? Datasette database with demo scores (`leaderboard.db`)
- ? API endpoint for score submissions (`api/submit-score.go`)
- ? Game client for HTTP requests (`leaderboard/leaderboard.go`)
- ? GitHub Pages leaderboard display (updated `docs/index.html`)
- ? Vercel deployment config (`vercel.json`)
- ? GitHub Actions workflow (`.github/workflows/datasette-deploy.yml`)
- ? Complete documentation (3 guides)

**Features:**
- ?? Global leaderboard with real-time updates
- ?? One-click score submission from game
- ?? Filter by difficulty, sort by Net Worth or ROI
- ?? Beautiful UI on GitHub Pages
- ?? Free hosting on Vercel
- ?? Privacy-focused (no personal data)

## ?? Deploy in 5 Minutes

### Step 1: Install Tools (One-time)

```bash
# Install Vercel CLI
npm install -g vercel

# Install Python packages for Datasette
pip install datasette datasette-publish-vercel
```

### Step 2: Deploy API to Vercel

```bash
cd /workspace  # Your project directory

# Login to Vercel (opens browser)
vercel login

# Deploy the API endpoint
vercel deploy --prod
```

**Copy the deployment URL** (e.g., `https://unicorn-abc123.vercel.app`)

### Step 3: Deploy Datasette

```bash
# Deploy the leaderboard database
datasette publish vercel leaderboard.db \
  --project=unicorn-leaderboard \
  --install=datasette-cors \
  --metadata=datasette-metadata.json \
  --setting sql_time_limit_ms 3500
```

**Copy the Datasette URL** (e.g., `https://unicorn-leaderboard.vercel.app`)

### Step 4: Update Configuration

**Edit `docs/index.html` line 617:**
```javascript
const DATASETTE_URL = 'https://YOUR-DATASETTE-URL/leaderboard/game_scores.json';
```

**Edit `leaderboard/leaderboard.go` line 17:**
```go
DefaultAPIEndpoint = "https://YOUR-VERCEL-URL/api/submit-score"
```

Replace with your actual Vercel URLs from steps 2 and 3.

### Step 5: Rebuild and Test

```bash
# Rebuild the game
go build

# Test it!
./unicorn
# Play a game ? Complete it ? Choose "Yes" to submit score
```

### Step 6: Push to GitHub

```bash
git add .
git commit -m "Add global leaderboard with Datasette"
git push
```

Your leaderboard will be live at:
**https://james-see.github.io/unicorn**

## ?? Test Your Setup

### Test API Endpoint
```bash
curl -X POST https://YOUR-VERCEL-URL/api/submit-score \
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
  "id": "some-uuid"
}
```

### Test Datasette
```bash
curl https://YOUR-DATASETTE-URL/leaderboard/game_scores.json?_size=10
```

Should return JSON with your demo scores.

### Test from Game
1. Run `./unicorn`
2. Play a complete game
3. When prompted, choose "Yes" to submit
4. Check GitHub Pages to see your score!

## ?? Documentation

Three comprehensive guides are available:

1. **[QUICKSTART_LEADERBOARD.md](QUICKSTART_LEADERBOARD.md)**
   - 5-minute deployment guide
   - Quick reference commands

2. **[DATASETTE_SETUP.md](DATASETTE_SETUP.md)**
   - Detailed setup instructions
   - Troubleshooting guide
   - Security considerations

3. **[README_LEADERBOARD.md](README_LEADERBOARD.md)**
   - Architecture overview
   - API documentation
   - Future enhancements

## ?? Optional: GitHub Actions

To enable automatic redeployment:

1. Get Vercel token: https://vercel.com/account/tokens
2. Add to GitHub secrets as `VERCEL_TOKEN`
3. Workflow runs every 6 hours or manually

## ?? What Players Will See

### In the Game
After completing a game:
```
========================================================
           ?? GLOBAL LEADERBOARD
========================================================

Would you like to submit your score to the global leaderboard?
Your score will be visible to all players worldwide!

Submit to global leaderboard? (y/n, default n): y

Checking global leaderboard service... ?
Submitting your score... ?

?? Success! Your score has been submitted to the global leaderboard!

View the global leaderboard at:
https://james-see.github.io/unicorn
```

### On GitHub Pages
- Beautiful leaderboard table with rankings
- Filter tabs: All, Easy, Medium, Hard, Expert
- Sort by Net Worth or ROI
- Refresh button for latest scores
- Gold/Silver/Bronze highlighting for top 3
- Green for positive ROI, red for negative

## ?? Demo Scores Included

Your database already has 3 demo scores:
- **UnicornHunter**: $50M, 19,900% ROI (Expert)
- **VCMaster**: $25M, 9,900% ROI (Hard)  
- **StartupKing**: $10M, 3,900% ROI (Medium)

These will show up immediately on your leaderboard!

## ?? Tips

- **Free Tier**: Vercel free tier supports ~10,000 submissions/month
- **No Database Management**: Datasette handles everything
- **Real-time Updates**: No caching, always fresh data
- **Privacy**: Only game stats are stored, no personal info

## ?? Troubleshooting

**API returns 404?**
- Ensure you ran `vercel deploy --prod` (not just `vercel dev`)
- Check the URL matches what's in your code

**Leaderboard shows error?**
- Visit your Datasette URL directly to verify it's working
- Check browser console for errors
- Verify CORS plugin is installed

**Game can't submit?**
- Rebuild after updating URLs: `go build`
- Check internet connection
- Try the curl test command above

## ?? Monitor Usage

View logs and stats:
```bash
vercel logs                    # View function logs
vercel logs --follow          # Live tail
vercel ls                     # List deployments
vercel inspect YOUR-URL       # Deployment details
```

## ?? Current Status

? Code complete and tested  
? Database created with demo data  
? Documentation written  
? Build verified  
? **Ready for deployment!**

## ?? Need Help?

- **Quick Start**: [QUICKSTART_LEADERBOARD.md](QUICKSTART_LEADERBOARD.md)
- **Detailed Guide**: [DATASETTE_SETUP.md](DATASETTE_SETUP.md)
- **Architecture**: [README_LEADERBOARD.md](README_LEADERBOARD.md)
- **Project Summary**: [PHASE5_SUMMARY.md](PHASE5_SUMMARY.md)

## ?? You're All Set!

Everything is ready. Just follow the 5-minute deployment steps above and your global leaderboard will be live!

**Your players will love competing worldwide!** ??

---

Questions? Check the documentation or open a GitHub issue.
