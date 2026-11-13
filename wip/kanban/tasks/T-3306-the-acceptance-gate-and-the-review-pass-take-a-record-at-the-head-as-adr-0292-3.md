---
id: T-3306
type: task
nature: feature
title: "The acceptance gate and the review pass take a record at the head as ADR-0292 §3 widens it"
status: done
parent: S-0662
owner: alex
created: 2026-11-13T13:43:45Z
updated: 2026-11-13T14:27:13Z
transitions:
  - to: ready
    at: 2026-11-13T13:56:32Z
    by: agent-S-0662
  - to: in-progress
    at: 2026-11-13T13:56:32Z
    by: agent-S-0662
  - to: done
    at: 2026-11-13T14:27:13Z
    by: agent-S-0662
stream: S-0662
tags: [cli]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/preview/accept.go, flai/internal/mergequeue/eligible.go]
after: [T-2875]
usage:
  source: log
  seconds: 1841
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 72
      output: 389
      cache_read: 2925431
      cache_write: 196722
      cost: 1.1671
---
# T-3306 The acceptance gate and the review pass take a record at the head as ADR-0292 §3 widens it

## Work

ADR-0292 §3 widens "at the head": a record is at the head when its commit is the head, or when every commit from the record's commit to the head changes only paths under the manifest's `design`, `docs`, and `wip` folders, so the issues a close-out commits after the deploy do not undo the evidence. `design/system/kluster.md` says S-0662 builds it for the review pass and the acceptance gate.

- `kluster.AcceptGate` in `flai/internal/kluster/deploy.go` judges the deploy and the checks at the head with T-2875's judgement in place of `r.Commit != head`. Its callers, `preview.AcceptWith` in `flai/internal/preview/accept.go` and `mergequeue` in `flai/internal/mergequeue/eligible.go`, pass what it needs.
- The review pass in `flai/internal/serve/kluster.go` (the check near `r.Commit != head` in the look that queues a review deploy) queues no deploy for a story whose record passed at the head as widened, its checks too while `checks` is on.
- `KlusterEnv` and a checks run's gate stay at the exact head: ADR-0292 names only the review row, the acceptance gate, the move, and the close-out.

It waits for T-2875, whose judgement it uses.

## Done when

- `deploy_test.go` shows the gate passing at a head reached by a `design/issues` commit after the record's commit, and refusing at a head reached by a commit outside the three folders.
- `kluster_test.go` in `serve` shows the review pass queuing no deploy at a head reached only by a `wip` or `design` commit, and queuing one at a head reached by a code commit.
- `flai test` passes on the changed files.

## Notes

- Not a criterion of S-0662 by number; `design/system/kluster.md` § What happens when names it as S-0662's.
