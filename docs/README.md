# Unicorn documentation

Docs are grouped by topic:

- **[ROADMAP.md](ROADMAP.md)** – Feature roadmap and priorities  
- **[QUICKSTART.md](QUICKSTART.md)** – Get the game running quickly  

## Datasette / leaderboard (Python)

Leaderboard and deployment docs live in **[datasette/](datasette/)**:

- Setup: [DATASETTE_SETUP.md](datasette/DATASETTE_SETUP.md)
- Quick deploy: [QUICKSTART_LEADERBOARD.md](datasette/QUICKSTART_LEADERBOARD.md)
- Architecture: [README_LEADERBOARD.md](datasette/README_LEADERBOARD.md)

The leaderboard runs as a small Python app (Datasette on Vercel). Dependencies are in **`pyproject.toml`** and **`requirements.txt`** at repo root; `requirements.txt` is used by Vercel. There is no need for a separate Pipfile.

## Development

**[development/](development/)** – Phase summaries, refactoring notes, integration and next steps.

## Features

**[features/](features/)** – Founder mode guides, reputation system, animations, upgrade proposals, and related feature docs.
