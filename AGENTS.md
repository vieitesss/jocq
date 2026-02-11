# AGENTS.md

## Project Summary

jocq is a Go-based Terminal User Interface for interactively querying JSON data using the full jq language. It renders a split-pane UI — raw input on the left, query results on the right, a query bar at the top, and a status bar at the bottom. It handles both static files and real-time piped streams, targeting datasets from a handful of objects up to millions. The project is work in progress.

## Architecture

- The **buffer** is the central shared dependency. The scanner (ingestion) writes to it, the TUI reads from it. Both are started from `main.go` and connected only through the buffer.
- Four packages live under `internal/`: `ingest` (JSON stream scanning), `buffer` (thread-safe data store), `query` (gojq wrapper + scheduler), and `tui` (Bubble Tea application).
- JSON is decoded once on ingestion. Queries run against pre-decoded Go values, never re-parsing raw bytes.
- The query engine is **gojq** — full jq compatibility, pure Go, no CGo.
- The TUI is built with **Bubble Tea** (Charmbracelet). It uses a view-based architecture: a root model routes between views, each view is a self-contained Bubble Tea model under `tui/views/`.
- Queries are debounced and cancellable. A new keystroke cancels any in-flight query.

## Preferences

<!-- Add your coding style preferences, conventions, and guidelines here -->

## Issues & Fixes

This section tracks problems encountered during development and how they were resolved.

- *(none yet — this section will be populated as development progresses)*
