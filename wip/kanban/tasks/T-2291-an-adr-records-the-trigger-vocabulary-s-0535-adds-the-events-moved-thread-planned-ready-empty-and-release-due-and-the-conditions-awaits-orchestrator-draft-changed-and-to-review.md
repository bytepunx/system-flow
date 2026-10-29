---
id: T-2291
type: task
nature: feature
title: "An ADR records the trigger vocabulary S-0535 adds: the events moved, thread, planned, ready-empty, and release-due, and the conditions awaits-orchestrator, draft-changed, and to-review"
status: done
parent: S-0535
owner: alex
created: 2026-10-28T21:37:18Z
updated: 2026-10-28T22:31:33Z
transitions:
  - to: ready
    at: 2026-10-28T22:31:28Z
    by: system-flow
  - to: in-progress
    at: 2026-10-28T22:31:29Z
    by: system-flow
  - to: done
    at: 2026-10-28T22:31:33Z
    by: agent-S-0535
stream: S-0535
tags: [flai, events, adr]
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 4
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 2
      output: 25
      cache_read: 206510
      cache_write: 1405
      cost: 0.0708
---
# T-2291 An ADR records the trigger vocabulary S-0535 adds: the events moved, thread, planned, ready-empty, and release-due, and the conditions awaits-orchestrator, draft-changed, and to-review

## Work

ADR-0207 §3 names the events `accepted`, `cancelled`, `ready`, `review`, `edited`, `reordered`, and `thread-opened`. ADR-0209 §1 writes the orchestrator's trigger with `moved`, `thread`, `planned`, `ready-empty`, and `{event: accepted, when: release-due}`. Both ADRs are accepted, so the words cannot change in place. This task writes one new ADR with `flai adr new`, refining ADR-0207 §3 and ADR-0209 §1, before any code reads the new names:

- the events `moved` (an item moved, with `to`), `thread` (a thread came to await an agent, with the thread, its item, and whom it awaits under ADR-0181), `planned` (a planner's run ended with items written, with the item planned and the IDs written), `ready-empty` (the ready column emptied while the in-progress limit has room, once per draining), and `release-due` (an acceptance changed what `flai release --evaluate` answers);
- whether `review` and `thread-opened` stay as they are or become `moved` with `to-review` and `thread`; recommend keeping both, so no role file breaks;
- `release-due` as an event of its own, as S-0535's criterion 3 says, rather than ADR-0209's condition on `accepted`;
- the trigger conditions `awaits-orchestrator` on `thread`, `draft-changed` on `edited`, and `to-review` on `moved`, and `to-backlog` qualifying `{event: moved}` as well as its calls, so that S-0536's trigger list passes `flai check`;
- where flai keeps the events it emits: a log under `.flai-cache` that `inbox.CatchUp` reads, as `edits.jsonl` and `overlaps.jsonl` are, and the `thread` event derived from the thread files.

It waits for nothing. Every other task follows the names it settles.

## Done when

- The ADR is in `design/adrs` with status accepted, refines ADR-0207 and ADR-0209, and is listed in `design/adrs/README.md`.
- It names each event and condition above with the fields it carries and who emits it.
- `flai check --strict` and the markdown lint are clean on it.

## Notes

The ADR's file name is not known until `flai adr new` allocates its number, so this task names only the index; `flai task done` widens its touches with the file.
