---
id: T-0638
type: task
nature: research
title: Measure the two delegating runs and comparable runs that did not delegate
status: in-progress
parent: S-0188
owner: arobson
created: 2026-10-01T08:31:38Z
updated: 2026-10-01T09:01:26Z
transitions:
  - to: ready
    at: 2026-10-01T08:31:52Z
    by: agent-S-0188
  - to: in-progress
    at: 2026-10-01T09:01:26Z
    by: agent-S-0188
stream: S-0188
tags: []
touches: [design/system/agent-context.md]
---
# T-0638 Measure the two delegating runs and comparable runs that did not delegate

## Work

Read each run's log as `agent-context.md § Sub-agents › Measured` reads S-0175 and S-0118 (ADR-0051: assistant events deduplicated by message ID, the `result` for cost and turns, `parent_tool_use_id` for a sub-agent's calls): cost, cache reads, model calls, and turns for the story's agent and its sub-agents, and wall time. Pick comparable runs that did not delegate by nature, size, and components touched.

## Done when

Each delegating run has its numbers and at least one comparable non-delegating run, with what makes them comparable.

## Notes
