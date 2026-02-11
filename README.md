![Release](https://img.shields.io/github/v/release/vieitesss/jocq?display_name=tag&sort=semver) ![License](https://img.shields.io/github/license/vieitesss/jocq)

# jocq

*JSON Operations, Control & Query*

> **The ultimate JSON control center for the terminal**: inspect, transform, generate commands, and monitor JSON streams in real time — all from a fast, keyboard-driven TUI.

⚠️ **Status:** Pre-release / work in progress. The application is not functional yet; this README describes the intended scope.

---

## What it will do

This project aims to become a high-performance **Terminal User Interface (TUI)** for working with JSON in both **static files** and **live streams**.

### 1) Analyze JSON files (static mode)
- **Tree explorer** for deeply nested JSON (collapsible navigation)
- **Search & filter** keys/values across large documents
- **Structured edits** (update values, rename keys, delete nodes)
- **Diff-friendly output** and optional formatting/pretty-print controls

### 2) Monitor JSON streams (real-time mode)
- **Follow** JSON lines / event streams (stdin, pipes, logs)
- **Live filtering** (by key paths, regex, value predicates)
- **Rate + throughput** visibility (basic observability for the stream)
- **Fast rendering** for high-volume updates

### 3) Generate commands & transformations (power mode)
- Generate **jq** (and/or other) commands from interactive selections
- Copy-to-clipboard / export of transformations and filters
- Saved “recipes” (reusable query + transform presets)
