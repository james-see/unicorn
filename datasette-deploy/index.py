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
asyncio.run(ds.invoke_startup())
app = ds.app()
