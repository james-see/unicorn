import asyncio
from datasette.app import Datasette
import json
import pathlib
import os

static_mounts = [
    (static, str((pathlib.Path(".") / static).resolve()))
    for static in []
]

metadata = dict()
try:
    metadata_path = os.path.join(os.path.dirname(__file__), 'datasette-metadata.json')
    metadata = json.load(open(metadata_path))
except Exception:
    pass

secret = os.environ.get("DATASETTE_SECRET")

# Connect to Vercel Postgres database
postgres_url = os.environ.get("POSTGRES_URL")
if not postgres_url:
    raise ValueError("POSTGRES_URL environment variable is required")

# Datasette can connect to Postgres databases using connection strings
# Format: postgres://user:password@host:port/database
ds = Datasette(
    [],  # No SQLite files
    [postgres_url],  # Postgres connection string
    static_mounts=static_mounts,
    metadata=metadata,
    secret=secret,
    cors=True,
    settings={}
)
asyncio.run(ds.invoke_startup())
app = ds.app()
