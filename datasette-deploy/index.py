from datasette.app import Datasette
import json
import os
from mangum import Mangum

metadata = dict()
try:
    metadata_path = os.path.join(os.path.dirname(__file__), 'datasette-metadata.json')
    metadata = json.load(open(metadata_path))
except Exception:
    pass

ds = Datasette(
    [os.path.join(os.path.dirname(__file__), 'leaderboard.db')],
    metadata=metadata,
    cors=True
)

app = ds.app()

# Mangum handles async initialization via lifespan events
handler = Mangum(app, lifespan="auto")
