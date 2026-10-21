---
id: S-0450
type: story
nature: improvement
title: task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator
status: done
owner: alex
created: 2026-10-20T18:53:36Z
updated: 2026-10-21T00:29:28Z
transitions:
  - to: ready
    at: 2026-10-20T19:07:01Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-20T19:07:26Z
    by: agent-S-0450
  - to: review
    at: 2026-10-20T22:40:52Z
    by: system-flow
  - to: done
    at: 2026-10-21T00:29:28Z
    by: alex
tags: [flai, template]
topics: [mcp, inbox]
touches: [design/adrs/README.md, design/system/agent-narrative.md, design/system/flai-cli.md, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/storystart/start.go, flai/internal/storystart/start_test.go, flai/internal/mcpserver/story_start.go, flai/internal/mcpserver/story_start_test.go, flai/cmd/story_start.go, flai/cmd/story_start_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/session-start.md, template/root/design/conventions/session-start.md, design/issues/I-0167-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md, design/issues/summary.md, design/adrs/0179-task-done-and-story-start-carry-an-inbox-that-lists-the-story-s-threads-and.md, design/issues/I-0171-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4562
  priced_by: harness
  turns:
    - day: 2026-10-20
      test_runs: 1
      hand_edits: 1
      work: 86
  models:
    - model: claude-opus-5-5
      input: 460
      output: 137374
      cache_read: 25932737
      cache_write: 1301720
      cost: 15.8561
  strategic:
    - kind: planner
      seconds: 324
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 44
          output: 30293
          cache_read: 1973159
          cache_write: 153118
          cost: 0.1845
        - model: claude-opus-5-5
          input: 497
          output: 35373
          cache_read: 84798597
          cache_write: 231655
          cost: 19.7464
    - kind: orchestrator
      seconds: 7381
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 106
          output: 2022
          cache_read: 31485628
          cache_write: 46134
          cost: 7.9048
        - model: claude-sonnet-5-5
          input: 22
          output: 120
          cache_read: 397871
          cache_write: 58537
          cost: 0.3736
cost_of_delay:
  inputs:
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-10-20T18:53:36Z
  value: 30
  by: planner-S-0450
  at: 2026-10-20T19:06:22Z
forecast:
  duration: 37m
  delivery: 2026-10-21T10:40:00Z
  basis: "flai forecast: median 87 s per unit of size over 78 done large improvement stories on claude-opus-5-5, times size 25 (2 criteria, 23 touches), 30th in the pull order under an in-progress limit of 5."
  by: planner-S-0450
  at: 2026-10-20T19:06:22Z
finalized:
  by: orchestrator
  at: 2026-10-20T19:06:55Z
---
# S-0450 task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator

## Goal

This story remediates [I-0167](../../../design/issues/I-0167-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md), "task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0167 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0167 is closed with `flai issue close I-0167 --reason` saying what fixed it

## Tasks
- T-1914 Record in an ADR and the design that task_done and story_start carry the story's threads in full and only a count of the rest, fitted to 40,000 bytes
- T-1915 Give the inbox package a story-scoped embedded read with a count of the other threads, and a fit to a byte limit
- T-1916 Make story_start, MCP and CLI --json, carry the embedded inbox and fit its whole answer to 40,000 bytes
- T-1917 Make task_done, MCP and CLI --json, carry the embedded inbox and fit its whole answer to 40,000 bytes
- T-1918 Document the embedded inbox and its fit in the user docs, the generated reference, and the session-start convention
- T-1919 Close I-0167 with the reason naming the fix

## Notes

Cost of delay inputs set by flai from I-0167. time_lost_per_cycle 12m: 4m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-19T15:28:45Z, 1.1 days before this story; under one cycle counts as one).

### Planning

Touches, all files; no folder touch is kept. The story declared none, so `flai touches suggest` needed seed paths; it was run from the three packages the explorer found build the answers.

| Touch | Source |
|-------|--------|
| `flai/internal/inbox/inbox.go`, `flai/internal/storystart/start.go`, `flai/internal/taskdone/taskdone.go` | layout: where the inbox is built and embedded |
| `flai/internal/inbox/inbox_test.go`, `flai/internal/taskdone/taskdone_test.go`, `flai/internal/mcpserver/task.go`, `flai/internal/mcpserver/task_test.go`, `flai/cmd/task_done.go`, `docs/users/flai-reference.md`, `design/system/flai-cli.md` | co-change with the seeds |
| `flai/internal/storystart/start_test.go`, `flai/internal/mcpserver/story_start.go`, `flai/internal/mcpserver/story_start_test.go`, `flai/cmd/story_start.go`, `flai/cmd/story_start_test.go`, `flai/cmd/task_done_test.go` | layout: the MCP and CLI forms of the two answers and their tests |
| `design/adrs/README.md`, `design/system/agent-narrative.md` | design: the ADR the goal's proposal needs and § What an agent is told through MCP |
| `docs/users/flai.md`, `design/conventions/session-start.md`, `template/root/design/conventions/session-start.md` | design: the user docs and the convention sentence on what `story_start` answers |
| the I-0167 issue file, `design/issues/summary.md` | criterion 2: `flai issue close` |

Left out of the co-change list: `flai/internal/mcpserver/server.go`, `folder.go`, `cursor.go`, and `messages.go`, since the standalone `inbox`, `wait_for_events`, and `wait_for_work` are outside the issue; and the `stream` and `branch` commands, which co-change with `taskdone` for the sync step only. The new ADR's file is named by `flai adr new`; `flai task done` widens the touches with it.

Forecast 37m, delivery 2026-10-21T10:40:00Z, as `flai forecast` gives them after the touches were set (4m before, at size 2 with no touches). Kept: the median comes from 78 done large improvement stories on the same model, and size 25 already counts the six tasks' 23 files.

Cost of delay 30 USD a week, as `flai cod` gives it from flai's input of 12m lost per 168h cycle at 150 USD an hour. Kept. The issue now records 4 instances where the input counted 3, which would give 16m and 40 USD a week; the input is the operator's to change, so it is raised on the plan's thread, not changed here.
