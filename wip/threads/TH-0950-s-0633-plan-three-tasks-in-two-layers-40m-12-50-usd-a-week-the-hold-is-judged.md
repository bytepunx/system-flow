---
id: TH-0950
title: "S-0633 plan: three tasks in two layers, 40m, 12.50 USD a week; the hold is judged again at every automatic start"
anchor:
  path: wip/kanban/stories/S-0633-flai-serve-starts-a-story-s-agent-for-a-ready-story-the-board-holds-no-touches.md
  item: S-0633
status: resolved
participants: [planner-S-0633, orchestrator]
created: 2026-11-09T22:02:49Z
updated: 2026-11-09T22:07:33Z
---

# TH-0950 S-0633 plan: three tasks in two layers, 40m, 12.50 USD a week; the hold is judged again at every automatic start

On wip/kanban/stories/S-0633-flai-serve-starts-a-story-s-agent-for-a-ready-story-the-board-holds-no-touches.md.

## Entries

### 2026-11-09T22:02:49Z planner-S-0633
Planned from scratch. The run prompt says it was asked to plan again on a thread "asked". No such thread exists, and S-0633 had no earlier plan: no touches, forecast, value, or tasks.

## Proposed solution

agent-S-0568 started at 2026-11-07T08:18:46Z "because it entered ready", with no touches and no tasks. Its first task came at 08:18:56Z and its first touches at 08:19:25Z. S-0562, S-0564, and S-0584 were in progress then. The look's gate in `flai/internal/serve/agents.go` already skips a story `holds.Of` holds. So the start came through a path or a race that the gate does not cover. I could not read the serve journal for that minute from the main checkout.

- Fix the start: find that path, reproduce it in a test, and judge the hold again at every automatic start.
- Fix the advice: when the story's own empty claim holds it, `story_start` names `flai touches` as what clears the hold, instead of "pull the first ready story … not held".

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2786 serve starts no agent while the hold stands, with a reproducing test | `flai/internal/serve/agents.go`, `flai/internal/serve/hold_test.go` | none |
| 1 | T-2787 story_start's no-touches refusal names `flai touches` | `flai/internal/storystart/start.go`, `flai/internal/storystart/start_test.go` | none |
| 2 | T-2789 workflow.md Held and Skip-ahead bullets, close I-0301 | `design/system/workflow.md`, I-0301, `design/issues/summary.md` | T-2786, T-2787 |

T-2786 and T-2787 share no path, so they can run together.

## Figures

- Forecast: 40m, up from flai's 18m. The cause is still to be traced, and S-0568, a similar remediation, spent about 31m in progress.
- Cost of delay: 12.50 USD a week from the operator's 5m per cycle input, unchanged.

## Assumptions

- The operator's own start (`flai/internal/serve/start.go`) keeps passing a hold with a warning. Only automatic starts wait on the hold.
- flai serve does not wait for the planner before it starts a story's agent. The hold is the one rule, as `design/system/workflow.md` § Skip-ahead says.
- If T-2786 finds a race rather than a path that skips the gate, the fix re-checks the hold just before the launch. It does not widen the look's lock.

Nothing to split, merge, or drop.

### 2026-11-09T22:06:05Z orchestrator
Approved. The tasks, touches, and figures fit S-0633.

- The plan proposes a fix, as the goal asks. T-2786 finds the start path, reproduces I-0301 in a test that fails before the change, and gates every automatic start on the hold, for criterion 1. T-2787 fixes the advice the agent got. T-2789 closes I-0301 for criterion 2.
- The layers are justified. T-2786 and T-2787 share no path, so they run together. T-2789 describes what both built.
- The touches name files only, with no folder touch. The story file itself is touched for the criteria ticks.
- 40m over flai's 18m gives its reason: the cause is still to be traced. 12.50 USD a week comes from the operator's 5m input.
- The assumptions stand. The operator's own start keeps passing a hold with a warning, as `design/system/workflow.md` § Skip-ahead says. Re-checking the hold just before launch, if T-2786 finds a race, keeps the look's lock as it is.

Note for the story's agent: S-0632 (in progress) also changes `flai/internal/serve/agents.go`, and S-0631 (ready) splits it. The trial merge judges any hold, and `flai stream sync` and the stories' conversation settle a conflict.

The operator has already finalized S-0633 and moved it to ready.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-09T22:07:33Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: tasks, touches, forecast, and value fit S-0633.
