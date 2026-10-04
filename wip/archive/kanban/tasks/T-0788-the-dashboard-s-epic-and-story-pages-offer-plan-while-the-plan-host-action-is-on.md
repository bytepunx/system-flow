---
id: T-0788
type: task
nature: feature
title: The dashboard's epic and story pages offer Plan while the plan host action is on
status: done
parent: S-0208
owner: alex
created: 2026-10-04T00:44:07Z
updated: 2026-10-04T03:23:49Z
transitions:
  - to: ready
    at: 2026-10-04T00:44:32Z
    by: agent-S-0208
  - to: in-progress
    at: 2026-10-04T03:16:40Z
    by: agent-S-0208
  - to: done
    at: 2026-10-04T03:23:49Z
    by: agent-S-0208
stream: S-0208
tags: []
touches: [flaiover/src]
after: [T-0786]
usage:
  source: log
  seconds: 429
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 57
      output: 16662
      cache_read: 4171042
      cache_write: 72359
      cost: 1.6063
---
# T-0788 The dashboard's epic and story pages offer Plan while the plan host action is on

## Work

- On an epic's and a story's page, a Plan action that calls `plan.run` with the item's ID through the dashboard's API, shown only while the `plan` host action is on for the project, and reporting the refusal (a run already going) as other host actions do.

Waits for T-0786 for `plan.run`; runs beside T-0787, with no path in common.

## Done when

- [x] A server route and its test, and a component test showing the action shown and hidden with the host action
- [x] `npm run check` and the tests for what changed pass in `flaiover/`

## Notes
