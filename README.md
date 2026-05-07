# maeve

Memory and agent evolution engine for coding agents.

`maeve` tracks session context, scores what matters, compresses it to a token
budget, and saves named snapshots that can be restored or diffed later. The
repo includes a Go CLI plus a VS Code extension that exposes the same workflow
inside the editor.

## Quick start

```bash
go install ./cmd/maeve
maeve init
maeve ingest file PRD.md
maeve status
maeve compress --budget 8000
maeve snapshot save "before refactor"
```

## VS Code extension

The extension adds a `maeve` activity bar view with Context Lens and Snapshots
panels, plus commands for the common workflow:

- initialize a session
- ingest the active file or selected Explorer files
- ingest the current Git diff
- compress context and copy it to the clipboard
- save and list snapshots
- build a workspace-local CLI at `.maeve/bin/maeve`

To run it locally:

```bash
npm install
npm run compile
```

Open this folder in VS Code, press `F5`, and use the `maeve` commands from the
Command Palette. If the `maeve` binary is not installed globally, run
`maeve: Build Workspace CLI` first.

## Current scope

This repository is building the v0.1 slice first:

- local CLI
- VS Code extension
- SQLite store
- context object ingestion
- importance scoring
- token budget compression
- snapshot creation and listing

HTTP, MCP, and tree-sitter parsing are laid out as follow-on packages.
