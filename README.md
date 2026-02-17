![Release](https://img.shields.io/github/v/release/vieitesss/jocq?display_name=tag&sort=semver) ![License](https://img.shields.io/github/license/vieitesss/jocq)

# jocq

*JSON Operations, Control & Query*

jocq is a Go-based terminal UI for interactively querying JSON files with full jq syntax (via gojq).

## Status

This project is a work in progress, but it is functional for file-based querying.

## Current Features

- File input through `--file` / `-f`
- Synchronous ingestion of a single JSON file before the TUI starts
- Input size guard (`<= 100 MB`) and explicit empty-file errors
- Split layout with a query bar, raw JSON pane, and query result pane
- Query execution with full jq language support through `gojq`
- Pretty-printed JSON output in the result pane
- Inline query error display in the result pane
- Keyboard focus navigation across panes

## Quick Start

Requirements:

- Go (as specified in `go.mod`)

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

- `Tab`: move focus forward (Query -> Raw JSON -> Query Result)
- `Shift+Tab`: move focus backward
- `Enter` (from query input): run current jq query
- `Enter` (from Raw JSON pane): collapse/expand object or array at cursor
- `Ctrl+C`: quit

## Notes and Limitations

- Current mode is file-only; stream/pipeline input is not implemented yet.
- Ingestion currently loads the full file into memory before launching the UI.
- Query execution runs against decoded JSON values kept in memory.

## Roadmap

Planned areas include:

- Better large-data ergonomics and progressive/async ingestion
- Richer JSON exploration UX
- Real-time stream mode
- Additional query workflow improvements
