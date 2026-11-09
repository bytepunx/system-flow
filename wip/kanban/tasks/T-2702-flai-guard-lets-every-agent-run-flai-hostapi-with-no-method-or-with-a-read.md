---
id: T-2702
type: task
nature: improvement
title: flai guard lets every agent run flai hostapi with no method or with a read method, and refuses the write methods
status: done
parent: S-0612
owner: alex
created: 2026-11-09T07:56:45Z
updated: 2026-11-09T08:27:55Z
transitions:
  - to: ready
    at: 2026-11-09T08:05:36Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T08:05:36Z
    by: system-flow
  - to: done
    at: 2026-11-09T08:27:55Z
    by: agent-S-0612
stream: S-0612
tags: [flai, guard]
touches: [flai/internal/guard/guard.go, flai/internal/guard/hostapi_test.go, flai/cmd/guard_test.go]
usage:
  source: log
  seconds: 1339
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 80
      output: 367
      cache_read: 10553458
      cache_write: 371849
      cost: 3.9719
---
# T-2702 flai guard lets every agent run flai hostapi with no method or with a read method, and refuses the write methods

## Work

I-0307: `flai guard` refuses a story's sub-agents `flai hostapi`, since `hostapi` is in neither `cliReads` nor `flagReads` in `flai/internal/guard/guard.go`, so the `flai hostapi --timing <method>` that `docs/users/flai.md` § Request timing names cannot be run by a task sub-agent. S-0596's three exploration sub-agents built in-process Go programs instead, about 20 minutes each.

Proposed fix, from the issue's instance: make `flai hostapi` a read in the guard's form, as `flai role` is (S-0568, `readsRoles`). Add `hostapi` to `flagReads` with a `form` that allows it with no method (it lists the methods) and with a method in a new exported list of the host API's read methods, `HostAPIReads` (name it as fits), and refuses every other method with a reason that names the method and says the write methods are the story's agent's (ADR-0060). `--timing` and `--json` change nothing. The guard cannot import `flai/internal/hostapi`, which imports it, so the list lives in the guard; the next task ties it to the method table.

Leave out of the list any method that writes anything, a cursor or kept state under `.flai-cache` included, such as a method that marks what the designer has seen: read each candidate's code in `flai/internal/hostapi` before listing it. Because `reads` is shared, the planner's, the orchestrator's, and the analyzer's own calls gain the reads too, as they did `flai role list`.

Tests, in a new `flai/internal/guard/hostapi_test.go`, reproduce I-0307: a task sub-agent's `flai hostapi --timing board.get '{"all":true}' >/dev/null` and a read-only sub-agent's `flai hostapi --timing inbox.designer 2>&1 >/dev/null` (the documented example) pass; `flai hostapi` alone passes; a write method, such as `item.move`, is refused to every sub-agent with the reason. Add the case to `flai/cmd/guard_test.go` as S-0568 did for `flai role`.

## Done when

- A task sub-agent and a read-only sub-agent may run `flai hostapi` alone and with any listed read method, with or without `--timing`.
- Every method not on the list is refused to every sub-agent with a reason naming it.
- The tests above fail on the guard before the change and pass after it, and `flai test flai/internal/guard flai/cmd/guard_test.go` passes.

## Notes
