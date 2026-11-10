---
id: S-0673
type: story
nature: improvement
title: after chains on planned stories hold ready stories for days while in-progress stays under two of five
status: ready
owner: alex
created: 2026-11-10T08:34:51Z
updated: 2026-11-10T08:54:23Z
transitions:
  - to: ready
    at: 2026-11-10T08:50:36Z
    by: alex
tags: [flai]
topics: [ready]
touches: [design/adrs, design/adrs/README.md, flai/internal/serve/readyempty.go, flai/internal/serve/readyempty_test.go, flai/internal/serve/agents_launch.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/inbox/events.go, flai/internal/mcpserver/server.go, flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go, flai/cmd/move.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/agent-narrative.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/system/workflow.md, design/issues/I-0285-after-chains-on-planned-stories-hold-ready-stories-for-days-while-in-progress.md, design/issues/summary.md]
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
          input: 33
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
    time_lost_per_cycle: 56h
    by: flai
    at: 2026-11-10T08:34:51Z
  value: 1039
  by: planner-S-0673
  at: 2026-11-10T08:48:47Z
forecast:
  duration: 1h
  delivery: 2026-11-10T10:57:00Z
  basis: "Its own forecast of 1h; 2nd in the pull order with an in-progress limit of 5, behind S-0680, S-0682, S-0684, S-0685, S-0687 and S-0674."
  by: flai
  at: 2026-11-10T08:54:23Z
finalized:
  by: alex
  at: 2026-11-10T08:50:30Z
---
# S-0673 after chains on planned stories hold ready stories for days while in-progress stays under two of five

## Goal

This story remediates [I-0285](../../../design/issues/I-0285-after-chains-on-planned-stories-hold-ready-stories-for-days-while-in-progress.md), "after chains on planned stories hold ready stories for days while in-progress stays under two of five", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0285 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0285 is closed with `flai issue close I-0285 --reason` saying what fixed it

## Tasks
- T-2928 An ADR records what I-0285's instances show and decides that ready-empty counts only ready stories not held on after and that a move to ready warns of an after hold
- T-2930 ready-empty counts only the ready stories not held on after, so held stories in ready no longer keep the orchestrator from promoting
- T-2931 A move of a story to ready warns when a story it names in after is not done
- T-2933 The user guide says what ready-empty counts and that a move to ready warns of an after hold, and I-0285 is closed

## Notes

The analyzer found I-0285: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0285. time_lost_per_cycle 56h carried over from I-0285's Impact section.

Strategic usage carried over from I-0285: analyzer 42s, 0.8156 USD.

### Planning

Planned by planner-S-0673 on 2026-11-10. The plan's thread is on S-0673.

**What the instances show.** The planner's proposal, for the story's agent to confirm or change in the ADR of the first task:

- The longest holds were on stories planned before ADR-0230 (2026-10-31) made `after` name only a story whose output it reads: S-0425 (326,792 s), S-0424 (133,349 s), S-0423, S-0422, and S-0420, all created 2026-10-16 in E-0024's chain, and S-0554 (74,959 s), created 2026-10-30. Of the stories that carried `after` and started since 2026-11-03, the longest hold is S-0579's 6,259 s, and `claims.weeks[]` gives W46 0 s under `after`. The planner's rule needs no change.
- Each long hold sat in ready because it was moved there by hand while a story it names in `after` was not done: S-0425 on 2026-10-30, S-0424 and S-0554 on 2026-11-01, by alex. `flai promote --candidates` already leaves such a story out, and `flai move` to ready says nothing of the hold; only the move to in-progress warns (`flai/internal/workitem/rules.go`).
- A story held on `after` in ready counts toward the ready column in the `ready-empty` event: `readyDrained` in `flai/internal/serve/readyempty.go` counts held stories, so with as many held stories in ready as the in-progress room, no `ready-empty` is told, the orchestrator is not started to promote, and lanes stand empty. This is how in-progress stayed at 1 to 4 of 5 while `held` reached 6 on 2026-11-01.

**Proposed remedy.** `ready-empty` counts only the ready stories not held on `after`, the hold that lasts until an acceptance; a story held on overlap alone still counts, since a trial merge may clear it at the next look. A move of a story to ready warns when a story it names in `after` is not done, as the move to in-progress warns of a hold. A warning, not a refusal: the move is the operator's.

**Touches.** No touches were declared, so `flai touches suggest` ran from `flai/internal/serve/readyempty.go` and `flai/internal/workitem/rules.go` (16 commits).

| Touch | Source | Why |
|-------|--------|-----|
| `design/adrs/` | layout | The first task writes a new ADR whose number and name are not known until `flai adr new` gives them; no task can name the file. Kept as a folder touch for that alone |
| `design/adrs/README.md` | co-change 19% | The generated ADR index |
| `flai/internal/serve/readyempty.go`, `readyempty_test.go`, `agents_launch.go` | layout | `readyDrained` and the count `look` passes it |
| `flai/internal/serve/orchestrate_reasons.go`, `orchestrate_reasons_test.go` | layout | The reason's words, "the ready column holds n" |
| `flai/internal/inbox/events.go`, `flai/internal/mcpserver/server.go` | layout, co-change 25% | The event's and the tools' descriptions of `ready-empty` |
| `flai/internal/workitem/rules.go`, `rules_test.go` | layout | `Move`'s warnings |
| `flai/cmd/move.go`, `docs/users/flai-reference.md` | co-change 12%, 19% | The command's help and the reference generated from it |
| `docs/users/flai.md`, `design/system/flai-cli.md`, `design/system/agent-narrative.md`, `design/system/strategic-agents.md`, `design/system/workflow.md` | design, co-change | Where `ready-empty` and the move's warnings are described |
| `design/issues/I-0285-…md`, `design/issues/summary.md` | design | Closing I-0285 |

Left out from `flai touches suggest`: `flai/internal/workitem/item.go`, `boardview.go`, `flai/cmd/items.go`, `board.go`, and `flai/internal/check/check.go`, which the change does not reach.

**Forecast.** `flai forecast` gave 37m once the touches were in: 101 s a unit over 151 large improvement stories, size 22 (2 criteria, 20 touches). Raised to 1h, because the size counts neither the ADR the first task writes nor the four tasks in three layers.

**Cost of delay.** `flai cod` gave 8,400 USD a week from the 56h input: W45's `after` hold at 150 USD an hour. The value is lowered to 1,039 USD a week, the rate since the planner's rule held: the stories with `after` started since 2026-11-03 were held 24,939 s (6.93 h) in all, and W46 0 s so far. The input is the operator's and stays as it is.

**After.** None: the story reads no other open story's output.
