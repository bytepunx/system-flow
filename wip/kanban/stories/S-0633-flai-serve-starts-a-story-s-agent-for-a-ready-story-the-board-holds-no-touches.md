---
id: S-0633
type: story
nature: improvement
title: flai serve starts a story's agent for a ready story the board holds no-touches, and story_start then refuses it
status: in-progress
owner: alex
created: 2026-11-07T23:41:06Z
updated: 2026-11-09T23:21:01Z
transitions:
  - to: ready
    at: 2026-11-09T22:05:17Z
    by: alex
  - to: in-progress
    at: 2026-11-09T22:12:26Z
    by: agent-S-0633
tags: [flai, serve]
touches: [flai/internal/serve/agents.go, flai/internal/serve/hold_test.go, flai/internal/storystart/start.go, flai/internal/storystart/start_test.go, design/system/workflow.md, design/issues/I-0301-flai-serve-starts-a-story-s-agent-for-a-ready-story-the-board-holds-no-touches.md, design/issues/summary.md, wip/kanban/stories/S-0633-flai-serve-starts-a-story-s-agent-for-a-ready-story-the-board-holds-no-touches.md, flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go, flai/internal/workitem/itemlock.go, design/adrs/0272-a-checked-creation-holds-a-checked-write-lock-from-its-write-until-it-keeps-or.md, design/adrs/README.md, design/issues/I-0328-a-checked-edit-of-a-work-item-writes-checks-and-restores-it-with-no-lock-so-a.md, design/issues/I-0314-testaformularolestepsendistoldtotheinstancesagent-fails-on-its-timing-under.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4102
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-09
      ceremony: 2
      hand_edits: 2
      work: 70
  models:
    - model: claude-opus-5-5
      input: 248
      output: 1747
      cache_read: 33756538
      cache_write: 2144377
      cost: 13.1272
  strategic:
    - kind: orchestrator
      seconds: 119
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 20
          output: 120
          cache_read: 1089370
          cache_write: 133221
          cost: 0.3253
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-07T23:41:06Z
  value: 12.5
  by: planner-S-0633
  at: 2026-11-09T22:02:17Z
forecast:
  duration: 40m
  delivery: 2026-11-09T22:58:00Z
  basis: "Its own forecast of 40m; 1st in the pull order with an in-progress limit of 5, behind S-0622, S-0623, S-0632 and S-0634."
  by: flai
  at: 2026-11-09T22:09:59Z
finalized:
  by: alex
  at: 2026-11-09T22:05:02Z
---
# S-0633 flai serve starts a story's agent for a ready story the board holds no-touches, and story_start then refuses it

## Goal

This story remediates [I-0301](../../../design/issues/I-0301-flai-serve-starts-a-story-s-agent-for-a-ready-story-the-board-holds-no-touches.md), "flai serve starts a story's agent for a ready story the board holds no-touches, and story_start then refuses it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0301 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0301 is closed with `flai issue close I-0301 --reason` saying what fixed it

## Tasks
- T-2786 flai serve starts no agent for a ready story while its hold stands, with a test that reproduces I-0301
- T-2787 story_start's refusal of a story held no-touches by its own empty claim says that declaring its touches clears the hold
- T-2789 workflow.md says an automatic start waits on the hold and story_start names the touches that clear a no-touches hold, and I-0301 is closed

## Notes

Cost of delay inputs set by flai from I-0301. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T08:44:09Z, 0.6 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0301's one instance. agent-S-0568 was started at 2026-11-07T08:18:46Z "because it entered ready". At that moment S-0568 had no touches and no tasks: T-2551 was created at 08:18:56Z, and the first touches edit came at 08:19:25Z. S-0562, S-0564, and S-0584 were in progress. `holds.Of` holds such a story `no-touches`, and the look's gate in `flai/internal/serve/agents.go` skips a held story, as `design/system/workflow.md` § Skip-ahead says. So the start came through a path or a race that the gate does not cover. The fix has two parts:

- T-2786 finds that path, reproduces it in a test, and judges the hold again at every automatic start.
- T-2787 changes what `story_start` says when the story's own empty claim holds it: declaring its touches clears the hold, so the agent is not sent to pull another story.

The serve journal for that minute could not be read from the planner's checkout. Naming the path is T-2786's first step.

Touches, where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/serve/agents.go` | layout: the look's hold gate and the start paths |
| `flai/internal/serve/hold_test.go` | layout: the serve hold tests (`TestOneLookStartsNoTwoThatOverlap`) |
| `flai/internal/storystart/start.go` | layout: `refuse`, the held refusal |
| `flai/internal/storystart/start_test.go` | layout: the held-refusal tests |
| `design/system/workflow.md` | design: the Held and Skip-ahead bullets |
| `design/issues/I-0301-…md`, `design/issues/summary.md` | declared by the goal and criteria |
| the story file | criteria ticks |

`flai touches suggest` from `agents.go` and `start.go` ranks `flai/internal/serve/agents_test.go` (49%) and `flai/cmd/serve_actions.go` (30%) highest. They are left out because the tests go in `hold_test.go` and no serve action changes. `restart.go` and `start.go` under `flai/internal/serve` are left out too: the operator's start passes a hold with a warning by design, and the restart paths start only stories already in progress. No folder touch and no `after`: the story reads nothing another open story writes.

Forecast: flai forecast gave 18m, which is raised to 40m with delivery 2026-11-09T23:48:00Z, flai's delivery plus the 22m added. The reason is that the cause is still to be traced, and S-0568, a three-task remediation of this size, spent about 31m in progress.

Cost of delay: 12.50 USD a week, as flai cod works it out from the operator's input of 5m lost per 168h cycle at 150 USD an hour. It stands as computed.
