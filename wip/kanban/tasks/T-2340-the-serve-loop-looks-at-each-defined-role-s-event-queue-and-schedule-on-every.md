---
id: T-2340
type: task
nature: feature
title: The serve loop looks at each defined role's event queue and schedule on every look of the project
status: done
parent: S-0418
owner: alex
created: 2026-10-29T16:49:16Z
updated: 2026-10-29T17:20:27Z
transitions:
  - to: ready
    at: 2026-10-29T17:14:08Z
    by: agent-S-0418
  - to: in-progress
    at: 2026-10-29T17:14:09Z
    by: agent-S-0418
  - to: done
    at: 2026-10-29T17:20:27Z
    by: agent-S-0418
stream: S-0418
tags: [cli]
touches: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go]
after: [T-1743, T-1744]
usage:
  source: log
  seconds: 376
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 46
      output: 213
      cache_read: 2804323
      cache_write: 127516
      cost: 1.2104
---
# T-2340 The serve loop looks at each defined role's event queue and schedule on every look of the project

## Work

- Wire the role event queue (T-1743) and the role scheduler (T-1744) into the serve loop in `serve.go`, beside the orchestrator and the analysis scheduler: made per project, looked at on the same changes, every minute, and when an agent ends, so that reasons waiting for a run start the next once it ends.
- Have the loop look again when a role's queued reasons come due between looks, as the orchestrator's wake does.
- A loop test that starts a fixture role once per trigger kind, event and schedule, and records each trigger.

Waits for T-1743 and T-1744, which it wires together.

## Done when

`flai test flai/internal/serve` passes.

## Notes
