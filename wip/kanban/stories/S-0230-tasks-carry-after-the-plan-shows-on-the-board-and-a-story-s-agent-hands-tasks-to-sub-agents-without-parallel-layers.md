---
id: S-0230
type: story
nature: improvement
title: Tasks carry after, the plan shows on the board, and a story's agent hands tasks to sub-agents without parallel layers
status: ready
owner: arobson
created: 2026-10-02T12:03:34Z
updated: 2026-10-02T16:27:00Z
transitions:
  - to: ready
    at: 2026-10-02T16:27:00Z
    by: alex
tags: [flai, dashboard, template]
touches: [flai/internal/workitem, flai/internal/harness, flai/internal/usage, design/conventions/, template/, flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0230 Tasks carry after, the plan shows on the board, and a story's agent hands tasks to sub-agents without parallel layers

## Goal

S-0176's experiment recommended adapt: keep a task's `after` and the plan's display, and keep handing a story's tasks to sub-agents, which is where its replays saved cost (15 and 28% on two of three), but drop the expectation that running a layer's tasks at once makes a story faster: on stories of three to five tasks of a few minutes each it did not. S-0176's branch built all of it as an experiment; this story adopts the parts that held, as the results document `design/experiments/S-0176-the-story-s-agent-plans-which-tasks-can-run-in-parallel-and-works-them-with-sub-agents.md` records.

## Acceptance criteria
- [ ] The host's flai carries a task's `after` and `flai.minimum` is raised so that no older flai refuses a planned task (S-0181)
- [ ] I-0054 is fixed: a task worked by a sub-agent is measured by its sub-agent's calls (its `parent_tool_use_id`), not by its time window, so tasks worked at once no longer each get the whole run's usage
- [ ] The conventions (`delegation.md`, `work-management.md`, here and in the template) and `harness.Prompt` say a story's agent plans its tasks with `after` and hands each to a task sub-agent, and that running a layer at once is optional and worth it only where the tasks are long beside the story's fixed costs
- [ ] The measurement is repeated once forked sub-agents are offered to a headless session, and `design/system/agent-context.md` records the result beside the first
- [ ] `design/system/agent-context.md` records what was adopted and links the experiment's results document

## Tasks

## Notes

Follows the experiment S-0176 (ADR-0066). Its own tasks did not carry `after` on `main` because the host's flai predated it; see S-0176's narrative.
