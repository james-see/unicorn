import asyncio
from datasette.app import Datasette
import json
import pathlib
import os

metadata = dict()
try:
    metadata = json.load(open("datasette-metadata.json"))
except Exception:
    pass

ds = Datasette(
    ["leaderboard.db"],
    metadata=metadata,
    cors=True
)

asyncio.run(ds.invoke_startup())
app = ds.app()

# For Vercel, we need to use mangum adapter
from mangum import Mangum
handler = Mangum(app, lifespan="off")
