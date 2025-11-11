# Datasette Postgres Deployment Fix - 2025-11-11

## Problem

The Datasette deployment in `/datasette-deploy` was failing because:

1. The `index.py` file requires `POSTGRES_URL` environment variable
2. The datasette-deploy is a separate Vercel project from the main API
3. The `POSTGRES_URL` environment variable needs to be set in the Vercel dashboard for the datasette deployment

## Context

You have TWO Vercel deployments:

### 1. Main API Deployment (Working ✅)
- **URL:** `https://unicorn-green.vercel.app`
- **Location:** Root `/workspace` directory
- **Features:**
  - Go serverless functions at `/api/submit-score` and `/api/get-leaderboard`  
  - Uses Postgres via `POSTGRES_URL` environment variable
  - Already provides Datasette-compatible JSON at `/leaderboard/game_scores.json`
- **Status:** Working perfectly

### 2. Datasette Deployment (Failing ❌)
- **Location:** `/workspace/datasette-deploy` directory
- **Purpose:** Actual Python Datasette instance for richer data exploration
- **Issue:** Missing `POSTGRES_URL` environment variable in this deployment's Vercel config

## Root Cause

The datasette-deploy directory is configured to connect to Postgres, but when deploying to Vercel as a separate project, the `POSTGRES_URL` environment variable from your main deployment is NOT automatically shared.

## Solution

You need to set the `POSTGRES_URL` environment variable in the Vercel dashboard for the datasette deployment:

### Steps to Fix:

1. **Get your Postgres URL** from the main deployment:
   - Go to [Vercel Dashboard](https://vercel.com/dashboard)
   - Select your main project (`unicorn`)
   - Go to **Settings** → **Environment Variables**
   - Copy the `POSTGRES_URL` value

2. **Set it in the Datasette deployment**:
   - Go to your Datasette project in Vercel (separate project)
   - Go to **Settings** → **Environment Variables**  
   - Add: `POSTGRES_URL` with the same value from step 1
   - Make sure to add it for all environments (Production, Preview, Development)

3. **Redeploy**:
   ```bash
   cd /workspace/datasette-deploy
   vercel --prod
   ```

## Updated Configuration

I've updated the files to properly connect to Postgres:

### `index.py`
- Connects to Postgres via `POSTGRES_URL` environment variable
- Uses `datasette-connectors` for Postgres support
- Loads metadata from `datasette-metadata.json`
- Enables CORS for API access

### `requirements.txt`
```
datasette>=0.65.0
datasette-connectors
psycopg2-binary
```

### `runtime.txt`
```
python-3.12
```

## Alternative: Use Existing API

**Do you actually need a separate Datasette deployment?**

Your main API deployment already provides:
- ✅ Datasette-compatible JSON API
- ✅ Filtering by difficulty
- ✅ Sorting by any column
- ✅ Pagination with `_size` parameter
- ✅ Uses the same Postgres database

The only advantage of actual Datasette is:
- Rich web UI for data exploration
- Built-in data export features
- SQL query interface

**Recommendation:** If you just need the API for your game leaderboard, stick with the existing Go API at `unicorn-green.vercel.app`. It's simpler, faster, and already working.

## Verification After Fix

Once you set the environment variable and redeploy, test:

```bash
# Test the Datasette deployment
curl "https://YOUR-DATASETTE-URL.vercel.app/leaderboard/game_scores.json?_size=5"
```

## Files Modified

1. `/workspace/datasette-deploy/index.py` - Updated Postgres connection
2. `/workspace/datasette-deploy/requirements.txt` - Added Postgres support
3. `/workspace/datasette-deploy/pyproject.toml` - Updated dependencies
4. `/workspace/datasette-deploy/runtime.txt` - Set Python 3.12

## Next Steps

**Choose one option:**

### Option A: Fix Datasette Deployment (for rich data exploration UI)
1. Set `POSTGRES_URL` in Vercel dashboard for datasette project
2. Redeploy datasette
3. Access full Datasette web UI

### Option B: Remove Datasette Deployment (simpler)
1. Delete `/workspace/datasette-deploy` directory
2. Use only the main API at `unicorn-green.vercel.app`
3. The Go API already provides all the JSON data you need

I recommend **Option B** unless you specifically want the Datasette web UI for data exploration.
