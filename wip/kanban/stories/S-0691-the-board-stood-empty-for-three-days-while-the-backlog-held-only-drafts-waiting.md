---
id: S-0691
type: story
nature: improvement
title: The board stood empty for three days while the backlog held only drafts waiting on operator decisions nobody asked for again
status: review
owner: alex
created: 2026-11-10T08:49:02Z
updated: 2026-11-11T08:00:48Z
transitions:
  - to: ready
    at: 2026-11-11T01:44:32Z
    by: alex
  - to: in-progress
    at: 2026-11-11T05:06:07Z
    by: agent-S-0691
  - to: review
    at: 2026-11-11T08:00:48Z
    by: agent-S-0691
blocked:
  - from: 2026-11-11T05:26:42Z
    until: 2026-11-11T07:42:24Z
    reason: pause for urgent fix
tags: [flai]
topics: [ready, threads]
touches: [design/adrs, design/adrs/README.md, flai/internal/inbox/starved.go, flai/internal/inbox/starved_test.go, flai/internal/serve/starved.go, flai/internal/serve/starved_test.go, flai/internal/serve/agents_launch.go, design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md, design/issues/I-0288-the-board-stood-empty-for-three-days-while-the-backlog-held-only-drafts-waiting.md, design/issues/summary.md, design/issues/I-0347-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md]
after: [S-0673]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3533
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 7
      work: 66
  models:
    - model: claude-opus-5-5
      input: 264
      output: 94764
      cache_read: 15680316
      cache_write: 623520
      cost: 9.4029
  strategic:
    - kind: orchestrator
      seconds: 68
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 14
          output: 125
          cache_read: 741371
          cache_write: 76384
          cost: 0.2255
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
    time_lost_per_cycle: 1h
    by: orchestrator
    at: 2026-11-10T10:55:47Z
  value: 150
  by: planner-S-0691
  at: 2026-11-10T10:56:15Z
forecast:
  duration: 1h
  delivery: 2026-11-11T03:59:00Z
  basis: "Its own forecast of 1h; 2nd in the pull order with an in-progress limit of 3, behind S-0690, S-0692 and S-0693."
  by: flai
  at: 2026-11-11T02:27:14Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:04:48Z
---
# S-0691 The board stood empty for three days while the backlog held only drafts waiting on operator decisions nobody asked for again

## Goal

This story remediates [I-0288](../../../design/issues/I-0288-the-board-stood-empty-for-three-days-while-the-backlog-held-only-drafts-waiting.md), "The board stood empty for three days while the backlog held only drafts waiting on operator decisions nobody asked for again", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0288 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0288 is closed with `flai issue close I-0288 --reason` saying what fixed it

## Tasks
- T-3043 An ADR records what I-0288's instance shows and decides that flai serve asks the operator on one thread when nothing can be made ready
- T-3046 inbox.Starved lists each backlog story that cannot be promoted and the operator's decision it waits on
- T-3050 flai serve opens one thread to the operator on the board when nothing can be made ready, updates it when the list changes, and resolves it once a story can start
- T-3056 The workflow and strategic-agents design and the user guide say that flai serve asks the operator on one thread when nothing can be made ready
- T-3057 I-0288 is closed with the reason that flai serve asks the operator on one thread when nothing can be made ready

## Notes

The analyzer found I-0288: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0288: analyzer 42s, 0.8156 USD.

### Planning

Planned by planner-S-0691 on 2026-11-10. The plan's thread is TH-1042, on S-0691.

**What the instance shows.** These are the planner's findings. The story's agent confirms or changes them in the ADR of the first task.

