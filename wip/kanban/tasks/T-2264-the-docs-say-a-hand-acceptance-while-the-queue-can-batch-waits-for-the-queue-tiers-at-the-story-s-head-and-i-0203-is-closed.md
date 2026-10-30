---
id: T-2264
type: task
nature: remediation
title: The docs say a hand acceptance while the queue can batch waits for the queue tiers at the story's head, and I-0203 is closed
status: done
parent: S-0529
owner: alex
created: 2026-10-28T18:17:20Z
updated: 2026-10-30T00:57:46Z
transitions:
  - to: ready
    at: 2026-10-30T00:52:56Z
    by: agent-S-0529
  - to: in-progress
    at: 2026-10-30T00:52:57Z
    by: agent-S-0529
  - to: done
    at: 2026-10-30T00:57:46Z
    by: agent-S-0529
stream: S-0529
tags: [cli]
touches: [docs/users/flai.md, docs/operators/index.md, design/system/workflow.md, design/system/flai-cli.md, design/issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md, design/issues/summary.md]
after: [T-2261, T-2263]
usage:
  source: log
  seconds: 289
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 69
      output: 22460
      cache_read: 4705091
      cache_write: 197893
      cost: 2.5671
---
# T-2264 The docs say a hand acceptance while the queue can batch waits for the queue tiers at the story's head, and I-0203 is closed

## Work

Replace each statement that a hand acceptance runs no test tier, or that the operator should run `flai verify --all` first, with what T-2263 built, linking T-2261's ADR:

- `docs/users/flai.md`: § The merge queue, under Accept and release, and § Tiers left to the merge queue, under Verify a story before review.
- `docs/operators/index.md`: the merge-queue host action's section, its paragraph on what it changes in a close-out.
- `design/system/workflow.md`: § What waits on an acceptance (S-0296).
- `design/system/flai-cli.md`: the `flai accept` row, with the blocker and the override flag, and the `flai verify` row's `--all`, which now reuses the other tiers.

Then close the issue with `flai issue close I-0203 --reason` in the story's worktree, saying that a hand acceptance while the queue can batch is refused until the queue tiers ran at the story's head, with the ADR's number and the test that reproduces it.

It waits for T-2261, whose ADR it links, and T-2263, whose flag and blocker it names.

## Done when

- No document under `docs/` or `design/system/` says a hand acceptance runs no tier while the queue can batch.
- I-0203 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` and the markdown lint pass.

## Notes
