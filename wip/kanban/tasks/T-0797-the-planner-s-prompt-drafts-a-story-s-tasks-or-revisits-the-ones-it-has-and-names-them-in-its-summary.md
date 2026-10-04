---
id: T-0797
type: task
nature: feature
title: The planner's prompt drafts a story's tasks or revisits the ones it has, and names them in its summary
status: done
parent: S-0255
owner: alex
created: 2026-10-04T04:03:35Z
updated: 2026-10-04T04:07:52Z
transitions:
  - to: ready
    at: 2026-10-04T04:04:26Z
    by: agent-S-0255
  - to: in-progress
    at: 2026-10-04T04:04:26Z
    by: agent-S-0255
  - to: done
    at: 2026-10-04T04:07:52Z
    by: agent-S-0255
stream: S-0255
tags: []
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
usage:
  source: log
  seconds: 206
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 35
      output: 211
      cache_read: 1012040
      cache_write: 52966
      cost: 0.4344
---

# T-0797 The planner's prompt drafts a story's tasks or revisits the ones it has, and names them in its summary

## Work

`planPrompt` in `flai/internal/harness/harness.go` tells the planner, for a story, to enrich it as it does today and then: with no tasks, draft the tasks that deliver the story's outcome, each with `## Work` and `## Done when`, a nature, tags, `touches`, and `after` between them, in the backlog (`item_new` type task, or `flai task new`), and open one thread on the story summarising the plan (the tasks, their order and layers, and the assumptions made); with tasks, revisit each one not done or cancelled: re-enrich its `touches` and `after`, add tasks the outcome lacks, and propose in the story's thread any task it would split, merge, or drop, never cancelling a task or rewriting a task's words it did not write without asking. Tasks carry no topics (TH-0098): the planner adds to the story's topics a topic a task reaches that the story lacks. The final summary names the tasks created and the tasks revisited by ID. No path of another task.

## Done when

- `planPrompt` for a story says each of the above, and the epic's text is unchanged
- a test in `harness_test.go` fails without it: a story's plan prompt names drafting tasks into the backlog, the plan thread, revisiting open tasks without cancelling or rewriting without asking, and a summary naming the tasks created and revisited
- `go test ./internal/harness/` passes

## Notes
