---
id: T-2146
type: task
nature: improvement
title: flai serve confirms, on its per-minute look while orchestrate is on, each orchestrator recommendation that waited confirm_after with no owner entry since, on a thread an agent opened and not marked for the operator alone, and logs it
status: done
parent: S-0502
owner: alex
created: 2026-10-26T04:45:44Z
updated: 2026-10-28T13:23:35Z
transitions:
  - to: ready
    at: 2026-10-28T13:05:54Z
    by: agent-S-0502
  - to: in-progress
    at: 2026-10-28T13:05:55Z
    by: agent-S-0502
  - to: done
    at: 2026-10-28T13:23:35Z
    by: agent-S-0502
stream: S-0502
tags: [flai]
touches: [flai/internal/serve/confirm.go, flai/internal/serve/confirm_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/internal/serve/activity.go, flai/internal/workitem/activity.go]
after: [T-2142, T-2143]
usage:
  source: log
  seconds: 1060
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 45
      output: 15079
      cache_read: 4371827
      cache_write: 151170
      cost: 2.23
---
# T-2146 flai serve confirms, on its per-minute look while orchestrate is on, each orchestrator recommendation that waited confirm_after with no owner entry since, on a thread an agent opened and not marked for the operator alone, and logs it

## Work

Points 2 and 5 of T-2141's ADR, in `flai/internal/serve`:

- `confirm.go`: a `confirmer` built per project beside `newOrchestrator` in `serve.go`, with a `look(ctx)` called in the per-minute block after `orch.look`, as `analyst.look` and `curation.look` are. Each look: return while the project's `orchestrate` host action is off, as `orchestrator.look` returns on `!cfg.Orchestrate`; read the manifest afresh, since the operator may change it while flai serve runs, and return while `ConfirmAfterDuration()` is zero; then, for each open thread whose `PendingRecommendation()` is by `orchestrator` (the orchestrator's run's author, ADR-0087), is not marked for the operator alone, has `at` plus `confirm_after` at or before now, has no entry by the manifest's `owner` since (`WroteSince`), and whose opener is an agent `flai serve` starts (`agent-S-nnnn`, `planner-S-nnnn`, `planner-E-nnnn`) rather than `flai serve`, `flai`, or the owner, and whose title is not the permission prompt's (`Allow … ?`, the digest's `permissionTitle`, moved or exported so both read one pattern): call `threads.ConfirmByTimeout(repo, id, "flai serve", now, waited)`, then `LogActivity` with kind `orchestrator`, the summary `Confirmed by timeout the recommendation on TH-nnnn after <duration> with no answer from the operator`, and the thread's story as its item, and record a journal entry by `flai serve` (`serve.confirm`) as `orchestrate.go` records its runs. The look keeps no state file: a confirmed recommendation is no longer pending, so a second look finds nothing. An error on one thread is logged and the look goes on to the next.
- `confirm_test.go`: a repository with a thread opened by `agent-S-0001` and a recommendation by `orchestrator` 49 hours old is confirmed with the ADR's line and logged in `wip/agents/orchestrator.md` with no seconds and no cost; nothing is confirmed when `confirm_after` is unset, when the recommendation is 47 hours old, when the owner wrote after it (the owner having opened the thread, and not), when it is marked for the operator alone, when the opener is `flai serve`, when the title is a permission prompt's, when the recommendation's author is not `orchestrator`, and when the `orchestrate` host action is off. A second look after a confirmation writes nothing more.
- `serve_test.go`: the confirmer is wired into the per-project loop and runs on the minute as the other looks do.

## Done when

- Each case above has a test, and the confirmation is in the thread, in `wip/agents/orchestrator.md`, and in the host's journal.
- The asking agent's story is started again on the confirmation as it is on the operator's, which the existing resume path gives once the thread reads `answered`; a test shows the thread reads `answered` after the look.
- `flai test` on the changed paths passes.

## Notes

Layer 3, waits for T-2142 (the permission) and T-2143 (`WroteSince`, `ConfirmByTimeout`, the mark). The docs task waits for it. S-0498's T-2131 adds `remind.go`'s look to the same block in `serve.go`; this story waits for S-0498, so the look lands beside it, not against it.
