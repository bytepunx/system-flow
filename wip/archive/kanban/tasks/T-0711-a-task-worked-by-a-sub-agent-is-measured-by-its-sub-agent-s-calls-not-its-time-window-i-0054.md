---
id: T-0711
type: task
nature: remediation
title: A task worked by a sub-agent is measured by its sub-agent's calls, not its time window (I-0054)
status: done
parent: S-0230
owner: arobson
created: 2026-10-02T17:15:23Z
updated: 2026-10-02T17:22:41Z
transitions:
  - to: ready
    at: 2026-10-02T17:16:14Z
    by: agent-S-0230
  - to: in-progress
    at: 2026-10-02T17:16:14Z
    by: agent-S-0230
  - to: done
    at: 2026-10-02T17:22:41Z
    by: agent-S-0230
stream: S-0230
tags: []
touches: [flai/internal/usage, flai/internal/serve]
usage:
  source: log
  seconds: 387
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 76
      output: 24962
      cache_read: 2989907
      cache_write: 97119
      cost: 1.718
---
# T-0711 A task worked by a sub-agent is measured by its sub-agent's calls, not its time window (I-0054)

## Work

Fix I-0054. Today `serve.Measure` gives each task `rec.Windows(spans)`: each session's reported totals in the share of its calls' weight that fall in the task's in-progress spans (ADR-0051), so two tasks in progress at once each get the whole window.

Claude Code's stream-json gives every call a sub-agent makes `parent_tool_use_id`, the ID of the story's agent's `Agent` (or `Task`) `tool_use` that started it, in the same session as the story's agent; the sub-agent's events also carry `task_description`, the `Agent` call's description. The story's agent's own calls have `parent_tool_use_id` null.

Measure each task as:

1. the calls of every sub-agent started for it, wherever they fall in time. A sub-agent is started for a task when its `Agent` call's description names one of the story's task IDs (the first that it names), or, when the description names none, its prompt does (the first of the story's task IDs the prompt names). `task_description` stands in for the description when the `tool_use` was not seen.
2. plus its share of every other call (the story's agent's own, and sub-agents started for no task) made while it was in progress, split evenly among the tasks in progress at that call, so that tasks in progress at once share those calls rather than each getting them all.

Weigh calls as `Window` does (input plus cache tokens against each session and model's whole), so a task's share of each session's reported totals is its calls' weight over the whole's, and uncovered calls are estimated as now. Seconds stay the time the task was in progress. The story's total does not change. A task worked with no sub-agent and alone in its window gets what it gets today.

Done in `flai/internal/usage` (parse `parent_tool_use_id`, `task_description`, and the `Agent`/`Task` `tool_use` inputs; a method that measures tasks from their spans and IDs) and `flai/internal/serve/usage.go` (`Measure` uses it). Keep the package comments true.

## Done when

- [x] Tests in `flai/internal/usage` and `flai/internal/serve` show two tasks in progress at once, each worked by its own sub-agent, get their own sub-agents' calls and half of the story's agent's calls in their overlap, and that the tasks sum to no more than the story
- [x] A test shows a sub-agent named for a task by its prompt only, and one named for no task, is measured as above
- [x] The existing usage tests still pass unchanged in what they assert for tasks worked one at a time

## Notes

Waits for nothing: layer 1, beside T-0712, with which it shares no path.
