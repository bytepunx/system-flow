---
id: T-1663
type: task
nature: feature
title: The orchestrator's convention, agent definition, and prompt tell it to enqueue with evidence while the queue runs, and to open no thread on a blocker the queue settles
status: backlog
parent: S-0403
owner: alex
created: 2026-10-16T05:02:24Z
updated: 2026-10-16T05:02:24Z
transitions: []
stream: S-0403
tags: [cli]
touches: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/root/.claude/agents/orchestrator.md, ".claude/agents/orchestrator.md", flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go]
after: [T-1649]
---
# T-1663 The orchestrator's convention, agent definition, and prompt tell it to enqueue with evidence while the queue runs, and to open no thread on a blocker the queue settles

## Work

- In `strategic-agents.md` § As the orchestrator, baseline, in the repository and the template, change step 3 of `accept_reviews`. While the queue runs, enqueue with `queue_add` and the evidence, and log each enqueue with `activity_log`, naming the story, its place in the queue, and the policy figure. In step 4, open no thread for a blocker the queue settles.
- Name `queue_add` among the orchestrator's tools in both copies of `orchestrator.md`. The `.claude/` writes go through the operator's permission prompt.
- In `serve/orchestrate.go`, tell the orchestrator in its prompt whether the queue runs for the project.
- Test the prompt with the queue on and off.
- Waits for T-1649, which defines `queue_add`.

## Done when

The convention, both agent definitions, and the prompt say the above, and the tests pass.

## Notes
