---
id: T-0592
type: task
nature: feature
title: "A story moves back a column: ready to backlog, in-progress to ready, cancelled to backlog"
status: done
parent: S-0167
owner: alex
created: 2026-09-29T23:33:30Z
updated: 2026-09-29T23:37:36Z
transitions:
  - to: ready
    at: 2026-09-29T23:33:34Z
    by: agent-S-0167
  - to: in-progress
    at: 2026-09-29T23:33:34Z
    by: agent-S-0167
  - to: done
    at: 2026-09-29T23:37:36Z
    by: agent-S-0167
stream: S-0167
tags: []
touches: [flai/internal/workitem, flai/internal/metrics, design/system/workflow.md, design/system/metrics.md, docs/users]
usage:
  source: log
  seconds: 242
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 61
      output: 19613
      cache_read: 5231130
      cache_write: 60882
      cost: 1.9258
---
# T-0592 A story moves back a column: ready to backlog, in-progress to ready, cancelled to backlog

## Work

- Allow `ready → backlog`, `in-progress → ready`, and `cancelled → backlog` in `flai/internal/workitem/rules.go`; `done` stays terminal.
- A reopened item is not completed: derive `completed` from the last transition when the item is `done` or `cancelled` now, not the first ever, in `flai/internal/metrics`.
- Record the rule change as an ADR (it changes a rule the tooling enforces and the metrics contract), and update `design/system/workflow.md`, `design/system/metrics.md`, and `docs/users/flai.md`.
- Behaviour tests for each new transition, for `done` staying terminal, and for the metrics of a reopened item.

## Done when

- `flai move` accepts the three backward moves and refuses any move out of `done`.
- A story cancelled and moved back to backlog has no `completed`, lead time, or cycle time until it closes again.
- `make test` and lint pass.

## Notes
