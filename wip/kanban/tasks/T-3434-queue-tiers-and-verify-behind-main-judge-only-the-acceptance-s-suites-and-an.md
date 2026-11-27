---
id: T-3434
type: task
nature: feature
title: queue_tiers and verify_behind_main judge only the acceptance's suites, and an acceptance waits for a suite its verify did not run
status: done
parent: S-0808
owner: alex
created: 2026-11-16T23:46:37Z
updated: 2026-11-27T23:24:04Z
transitions:
  - to: ready
    at: 2026-11-27T23:09:14Z
    by: agent-S-0808
  - to: in-progress
    at: 2026-11-27T23:09:14Z
    by: agent-S-0808
  - to: done
    at: 2026-11-27T23:24:04Z
    by: agent-S-0808
stream: S-0808
tags: [flai]
touches: [flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/cmd/accept.go, flai/cmd/accept_suites_test.go, docs/users/flai-reference.md]
after: [T-3433]
usage:
  source: log
  seconds: 890
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 102
      output: 525
      cache_read: 6004453
      cache_write: 309318
      cost: 2.4266
---
# T-3434 queue_tiers and verify_behind_main judge only the acceptance's suites, and an acceptance waits for a suite its verify did not run

## Work

ADR-0312 §3, criteria 1, 3, and 4. It waits for T-3433, whose `AcceptSuites` and suite-aware verify functions it calls.

- In `flai/internal/preview/accept.go`, `verifyBehindMainGate` and `queueTiersGate` read the acceptance's suites with `AcceptSuites` from `repo.Flow()` and pass them to `StaleTiers` and `QueueTiersNotRun`. When main's paths select no tier of those suites, `verify_behind_main` lists no blocker.
- Add the refusal of criterion 3: a story whose last verify did not run a suite the acceptance needs is refused. The words name the suite and the close-out that runs it. It gates the same acceptances `verify_behind_main` gates. Give it a blocker code of its own, listed by `flai accept --dry-run` and `accept.preview` like the others. On the baseline (`verify: all`) it never fires where no blocker fired before.
- In `flai/cmd/accept.go`, say in the help which suites the gates judge, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Write the tests for a workflow whose acceptance lists fewer suites in a new `flai/cmd/accept_suites_test.go`. Leave `accept_stale_verify_test.go`, `accept_queue_tiers_test.go`, and `accept_queue_test.go` unchanged, since criterion 4 rests on them passing as they are.

## Done when

- A test shows `verify_behind_main` sending back only tiers of the listed suites, and listing no blocker when main's paths select none of them.
- A test shows `queue_tiers` ignoring a `queue` tier of a suite the acceptance does not list.
- A test shows the refusal naming the missing suite and `flai verify S-nnnn --job`.
- `flai/cmd/accept_stale_verify_test.go`, `flai/cmd/accept_queue_tiers_test.go`, `flai/cmd/accept_queue_test.go`, and `flai/internal/preview/accept_test.go`'s existing cases pass with no change to them.
- `docs/users/flai-reference.md` matches the help.

## Notes
