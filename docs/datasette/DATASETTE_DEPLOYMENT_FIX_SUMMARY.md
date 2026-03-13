# Datasette Deployment Fix - Summary

## Issue Identified

The Datasette deployment keeps failing because the `POSTGRES_URL` environment variable is not set in the Vercel project for the datasette deployment.

## What Was Wrong

The `index.py` file expects `POSTGRES_URL` to be available, but you have two separate Vercel projects:
1. **Main API** (`unicorn-green.vercel.app`) - has `POSTGRES_URL` set ✅
2. **Datasette deployment** - missing `POSTGRES_URL` ❌

## Fix Applied

### 1. Updated `index.py` to properly connect to Postgres
- Connects to Postgres database using `POSTGRES_URL` environment variable  
- Provides clear error message if variable is missing
- Loads metadata from `datasette-metadata.json`

### 2. Updated `requirements.txt`
```
datasette>=0.65.0
datasette-connectors
psycopg2-binary
```

### 3. Added `runtime.txt`
```
python-3.12
```

## ⚠️ CRITICAL: Action Required

**You MUST set the environment variable in Vercel for this to work:**

1. Go to [Vercel Dashboard](https://vercel.com/dashboard)
2. Find your Datasette project (separate from main `unicorn` project)
3. Go to **Settings** → **Environment Variables**
4. Add `POSTGRES_URL` with the same value from your main project
5. Add for all environments: Production, Preview, Development
6. Redeploy

## How to Get the POSTGRES_URL

From your main `unicorn` project:
```bash
vercel env pull .env.local
cat .env.local | grep POSTGRES_URL
```

Or from Vercel Dashboard:
- Main project (`unicorn`) → Settings → Environment Variables → `POSTGRES_URL`

## Test After Deployment

```bash
# Should return JSON with leaderboard data
curl "https://YOUR-DATASETTE-URL.vercel.app/leaderboard/game_scores.json?_size=5"
```

## Alternative: Use Existing API (Recommended)

**You already have a working API!**

Your main deployment at `unicorn-green.vercel.app` already provides:
- `/api/get-leaderboard` - Datasette-compatible JSON
- `/leaderboard/game_scores.json` - Same endpoint with friendly URL  
- Full filtering, sorting, pagination support
- Connected to same Postgres database

**Recommendation:** Unless you need the Datasette web UI for data exploration, just use the existing Go API. It's simpler and already working.

## Files Modified

✅ `/workspace/datasette-deploy/index.py` - Postgres connection
✅ `/workspace/datasette-deploy/requirements.txt` - Dependencies
✅ `/workspace/datasette-deploy/pyproject.toml` - Project config
✅ `/workspace/datasette-deploy/runtime.txt` - Python version

## What Happens Without the Fix

Without `POSTGRES_URL` set:
```
ValueError: POSTGRES_URL environment variable is required. Set it in Vercel dashboard.
```

This is why your deployments are failing.

## Deploy Commands

```bash
cd /workspace/datasette-deploy
vercel --prod
```

## Summary

✅ **Code fixed** - Proper Postgres connection configured
❌ **Environment variable needed** - You must set `POSTGRES_URL` in Vercel
🤔 **Consider alternatives** - Your existing API might be sufficient
