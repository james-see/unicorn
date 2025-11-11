import asyncio
from datasette.app import Datasette
import json
import pathlib
import os
import logging

# Set up logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

static_mounts = [
    (static, str((pathlib.Path(".") / static).resolve()))
    for static in []
]

metadata = dict()
try:
    metadata_path = os.path.join(os.path.dirname(__file__), 'datasette-metadata.json')
    if os.path.exists(metadata_path):
        metadata = json.load(open(metadata_path))
        logger.info(f"Loaded metadata from {metadata_path}")
except Exception as e:
    logger.warning(f"Could not load metadata: {e}")

secret = os.environ.get("DATASETTE_SECRET")

# Determine database configuration
postgres_url = os.environ.get("POSTGRES_URL")
sqlite_db_path = os.path.join(os.path.dirname(__file__), 'leaderboard.db')

# Initialize Datasette with appropriate database
# Note: For Postgres support, Datasette requires the connection string to be passed
# via environment variables or plugins. For now, we'll use SQLite as the primary
# database source, with Postgres support via the API endpoints.
if postgres_url:
    logger.info("POSTGRES_URL detected, but Datasette will use SQLite for read-only access.")
    logger.info("Postgres writes should go through the API endpoints.")
    # Fall through to SQLite for Datasette read-only access
    postgres_url = None  # Don't use Postgres directly in Datasette

if os.path.exists(sqlite_db_path):
    logger.info(f"Using SQLite database from {sqlite_db_path}")
    ds = Datasette(
        [sqlite_db_path],
        static_mounts=static_mounts,
        metadata=metadata,
        secret=secret,
        cors=True,
        settings={
            "sql_time_limit_ms": 3500,
            "allow_download": False,
        }
    )
else:
    logger.warning("No database found. Creating empty SQLite database.")
    # Create empty database if neither exists
    import sqlite3
    conn = sqlite3.connect(sqlite_db_path)
    # Create table structure
    conn.execute("""
        CREATE TABLE IF NOT EXISTS game_scores (
            id TEXT PRIMARY KEY,
            player_name TEXT NOT NULL,
            final_net_worth INTEGER NOT NULL,
            roi REAL NOT NULL,
            successful_exits INTEGER NOT NULL,
            turns_played INTEGER NOT NULL,
            difficulty TEXT NOT NULL,
            played_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    """)
    conn.execute("CREATE INDEX IF NOT EXISTS idx_net_worth ON game_scores(final_net_worth DESC)")
    conn.execute("CREATE INDEX IF NOT EXISTS idx_roi ON game_scores(roi DESC)")
    conn.execute("CREATE INDEX IF NOT EXISTS idx_player ON game_scores(player_name)")
    conn.execute("CREATE INDEX IF NOT EXISTS idx_difficulty ON game_scores(difficulty)")
    conn.commit()
    conn.close()
    logger.info(f"Created empty database at {sqlite_db_path}")
    ds = Datasette(
        [sqlite_db_path],
        static_mounts=static_mounts,
        metadata=metadata,
        secret=secret,
        cors=True,
        settings={
            "sql_time_limit_ms": 3500,
            "allow_download": False,
        }
    )

asyncio.run(ds.invoke_startup())
app = ds.app()
