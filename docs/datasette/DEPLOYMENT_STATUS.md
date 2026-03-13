# Deployment Status & Next Steps

## ✅ Completed

1. **API Deployed to Vercel**
   - URL: `https://unicorn-g5xp1fzsv-james-campbells-projects-98ba50e1.vercel.app/api/submit-score`
   - Status: ✅ Deployed and ready
   - Note: Currently has authentication protection enabled

2. **Datasette Deployed to Vercel**
   - URL: `https://unicorn-leaderboard-pgima6xn0-james-campbells-projects-98ba50e1.vercel.app`
   - Status: ✅ Deployed (build had warning but deployment succeeded)
   - Note: Currently has authentication protection enabled

3. **Code Updated**
   - ✅ `leaderboard/leaderboard.go` - Updated with API URL
   - ✅ `docs/index.html` - Updated with Datasette URL
   - ✅ Game binary rebuilt with new URLs

## ⚠️ Important: Authentication Protection

Both deployments currently have **Vercel Authentication Protection** enabled, which means:
- The URLs require authentication to access
- Public users cannot access the API or Datasette endpoints
- This needs to be disabled for public access

### How to Disable Protection:

1. Go to [Vercel Dashboard](https://vercel.com/dashboard)
2. Select your project:
   - **API Project**: `unicorn`
   - **Datasette Project**: `unicorn-leaderboard`
3. Go to **Settings** → **Deployment Protection**
4. Disable protection or configure it for public access
5. Redeploy if needed

### Alternative: Use Production Domains

Vercel assigns production domains automatically. Check your Vercel dashboard for:
- Production URL for `unicorn` project
- Production URL for `unicorn-leaderboard` project

Then update the URLs in:
- `leaderboard/leaderboard.go` (line 16)
- `docs/index.html` (line 619)

## 📋 Next Steps

### Immediate Actions Needed:

1. **Disable Authentication Protection** (see above)
   - Or wait for production domains and update URLs

2. **Test the API** (once protection is disabled):
   ```bash
   curl -X POST https://YOUR-API-URL/api/submit-score \
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

3. **Test Datasette** (once protection is disabled):
   ```bash
   curl https://YOUR-DATASETTE-URL/leaderboard/game_scores.json?_size=3
   ```

4. **Test from Game**:
   ```bash
   ./unicorn
   # Play a game and submit a score
   ```

### Optional: Set Up Production Domains

1. In Vercel Dashboard, go to project settings
2. Add custom domain or use default production domain
3. Update URLs in code and rebuild:
   ```bash
   go build
   ```

### Optional: GitHub Actions Auto-Deployment

1. Get Vercel token: https://vercel.com/account/tokens
2. Add to GitHub secrets as `VERCEL_TOKEN`
3. The workflow at `.github/workflows/datasette-deploy.yml` will auto-deploy every 6 hours

## 🔗 Deployment URLs

**API Endpoint:**
- Preview: `https://unicorn-g5xp1fzsv-james-campbells-projects-98ba50e1.vercel.app/api/submit-score`
- Production: Check Vercel dashboard

**Datasette:**
- Preview: `https://unicorn-leaderboard-pgima6xn0-james-campbells-projects-98ba50e1.vercel.app`
- Production: Check Vercel dashboard

## 📝 Notes

- Both deployments are using preview URLs (they change with each deployment)
- For production use, set up production domains or use the stable production URLs
- The Datasette build showed a warning about `pip3.9` but deployment succeeded
- All code is updated and game is rebuilt with the new URLs

## ✨ Once Protection is Disabled

Your leaderboard will be fully functional:
- Players can submit scores from the game
- Leaderboard displays on GitHub Pages
- Real-time updates from Datasette
- Filter by difficulty, sort by Net Worth or ROI

