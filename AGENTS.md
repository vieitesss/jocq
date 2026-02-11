# AGENTS.md

## Project Summary

jocq is a Go-based Terminal User Interface for interactively querying JSON data using the full jq language. It renders a split-pane UI — raw input on the left, query results on the right, a query bar at the top, and a status bar at the bottom. It targets datasets from a handful of objects up to millions. The project is work in progress.

## Architecture

- The **buffer** is the central shared dependency. The scanner (ingestion) writes to it, the TUI reads from it. Both are started from `main.go` and connected only through the buffer.
- Four packages live under `internal/`: `ingest` (JSON ingestion), `buffer` (thread-safe data store), `query` (gojq wrapper + scheduler), and `tui` (Bubble Tea application).
- JSON is decoded once on ingestion. Queries run against pre-decoded Go values, never re-parsing raw bytes.
- The query engine is **gojq** — full jq compatibility, pure Go, no CGo.
- The TUI is built with **Bubble Tea** (Charmbracelet). It uses a view-based architecture: a root model routes between views, each view is a self-contained Bubble Tea model under `tui/views/`.
- Queries are debounced and cancellable. A new keystroke cancels any in-flight query.
- Current ingestion mode is synchronous full-file loading (up to 100MB) before starting the TUI.

## Preferences

- Prioritize throughput and low allocations in hot paths.
- `buffer.Data.Raw()` and `buffer.Data.RawRange()` return read-only shallow snapshots; callers must not mutate returned slices.
- Prefer synchronous ingestion in file mode until realtime/chunked ingestion is implemented.

## Issues & Fixes

This section tracks problems encountered during development and how they were resolved.

- Scanner runs synchronously in file mode and returns ingestion errors directly to the CLI.
- Ingestion currently enforces a max input size of 100MB (`<= 100MB` accepted) and returns explicit errors for empty files.
- TUI initialization receives `*buffer.Data` and fetches raw content from the buffer during explorer init.
