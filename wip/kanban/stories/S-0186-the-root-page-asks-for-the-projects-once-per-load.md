---
id: S-0186
type: story
nature: improvement
title: The root page asks for the projects once per load
status: ready
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:32:17Z
transitions:
  - to: ready
    at: 2026-10-01T08:32:17Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/routes/+page.svelte, flaiover/src/routes/+layout.svelte, flaiover/src/lib/project.svelte.ts, flaiover/src/routes/api/projects]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0186 The root page asks for the projects once per load

## Goal

I-0033: loading the root page fetches `/api/projects` up to three times. The layout calls `projectState.refresh()` (`flaiover/src/routes/+layout.svelte` ~38), and the root page calls `projectState.refresh()` and `loadGlances()` on mount (`flaiover/src/routes/+page.svelte` ~52-57), each its own `fetch('/api/projects')` (`lib/project.svelte.ts` ~73). Each call gathers every project's glance on a 3 s timeout (`routes/api/projects/+server.ts` ~30-40), and the calls appear to be served one after another, so a slow project delays the page several times over.

## Acceptance criteria
- [ ] One load of the root page makes one `/api/projects` request: `projectState` shares a single in-flight request, holds the glances, and `loadGlances` reads from it
- [ ] Whether the dashboard serves concurrent `/api/projects` calls one after another (for example `hub.ask` per connection) is measured, and fixed if so, or the design says why not
- [ ] A test counts the requests on a root page load
- [ ] `design/system/flaiover-dashboard.md` says how the projects and their glances are loaded
- [ ] I-0033 is closed with what fixed it

## Tasks

## Notes
