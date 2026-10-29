---
id: T-1743
type: task
nature: feature
title: serve queues a role whose trigger names an event through the orchestrator's reason queue, coalesced into one run per role and item
status: done
parent: S-0418
owner: alex
created: 2026-10-16T05:45:27Z
updated: 2026-10-29T17:11:27Z
transitions:
  - to: ready
    at: 2026-10-29T16:52:46Z
    by: agent-S-0418
  - to: in-progress
    at: 2026-10-29T16:52:47Z
    by: agent-S-0418
  - to: done
    at: 2026-10-29T17:11:27Z
    by: agent-S-0418
stream: S-0418
tags: [cli]
touches: [flai/internal/serve/role_events.go, flai/internal/serve/role_events_test.go, flai/internal/serve/orchestrate_reasons.go]
after: [T-1748]
usage:
  source: log
  seconds: 1120
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 81
      output: 33262
      cache_read: 6167550
      cache_write: 161449
      cost: 2.8801
---
# T-1743 serve queues a role whose trigger names an event through the orchestrator's reason queue, coalesced into one run per role and item

## Work

Queue a defined role whose trigger names an event through the reason queue S-0536 gave the orchestrator (`flai/internal/serve/orchestrate_reasons.go`: `orchestrateReason`, `reasonQueue`, `reasonsOf`), generalised where it must be so that a role a role file defines takes the same path. Leave `prime()` and `primeReasons()` as they are (MS-0232).

- A `roleEvents` per project, in `role_events.go`, with a `look(ctx)`. At each look, for each role a role file defines but the built-in ones, whose trigger names an event and whose host action is on, read the inbox's events since its last look as the role (`inbox.CatchUp` under a cursor of flai serve's per role, as `serveCursor` is the orchestrator's) and queue those its triggers match, under ADR-0207's and ADR-0211's events and the conditions `flai/internal/roles` implements. A thread awaits the role when it awaits the role's agent name, not only the orchestrator's.
- Keep a queue per project and role under the serve folder, beside `orchestrate-reasons.json`, so that it outlives one flai serve.
- Coalesce over `orchestration.coalesce`, the window S-0536 uses, since a role has no window of its own. A role that takes no item gets one run per batch; a role that takes items gets one run per item, with the reasons about that item, and reasons about an item of a kind it does not take are not queued.
- Start each batch through `startRole` (T-1748) with the trigger `<event> <ID>` per reason joined by semicolons. Reasons that arrive while a run goes start the next run once it ends. A start `startRole` refuses drops its reasons and logs the refusal; a start the budget holds keeps them queued.
- While the action is off, queue nothing, keep up with the events, and empty the queue.
- Test each event kind, a condition, the coalescing, an item role, a refusal, a held start, and the action off.

Waits for T-1748, whose `startRole` it calls. Runs beside T-1744: they share no path.

## Done when

`flai test flai/internal/serve` passes.

## Notes
