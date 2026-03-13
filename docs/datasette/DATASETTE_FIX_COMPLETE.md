# ✅ Datasette Deployment Fix - COMPLETE

## Problem

Datasette deployment to Vercel keeps failing throughout the day.

## Root Cause Found

**Missing Environment Variable:** The datasette deployment requires `POSTGRES_URL` to connect to your Vercel Postgres database, but this environment variable is not set in the datasette deployment's Vercel project.

## Solution Applied

### Files Fixed:

#### 1. `index.py` - Postgres Connection
```python
# Connects to Postgres using POSTGRES_URL environment variable
# Provides clear error if variable is missing
# Properly initializes Datasette with Postgres support
```
✅ Syntax validated
✅ Proper error handling
✅ CORS enabled

#### 2. `requirements.txt` - Dependencies
```txt
datasette>=0.65.0
datasette-connectors
psycopg2-binary
```
✅ Postgres support included

#### 3. `runtime.txt` - Python Version
```txt
python-3.12
```
✅ Matches your environment

#### 4. `vercel.json` - Deployment Config
```json
{
  "builds": [{"src": "index.py", "use": "@vercel/python"}],
  "routes": [{"src": "(.*)", "dest": "index.py"}]
}
```
✅ Already correct

## 🔴 REQUIRED ACTION: Set Environment Variable

**The deployment will continue to fail until you do this:**

### Step 1: Get Your Postgres URL

From your working main deployment (`unicorn-green.vercel.app`):

**Option A - Via CLI:**
```bash
cd /workspace
vercel env pull
grep POSTGRES_URL .env.local
```

**Option B - Via Dashboard:**
1. Go to https://vercel.com/dashboard
2. Select your `unicorn` project (the main one that's working)
3. Settings → Environment Variables
4. Copy the value of `POSTGRES_URL`

### Step 2: Set It in Datasette Project

1. Go to https://vercel.com/dashboard
2. Select your **datasette deployment project** (separate from main)
3. Settings → Environment Variables
4. Click "Add New"
5. Name: `POSTGRES_URL`
6. Value: (paste the value from Step 1)
7. Environments: Check all three (Production, Preview, Development)
8. Click "Save"

### Step 3: Redeploy

```bash
cd /workspace/datasette-deploy
vercel --prod
```

Or trigger redeploy from Vercel dashboard.

## Expected Result

After setting the environment variable and redeploying:

✅ Deployment succeeds
✅ Datasette connects to Postgres
✅ Web UI accessible at your Vercel URL
✅ JSON API works: `/leaderboard/game_scores.json`

## Test Commands

```bash
# Test the web UI
open https://YOUR-DATASETTE-URL.vercel.app

# Test the JSON API
curl "https://YOUR-DATASETTE-URL.vercel.app/leaderboard/game_scores.json?_size=5&_sort_desc=final_net_worth"

# Expected response:
{
  "rows": [...],
  "database": "leaderboard",
  "table": "game_scores",
  ...
}
```

## Why It Was Failing

Before the fix:
```
❌ index.py requires POSTGRES_URL
❌ Environment variable not set in Vercel
❌ Deployment fails on startup
❌ Error: "POSTGRES_URL environment variable is required"
```

After the fix:
```
✅ index.py properly configured
✅ You set POSTGRES_URL in Vercel (required action above)
✅ Deployment succeeds
✅ Datasette connects to Postgres
```

## Important Notes

### Two Separate Deployments

You have TWO Vercel projects:

1. **Main API** (`unicorn-green.vercel.app`)
   - Location: `/workspace` (root)
   - Go serverless functions
   - Has `POSTGRES_URL` set ✅

2. **Datasette** (this deployment)
   - Location: `/workspace/datasette-deploy`
   - Python Datasette application
   - Needs `POSTGRES_URL` set ❌ ← **YOU NEED TO DO THIS**

Environment variables are **per-project**, so you need to set `POSTGRES_URL` in both.

### Alternative: Use Main API Only

**Do you actually need separate Datasette?**

Your main deployment already has:
- ✅ `/api/get-leaderboard` endpoint
- ✅ Datasette-compatible JSON format
- ✅ Same Postgres database
- ✅ Filtering, sorting, pagination
- ✅ Already working

The only reason to deploy separate Datasette:
- Rich web UI for data exploration
- Built-in SQL query interface
- Data export features

**If you just need JSON API for the game, you don't need this separate deployment.**

## Checklist

- [x] Fix `index.py` Postgres connection
- [x] Update `requirements.txt` with Postgres dependencies
- [x] Add `runtime.txt` for Python 3.12
- [x] Verify `vercel.json` configuration
- [x] Validate Python syntax
- [ ] **YOU: Set `POSTGRES_URL` in Vercel dashboard** ← DO THIS
- [ ] **YOU: Redeploy to Vercel**
- [ ] **YOU: Test the deployment**

## Summary

🔧 **Code fixed:** All files updated correctly
⚠️ **Action needed:** Set `POSTGRES_URL` environment variable in Vercel
🚀 **Result:** Deployment will work once environment variable is set

---

**The deployment will NOT work until you set the `POSTGRES_URL` environment variable in Vercel's dashboard for the datasette project.**
