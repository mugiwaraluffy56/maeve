# maeve - Product Requirements Document

**Memory and Agent's Evolution Engine**

> "Keeps coding agents fast, focused, and cheap."

---

## Problem

Coding agents degrade over long sessions:

- Context windows bloat with noise: dead imports, old terminal output, stale diffs
- Token costs spike 3-10x mid-session vs fresh session
- Agents repeat themselves, lose coherence, hallucinate on stale context
- No tooling to inspect, prune, or snapshot agent context state
- Session state lost on window close - no continuity across work days

Nobody has built a "context operating system" for agents. This is the gap.

---

## Positioning

**"Git for AI context."**

Not another AI autocomplete. Not another editor. Infrastructure layer for AI cognition.

---

## Target Users

| User | Pain |
|------|------|
| Power dev using coding agents daily | Token bills, context rot, session loss |
| OSS maintainer running AI on large repos | Agent gets confused by repo size |
| Team using AI agents in CI/CD | Expensive, unpredictable per-run costs |
| Agent framework builder | No standard context lifecycle API |

---

## Tech Stack

### Language split

| Layer | Language | Reason |
|-------|----------|--------|
| Core engine, server, MCP, CLI | Go | Fast to ship, idiomatic for servers/CLI, solid tree-sitter bindings |
| VS Code extension | TypeScript | Required by VS Code API, kept as thin shell only |

All logic lives in Go. TypeScript = UI + event bridge only.

### Go libraries

| Purpose | Library |
|---------|---------|
| AST parsing | `github.com/smacker/go-tree-sitter` + language grammars |
| SQLite | `modernc.org/sqlite` (pure Go, no CGO) |
| HTTP router | `github.com/go-chi/chi/v5` |
| CLI | `github.com/spf13/cobra` |
| Config | `github.com/spf13/viper` |
| Token counting | `github.com/pkoukk/tiktoken-go` |
| Logging | `go.uber.org/zap` |
| File watching | `github.com/fsnotify/fsnotify` |
| Testing | stdlib + `github.com/stretchr/testify` |

---

## Repo Structure (wide, not deep)

```
maeve/
├── cmd/                  CLI + daemon entrypoints
│   ├── maeve/            main.go - maeve CLI (cobra)
│   └── maevd/            main.go - server daemon
├── compress/             scoring, dedup, truncation, packing
├── ingest/               file, diff, terminal ingestion
├── parse/                tree-sitter wrapper, language grammars, tokenizer
├── session/              session lifecycle, snapshot management
├── store/                SQLite driver, migrations
├── server/               HTTP router, handlers, WebSocket
├── mcp/                  MCP server (JSON-RPC 2.0, tools)
├── ext/                  VS Code extension (TypeScript)
├── maeve.toml.example
├── go.mod
├── go.sum
├── Makefile
└── PRD.md
```

Flat. No `internal/`, no `pkg/`. Each top-level dir is one concern.

---

## Core Concepts

### ContextObject

Atomic unit maeve tracks.

| Type | Description |
|------|-------------|
| `file_snapshot` | File content at point in time |
| `terminal_block` | Chunked shell session output |
| `conversation_turn` | Agent message pair |
| `diff_block` | Git diff segment |
| `symbol` | Function/class/struct extracted by tree-sitter |

### Session

Ordered log of ContextObjects with metadata: timestamps, importance scores, agent references.

### Snapshot

Named checkpoint of compressed session state. Portable, shareable, diffable.

### ImportanceScore

`float64 [0.0–1.0]` computed per ContextObject. Drives all eviction/pruning decisions.

---

## Importance Scoring

```
score = w1*recency + w2*access_freq + w3*symbol_type + w4*diff_coverage + w5*agent_refs
```

| Factor | Weight (default) | Notes |
|--------|-----------------|-------|
| `recency` | 0.30 | Exponential decay, half-life configurable |
| `access_freq` | 0.20 | How many times agent/user touched this |
| `symbol_type` | 0.25 | fn signature > type def > impl > comment > whitespace |
| `diff_coverage` | 0.15 | +0.3 bonus if symbol in open git diff |
| `agent_refs` | 0.10 | Did agent explicitly reference this object? |

All weights configurable in `maeve.toml`.

---

## Compression Engine

Given a token budget, produce optimal context payload.

### Step 1 - Structural Dedup