- From 2026-11-03 to 2026-11-05, ready, in progress, and review were all 0. Every backlog story was a draft, or was held by `after` behind a draft. The drafts S-0406 and S-0412 waited on TH-0532 and TH-0533, which awaited the operator with one entry each from 2026-10-20.
- `ready-empty` fired as designed (`flai/internal/serve/readyempty.go`), and the orchestrator was started for it. Its `ready` reason pack (`loadReady`, `flai/internal/context/reasons_figures.go`) holds the board and `flai promote --candidates`, and nothing names a draft's open thread. Its run of 2026-11-02 promoted nothing and said nothing new to the operator (`wip/agents/orchestrator.md`).
- The reminder of ADR-0205 (`flai/internal/serve/remind.go`, with `threads.remind_after: 4h` set since S-0498) would have reminded the operator of both threads, but only through a channel set with `flai serve notify set`. Without one it only logs. TH-0660 says no channel was set.
- So nothing in the repository or on the dashboard told the operator that the board had stopped, or which of their decisions would start it again.

**Proposed remedy.** When the ready column holds no story that can start and the in-progress limit has room, and `flai promote --candidates` lists no candidate, flai serve opens one thread to the operator on `wip/kanban/board.md`. The thread lists each backlog story that cannot be promoted and the operator's decision it waits on: a draft to finalize, with each open thread on it that awaits the operator; a story held by `after`, with the story at the root of its chain; or a blocked story, with its reason. It says when the host has no reminder channel. While the board stays starved, flai serve adds an entry only when that list changes. It resolves the thread once a story can start again. The thread is flai serve's, so the orchestrator leaves it alone (`answer_threads` passes over entries by neither a story's agent nor a planner), and while it is open the reminder of ADR-0205 covers it like any other.

**Touches.** No touches were declared. `flai touches suggest` ran from `flai/internal/serve/readyempty.go`, `flai/cmd/promote.go`, and `flai/internal/inbox/inbox.go` (18 commits).

| Touch | Source | Why |
|-------|--------|-----|
| `design/adrs/` | layout | The first task writes a new ADR whose number and name `flai adr new` gives, so no task can name its file yet. Kept as a folder touch for that alone |
| `design/adrs/README.md` | layout | The generated ADR index |
| `flai/internal/inbox/starved.go`, `starved_test.go` | layout | The list of what keeps each backlog story from ready. It joins `workitem.PromotionCandidates` with `inbox.Summarize`'s threads awaiting the operator. It lives in `inbox` because `threads` imports `workitem`, and `inbox` imports both |
| `flai/internal/serve/starved.go`, `starved_test.go` | layout | Opening, updating, and resolving the thread, as `agents_restart.go`, `forge.go`, and `kluster.go` open and resolve flai serve's other threads |
| `flai/internal/serve/agents_launch.go` | layout | `look()` calls the new step after `noteReadyEmpty` |
| `design/system/workflow.md`, `design/system/strategic-agents.md` | design | Cadence and Pull policy, and the `ready-empty` row and its reason, say what happens when nothing can be made ready |
| `docs/users/flai.md` | design | The `flai serve` section lists the threads flai serve opens |
| `design/issues/I-0288-…md`, `design/issues/summary.md` | criterion 2 | Closing I-0288 |

Left out from `flai touches suggest`: `flai/internal/inbox/inbox_test.go`, `events.go`, the `flai/internal/mcpserver` files, `flai/internal/storystart/start.go`, and `flai/cmd/serve_actions.go`. The change reaches none of them: it adds no event and no tool.

**After.** S-0673. This story reads S-0673's decision, ADR-0279, and its change to `readyDrained` in `flai/internal/serve/readyempty.go`: `ready-empty` counts only the ready stories not held on `after`. "Holds no story that can start" is built on that count. S-0673 was accepted on 2026-11-10, so the `after` holds nothing now.

**Forecast.** `flai forecast` gave 24m once the touches were in: 103 s a unit over 155 done large improvement stories, size 14 (2 criteria, 12 touches). It is raised to 1h, as S-0673's was for the same shape: the size counts neither the ADR the first task writes nor five tasks in four layers, two of which add Go files with tests. flai keeps the delivery current with the pull order.

**Cost of delay.** 150 USD a week, as `flai cod` works it out from the input the orchestrator set on TH-1042, `time_lost_per_cycle: 1h`, at 150 USD an hour. It stands: I-0288 measured no time lost, and the hour is the operator's time per cycle to notice an empty board and find the decisions it waits on. The board's idle lanes are not counted again here; the stories that waited carry their own cost of delay.
