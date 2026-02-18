![Release](https://img.shields.io/github/v/release/vieitesss/jocq?display_name=tag&sort=semver) ![License](https://img.shields.io/github/license/vieitesss/jocq)

# jocq

*JSON Operations, Control & Query*

jocq is a Go-based terminal UI for interactively exploring and querying JSON files with full jq syntax (via gojq).

## Status

This project is a work in progress, but currently functional for file-based JSON exploration and querying.

## Current Features

- File input through `--file` / `-f`
- Synchronous ingestion of a single JSON file before the TUI starts
- Input size guard (`<= 100 MB`) with explicit empty-file and invalid-JSON ingestion errors
- Split layout with a query bar, source explorer pane, result pane, and help bar
- Source explorer built from decoded JSON values (no raw-text renderer)
- Deterministic source tree ordering (object keys are sorted)
- Collapsible objects/arrays in the source explorer (`Enter` to toggle; collapsed containers render as `{...}` / `[...]`)
- Relative line numbers in the source gutter (current line is always `0`)
- Editor-like navigation in the source pane (`j`/`k`, arrows, `g`/`G`, `ctrl+u`/`ctrl+d`, numeric prefixes)
- `{` / `}` to jump to the previous/next container node in the source pane
- Pending motion count shown in the source pane title while typing a numeric prefix
- Full jq query execution through `gojq`
- Bounded LRU cache (256 entries) for compiled jq programs; both valid programs and parse errors are cached
- Pretty-printed JSON output in the result pane
- Inline query error display in the result pane
- Pane headers with progress metadata (source cursor percent and result scroll percent)
- Help footer with key hints and full-help toggle

## Quick Start

Requirements:

- Go (as specified in `go.mod`; currently `go 1.25.5`)

Run with the bundled example file:

```bash
go run ./cmd/jocq -f assets/example.json
```

Or with your own file:

```bash
go run ./cmd/jocq --file /path/to/data.json
```

If you use `just`:

```bash
just file
```

## Controls

- `Tab`: move focus forward (Query -> Source -> Result)
- `Shift+Tab`: move focus backward
- `Enter` (in query input): run current jq query
- `Enter` (in source pane): collapse/expand object or array at cursor
- `j` / `k` and `Up` / `Down` (source/result panes): line-by-line movement
- `[count]j` / `[count]k` (and arrow variants): counted line movement in source pane
- `g` / `G` or `Home` / `End` (source pane): jump to top / bottom
- `{` / `}` (source pane): jump to previous / next container node
- `ctrl+u` / `ctrl+d` (source/result panes): half-page movement
- `?` (source/result panes): toggle full help
- `Ctrl+C`: quit

## Notes and Limitations

- Current mode is file-only; stream/pipeline input is not implemented yet.
- Ingestion currently loads the full file into memory before launching the UI.
- Query execution runs against decoded JSON values kept in memory.
- Query execution is currently manual (`Enter`), not debounced/cancellable yet.

## Roadmap

Planned areas include:

- Better large-data ergonomics and progressive/async ingestion
- Richer JSON exploration UX
- Real-time stream mode
- Query workflow improvements (including reactive/debounced execution)
