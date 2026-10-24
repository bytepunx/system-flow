---
id: T-2101
type: task
nature: remediation
title: Close I-0201 with the reason naming the go-test tier's change
status: done
parent: S-0490
owner: alex
created: 2026-10-23T21:22:52Z
updated: 2026-10-24T00:04:45Z
transitions:
  - to: ready
    at: 2026-10-24T00:03:45Z
    by: agent-S-0490
  - to: in-progress
    at: 2026-10-24T00:03:47Z
    by: agent-S-0490
  - to: done
    at: 2026-10-24T00:04:45Z
    by: agent-S-0490
stream: S-0490
tags: [flai, issues]
touches: [design/issues/I-0201-flai-test-runs-go-test-with-short-so-a-real-git-test-a-task-adds-is-skipped-until-the-integration-tier-runs-in-the-merge-queue-s-batch.md, design/issues/summary.md]
after: [T-2099]
usage:
  source: log
  seconds: 58
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 1936
      cache_read: 752180
      cache_write: 15582
      cost: 0.3138
---
# T-2101 Close I-0201 with the reason naming the go-test tier's change

## Work

Waits for T-2099, since the reason names its change.

In the story's worktree, run `flai issue close I-0201 --reason "<reason>"`, the reason saying that the `go-test` tier of `system-flow.yaml` now runs `go test` without `-short` and with `-skip ^TestGateway` on the changed packages, so a real-git or monorepo test a task adds runs in `flai test` and in the close-out's `go-test` step while `integration` is left to the merge queue, and that `flai/internal/verify/select_test.go` guards the tier. flai rewrites the issue and `design/issues/summary.md`. Then tick the story's second criterion with `flai criteria tick S-0490 2`, and the first once T-2099's test is in and `flai test flai/internal/verify` has run the real-git tests.

## Done when

- I-0201's front matter reads `status: closed` with the reason, and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` is clean for the story.
- Both acceptance criteria of S-0490 are ticked.

## Notes

Issue commands write under the current checkout: run them in the story's worktree, not the main checkout.
