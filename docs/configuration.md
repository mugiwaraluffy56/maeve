# Configuration

`maeve` reads optional TOML configuration and `MAEVE_` environment variables.

```toml
[compression]
default_budget_tokens = 8000

[server]
port = 7432
mcp_port = 7433
```

CLI flags still win for command-specific values like `--budget` and
`--session`.

