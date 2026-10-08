---
id: T-1337
type: task
nature: research
title: Create the epic for the LiteLLM and OpenRouter adapters and link it from the finding
status: done
parent: S-0339
owner: alex
created: 2026-10-08T07:21:13Z
updated: 2026-10-08T08:34:04Z
transitions:
  - to: ready
    at: 2026-10-08T08:32:43Z
    by: agent-S-0339
  - to: in-progress
    at: 2026-10-08T08:32:43Z
    by: agent-S-0339
  - to: done
    at: 2026-10-08T08:34:04Z
    by: agent-S-0339
stream: S-0339
tags: [research, agents]
touches: [design/system/agent-adapters.md]
after: [T-1336]
usage:
  source: log
  seconds: 81
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 155
      output: 8504
      cache_read: 819560
      cache_write: 46101
      cost: 1.3515
---
# T-1337 Create the epic for the LiteLLM and OpenRouter adapters and link it from the finding

## Work

Create the epic with `flai epic new` and `--body-stdin`. Give it:

- a goal: agent adapters for LiteLLM and OpenRouter, built on the decided abstractions
- acceptance criteria as checkboxes, one per outcome the decision names, such as each adapter, the agent schema change, the guard and permission contract, and usage measurement
- notes that link `design/system/agent-adapters.md` and each new ADR, and name the order the recommendation gives.

Write no stories under it. The planner drafts them from the finding, as the story's third criterion says. Set no cost of delay inputs: they are the operator's. Add a `## Epic` section, or a line in `## Decision`, to `design/system/agent-adapters.md`, naming the epic by ID. In that section, list what the planner must read first. Waits for T-1336 because the epic follows the decision.

## Done when

- The epic exists in the backlog with goal, criteria, and notes linking the finding and the ADRs, and `flai check --strict` passes.
- `design/system/agent-adapters.md` names the epic.
- The story's first and third criteria are ticked.

## Notes