- SHA-256 hash each ContextObject
- Identical hash → drop duplicate, keep highest-scored version
- Near-duplicate: SimHash with threshold 0.85 → merge

### Step 2 - Semantic Truncation (tree-sitter)

Per language grammar:

| Object | Condition | Action |
|--------|-----------|--------|
| Function | score < 0.4 | Emit signature only, drop body |
| Struct/class | score < 0.3 | Emit field names only, drop method bodies |
| Comments | score < 0.5 | Drop (except doc comments on public API) |
| Import blocks | any | Deduplicate across files, emit once |

### Step 3 - Temporal Decay

Objects older than `decay_window` get score multiplied by `0.7^(hours_old)`.

### Step 4 - Diff Awareness

Files touched in current `git diff` get `+0.3` score bonus. Prevents pruning active work.

### Step 5 - Pack

Sort by score descending. Greedily fill budget. Stop when budget exhausted.

Output: compressed context payload + compression report (original tokens → compressed tokens, dropped objects, reasons).

---

## SQLite Schema

```sql
CREATE TABLE sessions (
  id         TEXT PRIMARY KEY,
  name       TEXT,
  repo_path  TEXT NOT NULL,
  branch     TEXT,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

CREATE TABLE context_objects (
  id           TEXT PRIMARY KEY,
  session_id   TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  type         TEXT NOT NULL,
  source_path  TEXT,
  content      BLOB NOT NULL,   -- zstd compressed
  importance   REAL NOT NULL DEFAULT 0.5,
  token_count  INTEGER NOT NULL,
  created_at   INTEGER NOT NULL,
  accessed_at  INTEGER NOT NULL
);

CREATE TABLE snapshots (
  id                TEXT PRIMARY KEY,
  session_id        TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  name              TEXT NOT NULL,
  compressed_payload BLOB NOT NULL,
  original_tokens   INTEGER NOT NULL,
  compressed_tokens INTEGER NOT NULL,
  created_at        INTEGER NOT NULL
);

CREATE INDEX idx_context_objects_session    ON context_objects(session_id);
CREATE INDEX idx_context_objects_importance ON context_objects(importance);
CREATE INDEX idx_snapshots_session          ON snapshots(session_id);
```

---

## HTTP API

Server runs on `localhost:7432`.

```
POST   /api/v1/sessions                      Create session
GET    /api/v1/sessions/:id                  Get session metadata
GET    /api/v1/sessions/:id/status           Token count, object count, score histogram
POST   /api/v1/sessions/:id/ingest           Ingest ContextObject
POST   /api/v1/sessions/:id/compress         Compress to budget → return payload + report
POST   /api/v1/sessions/:id/prune            Remove objects below importance threshold
POST   /api/v1/sessions/:id/snapshots        Create snapshot
GET    /api/v1/sessions/:id/snapshots        List snapshots
GET    /api/v1/snapshots/:id                 Get snapshot payload
POST   /api/v1/snapshots/diff                Diff two snapshots → structured delta
WS     /api/v1/sessions/:id/stream           Live token count + score distribution updates
```

---

## MCP Server

Runs on `localhost:7433`. Speaks JSON-RPC 2.0 over stdio or HTTP.
Exposes maeve to any MCP-compatible agent as callable tools.

### Tools

#### `maeve_compress`
Compress current session context to token budget.

```json
{
  "name": "maeve_compress",
  "description": "Compress current session context to a token budget. Returns optimized context payload and savings report.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "budget_tokens": { "type": "integer" },
      "session_id":    { "type": "string", "description": "Defaults to repo-detected session" }
    },
    "required": ["budget_tokens"]
  }
}
```

#### `maeve_snapshot`
Save named snapshot of current session state.

```json
{
  "name": "maeve_snapshot",
  "inputSchema": {
    "type": "object",
    "properties": {
      "name":       { "type": "string" },
      "session_id": { "type": "string" }
    },
    "required": ["name"]
  }
}
```

#### `maeve_status`
Get current token usage, compression ratio, object count.

```json
{
  "name": "maeve_status",
  "inputSchema": {
    "type": "object",
    "properties": {
      "session_id": { "type": "string" }
    }
  }
}
```

#### `maeve_prune`
Remove low-importance context objects below threshold.

```json
{
  "name": "maeve_prune",
  "inputSchema": {
    "type": "object",
    "properties": {
      "threshold":  { "type": "number", "minimum": 0, "maximum": 1 },
      "session_id": { "type": "string" }
    },
    "required": ["threshold"]
  }
}
```

