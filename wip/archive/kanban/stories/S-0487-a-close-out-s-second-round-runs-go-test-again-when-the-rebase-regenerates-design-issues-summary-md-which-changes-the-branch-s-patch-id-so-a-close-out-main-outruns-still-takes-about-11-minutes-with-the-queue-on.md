---
id: S-0487
type: story
nature: improvement
title: A close-out's second round runs go-test again when the rebase regenerates design/issues/summary.md, which changes the branch's patch id, so a close-out main outruns still takes about 11 minutes with the queue on
status: done
owner: alex
created: 2026-10-23T17:43:16Z
updated: 2026-10-23T19:33:47Z
transitions:
  - to: ready
    at: 2026-10-23T18:08:42Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-23T18:08:48Z
    by: agent-S-0487
  - to: review
    at: 2026-10-23T18:33:10Z
    by: agent-S-0487
  - to: done
    at: 2026-10-23T19:33:47Z
    by: orchestrator
tags: [cli]
topics: [cli]
touches: [flai/internal/verify/paths.go, flai/internal/verify/paths_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/cmd/verify.go, design/adrs, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0196-a-close-out-s-second-round-runs-go-test-again-when-the-rebase-regenerates-design-issues-summary-md-which-changes-the-branch-s-patch-id-so-a-close-out-main-outruns-still-takes-about-11-minutes-with-the-queue-on.md, design/issues/summary.md, design/issues/I-0201-flai-test-runs-go-test-with-short-so-a-real-git-test-a-task-adds-is-skipped-until-the-integration-tier-runs-in-the-merge-queue-s-batch.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1475
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 2
      test_runs: 4
      hand_edits: 1
      work: 40
  models:
    - model: claude-opus-5-5
      input: 154
      output: 38643
      cache_read: 8623567
      cache_write: 358163
      cost: 4.8137
  strategic:
    - kind: orchestrator
      seconds: 2247
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 241
          output: 4035
          cache_read: 55441400
          cache_write: 57971
          cost: 12.207
        - model: claude-sonnet-5-5
          input: 29
          output: 948
          cache_read: 458396
          cache_write: 96551
          cost: 0.383
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-10-23T17:43:16Z
  value: 15
  by: planner-S-0487
  at: 2026-10-23T18:08:09Z
forecast:
  duration: 45m
  delivery: 2026-10-24T12:45:00Z
  basis: "flai forecast gave 23m (89 s per unit of size over 93 done improvement stories, size 15); raised to 45m, the median agent time of the three verify remediations it resembles, S-0374 (35m), S-0466 (44m), and S-0405 (77m), each a verify change with an ADR, docs, and an issue close; delivery is flai's 2026-10-24T12:22Z, 28th in the pull order, moved by the 22m added."
  by: planner-S-0487
  at: 2026-10-23T18:08:09Z
finalized:
  by: orchestrator
  at: 2026-10-23T18:08:31Z
---
# S-0487 A close-out's second round runs go-test again when the rebase regenerates design/issues/summary.md, which changes the branch's patch id, so a close-out main outruns still takes about 11 minutes with the queue on

## Goal

This story remediates [I-0196](../../../design/issues/I-0196-a-close-out-s-second-round-runs-go-test-again-when-the-rebase-regenerates-design-issues-summary-md-which-changes-the-branch-s-patch-id-so-a-close-out-main-outruns-still-takes-about-11-minutes-with-the-queue-on.md), "A close-out's second round runs go-test again when the rebase regenerates design/issues/summary.md, which changes the branch's patch id, so a close-out main outruns still takes about 11 minutes with the queue on". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0196 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0196 is closed with `flai issue close I-0196 --reason` saying what fixed it

## Tasks
- T-2078 flai verify's patch id leaves out the paths flai regenerates at sync, so a round after a rebase that regenerated design/issues/summary.md still reuses the tiers main's paths do not select
- T-2079 An ADR refining ADR-0156 records that the patch id leaves out the generated paths, and flai-cli.md and workflow.md say so
- T-2080 docs/users/flai.md says the patch id leaves out the generated paths, and the flai reference is regenerated
- T-2081 I-0196 is closed with the reason naming the patch id that leaves out the generated paths

## Notes

Cost of delay inputs set by flai from I-0196. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-23T05:04:30Z, 0.5 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0487 on 2026-10-23. The solution proposed, from I-0196's instance: `flai verify` records the patch id of the branch's own change without the generated paths, the files flai regenerates at a sync's rebase (`issues.Generated`, `design/issues/summary.md` alone today, ADR-0098). The regenerated summary is a function of the issue files, which stay in the diff, and of main's issue files, which already select tiers through main's new paths, so the round after such a rebase reuses the tiers main's new paths do not select, go-test among them when main changed only `design/issues`. ADR-0156's one patch id stays the gate; only what it covers narrows.

Touches. `flai touches suggest` had nothing to start from (the story declared no touches), so every touch comes from the code layout and the tasks:

- `flai/internal/verify/paths.go`, `paths_test.go`, `story.go`, `story_test.go`, `flai/cmd/verify.go`: layout, where `BranchPatch`, `storyRun.reusable`, their tests, and the help text live (T-2078).
- `design/adrs`: a folder touch, kept because T-2079 writes a new ADR with `flai adr new`, whose number and file name no task can name until it is allocated; `design/adrs` is in the manifest's `claims.shared`, so it holds no story. `design/adrs/README.md` is the index `flai adr accept` writes.
- `design/system/flai-cli.md`, `design/system/workflow.md`: design, the two places that describe the patch id and the reuse (T-2079).
- `docs/users/flai.md`, `docs/users/flai-reference.md`: design, the user section on resuming across a clean rebase and the reference `make flai-reference` regenerates from the help text (T-2080).
- `design/issues/I-0196-…md`, `design/issues/summary.md`: the second criterion, `flai issue close` rewrites both (T-2081).

Forecast. `flai forecast` gave 23m: 89 s per unit of size over 93 done improvement stories on claude-opus-5-5, size 15 (2 criteria, 13 touches), delivery 2026-10-24T12:22Z as 28th in the pull order. Raised to 45m: the three verify remediations this one resembles, each a verify change with an ADR, docs, and an issue close, took 35m (S-0374), 44m (S-0466), and 77m (S-0405) of agent time, and this one has two git-fixture tests to write and four tasks closed one after another by one agent. Delivery moved by the 22m added, to 2026-10-24T12:45Z.

Cost of delay. `flai cod` gives 15.00 USD a week from the inputs flai set (6m per 168h cycle at 150 USD an hour), and the value stands at that: the inputs are not mine to change. They understate it: the issue says every story that records or closes an issue meets this whenever main moves during its close-out, and each occurrence is about six minutes of a close-out that also outlives the 600000 ms call.

### Accepted by the orchestrator

- Verified: ce13fa6ddc8c5092a450c9fb5d56964bd994e5f5
- Batch: 20261023-192227
- Enqueued verified: 5bb483d7454ecc73682bbaf93bf5aca58679d11b
- At: 2026-10-23T19:33:47Z

Verdict: pass. The patch id now leaves out only design/issues/summary.md, the path issues.Generated names, a real-git test reproduces I-0196, and every changed file is within the story's touches.
- 1: flai/internal/verify/paths.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/paths_test.go
- 2: design/issues/I-0196-a-close-out-s-second-round-runs-go-test-again-when-the-rebase-regenerates-design-issues-summary-md-which-changes-the-branch-s-patch-id-so-a-close-out-main-outruns-still-takes-about-11-minutes-with-the-queue-on.md, design/issues/summary.md
