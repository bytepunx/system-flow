---
id: T-0604
type: task
nature: feature
title: Record the stop in an ADR, the design, and the docs
status: done
parent: S-0170
owner: alex
created: 2026-09-30T00:48:33Z
updated: 2026-09-30T01:02:31Z
transitions:
  - to: ready
    at: 2026-09-30T00:48:49Z
    by: agent-S-0170
  - to: in-progress
    at: 2026-09-30T00:59:21Z
    by: agent-S-0170
  - to: done
    at: 2026-09-30T01:02:31Z
    by: agent-S-0170
stream: S-0170
tags: []
touches: [design, docs]
usage:
  source: log
  seconds: 190
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 58
      output: 21238
      cache_read: 5231864
      cache_write: 63927
      cost: 1.9828
---

# T-0604 Record the stop in an ADR, the design, and the docs

## Work

- An ADR: the agent host action lets the dashboard stop a story's agent.
- `design/system/flaiover-dashboard.md` and `design/system/flai-cli.md` (or wherever the agent host action's methods are described) name the stop.
- `docs/users/flai.md`, the generated reference, and `docs/operators/index.md`'s agent host action section say what Stop does.

## Done when

- `flai check --strict` passes and the reference test is green.

## Notes
