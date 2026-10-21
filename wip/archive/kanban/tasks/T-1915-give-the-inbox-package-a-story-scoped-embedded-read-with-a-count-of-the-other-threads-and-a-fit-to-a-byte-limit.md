---
id: T-1915
type: task
nature: improvement
title: Give the inbox package a story-scoped embedded read with a count of the other threads, and a fit to a byte limit
status: done
parent: S-0450
owner: alex
created: 2026-10-20T19:05:20Z
updated: 2026-10-20T19:15:26Z
transitions:
  - to: ready
    at: 2026-10-20T19:12:06Z
    by: agent-S-0450
  - to: in-progress
    at: 2026-10-20T19:12:07Z
    by: agent-S-0450
  - to: done
    at: 2026-10-20T19:15:26Z
    by: agent-S-0450
stream: S-0450
tags: [flai]
touches: [flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go]
after: [T-1914]
usage:
  source: log
  seconds: 199
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 17
      output: 5202
      cache_read: 982086
      cache_write: 49297
      cost: 0.6005
---
# T-1915 Give the inbox package a story-scoped embedded read with a count of the other threads, and a fit to a byte limit

## Work

Build in `flai/internal/inbox/inbox.go` what T-1914's ADR decides, for `task_done` and `story_start` to share:

- A read for an embedded inbox: with `Options.Story` set, the threads on the story and its tasks are listed as today, and the other threads that would have been listed are counted in a new field of `Inbox` (its JSON name as the ADR gives it) instead of listed. The standalone `inbox` tool keeps its answer.
- A fit that cuts an `Inbox` to a byte budget: trims `Changes` from the oldest, adding to `Omitted`, then shortens the listed threads' last entries, and counts what it cut, so the answer says what was left out.

Unit-test both in `flai/internal/inbox/inbox_test.go`, including the reproduction of I-0167: 33 threads awaiting the agent, none on the story, each with a long last entry, and 50 changes, whose embedded inbox lists none of those threads, counts 33, and fits the budget.

It waits for T-1914, because the ADR fixes the field names and the order of the cuts.

## Done when

- The embedded read lists only the story's threads and counts the rest; the standalone read is unchanged, and its tests pass as they were.
- The fit brings an inbox over the budget under it and records what it cut.
- The I-0167 reproduction test passes, and `flai test flai/internal/inbox` passes.

## Notes
