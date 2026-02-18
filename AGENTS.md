# AGENTS.md

## Project Summary

jocq is a Go-based Terminal User Interface for interactively querying JSON data using the full jq language. It renders a split-pane UI: a navigable source JSON explorer on the left, query results on the right, a query bar at the top, and a status/help bar at the bottom. It targets datasets from a handful of objects up to millions. The project is work in progress.

## Architecture

- The **buffer** is the central shared dependency. The scanner (ingestion) writes to it, the TUI reads from it. Both are started from `main.go` and connected only through the buffer.
- Core packages live under `internal/`: `ingest` (JSON ingestion), `buffer` (thread-safe data store), and `tui` (Bubble Tea application, including query scheduling and execution).
- JSON is decoded once on ingestion. Queries run against pre-decoded Go values, never re-parsing raw bytes.
- The query engine is **gojq** — full jq compatibility, pure Go, no CGo.
- The TUI is built with **Bubble Tea** (Charmbracelet). It uses a view-based architecture: the root application in `internal/tui/app.go` routes between views, each view is a self-contained Bubble Tea model under `internal/tui/views/...`.
- Explorer source pane architecture:
  - `internal/tree` flattens decoded JSON into line-addressable nodes (path, depth, type, value metadata). Each node carries `Collapsible`, `Collapsed`, and `ChildCount` fields for collapse/expand support.
  - `internal/tui/components/treevp` is a custom cursor viewport. It maintains `nodes` (full flat list), `visible` (index list of non-hidden nodes), `cursor` (position in `visible`), and `offset` (scroll position). `rebuildVisible` is called whenever collapse state changes.
  - Source nodes are initialized once from decoded data and reused while navigating. `SetNodes` takes ownership of the provided slice — callers must not mutate it after the call.
- Tree nodes support collapse/expand: `enter` toggles a collapsible container (ObjectOpen or ArrayOpen). Hidden children are skipped by `rebuildVisible` using a depth-tracking skip range.
- `internal/query` package owns query compilation and execution, separate from the TUI. Compiled gojq programs are cached in a bounded LRU cache (default capacity 256) keyed by trimmed query string, protected by a mutex. Both successfully compiled programs and parse errors are cached, so repeated malformed queries short-circuit without re-parsing.
- Query execution is **manual and synchronous**: `Enter` in the query input triggers `RunQuery`, which calls `query.Execute` directly on the Bubble Tea event loop — no goroutine, no debounce, no cancellation. Debounced/cancellable execution is planned but not yet implemented.
- Current ingestion mode is synchronous full-file loading (up to 100MB) before starting the TUI.

## Preferences

- Prioritize throughput and low allocations in hot paths.
- `buffer.Data.Raw()` and `buffer.Data.RawRange()` return read-only shallow snapshots; callers must not mutate returned slices.
- Prefer synchronous ingestion in file mode until realtime/chunked ingestion is implemented.
- Keep source tree rendering stable and deterministic: object keys are sorted when flattened.
- `rebuildVisible` must be called after any collapse state change; never mutate `nodes` without rebuilding `visible` afterward.
- Cursor navigation behavior should feel editor-like:
  - There is always a highlighted current line in the source viewport.
  - Relative line numbers are shown in the source gutter; the current line is always `0`.
  - Numeric prefixes apply to `j/k` motions only (for example, `5j`, `12k`).
  - `j/k` and arrows move line-by-line.
  - `g` and `G` (or `Home`/`End`) jump to top/bottom.
  - `ctrl+u` / `ctrl+d` move half-page up/down in both the source and result panes.
  - `{`/`}` jump to the previous/next `ObjectOpen` node; `[`/`]` jump to the previous/next `ArrayOpen` node. All four are only active when the source pane is focused.
- Pane focus cycles circularly: `Tab` advances `InputPane → InPane → OutPane → InputPane`; `Shift+Tab` reverses. `RunQuery` (`enter`) and `ToggleNode` (`enter`) share the same key but are mutually exclusive via `SetFocusMode` enabling/disabling each binding per pane.
- Mouse events are received and routed to the focused pane (`tea.MouseMsg` handled in the explorer view).
- Use ANSI-aware width handling in tree rendering to avoid wrapped visual lines in the source pane.

## Issues & Fixes

This section tracks problems encountered during development and how they were resolved.

- Scanner runs synchronously in file mode and returns ingestion errors directly to the CLI.
- Ingestion currently enforces a max input size of 100MB (`<= 100MB` accepted) and returns explicit errors for empty files.
- TUI explorer initialization reads decoded content from `*buffer.Data` and builds the source tree from decoded values (no raw-string source pane rendering).
- Source pane line wrapping/highlighting fixes:
  - Long lines are ANSI-truncated and padded to viewport width to keep one logical node per visible row.
  - Full-line highlight is applied per rendered segment so ANSI resets from syntax colors do not break cursor background.
- Pane sizing consistency fix: source/result panes are bounded to consistent heights in layout rendering.
- Explorer pane headers show right-aligned progress percentages (source cursor progress and result scroll progress).
- Width overflow fix: source/output panes, query input, and help footer are clipped/padded to terminal width to prevent horizontal overflow.
