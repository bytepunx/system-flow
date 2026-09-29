---
id: T-0593
type: task
nature: feature
title: flai board limit sets a column's WIP limit, and the dashboard can run it through flai on the host
status: done
parent: S-0167
owner: alex
created: 2026-09-29T23:33:31Z
updated: 2026-09-29T23:40:08Z
transitions:
  - to: ready
    at: 2026-09-29T23:33:34Z
    by: agent-S-0167
  - to: in-progress
    at: 2026-09-29T23:37:36Z
    by: agent-S-0167
  - to: done
    at: 2026-09-29T23:40:08Z
    by: agent-S-0167
stream: S-0167
tags: []
touches: [flai/cmd, flai/internal/workitem, flai/internal/hostapi, docs/users, design/conventions/work-management.md]
usage:
  source: log
  seconds: 152
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 41
      output: 13148
      cache_read: 3506715
      cache_write: 40812
      cost: 1.291
---
# T-0593 flai board limit sets a column's WIP limit, and the dashboard can run it through flai on the host

## Work

- Add `flai board limit <column> <n>` writing `wip_limits` in `wip/kanban/board.md`, the one file flai, flai serve, `flai check`, and the dashboard read the limits from; `0` means no limit. Only `ready`, `in-progress`, and `review` carry limits.
- Add the `board.limit` write to the host channel (`flai/internal/hostapi`), so the dashboard reaches it only through flai on the host (ADR-0029).
- Stop `design/conventions/work-management.md`'s project addition restating the numbers, so a change made from the board leaves nothing stale.
- Document the command in `docs/users/flai.md` and the reference.

## Done when

- `flai board limit in-progress 3` changes `board.md` and nothing else, and `flai board` shows the new limit.
- The host channel runs it and refuses a column without limits or a negative number.
- `make test` and lint pass.

## Notes