#### `maeve_ingest_diff`
Ingest a git diff for importance-scoring and context tracking.

```json
{
  "name": "maeve_ingest_diff",
  "inputSchema": {
    "type": "object",
    "properties": {
      "diff_content": { "type": "string" },
      "session_id":   { "type": "string" }
    },
    "required": ["diff_content"]
  }
}
```

---

## CLI (`maeve`)

```bash
maeve init                           # Init session for current repo
maeve status                         # Session stats: objects, tokens, score dist
maeve compress --budget 8000         # Compress to token budget, print payload
maeve snapshot save "before refactor"
maeve snapshot restore <id>
maeve snapshot list
maeve diff <snap1> <snap2>           # Structured diff between two snapshots
maeve watch                          # Daemon mode: auto-ingest file changes
maeve prune --threshold 0.3          # Drop objects with score < 0.3
maeve serve                          # Start HTTP + MCP servers
```

---

## VS Code Extension

Thin TypeScript shell. All logic in Go server. Extension = UI + event bridge only.

### Panels

| Panel | Description |
|-------|-------------|
| Context Lens | Sidebar: live token budget bar, object list with scores, prune button |
| Snapshot Manager | List snapshots, restore, side-by-side diff |
| Compression Report | After compress: what dropped, why, savings % |

### Events → maeve-server

| VS Code event | maeve action |
|---------------|-------------|
| `onDidSaveTextDocument` | `POST /ingest` file snapshot |
| `onDidChangeGitBranch` | New session or switch active session |
| Terminal data | `POST /ingest` terminal block |

### Status Bar

```
maeve: 4.2k / 8k tokens  [compress]  [snapshot]
```

---

## Configuration (`maeve.toml`)

```toml
[session]
auto_init          = true
decay_window_hours = 2
decay_factor       = 0.7

[compression]
default_budget_tokens  = 8000
dedup_threshold        = 0.85
truncation_threshold   = 0.4

[scoring]
weight_recency        = 0.30
weight_access_freq    = 0.20
weight_symbol_type    = 0.25
weight_diff_coverage  = 0.15
weight_agent_refs     = 0.10

[server]
port     = 7432
mcp_port = 7433

[tokenizer]
model = "cl100k_base"   # or "o200k_base"
```

---

## Data Flow

```
File save / terminal output / git diff
          ↓
  maeve-server POST /ingest
          ↓
  tree-sitter parse → symbol extraction
          ↓
  importance scoring
          ↓
  SQLite store (zstd compressed blob)
          ↓
  WebSocket push → VS Code extension (live token bar)

Agent calls maeve_compress (MCP) or POST /compress (HTTP):
          ↓
  Load all session objects from SQLite
          ↓
  Score → dedup → truncate → pack
          ↓
  Return compressed payload + report
          ↓
  Agent uses payload as context
```

---

## Milestones

| Phase | Scope | Goal |
|-------|-------|------|
| v0.1 | Core compress engine + SQLite store + `maeve compress` / `maeve snapshot` CLI | Working compression end-to-end |
| v0.2 | `maeve-server` HTTP API + tree-sitter integration + `maeve watch` file watcher | Live ingestion |
| v0.3 | MCP server for agent integration | Agent-callable |
| v0.4 | VS Code extension with Context Lens panel and status bar | Visual feedback |
| v0.5 | Compression reports, snapshot diff viewer, token visualization | Full observability |
| v1.0 | Polish, docs, VS Code marketplace publish | Public launch |

---

## Success Metrics

| Metric | Target |
|--------|--------|
| Token reduction per compress | 60–80% on typical 30min session |
| Compression latency | < 100ms for sessions up to 50k tokens |
| Time to first compress (ext install → compress) | < 2 minutes |
| Agent coherence (dogfood on maeve's own repo) | Subjective, measurably fewer repetitions |

---

## Competitive Landscape

| Tool | Gap maeve fills |
|------|----------------|
| Editor rules / ignore files | Static, manual, no scoring |
| Built-in compact commands | One-shot, no persistence, no agent API |
| Closed context windows | No visibility, no control |
| MemGPT | NLP-heavy, not repo-aware, no tree-sitter |
| None | Structured, repo-aware, agent-agnostic, open source |

**Strongest wedge**: MCP integration. Every MCP-compatible agent gets maeve for free via tool call. Ship MCP server early, let ecosystem drive adoption.
