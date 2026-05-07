# Compression Pipeline

The v0.1 compressor is deterministic:

```text
load objects
    |
    v
exact SHA-256 dedup
    |
    v
sort by importance desc
    |
    v
greedy token budget pack
    |
    v
payload + report
```

Later versions can add near-duplicate detection and tree-sitter truncation
without changing the external CLI contract.

