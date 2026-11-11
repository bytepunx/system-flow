---
id: S-0674
type: story
nature: improvement
title: Declared touches match what a story changes in 16 percent of completed stories, so claims mislead the pull hold, the trial merge, and the conversations
status: ready
owner: alex
created: 2026-11-10T08:34:52Z
updated: 2026-11-11T01:41:10Z
transitions:
  - to: ready
    at: 2026-11-10T08:50:34Z
    by: alex
tags: [flai]
topics: [metrics, claims]
touches: [design/adrs, flai/internal/storygit/committed.go, flai/internal/storygit/committed_test.go, flai/internal/metrics/claims.go, flai/internal/metrics/claims_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, design/system/metrics.md, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/itemedit/widen.go, flai/internal/itemedit/widen_test.go, docs/users/flai.md, design/system/flai-cli.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0286-declared-touches-match-what-a-story-changes-in-16-percent-of-completed-stories.md, design/issues/summary.md]
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
      seconds: 67
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 73
          output: 14114
          cache_read: 4735506
          cache_write: 105255
          cost: 1.9583
    - kind: orchestrator
      seconds: 58
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 75
          cache_read: 289596
          cache_write: 14309
          cost: 0.0827
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13855
          cost: 0.278
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4h
    by: flai
    at: 2026-11-10T08:34:52Z
  value: 600
  by: planner-S-0674
  at: 2026-11-10T08:49:36Z
forecast:
  duration: 2h30m
  delivery: 2026-11-11T05:05:00Z
  basis: "Its own forecast of 2h30m; 3rd in the pull order with an in-progress limit of 3, behind S-0676 and S-0678."
  by: flai
  at: 2026-11-11T01:41:10Z
finalized:
  by: alex
  at: 2026-11-10T08:50:26Z
---
# S-0674 Declared touches match what a story changes in 16 percent of completed stories, so claims mislead the pull hold, the trial merge, and the conversations

## Goal

This story remediates [I-0286](../../../design/issues/I-0286-declared-touches-match-what-a-story-changes-in-16-percent-of-completed-stories.md), "Declared touches match what a story changes in 16 percent of completed stories, so claims mislead the pull hold, the trial merge, and the conversations", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0286 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0286 is closed with `flai issue close I-0286 --reason` saying what fixed it

## Tasks
- T-2936 ADR: the touches drift leaves out the folders flai writes, and flai task done narrows a closing task's unchanged touches
- T-2939 The touches drift leaves out the digest, analysis, issues, and wip folders on both sides, with a test that reproduces I-0286
- T-2940 flai task done narrows the closing task's touches, and the story's, to the files the branch changes
- T-2941 Measure exact_share again and close I-0286 with what fixed it

## Notes

The analyzer found I-0286: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0286. time_lost_per_cycle 4h carried over from I-0286's Impact section.

Strategic usage carried over from I-0286: analyzer 42s, 0.8156 USD.

### Planning

#### What the instances show

`flai stats --json` on 2026-11-10, `claims.drift[]` of the 115 stories from S-0560 on:

| Measure | Stories |
|---------|---------|
| exact as flai counts it today | 25 |
| with a file `outside` the touches | 80 |
| every `outside` file under `design/digest/`, `design/issues/`, or `design/analysis/` | 80 |
| exact once those folders and `wip/` are left out of both sides | 97 |
| still inexact: 1 to 4 predicted files the branch never changed | 18 |

The `outside` half is attribution, not missed claims. `storygit.Committed` (`flai/internal/storygit/committed.go`) gives a story every file of each commit whose subject names it in brackets. Three of those commits are flai's bookkeeping. `chore: [S-nnnn] accept and archive` sweeps the main checkout's uncommitted digest entries into the acceptance commit (ADR-0248). `chore: [S-nnnn] create story` writes the issue the story came from. `docs: [S-nnnn] record I-nnnn` is the close-out recording a finding outside the story. flai writes all three folders, as it writes `wip/`, which `Committed` already leaves out (S-0205). The `unchanged` half is partly the same mistake: a touch under `wip/` is never matched, because its commits are left out. The rest is real. A task's predicted file that its work never changed stays claimed, since `flai task done` only widens (ADR-0144).

#### Proposed remedy

1. The touches drift leaves out, on both sides, the folders flai writes: `wip/`, and the `digest`, `analysis`, and `issues` folders of the design folder.
2. `flai task done` also narrows. It drops from the closing task's touches each file the branch does not change against main. It drops the same file from the story's touches when no other open task names it.

An ADR refining ADR-0144 and the S-0205 definition in `design/system/metrics.md` § Claims and touches records both. The metrics contract changes only with an ADR.

#### Tasks and layers

| Layer | Task | Waits for |
|-------|------|-----------|
| 1 | T-2936 the ADR | none |
| 2 | T-2939 the drift leaves out flai's folders, with the reproducing test | T-2936, whose decision it builds |
| 2 | T-2940 task done narrows | T-2936, whose decision it builds |
| 3 | T-2941 measure and close I-0286 | T-2939 and T-2940, whose effect it measures |

T-2939 and T-2940 share no path and run together.

#### Touches

| Touch | Source | Why |
|-------|--------|-----|
| `design/adrs` | design | The new ADR, whose file name no task can give before `flai adr new` numbers it. The folder is in `claims.shared`, so it holds no story. Kept as a folder touch for that reason. |
| `flai/internal/storygit/committed.go`, `committed_test.go` | layout, co-change | Where a story's committed files are read (50% co-change for the test). |
| `flai/internal/metrics/claims.go`, `claims_test.go` | layout, co-change | `touchDrift` reads the touches side (60% co-change for the test). |
| `flai/internal/metrics/metrics.go` | co-change | `Options` documents `Commits` (40% co-change). |
| `flai/internal/statsread/statsread.go` | co-change | Hands `Commits` and the layout to the metrics (30% co-change). `statsread_test.go` (40%) is left out until a task changes it. Task done widens it if one does. |
| `design/system/metrics.md` | design | § Claims and touches defines a story's commits and touches. |
| `flai/internal/taskdone/taskdone.go`, `taskdone_test.go` | layout | The touches step, `touches()` and `branchChanges()`. |
| `flai/internal/itemedit/widen.go`, `widen_test.go` | layout | `Widen` and `WidenStory`. The narrowing sits beside them. |
| `docs/users/flai.md`, `design/system/flai-cli.md` | design | § Closing a task, step 6, and the `flai task done` row. Both are in `claims.shared`. |
| `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` | design | The task-done rule says what it does to touches. Conventions land in `design/` and then `template/` in one story. |
| `design/issues/I-0286-...md`, `design/issues/summary.md` | declared by criterion 2 | `flai issue close` writes them. They are in `claims.shared`. |

#### Figures

- **Forecast: 2h30m, delivery 2026-11-10T16:30Z.** `flai forecast` gave 36m once the 19 touches were written (101 s per unit times size 21, over 151 large improvement stories). Raised because size counts touches but not the work behind them: an ADR, two Go changes with tests, docs across seven files, a measured issue close, and close-out tiers that run integration and smoke for a story touching `flai/` and `template/`. The delivery keeps flai's start, 16th in the pull order, and adds the longer duration.
- **Cost of delay: 600 USD a week.** `flai cod` gives this: 4h lost per 168h cycle at 150 USD an hour, from the operator's `time_lost_per_cycle`. The value stands as flai works it out. The instances suggest that the issue's 4h, the send-backs at acceptance, comes less from wrong claims than the issue assumed. The `outside` files were never missed claims. The plan's thread raises this for the operator, who owns the input.

#### after

None. The story reads no output of another open story.
