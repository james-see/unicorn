# Datasette Deployment Fix - 2025-11-11

## Problem

The Datasette deployment to Vercel was failing throughout the day with the following issues:

1. **Missing Environment Variable**: The `index.py` file was expecting a `POSTGRES_URL` environment variable that was not configured in Vercel
2. **Wrong Database Type**: Code was configured to use PostgreSQL, but a SQLite database (`leaderboard.db`) was present in the deployment directory
3. **Unnecessary Dependencies**: Requirements included PostgreSQL-specific packages (`psycopg2-binary`) that weren't needed

## Root Cause

The deployment configuration was attempting to connect to a PostgreSQL database via environment variable, but:
- No PostgreSQL database was provisioned in Vercel
- No `POSTGRES_URL` environment variable was set
- A SQLite database file already existed and contained the data

This caused the deployment to fail immediately on startup with:
```
ValueError: POSTGRES_URL environment variable is required
```

## Solution

Changed the deployment to use SQLite database (which is the standard and recommended approach for Datasette on Vercel):

### 1. Updated `index.py`
**Before:**
```python
# Connect to Vercel Postgres database
postgres_url = os.environ.get("POSTGRES_URL")
if not postgres_url:
    raise ValueError("POSTGRES_URL environment variable is required")

ds = Datasette(
    [],  # No SQLite files
    [postgres_url],  # Postgres connection string
    ...
)
```

**After:**
```python
# Use SQLite database included in deployment
db_path = os.path.join(os.path.dirname(__file__), 'leaderboard.db')

ds = Datasette(
    [db_path],  # SQLite database file
    static_mounts=static_mounts,
    metadata=metadata,
    secret=secret,
    cors=True,
    settings={
        "sql_time_limit_ms": 3500,
        "allow_download": False
    }
)
```

### 2. Cleaned up `requirements.txt`
Removed unnecessary PostgreSQL dependencies:
```
datasette>=0.65.0
datasette-cors>=0.2.0
```

### 3. Updated `pyproject.toml`
Removed PostgreSQL dependencies from project dependencies.

### 4. Added `runtime.txt`
Ensured Python 3.12 is used:
```
python-3.12
```

## Benefits of SQLite Approach

1. **Simpler**: No need for external database configuration
2. **Faster**: Database is included in deployment bundle
3. **Cost-effective**: No external database costs
4. **Standard**: This is how Datasette is typically deployed on Vercel
5. **Reliable**: No network dependencies or connection issues

## Verification

- ✅ Python syntax validated
- ✅ Database schema confirmed (3 rows present)
- ✅ All required files in place
- ✅ Dependencies simplified

## Next Steps

1. Deploy to Vercel: `vercel deploy --prod` from the `datasette-deploy` directory
2. Test the deployment works
3. Update any documentation referencing PostgreSQL

## Note on Data Updates

Since SQLite database is bundled with deployment:
- To update leaderboard data, you'll need to:
  1. Update the `leaderboard.db` file locally
  2. Redeploy to Vercel
- Or: Keep the API endpoint (`submit-score`) which should update a persistent database
- Consider: Setting up a persistent database solution if real-time updates from the game are needed

## Alternative: Persistent Storage

If you need real-time score submissions without redeployment:
1. Use Vercel Postgres with the API endpoints
2. Keep Datasette deployment separate (read-only from static SQLite snapshots)
3. Periodically sync Postgres → SQLite for Datasette views
