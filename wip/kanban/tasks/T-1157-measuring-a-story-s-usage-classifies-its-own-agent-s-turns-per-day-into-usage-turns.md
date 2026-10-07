---
id: T-1157
type: task
nature: feature
title: Measuring a story's usage classifies its own agent's turns per day into usage.turns
status: done
parent: S-0293
owner: alex
created: 2026-10-07T09:28:11Z
updated: 2026-10-07T09:36:52Z
transitions:
  - to: ready
    at: 2026-10-07T09:28:44Z
    by: agent-S-0293
  - to: in-progress
    at: 2026-10-07T09:28:46Z
    by: agent-S-0293
  - to: done
    at: 2026-10-07T09:36:52Z
    by: agent-S-0293
stream: S-0293
tags: []
touches: [flai/internal/usage/turns.go, flai/internal/usage/turns_test.go, flai/internal/usage/log.go, flai/internal/usage/log_test.go, flai/internal/usage/usage.go, flai/internal/usage/usage_test.go, flai/internal/workitem/usage.go, flai/internal/workitem/usage_test.go, flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/fields_test.go, flai/internal/workitem/store_test.go]
usage:
  source: log
  seconds: 486
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 65
      output: 30380
      cache_read: 4310255
      cache_write: 127872
      cost: 2.3022
---
# T-1157 Measuring a story's usage classifies its own agent's turns per day into usage.turns

## Work

`usage.Read` already pairs the story agent's own calls with their results. Extend it to record each turn of the session's own agent: an assistant message, by its message id, without `parent_tool_use_id`, with at least one tool call. Keep each call's tool name, its Bash `command`, and its Edit/Write/MultiEdit `file_path`, and the turn's timestamp. When the record is totalled, put each turn in one class, first match wins: test run, hand edit, empty wake, ceremony, work (rules in the story's narrative, `## Decisions`). Count them per UTC day into `Usage.Turns`, written to front matter as `usage.turns`, a list of days, each with `day` and the counts of its classes, omitted when empty. A task carries none; an epic's roll-up sums its stories' by day. Waits for nothing: first layer.

## Done when

- `Record.Total` gives a story's usage its turns per day and class, and a test with a stream-json fixture shows each class, a sub-agent's calls left out, and a turn that mixes classes taking the first that matches.
- `usage.turns` round-trips through front matter, is listed in `front-matter-fields.txt`, sums by day in `Add`, is copied by `Clone`, and compared by `Same`.

## Notes
