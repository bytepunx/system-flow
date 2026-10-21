---
id: T-1976
type: task
nature: remediation
title: Record the batch's scoped checks in an ADR, the design, and the docs, and close I-0174
status: done
parent: S-0464
owner: alex
created: 2026-10-21T07:30:46Z
updated: 2026-10-21T08:00:45Z
transitions:
  - to: ready
    at: 2026-10-21T07:59:45Z
    by: agent-S-0464
  - to: in-progress
    at: 2026-10-21T07:59:47Z
    by: agent-S-0464
  - to: done
    at: 2026-10-21T08:00:45Z
    by: agent-S-0464
stream: S-0464
tags: [flai]
touches: [design/adrs/README.md, design/system/flai-cli.md, design/system/devex.md, docs/users/flai.md, design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md, design/issues/summary.md, design/adrs/0183-the-merge-queue-s-batch-and-each-bisection-step-run-their-tiers-with-merge.md, template/CHANGELOG.md]
after: [T-1974, T-1975]
usage:
  source: log
  seconds: 58
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 10093
      cache_read: 1958648
      cache_write: 70772
      cost: 1.0381
---
# T-1976 Record the batch's scoped checks in an ADR, the design, and the docs, and close I-0174

## Work

The batch's verification no longer "sets no CLOSE_OUT_STORY" alone: it names its stories, and the checks inside its tiers scope to them. Record it where the design and the docs say otherwise.

- Write an ADR with `flai adr new`: the merge queue's batch and each bisection step run their tiers with `MERGE_QUEUE_STORIES` naming their stories, and the repository check, the markdown lint, and the integration checks scope to them as a close-out's scope to its story, recording no issue (refines ADR-0085 and ADR-0174). `flai task done` widens the touches with the new ADR's file.
- `design/system/flai-cli.md`: the merge queue pass's verification, and the `flai check` row's `--story`.
- `docs/users/flai.md`: `flai check --story` with several stories, and the merge queue section on what its tiers see.
- `design/system/devex.md`: the close-out row, where it says what the tiers are given.
- Close I-0174 with `flai issue close I-0174 --reason`, naming T-1974's and T-1975's fix and the test that reproduces the cause.

It waits for T-1974 and T-1975, since it records what they built and closes the issue only once the fix is in; it shares `docs/users/flai.md` with no earlier task.

## Done when

- The ADR is accepted and listed in `design/adrs/README.md`.
- `design/system/flai-cli.md` and `docs/users/flai.md` no longer say the batch's tiers are given no scope.
- I-0174 is closed and `design/issues/summary.md` no longer lists it.
- `flai check --strict` and the markdown lint pass on the changed paths.

## Notes
