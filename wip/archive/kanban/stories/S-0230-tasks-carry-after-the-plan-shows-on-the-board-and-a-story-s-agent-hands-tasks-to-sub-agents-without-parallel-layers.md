---
id: S-0230
type: story
nature: improvement
title: Tasks carry after, the plan shows on the board, and a story's agent hands tasks to sub-agents without parallel layers
status: done
owner: arobson
created: 2026-10-02T12:03:34Z
updated: 2026-10-02T23:22:07Z
transitions:
  - to: ready
    at: 2026-10-02T16:27:00Z
    by: alex
  - to: in-progress
    at: 2026-10-02T17:11:19Z
    by: agent-S-0230
  - to: review
    at: 2026-10-02T17:29:51Z
    by: agent-S-0230
  - to: done
    at: 2026-10-02T23:22:07Z
    by: alex
tags: [flai, dashboard, template]
touches: [flai/internal/usage, flai/internal/serve, flai/internal/harness, flai/cmd, design/conventions, template, design/adrs, design/issues, design/system, docs]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1161
  models:
    - model: claude-haiku-4-5-20251001
      input: 114
      output: 4979
      cache_read: 413332
      cache_write: 43845
      cost: 0.1211
    - model: claude-opus-5-5
      input: 254
      output: 83945
      cache_read: 10054955
      cache_write: 326607
      cost: 5.7776
    - model: claude-sonnet-5-5
      input: 24
      output: 8184
      cache_read: 513676
      cache_write: 68036
      cost: 0.3547
---
# S-0230 Tasks carry after, the plan shows on the board, and a story's agent hands tasks to sub-agents without parallel layers

## Goal

S-0176's experiment recommended adapt: keep a task's `after` and the plan's display, and keep handing a story's tasks to sub-agents, which is where its replays saved cost (15 and 28% on two of three), but drop the expectation that running a layer's tasks at once makes a story faster: on stories of three to five tasks of a few minutes each it did not. S-0176's branch built all of it as an experiment; this story adopts the parts that held, as the results document `design/experiments/S-0176-the-story-s-agent-plans-which-tasks-can-run-in-parallel-and-works-them-with-sub-agents.md` records.

## Acceptance criteria
- [x] The host's flai carries a task's `after` and `flai.minimum` is raised so that no older flai refuses a planned task (S-0181)
- [x] I-0054 is fixed: a task worked by a sub-agent is measured by its sub-agent's calls (its `parent_tool_use_id`), not by its time window, so tasks worked at once no longer each get the whole run's usage
- [x] The conventions (`delegation.md`, `work-management.md`, here and in the template) and `harness.Prompt` say a story's agent plans its tasks with `after` and hands each to a task sub-agent, and that running a layer at once is optional and worth it only where the tasks are long beside the story's fixed costs
- [x] `design/system/agent-context.md` records that forked sub-agents are still not offered to a headless session, and the repeat measurement is filed as its own story (S-0241) to run once they are (TH-0069)
- [x] `design/system/agent-context.md` records what was adopted and links the experiment's results document

## Tasks
- T-0711 A task worked by a sub-agent is measured by its sub-agent's calls, not its time window (I-0054)
- T-0712 The prompt and delegation.md say to name the task's ID in a task sub-agent's description
- T-0713 The design records what S-0230 adopted and how a task's usage is measured

## Notes

Follows the experiment S-0176 (ADR-0066). Its own tasks did not carry `after` on `main` because the host's flai predated it; see S-0176's narrative.
