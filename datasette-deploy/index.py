import asyncio
from datasette.app import Datasette
import json
import os

# Load metadata
metadata = dict()
try:
    metadata_path = os.path.join(os.path.dirname(__file__), 'datasette-metadata.json')
    with open(metadata_path) as f:
        metadata = json.load(f)
except Exception:
    metadata = {"title": "Unicorn Leaderboard"}

secret = os.environ.get("DATASETTE_SECRET", "default-secret-change-in-production")

# Connect to Vercel Postgres database
postgres_url = os.environ.get("POSTGRES_URL")
if not postgres_url:
    raise ValueError("POSTGRES_URL environment variable is required. Set it in Vercel dashboard.")

# Datasette supports Postgres via connection strings
# The connection string should be in format: postgresql://user:pass@host:port/db
ds = Datasette(
    files=[],  # No SQLite files
    immutables=[],
    databases={"leaderboard": {"url": postgres_url}},
    metadata=metadata,
    secret=secret,
    cors=True,
    settings={
        "sql_time_limit_ms": 3500,
        "allow_download": False
    }
)

# Initialize datasette
asyncio.run(ds.invoke_startup())
app = ds.app()
