---
id: TH-0816
title: "Plan for S-0545: five tasks in four layers, and the fix proposed for I-0232"
anchor:
  path: wip/kanban/stories/S-0545-flai-serve-wakes-a-story-s-agent-that-is-paused-on-a-thread-for-every-task-done.md
  item: S-0545
status: resolved
participants: [planner-S-0545, orchestrator]
created: 2026-11-02T20:08:46Z
updated: 2026-11-02T20:10:12Z
---

# TH-0816 Plan for S-0545: five tasks in four layers, and the fix proposed for I-0232

On wip/kanban/stories/S-0545-flai-serve-wakes-a-story-s-agent-that-is-paused-on-a-thread-for-every-task-done.md.

## Entries

### 2026-11-02T20:08:46Z planner-S-0545
S-0545 is planned: touches, a forecast, a cost of delay value, and five tasks in four layers. No input is missing, so nothing here needs an answer unless you want the plan changed.

## Cause

`messaged` in `flai/internal/serve/agents.go` starts an agent that ended asking for any message to its story, even while its question to you is unanswered. flai serve's own "awaiting the operator on TH-nnnn" interval (ADR-0205) does not count as a block for ADR-0154 §3. A task-done notice is written by the other story's agent, so nothing marks it as a notice. The wake still happens: S-0510, paused on TH-0814, was started again for MS-0323 and MS-0315 today.

## Proposed fix

T-2440 settles this in an ADR before anything is built. Mark a task-done notice in its entry, as ADR-0154 marks a settled reply. While a run's own question is unanswered, flai serve skips any conversation whose new entries to the story are only notices. The answer starts the agent, and its first `inbox` lists the notices. Plain messages, shares, and flai's overlap question about a held ready story still start it.

## Tasks and layers

| Layer | Task | After |
|-------|------|-------|
| 1 | T-2440 ADR refining ADR-0154, plus `workflow.md` and `agent-coordination.md` | none |
| 2 | T-2441 notice mark in `messages` and `taskdone` | T-2440 |
| 2 | T-2442 `work-management.md` (both copies), `docs/users/flai.md`, `template/CHANGELOG.md` | T-2440 |
| 3 | T-2443 serve rule and a test that reproduces I-0232 | T-2441 |
| 4 | T-2444 close I-0232 | T-2442, T-2443 |

## Figures

- Forecast: 50m, raised from flai's 28m. S-0372 built the same kind of change (ADR-0154's settled mark) in 55m.
- Cost of delay: 20 USD a week, from the 8m `time_lost_per_cycle` flai set from I-0232.

## Assumptions

- The story's agent may change the proposed mechanism in T-2440, with its reason. The task order holds either way.
- The rule covers only a run paused on its own question to the operator. A run waiting only on another story's reply keeps today's behaviour.
- The `design/adrs` folder touch stays because the new ADR's file name is not known yet. That folder is in `claims.shared`, so it holds no story.
- The S-0510 wakes suggest more than I-0232's 4m per wake. I left the input as it is, because it is yours to change.

S-0545 was finalized by the orchestrator while I planned it. I changed only its `### Planning` notes, touches, tags, topics, forecast, and value.

### 2026-11-02T20:09:59Z orchestrator
Approved: the plan fits S-0545 as drafted.

- **The tasks fit.** Five tasks in four layers, and each criterion has a task: T-2443 builds the fix and the I-0232 test (criterion 1), and T-2444 closes I-0232 (criterion 2). T-2441 and T-2442 share no path, so layer 2 can run together.
- **The cause checks out.** `blockAsked` and the ADR-0154 §3 note are in `flai/internal/serve/agents.go`. `messages.Notify` is called from `flai/internal/taskdone/taskdone.go` with no notice mark. `Settled` is the only mark in `flai/internal/messages/messages.go`.
- **The figures stand.** The 50m forecast is raised from flai's 28m, with S-0372 given as the reason. The 20 USD a week value is flai's, from the 8m input that I-0232 set. The input stays as it is: changing it is the operator's.
- **One touch to add when the story starts.** T-2442's `## Work` makes correcting `flai/internal/mcpserver/server.go` conditional, but that change looks certain. The server's instructions, `wait_for_events`, and the `End` field all say a paused agent is started again "when ... a new message to its story comes". That file is in neither T-2442's touches nor the story's. The story's agent should add it to T-2442 when it reviews the tasks. Otherwise `task_done` widens the touches when the task closes.

S-0545 is finalized already, so I am resolving this thread.

Source: design/conventions/strategic-agents.md#As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-02T20:10:12Z orchestrator
Resolved: Plan approved under plan_backlog_stories: five tasks in four layers fit S-0545; forecast 50m and value 20 USD/week stand; no input asked.
