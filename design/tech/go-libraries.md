---
title: Go libraries
updated: 2026-09-15
status: active
---

# Go libraries used by flai

Versions are pinned in `flai/go.mod`; the ones here are the majors we track.

| Library | Version | Purpose | Why this one |
|---------|---------|---------|--------------|
| `github.com/spf13/cobra` | v1.10.2 | Command tree, flags, help, completion | The de facto Go CLI framework, agents know it well |
| `github.com/goccy/go-yaml` | v1.19 | Parse and write front matter, `system-flow.yaml`, `template.yaml` | Actively maintained, preserves comments better than `gopkg.in/yaml.v3`, which is archived |
| `github.com/charmbracelet/huh` | v1.0 | Interactive prompts for `new` and `import` | Composable forms, accessible mode, works non-interactively when answers are given by flags |
| `github.com/charmbracelet/lipgloss` | v1 | Table and board rendering | Same ecosystem as huh |
| `github.com/modelcontextprotocol/go-sdk` | v1.8.0 | `flai mcp`: MCP server, stdio transport, typed tools with inferred JSON schemas, in-memory transports for tests | The official SDK, maintained with the specification; hand-rolling JSON-RPC and schema inference would be the alternative |
| `golang.org/x/term` | v0.46 | Detect whether stdin is a terminal, to decide between prompting and defaults | Standard extended library |
| `github.com/adrg/xdg` | v0.5 | Not used: config is fixed at `~/.flai` by decision | Listed so nobody adds it |

Deliberately not used:

- `viper`: config is one small JSON file, `encoding/json` is enough and keeps behaviour obvious.
- `go-git`: cloning with branches, submodules, and credentials is more reliable by shelling out to the user's `git`. See [ADR 0010](../adrs/0010-shell-out-to-git-and-docker.md).
- Docker SDK: same reasoning, `docker` on `PATH` is required and invoked as a subprocess.
- Markdown parsers: `flai` only needs front matter and headings, a small hand-written splitter is enough.
