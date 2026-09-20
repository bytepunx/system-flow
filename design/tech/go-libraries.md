---
title: Go libraries
updated: 2026-09-20
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
| `github.com/coder/websocket` | v1.8 | `flai serve`: the WebSocket flai opens to the dashboard's agent endpoint (ADR-0029, S-0072); also the stand-in dashboard in `internal/channel/channeltest` | Maintained (ISC), no dependencies, context-aware reads and writes, concurrent writers, a ping API; `gorilla/websocket` is stable but slow-moving and deadline-driven |
| `github.com/modelcontextprotocol/go-sdk` | v1.8.0 | `flai mcp`: MCP server, stdio transport and, since S-0076, its Streamable HTTP handler in both modes (`internal/mcphttp`); typed tools with inferred JSON schemas; in-memory transports and its own client in tests | The official SDK, maintained with the specification; hand-rolling JSON-RPC and schema inference would be the alternative |
| `golang.org/x/term` | v0.46 | Detect whether stdin is a terminal, to decide between prompting and defaults | Standard extended library |
| `github.com/adrg/xdg` | v0.5 | Not used: config is fixed at `~/.flai` by decision | Listed so nobody adds it |

Deliberately not used:

- `viper`: config is one small JSON file, `encoding/json` is enough and keeps behaviour obvious.
- `go-git`: cloning with branches, submodules, and credentials is more reliable by shelling out to the user's `git`. See [ADR 0010](../adrs/0010-shell-out-to-git-and-docker.md).
- Docker SDK: same reasoning, `docker` on `PATH` is required and invoked as a subprocess.
- Markdown parsers: `flai` only needs front matter and headings, a small hand-written splitter is enough.

No library watches files. `flai serve` tells a dashboard which files changed (S-0073) by polling: `internal/watch` stats the manifest's three folders and the manifest every 300 ms and reports a file once it has looked the same for one tick. `fsnotify` was the alternative; inotify, kqueue, and Windows each have limits of their own (watches per user, a descriptor per file, no recursion), and a walk of a few hundred Markdown files costs less than the difference is worth. `wait_for_events` in `flai mcp` polls for the same reason.

## MCP revisions served (S-0076)

`flai mcp` speaks whatever the SDK negotiates; what differs by revision is the HTTP transport, and flai serves both generations at one address ([ADR-0030](../adrs/0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md)). SDK 1.8.0 knows 2024-11-05, 2025-03-26, 2025-06-18, 2025-11-25, and 2026-07-28.

| | Up to 2025-11-25 | 2026-07-28 |
|-|-|-|
| Start | `initialize`, then `notifications/initialized` | None. A client may ask `server/discover`; each call carries the revision, the client's info, and its capabilities in `_meta`, and the revision again in `Mcp-Protocol-Version` |
| Sessions | `Mcp-Session-Id` from `initialize`, repeated on every request; `DELETE` ends one; an unknown one is 404 | Removed. Every request stands alone |
| GET stream | A standalone event stream for messages the server starts | Removed, with resumption (`Last-Event-ID`) and requests from server to client (sampling, elicitation, roots) |
| In the SDK | A handler with sessions. It refuses a 2026-07-28 call with "unsupported protocol version" and the revisions it does serve, and a client falls back | Served over HTTP only by a handler with `Stateless: true`, which in turn gives older clients no session |
| In flai | Sessions that idle out and are capped; the agent named at `initialize` | No sessions; the agent named on each request, by header or by `_meta` clientInfo |

Tried on 2026-09-20 with the SDK's own client, which asks for 2026-07-28 first and falls back to `initialize`: both paths against `internal/mcphttp`, named by header and by client, and raw HTTP for the session path as a 2025-06-18 client sends it. flai uses nothing the newer revision removes: it sends nothing unprompted, and an agent waits by holding `wait_for_events`. When clients stop speaking the older revisions, the session half goes and nothing an agent configured changes.
