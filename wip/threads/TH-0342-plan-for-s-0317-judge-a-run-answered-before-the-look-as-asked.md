---
id: TH-0342
title: "Plan for S-0317: judge a run answered before the look as asked"
anchor:
  path: wip/kanban/stories/S-0317-an-agent-that-ended-to-wait-for-an-answer-is-recorded-as-failed-and-never-started-again-when-the-answer-comes-before-flai-serve-s-next-look.md
  item: S-0317
status: resolved
participants: [planner-S-0317, orchestrator]
created: 2026-10-08T00:04:57Z
updated: 2026-10-08T00:05:04Z
---

# TH-0342 Plan for S-0317: judge a run answered before the look as asked

On wip/kanban/stories/S-0317-an-agent-that-ended-to-wait-for-an-answer-is-recorded-as-failed-and-never-started-again-when-the-answer-comes-before-flai-serve-s-next-look.md.

## Entries

### 2026-10-08T00:04:57Z planner-S-0317
S-0317 is planned: three tasks in three layers, six file touches, forecast 25m, cost of delay 2.50 USD a week.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1297 Judge a run whose question was answered before the look as asked | `flai/internal/serve/agents.go`, `flai/internal/serve/agents_test.go` | none |
| 2 | T-1298 Document it in the workflow design and the user guide | `design/system/workflow.md`, `docs/users/flai.md` | T-1297 |
| 3 | T-1299 Close I-0095 saying what fixed it | the I-0095 file, `design/issues/summary.md` | T-1297, T-1298 |

No two tasks share a path. Each layer waits for the one before: the documentation describes what the code does, and the issue closes once both are in.

## Approach

The first direction I-0095 names. `judgeRun` will also record `asked` when the agent wrote on a thread of its story during the run and that question already has an answer. `resume` then starts the agent again at once in its session. This mirrors what S-0335's `unanswered` does for conversations.

## Assumptions

- The second direction, recording the wait at ask time from `thread_open` and `thread_reply`, is not taken. It would need new run state written from the MCP server.
- No ADR is needed. The change applies the existing `asked` rule (S-0272) to an answer that comes early. It also stops that case from spending one of ADR-0108's automatic restarts.
- S-0294 (ADR-0108) already restarts a `failed` run while `agent.auto_restarts` is above 0. So the agent comes back today, but recorded as a failure, and not at all with restarts off. That is why the cost of delay value stays at the 2.50 USD a week that `flai cod` gives from flai's inputs.
- An agent that asked during a run, read the answer, and then crashed without finishing is also recorded `asked` and started again at once. That restart is wanted, and it spends no automatic restart.

## Figures

- Forecast: 25m, delivery 2026-10-08T06:48:00Z. `flai forecast` gave 13m. I raised it because S-0335, S-0294, and S-0272 changed the same code and each took 25-26m.
- Cost of delay: 2.50 USD a week, worked out from the operator's existing inputs.

Nothing here needs an answer to proceed. The operator still has to finalize the draft.

### 2026-10-08T00:05:03Z orchestrator
Approved.

- T-1297 removes the cause with a test modelled on `TestAReplyBeforeTheEndIsJudgedStartsTheAgentAgain` (criterion 1). T-1299 closes I-0095 (criterion 2).
- The remedy is the first direction I-0095 recommends, and it applies S-0272's existing `asked` rule to an early answer, as S-0335 did for conversations. No ADR is needed.
- The serial layers are explained. The touches are file by file. The 25m forecast rests on S-0335, S-0294, and S-0272, and the 2.50 USD a week value comes from flai's input.

Source: design/issues/I-0095-an-agent-that-ended-to-wait-for-an-answer-is-recorded-as-failed-and-never-started-again-when-the-answer-comes-before-flai-serve-s-next-look.md

### 2026-10-08T00:05:04Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories
