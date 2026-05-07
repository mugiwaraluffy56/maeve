# Verify maeve

Run the standard checks from the repository root:

```bash
go test ./...
go vet ./...
go build ./cmd/maeve ./cmd/maevd
```

Smoke-test the CLI with a temporary store:

```bash
STORE=/tmp/maeve-smoke.db

go run ./cmd/maeve --store "$STORE" init --name smoke
go run ./cmd/maeve --store "$STORE" ingest file PRD.md
go run ./cmd/maeve --store "$STORE" status
go run ./cmd/maeve --store "$STORE" compress --budget 100
go run ./cmd/maeve --store "$STORE" compress --budget 100 --report
go run ./cmd/maeve --store "$STORE" snapshot save "smoke test"
go run ./cmd/maeve --store "$STORE" snapshot list
go run ./cmd/maeve --store "$STORE" prune --threshold 0.3
```

Expected signs:

- `go test ./...` passes.
- `go vet ./...` prints nothing and exits `0`.
- `init` prints `session sess_...`.
- `ingest file PRD.md` prints `object obj_... tokens=... score=...`.
- `status` shows object and token counts.
- `compress --report` prints JSON with `original_tokens`,
  `compressed_tokens`, `kept_objects`, and `dropped_objects`.

Smoke-test the daemon starter:

```bash
go run ./cmd/maevd
```

In another terminal:

```bash
curl http://localhost:7432/healthz
curl http://localhost:7432/api/v1/version
```

Expected output:

```text
ok
{"version":"dev"}
```

