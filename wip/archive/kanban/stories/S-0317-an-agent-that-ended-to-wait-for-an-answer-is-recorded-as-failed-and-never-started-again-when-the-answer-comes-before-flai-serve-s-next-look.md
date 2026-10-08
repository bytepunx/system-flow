---
id: S-0317
type: story
nature: remediation
title: An agent that ended to wait for an answer is recorded as failed, and never started again, when the answer comes before flai serve's next look
status: done
owner: alex
created: 2026-10-07T18:59:47Z
updated: 2026-10-08T00:31:19Z
transitions:
  - to: ready
    at: 2026-10-08T00:05:12Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T00:10:22Z
    by: agent-S-0317
  - to: review
    at: 2026-10-08T00:30:33Z
    by: agent-S-0317
  - to: done
    at: 2026-10-08T00:31:19Z
    by: orchestrator
tags: [flai]
topics: [cli]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, design/system/workflow.md, docs/users/flai.md, design/issues/I-0095-an-agent-that-ended-to-wait-for-an-answer-is-recorded-as-failed-and-never-started-again-when-the-answer-comes-before-flai-serve-s-next-look.md, design/issues/summary.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1231
  turns:
    - day: 2026-10-08
      ceremony: 3
      hand_edits: 1
      work: 29
  models:
    - model: claude-opus-5-5
      input: 134
      output: 35469
      cache_read: 5971681
      cache_write: 281701
      cost: 3.7974
  strategic:
    - kind: orchestrator
      seconds: 371
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 58
          output: 877
          cache_read: 15897476
          cache_write: 34512
          cost: 3.9253
        - model: claude-sonnet-5-5
          input: 6
          output: 61
          cache_read: 39763
          cache_write: 29134
          cost: 0.0594
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1m
    by: flai
    at: 2026-10-07T18:59:47Z
  value: 2.5
  by: planner-S-0317
  at: 2026-10-08T00:04:36Z
forecast:
  duration: 25m
  delivery: 2026-10-08T06:48:00Z
  basis: "flai forecast gave 13m (94 s per unit over 17 remediation stories, size 8); raised to 25m because S-0335, S-0294, and S-0272, which changed the same run judging in flai/internal/serve/agents.go, each took 25-26m of agent time; delivery is flai's 06:36Z, 21st in the pull order, plus the 12m added."
  by: planner-S-0317
  at: 2026-10-08T00:04:36Z
finalized:
  by: orchestrator
  at: 2026-10-08T00:05:06Z
---
# S-0317 An agent that ended to wait for an answer is recorded as failed, and never started again, when the answer comes before flai serve's next look

## Goal

This story remediates [I-0095](../../../design/issues/I-0095-an-agent-that-ended-to-wait-for-an-answer-is-recorded-as-failed-and-never-started-again-when-the-answer-comes-before-flai-serve-s-next-look.md), "An agent that ended to wait for an answer is recorded as failed, and never started again, when the answer comes before flai serve's next look". The issue recommends this solution:

Directions to weigh: when a run ends, read whether its story has a thread whose last answer came after the agent's last entry on it, and treat that as an answer to act on rather than as no question; or record at the moment an agent asks, through `thread_open` or `thread_reply` on its own story, that it is waiting, so that the end is `asked` whatever the thread says by the time of the look. Either way, a story in progress whose agent has ended, with no thread awaiting the operator and no block, could be restarted once by flai serve, which is the second remediation of I-0084.

## Acceptance criteria
- [x] The cause I-0095 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0095 is closed with `flai issue close I-0095 --reason` saying what fixed it

## Tasks
- T-1297 Judge a run whose question was answered before the look as asked, so flai serve starts it again on the answer
- T-1298 Say in the workflow design and the user guide that an answer before flai serve's look still starts the agent again
- T-1299 Close I-0095 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0095. time_lost_per_cycle 1m: 1m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T21:01:06Z, 0.9 days before this story; under one cycle counts as one).

### Planning

Touches, each a file; no folder touch is kept:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/serve/agents.go` | layout | `judgeRun`, `asking`, `answered`, and `threadAnswered` decide `asked` or `failed`; the fix goes beside `asking` |
| `flai/internal/serve/agents_test.go` | layout | the run-judging tests, such as `TestAReplyBeforeTheEndIsJudgedStartsTheAgentAgain`, the new test's model |
| `design/system/workflow.md` | design, co-change | the **An answered agent** bullet says when a run is `asked` |
| `docs/users/flai.md` | design, co-change | the S-0272 and S-0335 paragraphs tell the operator when an agent starts again |
| `design/issues/I-0095-…-next-look.md` | design | criterion 2 closes it |
| `design/issues/summary.md` | co-change | `flai issue close` regenerates it |

Left out: `design/system/flai-cli.md` (68% co-change), since its `flai serve agent` row names no judging rule this changes; `flai/internal/mcpserver/server.go`, since the second direction, recording the wait at ask time, is not taken; and `docs/operators/index.md`, whose automatic restart paragraph stays true.

Approach: the first direction I-0095 names. S-0335 already counts a message that came after the run started and has no reply (`unanswered`); T-1297 does the same for a thread the agent asked on during the run and that has an answer since. S-0294 (ADR-0108) already restarts a `failed` run, so today the agent comes back, but as a failure, spending one of `agent.auto_restarts`, and not at all with restarts off.

Figures:

- Forecast 25m, raised from flai's 13m: S-0335, S-0294, and S-0272 changed the same code and each took 1526-1576 s of agent time. Delivery 2026-10-08T06:48:00Z is flai's 06:36Z plus the 12m added.
- Cost of delay 2.50 USD a week, as `flai cod` gives it from the inputs flai set from I-0095 (1m lost per 168h cycle at 150 USD an hour). It stands: one occurrence, 35 seconds lost, and S-0294 already restarts the agent while restarts are on.

### Accepted by the orchestrator

- Verified: 337b378edfbb8665f0338147fe104e1e03aa0347
- At: 2026-10-08T00:31:19Z

Verdict: accept; both criteria met (verifier at 337b378edfbb8665f0338147fe104e1e03aa0347; flai verify passed every step at that commit). The I-0079 bump is the close-out's own record.

- 1: flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, design/system/workflow.md, docs/users/flai.md
- 2: design/issues/I-0095-an-agent-that-ended-to-wait-for-an-answer-is-recorded-as-failed-and-never-started-again-when-the-answer-comes-before-flai-serve-s-next-look.md, design/issues/summary.md
