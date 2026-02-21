# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Spec Viewer is a Go-based live-preview tool for markdown technical specifications. It serves a web UI that renders `.md` files with real-time reload via WebSocket when files change on disk. Part of the Spec Kit ecosystem for Spec Driven Development (SDD).

## Build & Development Commands

```bash
# Development with hot-reload (recommended)
air

# One-shot development run with example specs
make run

# Production build → bin/spec-viewer
make build

# Run built binary
./bin/spec-viewer serve --port 9091 --folder ./specs
```

No test suite exists yet.

## Architecture

**CLI → Server → WebSocket + File Watcher** pipeline:

1. **Entry point**: `cmd/spec-viewer/` — Cobra CLI with `serve` subcommand that starts the HTTP server, WebSocket hub, and file watcher concurrently
2. **HTTP routing** (`internal/server/`): Gorilla Mux with routes: `/` (home), `/view?file=<path>` (markdown viewer), `/ws` (WebSocket), `/public/*` (static assets)
3. **Handlers** (`internal/handlers/`): Each route has its own handler. `ViewSpecHandler` converts markdown→HTML via Goldmark with directory traversal protection
4. **WebSocket hub** (`internal/socket/`): Thread-safe client registry with broadcast. File watcher sends "reload" → all clients call `window.location.reload()`
5. **File watcher** (`internal/watcher/`): fsnotify-based recursive directory watcher that triggers WebSocket broadcasts on file changes
6. **Spec discovery** (`internal/spec/`): Recursive directory scanner that builds a tree of `.md` files (excluding hidden files), rebuilt on each request for freshness
7. **Templates** (`internal/templates/`): Cached Go `html/template` with base layout + component partials pattern using template cloning

**Web assets** (`web/`) are embedded into the binary via `go:embed`. Templates live in `web/templates/`, static files in `web/public/`.

**Frontend stack**: Tailwind CSS (CDN) + Alpine.js + Basecoat CSS components. Dark/light theming via CSS variables and class-based Tailwind dark mode, persisted in localStorage.

## Key Packages

- `internal/` — Private packages: handlers, server, socket, spec, templates, watcher
- `pkg/` — Public packages: `logger` (slog+tint colored output), `ui` (formatted console messages)
- `web/` — Embedded templates and static assets (`efs.go` defines the embed.FS)

## Conventions

- Go module: `github.com/SantiagoBobrik/spec-viewer`
- Requires Go 1.24+
- Commit messages follow conventional commits (`feat:`, `fix:`, `refactor:`, `docs:`)
- Code formatted with `go fmt`
