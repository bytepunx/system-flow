---
id: T-0460
type: task
nature: improvement
title: Thread entries show the operator's on the right and agents' on the left, in different colours
status: done
parent: S-0127
owner: alex
created: 2026-09-26T07:41:46Z
updated: 2026-09-26T07:44:45Z
transitions:
  - to: ready
    at: 2026-09-26T07:41:49Z
    by: agent-S-0127
  - to: in-progress
    at: 2026-09-26T07:41:49Z
    by: agent-S-0127
  - to: done
    at: 2026-09-26T07:44:45Z
    by: agent-S-0127
stream: S-0127
tags: []
touches: [flaiover/src]
---
# T-0460 Thread entries show the operator's on the right and agents' on the left, in different colours

## Work

- `GET /api/threads` marks each entry `operator: true` when its author is the manifest's owner (`designer` when there is none), the same person flai counts as the designer for its inbox and writes threads as from the dashboard.
- `Threads.svelte` lays an operator entry out on the right in a primary-tinted bubble and an agent entry on the left in a neutral bubble, each still naming its author and time.
- Behavior tests: the endpoint marks entries from the owner and only those; the component places and colours the two kinds apart.
- `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` say how entries are told apart.

## Done when

- On any page with threads, an operator entry and an agent entry are told apart at a glance by side and colour.
- `scripts/flaiover-test.sh` passes.

## Notes
