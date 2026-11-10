---
id: S-0702
type: story
nature: improvement
title: "flai check finds `issues.summary` outside the story at close-out"
status: backlog
owner: alex
created: 2026-11-10T08:55:37Z
updated: 2026-11-10T19:26:06Z
transitions: []
tags: [flai, issues, git]
topics: [git]
touches: [design/adrs, design/adrs/README.md, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/issues/generated.go, flai/internal/issues/generated_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, design/system/continuous-improvement.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, design/issues/I-0266-flai-check-finds-issues-summary-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 52
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 24
          output: 5841
          cache_read: 790028
          cache_write: 63601
          cost: 0.7013
    - kind: orchestrator
      seconds: 224
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 40
          output: 371
          cache_read: 3216816
          cache_write: 232329
          cost: 0.9428
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-11-10T10:55:38Z
  value: 12.5
  by: planner-S-0702
  at: 2026-11-10T11:07:46Z
forecast:
  duration: 60m
  delivery: 2026-11-11T11:08:00Z
  basis: "Its own forecast of 1h; 31st in the pull order with an in-progress limit of 3, behind S-0666, S-0667, S-0708, S-0688, S-0677, S-0679, S-0675, S-0676, S-0678, S-0674, S-0704, S-0681, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0686, S-0689, S-0690, S-0691, S-0692, S-0693, S-0694, S-0695, S-0696, S-0697, S-0698, S-0699, S-0700 and S-0701."
  by: flai
  at: 2026-11-10T19:26:06Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:10:11Z
---
# S-0702 flai check finds `issues.summary` outside the story at close-out

## Goal

This story remediates [I-0266](../../../design/issues/I-0266-flai-check-finds-issues-summary-outside-the-story-at-close-out.md), "flai check finds `issues.summary` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0266 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0266 is closed with `flai issue close I-0266 --reason` saying what fixed it

## Tasks
- T-3002 Decide in an ADR that a rebase which wrote a generated file again at a stop writes it again at the branch's head and commits it
- T-3005 flai stream sync writes the generated files again at the branch's head after a rebase that wrote one at a stop, and commits what changed
- T-3009 The merge queue's batch rebase writes the generated files again at each story's rebased head after it wrote one at a stop
- T-3011 The design and user docs say a rebase writes the generated files again at the branch's head
- T-3015 Close I-0266 with the reason the sync and the batch now write the generated files again at the head

## Notes

### Planning

The cause, found from I-0266's one instance and main's history. The instance is S-0576's close-out at 2026-11-07T00:11:02Z. Main's `design/issues/summary.md` lacked the row of I-0264, whose file was on main already. S-0576's close-out commit `2ccd15b26` added that row with I-0266's. I-0264 reached main from S-0520's branch:

- `66a374f21` closed I-0218 and rewrote `summary.md`.
- `44beb162d`, the next commit, added the I-0264 file and left `summary.md` alone.

`66a374f21`'s summary carries `updated: 2026-11-06T23:25:22Z`, later than the commit itself. So a sync's rebase stopped on `summary.md` at `66a374f21`, and `ContinueOverResolved` (`flai/internal/storygit/sync.go`) wrote it again. It wrote it from the issue files as the worktree was at that stop, as [ADR-0098](../../../design/adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md) says. `44beb162d` was not yet applied there, so the summary had no I-0264 row. `44beb162d` then replayed cleanly and left the branch's head with a summary its issue files do not render. Every worktree made from main after the acceptance found `issues.summary` outside its story, until a close-out wrote the summary again.

The ADR index (`design/adrs/README.md`) is the other generated file and is written again at a stop in the same way. A later commit of the branch that adds an ADR file leaves it stale too.

The explorer found two more paths that write issue files and can skip the summary, both only on a failure:

- a fold that fails after it wrote some files (`flai/internal/issues/fold.go`, `sync.go`'s `fold`);
- `recordOutside`'s early error return (`flai/cmd/check.go`).

Neither fits this instance's evidence, so neither is in scope. The fix below does not cover them: it runs only after a rebase that stopped.

Proposed fix, for the story's agent to settle in an ADR that refines ADR-0098 before building it:

- After a rebase that wrote a generated file again at a stop, write each generated file again at the branch's head.
- When what it renders differs from the head's, with the summary's `updated` line left out as the `issues.summary` check leaves it out, commit it on the branch.
- Do this in `flai stream sync`, which `flai task done`, the close-out, and `flai accept` all run, and in the merge queue's batch rebase, which calls `ContinueOverResolved` itself.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `design/adrs`, `design/adrs/README.md` | design: the new ADR and its generated index |
| `flai/internal/storygit/sync.go`, `sync_test.go` | layout: `Sync` and `ContinueOverResolved`, where the cause is, and the test that reproduces it |
| `flai/internal/issues/generated.go`, `generated_test.go` | layout: `Generated` writes the two files; the comparison that leaves the `updated` line out belongs beside it |
| `flai/internal/mergequeue/batch.go`, `batch_test.go` | layout: the batch's `rebase` settles stops with `ContinueOverResolved` |
| `design/system/continuous-improvement.md` | design: § Summary says how the sync writes the summary again |
| `design/system/flai-cli.md` | design: the `flai stream sync` row names what the sync writes again |
| `design/system/workflow.md` | design: § Branches and collisions names the files sync regenerates |
| `docs/users/flai.md` | design: "The issue summary" bullet says when the summary is written again |
| `design/issues/I-0266-…md`, `design/issues/summary.md` | declared by criterion 2 |

`flai touches suggest S-0702` listed nothing at first, since the story declared no touches. Run again on 2026-11-10 with the touches above, it listed only co-changes at 10% or less, such as `docs/operators/index.md` and `docs/users/flai-reference.md`. The fix changes no flag, setting, or convention, so none of them is added. The folder touch `design/adrs` stays because the ADR's number and slug are not known until it is written. ADR-0279 was taken by S-0673 after the first plan, so T-3002 now names `0280-…` as its guess, and `flai task done` widens the touches with the file it writes.

Figures:

- Forecast: 60m. `flai forecast` gives 28m: 103 s per unit of size times size 16, over 157 done large improvement stories. Raised because the fix needs real-git tests in both `storygit` and the merge queue, plus an ADR, and a change under `flai/` runs the integration and smoke tiers at close-out. S-0582, a smaller fix of the same kind with an ADR, took 47m. The orchestrator approved 60m on TH-1033.
- Cost of delay: 12.50 USD a week, as `flai cod` works it out from the operator's input set on TH-1013: 5m of time lost per 168h cycle at 150 USD an hour, one cycle a week. It stands unadjusted: the input already weighs that each stale summary costs one round of close-outs.

Layers, approved on TH-1033:

1. T-3002 (the ADR).
2. T-3005 (the sync) and T-3011 (the docs), together, with no file in common.
3. T-3009 (the merge queue), after T-3005, whose helper it calls.
4. T-3015 (close I-0266), last.

No `after`: S-0702 reads nothing another open story writes. S-0710 (`adr.index` outside the story) has a different cause, an ADR file left uncommitted on main. S-0709 and S-0712 also change `flai/internal/storygit/sync.go`, for `undoes_main`. The trial merge judges that overlap, and the pull order decides which goes first.
