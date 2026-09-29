---
id: T-0539
type: task
nature: feature
title: flai measures an agent's tokens and cost from its log
status: done
parent: S-0143
owner: alex
created: 2026-09-29T06:08:15Z
updated: 2026-09-29T06:11:04Z
transitions:
  - to: ready
    at: 2026-09-29T06:08:36Z
    by: agent-S-0143
  - to: in-progress
    at: 2026-09-29T06:08:36Z
    by: agent-S-0143
  - to: done
    at: 2026-09-29T06:11:04Z
    by: agent-S-0143
stream: S-0143
tags: []
touches: [flai/internal/usage]
---
# T-0539 flai measures an agent's tokens and cost from its log

## Work

- New package `flai/internal/usage` reads Claude Code stream-json logs: per session (its `session_id`), the newest `result` event's `modelUsage` is the reported total per model (cumulative across resumed runs of that session, as the logs show); assistant messages, deduplicated by message id, carry each call's timestamp and input, cache read, and cache write tokens.
- A total over a story's logs: the reported totals of each session, plus what came after a session's last result (a run still going or one that died), estimated.
- A window `[from, to)` of those logs: each session's reported totals apportioned by the share of its calls' input and cache tokens that fall in the window; cost marked estimated.
- An estimate prices tokens at the model's blended rate (reported cost over reported tokens) across the logs read, and leaves cost unset when the logs report none for that model.
- Working seconds: each run's first to last event, and for a window its overlap.

## Done when

- Behaviour tests over fixture logs cover a single run, a resumed session over two logs, a run with no result, a window, subagent messages of another model, and lines that are not JSON.
- `make test` and lint pass.

## Notes
