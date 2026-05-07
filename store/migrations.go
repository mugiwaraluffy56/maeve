package store

const schemaSQL = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS sessions (
  id         TEXT PRIMARY KEY,
  name       TEXT,
  repo_path  TEXT NOT NULL,
  branch     TEXT,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS context_objects (
  id           TEXT PRIMARY KEY,
  session_id   TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  type         TEXT NOT NULL,
  source_path  TEXT,
  content      BLOB NOT NULL,
  importance   REAL NOT NULL DEFAULT 0.5,
  token_count  INTEGER NOT NULL,
  created_at   INTEGER NOT NULL,
  accessed_at  INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS snapshots (
  id                 TEXT PRIMARY KEY,
  session_id         TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  name               TEXT NOT NULL,
  compressed_payload BLOB NOT NULL,
  original_tokens    INTEGER NOT NULL,
  compressed_tokens  INTEGER NOT NULL,
  created_at         INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_context_objects_session
  ON context_objects(session_id);

CREATE INDEX IF NOT EXISTS idx_context_objects_importance
  ON context_objects(importance);

CREATE INDEX IF NOT EXISTS idx_snapshots_session
  ON snapshots(session_id);
`
