---
title: flaiover dashboard
updated: 2026-09-15
status: active
---

# flaiover dashboard

`flaiover` is the web view of a conforming repo. It renders every markdown document, searches `design/` and `wip/`, shows the kanban board, and charts the flow metrics defined in [metrics.md](metrics.md). It runs locally in Docker with the repo mounted read-write and is started by `flai dashboard`.

Source: `./flaiover` in this monorepo. Image: `ghcr.io/bytepunx/flaiover`.

## Shape

SvelteKit with `adapter-node`. The browser side is a single-page app (`ssr = false` at the root layout) so navigation is instant and state lives in the client. The server side is a thin set of `+server.ts` endpoints under `/api` that read and write the mounted repo. This is the only way a browser app can reach the filesystem, so "SPA" here means the UI, not the deployment. See ADR 0007.

```mermaid
flowchart LR
    B[Browser SPA<br/>Svelte 5 + Tailwind] -->|/api/*| S[SvelteKit server<br/>adapter-node]
    S -->|read, watch, write| M[/project mount]
    S -->|index| I[(in-memory index<br/>front matter, search)]
```

## Views

| Route | View | Data |
|-------|------|------|
| `/` | Overview: WIP by column, throughput this week, aging items, epics burn-up sparkline | `/api/stats` |
| `/board` | Kanban board: one column per state with WIP count against the limit, story cards (epics and tasks on request) with age in column, nature, blocked flag; drag to transition, refused moves show the rule; refreshes on SSE (S-013) | `/api/board`, `POST /api/items/:id/move` |
| `/items/:id` | Item detail: front matter, rendered body, children, transitions timeline, blocked intervals, narrative link, and actions for allowed moves, block, unblock, and a narrative log entry (S-013) | `/api/items/:id`, `POST .../move`, `.../block`, `.../unblock`, `POST /api/streams/:id/log` |
| `/charts/cycle-time` | Cycle time scatter by nature with p50 and p85 lines (S-014) | `/api/stats` |
| `/charts/burn-up` | Scope and done, per epic or total | `/api/stats` |
| `/charts/cfd` | Cumulative flow diagram | `/api/stats` |
| `/charts/time-in-state` | Stacked bars per completed item and the share bar | `/api/stats` |
| `/charts/throughput` | Weekly bars by nature | `/api/stats` |
| `/charts/aging` | Aging WIP against p85 | `/api/stats` |
| `/charts/estimates` | Estimate versus actual with the perfect-estimate diagonal | `/api/stats` |
| `/docs/<path>` | Documentation explorer: collapsible tree of `design/` (including `conventions/` and `issues/`), `docs/`, and `wip/`; rendered markdown with Mermaid, highlighted code, task lists, heading anchors, rewritten links; front matter panel (S-012) | `/api/docs/tree`, `/api/docs/file` |
| `/conventions` | The conventions in read order with project additions highlighted; the same set `flai prime` prints | `/api/conventions` |
| `/adrs` | ADR list with status, date, and supersession chain linking into the explorer (S-012) | `/api/docs/adrs` |
| `/streams` | Active narratives with current state and next steps | `/api/streams` |
| `/search` | Search across `design/` and `wip/`, `docs/` on request, with snippets and routes (S-012) | `/api/search?q=&docs=` |

Every chart has the same filter bar: window, type, and epic where the chart supports it; a summary strip shows completed, cancelled, WIP, throughput, and cycle time percentiles; a table view sits under every chart.

## API (S-011)

| Route | Returns |
|-------|---------|
| `GET /api/manifest` | `system-flow.yaml` parsed |
| `GET /api/items?type=&status=&archived=` | Every item from kanban and archive without bodies, sorted by ID; timestamps as `YYYY-MM-DDTHH:MM:SSZ` strings |
| `GET /api/items/:id` | `{ item, children }` with the body |
| `GET /api/docs/tree` | Three trees (design, docs, wip) of `{ name, path, kind, title?, frontMatter?, children? }` |
| `GET /api/docs/file?path=` | `{ path, frontMatter, body, raw }` for one markdown file; paths outside the repo or non-markdown are 400, missing 404 |
| `GET /api/events` | Server-sent events: `ready` once, then `change` with `{ path }` per changed file under design, docs, wip, or the manifest |
| `GET /api/search?q=&docs=` | `{ query, indexed, hits[] }`; each hit has `path`, `kind`, `itemId?`, `title`, `scope`, `status?`, `type?`, `score`, `snippet`, `route` |
| `GET /api/docs/adrs` | ADR front matter: `id`, `title`, `status`, `date`, `supersedes[]`, `supersededBy[]`, `path` |
| `GET /api/board` | `{ wip_limits, order, writable, columns: { <state>: card[] } }`; a card has `id`, `type`, `title`, `nature`, `parent`, `status`, `blocked`, `age_seconds`, `entered_at` |
| `POST /api/items/:id/move` `{ to, reason?, by? }` | flai move; `{ id, status, warnings[] }` or 400 `{ error }` with the rule |
| `POST /api/items/:id/block` `{ reason }`, `POST /api/items/:id/unblock` | flai block and unblock |
| `POST /api/streams/:id/log` `{ entry }` | flai stream log |
| `GET /api/stats?since=&type=&by=` | `flai stats --json` verbatim (see metrics.md), cached per query and cleared on change; bad arguments are 400 |
| `GET /_health`, `GET /_ready`, `GET /metrics` | Liveness, readiness with named dependency checks, Prometheus metrics (S-032, see docs/operators) |

Errors are `{ error }` with the status. The reader caches by path and mtime and is invalidated by the watcher.

## Writes

The mount is read-write so the board can be operated from the browser. Every write (move, block, unblock, narrative log) is performed by invoking the bundled `flai` binary with `--json` inside the project ([ADR-0016](../adrs/0016-dashboard-delegates-to-flai.md)); a refused write returns flai's rule text. The server finds the binary through `FLAI_BIN`, then `PATH`, and runs it with `FLAI_CONFIG` and `FLAI_CACHE_DIR` under the project's `.flai-cache` and `FLAI_AGENT=flaiover`. Without a binary the dashboard is read-only and the board says so. Writes are ordinary file edits, so they show up in `git status` for the human to commit. Metrics (S-014) come from `flai stats --json` the same way.

## Search

Server builds a MiniSearch index over title, tags, ID, headings, and body text of every markdown file under `design/` (conventions and issues included) and `wip/`, rebuilt on file change. Results link to the docs explorer or the item page.

## Runtime

- Container listens on `3000`. `flai dashboard` maps it to the configured host port on `127.0.0.1`, default `4242`, and runs the container as the host user (`--user uid:gid`), so the image must work as an arbitrary non-root UID: no privileged ports, no writes outside `/project` and `/tmp`, and a writable working directory is not assumed.
- Repo mounted at `/project`. `PROJECT_DIR` overrides for development outside Docker.
- File watching with `chokidar`, debounced, invalidates the index and pushes updates to open tabs with server-sent events.
- No authentication. It is a local tool bound to localhost by `flai`. Operators exposing it further are told not to in `docs/operators`.

## Internal structure

```text
flaiover/
├── src/
│   ├── lib/
│   │   ├── server/       # repo reader, front matter, metrics port, search index, watcher
│   │   ├── components/   # board, cards, charts, doc tree, markdown renderer
│   │   └── stores/       # client state
│   └── routes/
│       ├── +layout.ts    # ssr = false
│       ├── api/          # +server.ts endpoints
│       └── ...           # views above
├── static/
├── Dockerfile            # multi-stage, node:24-alpine runtime
└── tests/                # vitest unit, playwright e2e against the template sample repo
```
