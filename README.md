# maeve

Memory and agent evolution engine for coding agents.

`maeve` tracks session context, scores what matters, compresses it to a token
budget, and saves named snapshots that can be restored or diffed later.

## Quick start

```bash
go install ./cmd/maeve
maeve init
maeve ingest file PRD.md
maeve status
maeve compress --budget 8000
maeve snapshot save "before refactor"
```

## Current scope

This repository is building the v0.1 slice first:

- local CLI
- SQLite store
- context object ingestion
- importance scoring
- token budget compression
- snapshot creation and listing

HTTP, MCP, tree-sitter parsing, and the VS Code extension are laid out as
follow-on packages.

