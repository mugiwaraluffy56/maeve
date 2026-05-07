# Storage Reference

The default store lives under `.maeve/maeve.db` in the current repository.

SQLite is opened with one writer connection to match SQLite's local-first write
model and avoid accidental lock contention in the CLI.

