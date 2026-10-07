---
id: T-1173
type: task
nature: feature
title: The orchestrator's prompt, agent definition, and convention tell it to plan candidate stories and to settle the planner's threads
status: cancelled
parent: S-0328
owner: alex
created: 2026-10-07T19:40:07Z
updated: 2026-10-07T19:42:16Z
transitions:
  - to: cancelled
    at: 2026-10-07T19:42:16Z
    by: agent-S-0328
stream: S-0328
tags: [flai, orchestrator, conventions, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md]
after: [T-1171, T-1172]
---
# T-1173 The orchestrator's prompt, agent definition, and convention tell it to plan candidate stories and to settle the planner's threads

## Work

Tell the orchestrator what the new permission lets it do, in `orchestratePrompt` in `flai/internal/harness/harness.go`, in `.claude/agents/orchestrator.md` and its template copy, and under "As the orchestrator" in `strategic-agents.md` and its template copy:

- While `plan_backlog_stories` is on, run `flai plan --candidates` and start the planner with `plan` for each story it lists, one at a time, as for epics under `plan_backlog_epics`.
- Take the planner's threads with the threads awaiting the operator, whose last entry is now a story's agent's or a planner's. Answer a plan thread with an approval when the plan's tasks, touches, and forecast describe the story's criteria, or say what is inconsistent. Answer a planner's cost of delay question by choosing the inputs, or taking the planner's recommendation, and writing them as T-1172 allows, citing the planner's evidence; post a recommendation instead while a permission T-1168's ADR names is off.
- Narrow the rule that money escalates to the operator as the ADR says, in the prompt and in both copies of the convention.

Add a `template/CHANGELOG.md` entry for the convention and agent definition. The story's agent writes the two `.claude/agents/orchestrator.md` files itself, since a sub-agent is refused writes under `.claude/` (ADR-0102).

It waits for T-1171 and T-1172, whose behaviour it describes. It shares no path with T-1174 and runs beside it.

## Done when

- The prompt, both agent definitions, and both copies of the convention say the same, and `harness_test.go` checks the prompt names `plan_backlog_stories` and the planner's threads.
- `flai prime --role orchestrate` shows the new rules.
- `flai test` and the markdown lint pass on the paths changed.

## Notes

A branch that changes `.claude/` is accepted by the operator only (ADR-0106).
- 2026-10-07T19:42:16Z: moved to cancelled: duplicates the story agent's plan (T-1174, T-1176 to T-1179), written at start while the planner ran; its refinements are folded into T-1176 to T-1178
