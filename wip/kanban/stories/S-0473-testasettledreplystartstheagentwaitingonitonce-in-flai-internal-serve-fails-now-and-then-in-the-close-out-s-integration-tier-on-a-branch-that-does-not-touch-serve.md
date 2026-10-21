---
id: S-0473
type: story
nature: remediation
title: TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve
status: in-progress
owner: alex
created: 2026-10-21T11:58:53Z
updated: 2026-10-21T12:40:57Z
transitions:
  - to: ready
    at: 2026-10-21T12:05:25Z
    by: alex
  - to: in-progress
    at: 2026-10-21T12:06:08Z
    by: agent-S-0473
tags: []
touches: [flai/internal/serve/agents_test.go, design/issues/I-0158-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md, design/issues/summary.md, design/issues/I-0183-flai-check-finds-item-chronology-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 286
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-21
      hand_edits: 2
      work: 19
  models:
    - model: claude-opus-5-5
      input: 42
      output: 356
      cache_read: 2206353
      cache_write: 127568
      cost: 0.7901
  strategic:
    - kind: orchestrator
      seconds: 527
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 133
          output: 2167
          cache_read: 25881209
          cache_write: 41873
          cost: 5.7241
cost_of_delay:
  inputs:
    time_lost_per_cycle: 45m
    by: flai
    at: 2026-10-21T11:58:53Z
  value: 112.5
  by: planner-S-0473
  at: 2026-10-21T12:08:26Z
forecast:
  duration: 45m
  delivery: 2026-10-21T12:55:00Z
  basis: "flai forecast's 11m (131 s per size unit over 8 small remediation stories, size 5) raised to 45m, because the close-out's integration tier alone took 6m19s on S-0463 and the fix runs the test 200 times; delivered 45m after its start at 12:06Z."
  by: planner-S-0473
  at: 2026-10-21T12:08:26Z
finalized:
  by: alex
  at: 2026-10-21T12:04:09Z
---
# S-0473 TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve

## Goal

This story remediates [I-0158](../../../design/issues/I-0158-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md), "TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0158 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0158 is closed with `flai issue close I-0158 --reason` saying what fixed it

## Tasks
- T-1986 Pin both runs of TestASettledReplyStartsTheAgentWaitingOnItOnce to one second so the first look never finds the message new to the other story
- T-1989 Close I-0158 with the reason the pinned run times fixed it
- T-1992 TestASettledReplyStartsTheAgentWaitingOnItOnce writes its message before the other story's run started, so a turning second cannot make it new

## Notes

Cost of delay inputs set by flai from I-0158. time_lost_per_cycle 45m: 9m per occurrence × 5 occurrences ÷ 1 cycle of 168h (first reported 2026-10-17T04:54:11Z, 4.3 days before this story; under one cycle counts as one).

### Planning

Proposed cause, from I-0158's five instances and the code: the test's two `putEnded` calls stamp the runs with the wall-clock second. When they straddle a second boundary, the message from `id` to `other` falls strictly after `other`'s run (`messaged` in `flai/internal/serve/agents.go`, `messages.ToSince` comparing with `After`). It then reads as new to `other` and owes a reply, so the first look starts `other`'s agent (S-0001 on MS-0001), which fails at `started before the reply`. The fix is in the test, not in `messaged`: a message that comes after a run ended is rightly new to it. No other test in `agents_test.go` makes two `putEnded` runs in one test.

The planner drafted T-1986 (fix) and T-1989 (close the issue, after T-1986). The story's agent then wrote T-1992, which reaches the same cause and covers both. The plan's thread proposes keeping one of the two plans.

Touches, all files, no folder:

| Touch | Source | Why |
| --- | --- | --- |
| `flai/internal/serve/agents_test.go` | design (the issue names the test) | The test is fixed here. |
| `design/issues/I-0158-…md` | layout | `flai issue close` writes it. |
| `design/issues/summary.md` | layout, co-change (2 of 31 commits) | `flai issue close` regenerates it. |

Left out: `flai/internal/serve/agents.go` co-changes with the test in 29 of 31 commits, but the fix needs no production change. T-1986 tells the agent to plan again and widen its touches if the reproduction points there. S-0406, a backlog draft, also touches `agents_test.go`. It holds nothing until it is in progress.

Figures:

- Forecast 45m, up from `flai forecast`'s 11m (131 s per size unit × size 5). The close-out's integration tier alone took 6m19s in S-0463's instance, and the fix runs the test 200 times. Delivery is 45m after the story started at 12:06Z.
- Cost of delay value 112.50 USD a week, as `flai cod` gives it: 45m lost per 168h cycle at 150 USD an hour. It stands on the operator's inputs. The instances suggest they understate: 5 occurrences in 4.3 days is about 8 a week, and the last one cost a story its place in the merge queue too. That is raised on the plan's thread, not changed here.
