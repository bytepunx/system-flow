---
id: S-0809
type: story
nature: feature
title: The merge queue's batch runs the suites verify_at.batch lists over the merge of its stories
status: done
parent: E-0034
owner: alex
created: 2026-11-16T23:42:05Z
updated: 2026-11-21T17:20:21Z
transitions:
  - to: ready
    at: 2026-11-17T08:26:02Z
    by: alex
  - to: in-progress
    at: 2026-11-18T08:29:51Z
    by: agent-S-0809
  - to: review
    at: 2026-11-18T09:19:01Z
    by: agent-S-0809
  - to: in-progress
    at: 2026-11-18T13:04:38Z
    by: flai
  - to: review
    at: 2026-11-18T13:21:31Z
    by: agent-S-0809
  - to: done
    at: 2026-11-21T17:20:21Z
    by: alex
tags: [flai]
topics: [verify, workflow]
touches: [flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/mergequeue/state.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/workflow.md, design/system/project-manifest.md, docs/users/flai.md, docs/operators/index.md, flai/internal/mergequeue/land_test.go, design/issues/I-0407-flai-stream-sync-panics-while-continuing-a-rebase-whose-stop-the-agent-resolved.md, design/issues/summary.md]
after: [S-0806]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3963
  priced_by: harness
  turns:
    - day: 2026-11-18
      ceremony: 4
      work: 89
  models:
    - model: claude-opus-5-5
      input: 384
      output: 100123
      cache_read: 33647314
      cache_write: 1169505
      cost: 16.7225
  strategic:
    - kind: planner
      seconds: 34
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 42
          output: 4830
          cache_read: 4486202
          cache_write: 73387
          cost: 1.6333
    - kind: orchestrator
      seconds: 277
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 68
          output: 556
          cache_read: 3608019
          cache_write: 359935
          cost: 1.232
        - model: claude-sonnet-5-5
          input: 8
          output: 61
          cache_read: 86259
          cache_write: 42238
          cost: 0.1044
cost_of_delay:
  value: 64.2
  by: planner-S-0809
  at: 2026-11-16T23:48:51Z
forecast:
  duration: 41m
  delivery: 2026-11-18T09:03:00Z
  basis: "Its own forecast of 41m; 1st in the pull order with an in-progress limit of 3, behind S-0811, S-0834 and S-0842."
  by: flai
  at: 2026-11-18T08:13:57Z
finalized:
  by: orchestrator
  at: 2026-11-16T23:56:25Z
---
# S-0809 The merge queue's batch runs the suites verify_at.batch lists over the merge of its stories

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §6. Today a batch runs every tier its paths select.

## Acceptance criteria

- [x] A batch and its bisection run the tiers of the suites `verify_at.batch` lists that the batch's paths select, and every such tier under `all`.
- [x] `flai queue` and the batch's record name the suites a batch ran.
- [x] On the baseline a batch runs what it ran before this story.

## Tasks
- T-3451 verify.Tiers keeps the tiers of the suites it is given among those the paths select, and every one under all
- T-3453 The batch and each bisection step run the suites verify_at.batch lists, and the batch's record names them
- T-3456 flai queue names the suites the batch runs, in its listing and its --json
- T-3458 The design and the guides say the batch runs the suites verify_at.batch lists, and the baseline runs every one

## Notes

- After S-0806: the batch reads verify_at.batch, which S-0806 parses.

### Planning

- After: S-0806 writes the parse of `verify_at.batch` in `flai/internal/workflow`, which T-3453 reads. S-0806 waits for S-0805, which gives `verify.Tier` the `suite` that T-3451 filters on, so S-0805 is reached through S-0806 and is not named again.
- Touches came from layout and design. `flai touches suggest` had nothing to start from, since the story declared none. An explorer then read the code. The batch picks its tiers in `CheckoutVerifier` (`flai/internal/mergequeue/batch.go`), which calls `verify.Tiers` and `SelectStory` (`flai/internal/verify/story.go`, `select.go`). Bisection calls the same verifier (`bisect.go`). The batch's record is `Batch` in `state.go`. `flai serve` builds the queue's `Deps` in `flai/internal/serve/mergequeue.go`. `flai queue` prints the batch in `flai/cmd/queue.go`.
- Touches came from design for the documents: `design/system/flai-cli.md` § The merge queue pass, `design/system/workflow.md`, and `docs/users/flai.md`, each of which describes the batch's tiers.
- Touches came from co-change for three more. Once the code touches were written, the suggestion listed `docs/users/flai-reference.md` at 22% and `docs/operators/index.md` at 25%. The reference is regenerated from `flai queue`'s help. The operators guide and `design/system/project-manifest.md` both say the batch runs every tier its paths select. The rest of the co-change list names nothing this story changes.
- The touches name only files, with no folder kept. The `verified` and `step_verified` lines of `events.jsonl` are left out of the change, since their fields are in `design/system/metrics.md`, which changes only with an ADR. So the batch's record is `state.json`, and `metrics.md` is not touched.
- The forecast stands at flai's figure: 41m, delivered 2026-11-17T05:47:00Z. That is the median of 110 s per unit of size over 142 done feature stories on this model in the large band, times size 22 (3 criteria, 19 touches). No evidence was found to move it. The story's work is a filter, a parameter carried through, a field, a printed line, and five documents.
- The cost of delay value, 64.20 USD a week, is flai's figure, left unchanged. It is the story's share of E-0034's 1500.00 USD a week, which comes from the operator's input `time_lost_per_cycle: 10h`. The share is 41m of the 15h58m forecast over the epic's 20 open stories that have no inputs of their own.
- Layers: T-3451 first. T-3453 next, since it passes T-3451's suites. Then T-3456 and T-3458 together, since both read T-3453's field and behaviour and share no file.
- 2026-11-18T13:04:38Z: moved to in-progress: rebase and resolve job S-0809-20261118T130438Z: the rebase of story/S-0809 onto main stopped on conflicts in flai/internal/verify/select.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, which flai aborted, so the branch is as it was; its agent settles it in the worktree and syncs it with flai stream sync S-0809
