---
id: S-0623
type: story
nature: improvement
title: flai guard refuses the orchestrator git tag --list and git worktree list, which only read
status: backlog
owner: alex
created: 2026-11-07T23:40:02Z
updated: 2026-11-09T21:56:35Z
transitions: []
tags: []
topics: [guard]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, docs/users/flai-reference.md, design/system/flai-cli.md, design/issues/I-0289-flai-guard-refuses-the-orchestrator-git-tag-list-and-git-worktree-list-which.md, design/issues/summary.md]
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
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3103
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-11-07T23:40:02Z
  value: 75
  by: planner-S-0623
  at: 2026-11-09T21:55:36Z
forecast:
  duration: 1h30m
  delivery: 2026-11-10T17:00:00Z
  basis: "flai forecast gave 3m on size 2 with no touches; raised for a guard change with role table tests, a regenerated reference, and the integration and smoke tiers a flai/ close-out runs"
  by: planner-S-0623
  at: 2026-11-09T21:55:36Z
finalized:
  by: alex
  at: 2026-11-09T21:56:35Z
---
# S-0623 flai guard refuses the orchestrator git tag --list and git worktree list, which only read

## Goal

This story remediates [I-0289](../../../design/issues/I-0289-flai-guard-refuses-the-orchestrator-git-tag-list-and-git-worktree-list-which.md), "flai guard refuses the orchestrator git tag --list and git worktree list, which only read", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0289 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0289 is closed with `flai issue close I-0289 --reason` saying what fixed it

## Tasks
- T-2766 The guard passes git tag, worktree, and branch in the forms that only read
- T-2767 Document the git forms flai guard passes as reads
- T-2768 Close I-0289 naming the guard's read forms of git tag, worktree, and branch

## Notes

The analyzer found I-0289: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0289. time_lost_per_cycle 30m carried over from I-0289's Impact section.

Strategic usage carried over from I-0289: analyzer 42s, 0.8156 USD.

### Planning

Proposed solution, from I-0289's instances: `flai/internal/guard/guard.go` judges a git command by its subcommand alone, against `gitReads`, which has no `tag`, `worktree`, or `branch`; and the orchestrator's git rule puts `tag` under `gitPublishes`. All three write in some forms and read in others, so the guard learns their read forms, as `flagReads` does for flai's commands: `git tag` that lists (no positional word, or `-l`/`--list`, with `--sort`, `--contains`, `--points-at`, `--merged`, `-n` and the like, and none of `-a`, `-s`, `-u`, `-m`, `-F`, `-d`, `-f`, `-e`, `--annotate`, `--sign`, `--delete`, `--force`); `git worktree list`; and `git branch` that lists (no positional word, or `--list`/`-l`, `--show-current`, `-a`, `-r`, `-v`, `--contains`, `--merged`, and none of `-d`, `-D`, `-m`, `-M`, `-c`, `-C`, `-f`, `-u`, `--set-upstream-to`, `--unset-upstream`, `--edit-description`). Every other form keeps its refusal, and the orchestrator's `git tag v1.2.3` keeps `publishesThrough`. The rule is shared by every role the guard judges git for (sub-agents, the planner, the orchestrator, the analyzer, the curator), so all of them gain the reads. `git branch` is beyond what the instances show; it is added because the analysis names it with the other two and it has the same read/write split.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/guard/guard.go` | design: I-0289's Impact names `gitReads` and `gitPublishes` there |
| `flai/internal/guard/guard_test.go` | co-change 64% with `guard.go`; the orchestrator's git tests are there |
| `flai/cmd/guard.go` | co-change 51%; its help lists the git reads |
| `docs/users/flai-reference.md` | co-change 34%; generated from that help by `make flai-reference` |
| `design/system/flai-cli.md` | design: its `flai guard` row lists the git reads |
| `design/issues/I-0289-...md`, `design/issues/summary.md` | declared by criterion 2, `flai issue close` |

No folder touch. `flai/cmd/guard_test.go` (co-change 34%) is left out: it is needed only if a test there pins the help text, which the agent widens to when it finds one.

Forecast: `flai forecast` gave 3m on size 2 with no touches. Adjusted to 1h30m: a guard change with table tests across roles, a regenerated reference, and a close-out that runs the integration and smoke tiers for a story under `flai/`. Delivery set a working day out, behind the three stories ahead of it in the pull order.

Cost of delay: 75 USD a week, as `flai cod` computes it from the 30m `time_lost_per_cycle` the issue carries, at 150 USD an hour, one cycle a week. It stands: the 22 refusals over 2026-10-25 to 2026-11-07 are about one a day, each costing the orchestrator a retry or a detour, which 30m a week fits.
