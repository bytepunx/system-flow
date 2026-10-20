---
id: T-1663
type: task
nature: feature
title: The orchestrator's convention, agent definition, and prompt tell it to enqueue with evidence while the queue runs, and to open no thread on a blocker the queue settles
status: done
parent: S-0403
owner: alex
created: 2026-10-16T05:02:24Z
updated: 2026-10-20T00:24:54Z
transitions:
  - to: ready
    at: 2026-10-20T00:02:25Z
    by: agent-S-0403
  - to: in-progress
    at: 2026-10-20T00:02:26Z
    by: agent-S-0403
  - to: done
    at: 2026-10-20T00:24:54Z
    by: agent-S-0403
stream: S-0403
tags: [cli]
touches: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/serve/orchestrate_hold.go]
after: [T-1649]
usage:
  source: log
  seconds: 1348
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 132
      output: 46012
      cache_read: 9653019
      cache_write: 350206
      cost: 5.0241
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
