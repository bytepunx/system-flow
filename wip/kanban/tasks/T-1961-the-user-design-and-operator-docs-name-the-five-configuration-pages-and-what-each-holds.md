---
id: T-1961
type: task
nature: improvement
title: The user, design, and operator docs name the five Configuration pages and what each holds
status: done
parent: S-0459
owner: alex
created: 2026-10-21T07:03:25Z
updated: 2026-10-21T09:16:42Z
transitions:
  - to: ready
    at: 2026-10-21T09:12:15Z
    by: system-flow
  - to: in-progress
    at: 2026-10-21T09:12:16Z
    by: system-flow
  - to: done
    at: 2026-10-21T09:16:42Z
    by: agent-S-0459
stream: S-0459
tags: [dashboard]
touches: [docs/users/flaiover.md, design/system/flaiover-dashboard.md, docs/operators/settings.md]
after: [T-1960]
usage:
  source: log
  seconds: 266
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 83
      output: 33517
      cache_read: 4479150
      cache_write: 151689
      cost: 2.482
---
# T-1961 The user, design, and operator docs name the five Configuration pages and what each holds

## Work

- In `docs/users/flaiover.md` § Settings, replace the one page with the five Configuration pages. Name each page's sections and what each changes.
- In `design/system/flaiover-dashboard.md` § Views, give `/settings/:page` its row. Describe the redirect from `/settings` and its anchors, and `settingspages.ts`.
- In `docs/operators/settings.md`, each mention of the dashboard's Settings page names the Configuration page that holds the setting: Host actions under Features, Shared paths under File System, and so on.
- Waits for T-1960, so that the docs name the routes as built.

## Done when

No document in the three describes Settings as one page, and the markdown tier passes.

## Notes
