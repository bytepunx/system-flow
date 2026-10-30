---
id: T-2263
type: task
nature: remediation
title: flai accept by hand while the queue can batch refuses a story whose queue tiers did not run at its head, lists it as a blocker in --dry-run, and lands it with an override flag
status: done
parent: S-0529
owner: alex
created: 2026-10-28T18:17:04Z
updated: 2026-10-30T00:52:20Z
transitions:
  - to: ready
    at: 2026-10-30T00:23:55Z
    by: agent-S-0529
  - to: in-progress
    at: 2026-10-30T00:23:56Z
    by: agent-S-0529
  - to: done
    at: 2026-10-30T00:52:20Z
    by: agent-S-0529
stream: S-0529
tags: [cli]
touches: [flai/cmd/accept.go, flai/cmd/verify.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/cmd/accept_queue_tiers_test.go, flai/cmd/move_refused_test.go, docs/users/flai-reference.md, docs/operators/settings.md, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go, flai/internal/preview/orchestrator.go]
after: [T-2262]
usage:
  source: log
  seconds: 1704
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 168
      output: 54576
      cache_read: 11432730
      cache_write: 480853
      cost: 6.2378
---
# T-2263 flai accept by hand while the queue can batch refuses a story whose queue tiers did not run at its head, lists it as a blocker in --dry-run, and lands it with an override flag

## Work

In `preview.AcceptWith` (`flai/internal/preview/accept.go`), add a blocker, such as `queue_tiers`, for an acceptance that is neither the orchestrator's nor the queue's landing (`--batch`), while the queue can batch (`queueBatches`, `flai/cmd/verify.go:412`, passed in through the accept options): it holds when the story's last report, read with `verify.LastReport`, leaves a queue tier unrun at the branch's head by T-2262's helper. Its text names the tiers and `flai verify <S-nnnn> --all`, and the flag below.

In `flai/cmd/accept.go`, refuse on that blocker before `mergeStoryBranch` changes anything, so main and the story stay as they were. Add one flag, such as `--without-queue-tiers`, that lands the story anyway; the acceptance records in the story's Notes, beside what the orchestrator's acceptance writes there, that the queue tiers did not run, naming them. Update the command's help, which says the queue's rules, and regenerate `docs/users/flai-reference.md` with `make flai-reference`, or `TestReferenceIsCurrent` fails.

`flai move <S> done`, and the host API's `item.move` and `accept.run`, reach `acceptItem`, so they refuse too; the host API needs no change for the refusal, and the override stays a CLI flag.

It waits for T-2262, whose helper it reads.

## Done when

- A test in the new `flai/cmd/accept_queue_tiers_test.go` reproduces I-0203: with the queue batching (a `FLAI_CONFIG` turning on the `merge-queue` host action, as `cmd/verify_test.go`'s `queueFixture` does, `accept_reviews` on, local acceptance), a story whose stored report left `integration` queued at its head is refused by `flai accept <S> --by alex`, main unchanged and the story in review.
- The same story with a report at its head that ran every tier is accepted; so is one under the override flag, with the Notes naming the tiers that did not run.
- With the queue off, a hand acceptance is unchanged.
- `--dry-run` lists the blocker, tested in `flai/internal/preview/accept_test.go`; `flai move <S> done` is refused, tested in `flai/cmd/move_refused_test.go`.
- The orchestrator's and the queue landing's existing tests pass unchanged.
- `flai test` on the changed paths passes, `TestReferenceIsCurrent` included.

## Notes
