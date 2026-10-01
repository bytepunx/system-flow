---
id: S-0186
type: story
nature: improvement
title: The root page asks for the projects once per load
status: review
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T09:28:44Z
transitions:
  - to: ready
    at: 2026-10-01T08:32:17Z
    by: alex
  - to: in-progress
    at: 2026-10-01T09:17:23Z
    by: agent-S-0186
  - to: review
    at: 2026-10-01T09:28:44Z
    by: agent-S-0186
tags: [dashboard]
touches: [flaiover/src/routes/+page.svelte, flaiover/src/routes/+layout.svelte, flaiover/src/lib/project.svelte.ts, flaiover/src/routes/api/projects, flaiover/src/lib/project.svelte.test.ts, flaiover/src/routes/overview.svelte.test.ts, design/system/flaiover-dashboard.md, design/issues/I-0033-the-root-page-asks-api-projects-twice-on-one-load-and-the-dashboard-serves-the-two-calls-serially.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 708
  models:
    - model: claude-opus-5-5
      input: 288
      output: 63825
      cache_read: 14407493
      cache_write: 279409
      cost: 6.0438
---
# S-0186 The root page asks for the projects once per load

## Goal

I-0033: loading the root page fetches `/api/projects` up to three times. The layout calls `projectState.refresh()` (`flaiover/src/routes/+layout.svelte` ~38), and the root page calls `projectState.refresh()` and `loadGlances()` on mount (`flaiover/src/routes/+page.svelte` ~52-57), each its own `fetch('/api/projects')` (`lib/project.svelte.ts` ~73). Each call gathers every project's glance on a 3 s timeout (`routes/api/projects/+server.ts` ~30-40), and the calls appear to be served one after another, so a slow project delays the page several times over.

## Acceptance criteria
- [x] One load of the root page makes one `/api/projects` request: `projectState` shares a single in-flight request, holds the glances, and `loadGlances` reads from it
- [x] Whether the dashboard serves concurrent `/api/projects` calls one after another (for example `hub.ask` per connection) is measured, and fixed if so, or the design says why not
- [x] A test counts the requests on a root page load
- [x] `design/system/flaiover-dashboard.md` says how the projects and their glances are loaded
- [x] I-0033 is closed with what fixed it

## Tasks
- T-0656 projectState shares one in-flight /api/projects request and holds the glances, and the root page reads them
- T-0657 Measure whether the dashboard serves concurrent /api/projects calls one after another, and fix it or say why not
- T-0658 The dashboard design says how the projects and their glances are loaded, and I-0033 is closed

## Notes

- Criterion 1: `loadGlances()` was removed rather than made to read from `projectState`. The root page's list reads `projectState.list`, which holds the glances, so there is nothing left for `loadGlances` to do.
- Criterion 2: the dashboard and flai serve concurrent calls side by side. The serial calls I-0033 measured were Chromium's HTTP cache lock, measured with a scratch Playwright script (narrative log, 09:23:06Z), and fixed by fetching with `cache: 'no-store'`.
