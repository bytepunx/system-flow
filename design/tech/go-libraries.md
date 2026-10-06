---
title: Go libraries
updated: 2026-10-06
status: active
topics: [go]
---

# Go libraries used by flai

Versions are pinned in `flai/go.mod`; the ones here are the majors we track.

| Library | Version | Purpose | Why this one |
|---------|---------|---------|--------------|
| `github.com/spf13/cobra` | v1.10.2 | Command tree, flags, help, completion | The de facto Go CLI framework, agents know it well |
| `github.com/goccy/go-yaml` | v1.19 | Parse and write front matter, `system-flow.yaml`, `template.yaml` | Actively maintained, preserves comments better than `gopkg.in/yaml.v3`, which is archived |
| `github.com/coder/websocket` | v1.8 | `flai serve`: the WebSocket flai opens to the dashboard's agent endpoint (ADR-0029, S-0072); also the stand-in dashboard in `internal/channel/channeltest` | Maintained (ISC), no dependencies, context-aware reads and writes, concurrent writers, a ping API; `gorilla/websocket` is stable but slow-moving and deadline-driven |
| `github.com/modelcontextprotocol/go-sdk` | v1.8.0 | `flai mcp`: MCP server, stdio transport and, since S-0076, its Streamable HTTP handler in both modes (`internal/mcphttp`); typed tools with inferred JSON schemas; in-memory transports and its own client in tests | The official SDK, maintained with the specification; hand-rolling JSON-RPC and schema inference would be the alternative |
| `golang.org/x/term` | v0.46 | Detect whether stdin is a terminal, to decide between prompting and defaults | Standard extended library |
| `github.com/adrg/xdg` | v0.5 | Not used: config is fixed at `~/.flai` by decision | Listed so nobody adds it |

Deliberately not used:

- `viper`: config is one small JSON file, `encoding/json` is enough and keeps behaviour obvious.
- `go-git`: cloning with branches, submodules, and credentials is more reliable by shelling out to the user's `git`. See [ADR 0010](../adrs/0010-shell-out-to-git-and-docker.md).
- Docker SDK: same reasoning, `docker` on `PATH` is required and invoked as a subprocess.
- Markdown parsers: `flai` only needs front matter and headings, a small hand-written splitter is enough.
- Glob libraries (`bmatcuk/doublestar`): the patterns of `claims.shared` (S-0295, [ADR-0096](../adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md)) are matched a segment at a time with `path.Match`, and `**` is handled by hand in a few lines of `internal/manifest/shared.go`. The library would add a dependency for the one thing the standard library lacks, and the rule that a folder lies inside a pattern only when everything below it does would still be flai's own.
- Prompt and terminal UI libraries (`charmbracelet/huh`, `bubbletea`, `lipgloss`): removed in S-0160. See below.

## Prompts (S-0160)

flai asks few questions: the template's variables and the layout in `new` and `import`, whether to apply an import and move its folders, where a loose markdown file belongs, how to settle an upgrade conflict, and whether to cancel what is open under an item. `internal/prompt` asks them a line at a time on the terminal, with the standard library: a value with a default and a check, yes or no with a default, or one of a few numbered options. Every question has a flag or `--yes` that answers it without a terminal.

Until S-0160 they were `charmbracelet/huh` forms. huh brings in `bubbles/textarea`, and through it `atotto/clipboard`, which looks for clipboard programs on `PATH` when its package loads, whether flai prompts or not. On a host whose `PATH` has 54 entries, 17 of them Windows folders under `/mnt`, that took 145 ms of every flai process ([cause 5](../system/server-performance.md#what-it-says)). The alternatives were huh fields that do not use `textarea`, which does not help because huh imports it whichever fields are used; a `replace` of `atotto/clipboard` with an empty module, which `go install` of flai refuses; and huh behind a build tag or a second binary, which keeps a dependency tree of 25 modules for six questions. Dropping huh took those 25 modules out of `go.mod`. What was lost is arrow-key selection and the form styling. The decision is [ADR-0052](../adrs/0052-flai-asks-its-questions-a-line-at-a-time-with-its-own-prompt-package-not-with.md), which refines [ADR-0006](../adrs/0006-go-for-the-cli.md).

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
