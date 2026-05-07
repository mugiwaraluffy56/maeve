# Architecture

`maeve` is organized as wide, shallow Go packages. Each top-level directory owns
one product concern.

```text
CLI command
   |
   v
app wiring
   |
   +--> session service
   +--> ingest service
   +--> compression engine
   +--> snapshot service
   |
   v
SQLite store
```

The CLI should stay thin. Business behavior belongs in service packages so the
same code can later be exposed through HTTP and MCP.

