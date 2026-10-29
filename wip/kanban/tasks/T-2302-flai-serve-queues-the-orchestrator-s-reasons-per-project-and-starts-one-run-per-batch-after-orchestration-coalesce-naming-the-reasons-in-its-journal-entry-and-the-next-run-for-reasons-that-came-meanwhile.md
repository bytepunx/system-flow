---
id: T-2302
type: task
nature: feature
title: flai serve queues the orchestrator's reasons per project and starts one run per batch after orchestration.coalesce, naming the reasons in its journal entry, and the next run for reasons that came meanwhile
status: in-progress
parent: S-0536
owner: alex
created: 2026-10-28T21:45:42Z
updated: 2026-10-29T07:25:42Z
transitions:
  - to: ready
    at: 2026-10-29T07:25:41Z
    by: agent-S-0536
  - to: in-progress
    at: 2026-10-29T07:25:42Z
    by: agent-S-0536
stream: S-0536
tags: [flai, serve, orchestrator]
touches: [flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/orchestrate_hold_test.go, flai/internal/serve/serve.go, flai/internal/harness/harness.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go]
after: [T-2301]
---
# T-2302 flai serve queues the orchestrator's reasons per project and starts one run per batch after orchestration.coalesce, naming the reasons in its journal entry, and the next run for reasons that came meanwhile

## Work

ADR-0209 §2, §7, criterion 2. Today `orchestrator.look` in `flai/internal/serve/orchestrate.go` starts a run whenever the `orchestrate` action is on and none is going, and again a minute after a failure (`orchestrateRetry`). Replace that with a reason queue per project, in a new `flai/internal/serve/orchestrate_reasons.go`:

- At each look, read the events since the queue's cursor (the work item changes, the threads, and the event log S-0535 added, as `inbox.CatchUp` gives them) and queue each that a trigger of the orchestrator's role file matches, with its item or thread. `orchestration.schedule`'s ticks are reasons too. A file change with no event, a task move, and a replan's forecast edits match nothing.
- Keep the queue in the serve folder beside `agents.json`, so that a restart of flai serve loses nothing and the MCP server can read it (the end task).
- Start one run no sooner than `orchestration.coalesce` after the first queued reason, with every reason queued meanwhile, and move them from the queue to the run's record in `agents.json`. Reasons that come while a run is going stay queued and start the next run when it ends.
- With the action off, queue nothing; with the run held stopped (`flai serve orchestrate stop`), keep queuing and start nothing, as today. ADR-0087's guard, permissions, activity log, and `orchestrate` journal action stay.
- The run's start entry in the journal names its reasons. Pass them to the start through a new `Reasons` field of `harness.Request` in `flai/internal/harness/harness.go`; the prompt that lists them is the prompt task's.
- An `asked` reason is the operator's `flai serve orchestrate start` and the dashboard's start, which queue a reason and start at once.

Waits for T-2301, which gives `orchestration.coalesce` and `orchestration.schedule`.

## Done when

- A test in `orchestrate_reasons_test.go` drives three matching events inside one window into one run whose record and journal entry name all three, and a fourth event during that run into the next run once it ends.
- Tests show a non-matching change starts nothing, the queue survives a restart, and the action off queues nothing.
- `orchestrate_test.go`'s tests of the old loop are rewritten to the queue, and `flai test` passes on the changed paths.

## Notes
