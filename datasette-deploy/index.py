import asyncio
from datasette.app import Datasette
import json
import os

metadata = dict()
try:
    metadata_path = os.path.join(os.path.dirname(__file__), 'datasette-metadata.json')
    metadata = json.load(open(metadata_path))
except Exception:
    pass

ds = Datasette(
    [os.path.join(os.path.dirname(__file__), 'leaderboard.db')],
    metadata=metadata,
    cors=True,
    sql_time_limit_ms=3500
)

asyncio.run(ds.invoke_startup())
app = ds.app()

# For Vercel, we need to use mangum adapter
from mangum import Mangum
handler = Mangum(app, lifespan="off")
