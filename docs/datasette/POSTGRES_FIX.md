# Postgres Driver Fix - Resolution

## Problem

The Postgres driver was not being recognized in Vercel deployment, causing this error:
```
Database connection error: sql: unknown driver "postgres" (forgotten import?)
```

**Symptoms:**
- Local builds worked fine
- `submit-score` API worked on Vercel
- `get-leaderboard` API failed with driver error
- Both had identical code and imports: `_ "github.com/lib/pq"`

## Root Cause

**Vercel's Go function builder had issues with subdirectory structure:**
- Structure: `api/submit-score/main.go` and `api/get-leaderboard/main.go`
- Vercel generates wrapper code (`main__vc__go__.go`) during build
- The wrapper wasn't properly linking the Postgres driver for subdirectories
- Build logs showed: `Error: Command failed: go build -ldflags -s -w -o /tmp/.../bootstrap`

## Solution

**Moved to single-file structure (Vercel's recommended format):**

**Before:**
```
api/
├── submit-score/
│   ├── main.go
│   ├── go.mod
│   └── go.sum
└── get-leaderboard/
    ├── main.go
    ├── go.mod
    └── go.sum
```

**After:**
```
api/
├── submit-score.go
├── get-leaderboard.go
├── go.mod
└── go.sum
```

**Key Changes:**
1. Consolidated `go.mod` at `api/` level (not per-function)
2. Single `.go` files for each endpoint
3. Both use `package handler` with `func Handler(w, r)`
4. Postgres driver import works in both

## Verification

### Submit Score API ✅
```bash
curl -X POST https://unicorn-green.vercel.app/api/submit-score \
  -H "Content-Type: application/json" \
  -d '{"player_name":"Test","final_net_worth":12000000,"roi":2500.0,"successful_exits":9,"turns_played":120,"difficulty":"Expert"}'

Response: {"success":true,"message":"Score submitted successfully!","id":"89a9e387-9a00-40f2-baf7-a6cc4e3f880c"}
```

### Get Leaderboard API ✅
```bash
curl "https://unicorn-green.vercel.app/leaderboard/game_scores.json?_size=5&_sort_desc=final_net_worth"

Response: {
  "rows": [
    ["id1", "FixTest", 12000000, 2500.0, 9, 120, "Expert", "2025-11-01T..."],
    ["id2", "DirectTest1", 10500000, 2300.0, 7, 120, "Expert", "2025-11-01T..."],
    ...
  ],
  "filtered_table_rows_count": 5,
  "database": "leaderboard",
  "table": "game_scores",
  "columns": ["id", "player_name", "final_net_worth", "roi", "successful_exits", "turns_played", "difficulty", "played_at"]
}
```

### Difficulty Filtering ✅
```bash
curl "https://unicorn-green.vercel.app/api/get-leaderboard?_size=10&difficulty=Expert"

Response: Successfully returns only Expert difficulty scores
```

## Deployment URLs

- **Production:** https://unicorn-green.vercel.app
- **Submit Score:** https://unicorn-green.vercel.app/api/submit-score
- **Get Leaderboard:** https://unicorn-green.vercel.app/api/get-leaderboard
- **Datasette Route:** https://unicorn-green.vercel.app/leaderboard/game_scores.json

## Lessons Learned

1. **Vercel Go functions prefer single files** in `api/` directory
2. **Subdirectories can work** but may cause build wrapper issues
3. **Always test both endpoints** when using shared dependencies
4. **Check build logs** for wrapper generation errors
5. **Single go.mod at api/ level** is simpler than per-function modules

## References

- [Vercel Go Runtime](https://vercel.com/docs/functions/runtimes/go)
- [lib/pq Postgres Driver](https://github.com/lib/pq)
- Vercel expects: `api/filename.go` with `package handler` and `func Handler(w, r)`

## Date

Fixed: November 1, 2025
Deployed: https://unicorn-green.vercel.app
Status: ✅ All endpoints operational

