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
| `/board` | Kanban board, columns from `board.md`, cards with age, nature, blocked flag, drag to transition | `/api/items`, `/api/items/:id/move` |
| `/items/:id` | Item detail: rendered body, front matter, children, transitions timeline, narrative link | `/api/items/:id` |
| `/charts/cycle-time` | Cycle time scatter with percentiles | `/api/stats/cycle-time` |
| `/charts/burn-up` | Per epic and total | `/api/stats/burn-up` |
| `/charts/cfd` | Cumulative flow diagram | `/api/stats/cfd` |
| `/charts/time-in-state` | Stacked bars and aggregate share | `/api/stats/time-in-state` |
| `/charts/throughput` | Weekly bars by nature | `/api/stats/throughput` |
| `/charts/aging` | Aging WIP | `/api/stats/aging` |
| `/docs` | Documentation explorer: tree of `design/` (including `conventions/` and `issues/`), `docs/`, and `wip/` with rendered markdown, Mermaid, and front matter panel | `/api/docs/tree`, `/api/docs/file` |
| `/conventions` | The conventions in read order with project additions highlighted; the same set `flai prime` prints | `/api/conventions` |
| `/adrs` | ADR list with status and supersession chain | `/api/docs/adrs` |
| `/streams` | Active narratives with current state and next steps | `/api/streams` |
| `/search` | Full text and front matter search across `design/` and `wip/` | `/api/search` |

Every chart has the same filter bar: window, nature, epic.

## Writes

The mount is read-write so the board can be operated from the browser. Writes are limited to what `flai` also does: transitions, block and unblock, log entries, and editing an item's body. Every write goes through the same validation rules as `flai check`; the server ports the rules from the Go reference and the fixture test keeps them aligned. Writes are ordinary file edits, so they show up in `git status` for the human to commit.

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
