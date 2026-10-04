---
id: S-0255
type: story
nature: feature
title: The planner drafts a story's tasks into the backlog and revisits the children it already has
status: done
parent: E-0016
owner: alex
created: 2026-10-04T00:51:25Z
updated: 2026-10-04T19:34:24Z
transitions:
  - to: ready
    at: 2026-10-04T00:51:31Z
    by: alex
  - to: in-progress
    at: 2026-10-04T04:00:33Z
    by: agent-S-0255
  - to: review
    at: 2026-10-04T19:26:45Z
    by: agent-S-0255
  - to: done
    at: 2026-10-04T19:34:24Z
    by: alex
tags: [cli]
topics: [planner-agent]
touches: [flai/internal/harness, flai/internal/guard, flai/cmd/guard.go, flai/internal/itemnew, flai/internal/mcpserver/items_write_test.go, flai/internal/serve/plan_test.go, design/conventions, template, design/system/strategic-agents.md, docs/users/flai.md, design/system/flai-cli.md, docs/users/flai-reference.md, docs/users/flaiover.md, flai/cmd/plan.go, design/system/work-hierarchy.md, design/system/workflow.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md, design/issues/summary.md]
after: [S-0209]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 12611
  models:
    - model: claude-haiku-4-5-20251001
      input: 218
      output: 7788
      cache_read: 1182326
      cache_write: 65555
      cost: 0.2393
    - model: claude-opus-5-5
      input: 692
      output: 132699
      cache_read: 57444200
      cache_write: 910491
      cost: 20.5188
    - model: claude-sonnet-5
      input: 132
      output: 31540
      cache_read: 4503093
      cache_write: 213940
      cost: 1.7511
---
# S-0255 The planner drafts a story's tasks into the backlog and revisits the children it already has

## Goal

Given a story, the planner writes the tasks that deliver its outcome, in the backlog, and when the story already has tasks it revisits each open one: re-enriches it, proposes splits or merges, and reports what changed.

## Acceptance criteria
- [x] For a story with no tasks, the planner writes tasks with work, done when, nature, tags, touches, and `after` between them, each in the backlog, and opens one thread on the story summarizing the plan: the tasks, their order, and the assumptions it made
- [x] For a story with tasks, it revisits every task not done or cancelled: re-enriches it (the enrichment story), and proposes in the story's thread any task it would split, merge, add, or drop, creating tasks for additions but never cancelling or rewriting a finalized task's words without asking
- [x] Tasks it writes pass `flai check --strict` and the markdown lint before they are kept
- [x] Its log entry names the tasks created and revisited and the run's cost
- [x] Measured on E-0016 itself or a comparable epic: the designer judges the drafts on a thread and the result is recorded in `design/system/strategic-agents.md`

## Tasks
- T-0797 The planner's prompt drafts a story's tasks or revisits the ones it has, and names them in its summary
- T-0798 The guard lets the planner run flai task new
- T-0799 A task the planner writes that fails flai check --strict or the markdown lint is refused and leaves nothing
- T-0800 The conventions, the design, the docs, and the template say the planner drafts and revisits a story's tasks
- T-0801 The planner is run on a story, the designer judges its tasks on a thread, and the result is recorded in strategic-agents.md
- T-0802 Once S-0209 is accepted, item_new refuses a task that fails the check or the lint, and a story's planner entry names its tasks

## Notes

Criterion 1 named topics for tasks; tasks carry none, and the designer had it dropped (TH-0098): the planner puts a task's topic on the story. Criteria 3 and 4 for the MCP path rest on S-0209's checked `item_new` and activity items, with task cases added here (T-0802). Criterion 5: the operator ran the planner on S-0217 with this branch's prompt over flai 1.30.0 (TH-0108, six tasks in four layers) and resolved TH-0108 and TH-0102 without asking for a change; recorded in `design/system/strategic-agents.md` under "Measured on S-0217 (S-0255)". The revisit run was not made. Template 1.0.44.
