from datasette.app import Datasette
import json
import os
from mangum import Mangum

# Load metadata
metadata_path = os.path.join(os.path.dirname(__file__), 'datasette-metadata.json')
with open(metadata_path) as f:
    metadata = json.load(f)

# Create Datasette instance
db_path = os.path.join(os.path.dirname(__file__), 'leaderboard.db')
ds = Datasette(
    [db_path],
    metadata=metadata,
    cors=True,
    sql_time_limit_ms=3500
)

# Wrap with Mangum for Vercel compatibility
handler = Mangum(ds.app())
