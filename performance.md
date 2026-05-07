# Performance Notes

The hot path is:

```text
store read -> dedup -> sort -> pack
```

Avoid adding network calls or parser work to the compression hot path unless
the result is cached before compression starts.

