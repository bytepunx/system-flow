---
id: T-3355
type: task
nature: feature
title: The dashboard's log page and its LogViewer, LogLink, and LogText components show a log a message names
status: done
parent: S-0785
owner: alex
created: 2026-11-14T10:28:58Z
updated: 2026-11-14T10:50:32Z
transitions:
  - to: ready
    at: 2026-11-14T10:29:11Z
    by: agent-S-0785
  - to: in-progress
    at: 2026-11-14T10:29:11Z
    by: agent-S-0785
  - to: done
    at: 2026-11-14T10:50:32Z
    by: agent-S-0785
stream: S-0785
tags: []
touches: [flaiover/src/routes/api/log/+server.ts, flaiover/src/routes/api/log/log.test.ts, flaiover/src/lib/server/logs.ts, flaiover/src/lib/logs.ts, flaiover/src/lib/logs.test.ts, flaiover/src/lib/components/LogViewer.svelte, flaiover/src/lib/components/LogViewer.svelte.test.ts, flaiover/src/lib/components/LogLink.svelte, flaiover/src/lib/components/LogText.svelte, flaiover/src/lib/components/LogText.svelte.test.ts, flaiover/src/routes/ops/log/+page.svelte, flaiover/src/lib/sitemenu.ts, design/system/flaiover-dashboard.md, design/components/flaiover/src/lib/components/LogViewer.examples.json, design/components/flaiover/src/lib/components/LogLink.examples.json, design/components/flaiover/src/lib/components/LogText.examples.json, flaiover/src/lib/components/LogFeed.svelte]
usage:
  source: log
  seconds: 1281
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 113
      output: 38350
      cache_read: 7512823
      cache_write: 185009
      cost: 3.3642
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 20792
          cache_write: 579
          cost: 0.0063
---
# T-3355 The dashboard's log page and its LogViewer, LogLink, and LogText components show a log a message names

## Work

Give the dashboard a page that shows a log's contents, and the components that link a log path to it.

- `GET /api/log?path=&offset=`: flai's `log.read` on the host (T-3354), answered as `{ path, exists, offset, next, size, text, reset? }`; a missing path or an offset that is not a whole number of zero or more is refused with 400 before flai is asked; flai's refusal of a path outside what it reads keeps its message; 501 when the host flai is older than the method. Follow `routes/api/images/+server.ts` and `$lib/server/images.ts`.
- `LogViewer.svelte`: shows one log's contents by path, reading it through `/api/log` with `LogFeed.svelte`'s `read` prop, following it live while it grows, with the path shown and a note when the log is not there.
- The page `/ops/log?path=<path>`, which holds `LogViewer` set to the path, and says what to do when no path is given. Place it in `sitemenu.ts` under Ops so `locate` finds it; it need not be a menu entry.
- `LogLink.svelte`: a log path rendered as a link to `/ops/log?path=<encoded path>`, keeping the project the page is on as other links do, with the path as its text.
- `LogText.svelte` with a pure function in `$lib/logs.ts`: free text from flai (a status, a why, an error) rendered with each log path in it as a `LogLink` and the rest as text. A log path is an absolute, `~/`, or `.flai-cache/` path ending in `.log`, cut at whitespace, quotes, and backticks, without trailing punctuation.
- Tests for the route, the function, and each component; examples under `design/components/` for each new component as the components check asks.
- Document the page in the `## Views` table and the route in the `## API (S-0011)` table of `design/system/flaiover-dashboard.md` (the design-routes test reads both).
- Waits for nothing: the route's contract is fixed in T-3354, and the `REQUIRED_METHODS` entry is T-3354's, so it runs beside it.

## Done when

- The page shows a log's contents from a link, follows a growing log, and the route, function, and components have tests that pass with `flai test`.

## Notes
