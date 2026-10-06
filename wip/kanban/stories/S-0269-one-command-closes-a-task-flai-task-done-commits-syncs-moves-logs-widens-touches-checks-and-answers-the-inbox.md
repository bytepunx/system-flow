---
id: S-0269
type: story
nature: improvement
title: "One command closes a task: flai task done commits, syncs, moves, logs, widens touches, checks, and answers the inbox"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:28Z
updated: 2026-10-06T18:12:48Z
transitions: []
tags: [cli, mcp]
topics: [automation, mcp, hostapi, conventions]
touches: [flai/cmd/items.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/cmd/move.go, flai/cmd/touches.go, flai/cmd/check.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/work-management.md, design/conventions/git.md, template/root/design/conventions/work-management.md, template/root/design/conventions/git.md, template/CHANGELOG.md, design/adrs, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  value: 427
  by: planner-E-0017
  at: 2026-10-06T11:36:15Z
forecast:
  duration: 55m
  delivery: 2026-10-07T07:46:00Z
  basis: "Its own forecast of 55m; 31st in the pull order with an in-progress limit of 3, behind S-0226, S-0295, S-0296, S-0284, S-0278, S-0223, S-0224, S-0227, S-0229, S-0212, S-0213, S-0214, S-0215, S-0216, S-0228, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254, S-0261, S-0264 and S-0265."
  by: flai
  at: 2026-10-06T18:12:48Z
---
# S-0269 One command closes a task: flai task done commits, syncs, moves, logs, widens touches, checks, and answers the inbox

## Goal

Closing a task today is a run of six to nine single-purpose model turns: tick its criteria, `git add -A && git commit`, `flai stream sync`, `flai move T-nnnn done`, `flai stream log`, `flai touches` for the paths the diff added, `flai check`, and `inbox`; about 2,900 such turns across 108 story runs, each re-reading the agent's whole context. `flai task done T-nnnn -m "<message>"`, and `task_done` over MCP and `task.done` on the host channel, does them in one call from the story's worktree, in that order, stopping at the first that fails with its findings, and answers what changed: the commit, the sync's result, the move, the touches it widened, the check's findings, and the inbox. The log entry is the message unless `--log` gives another. `flai stream sync` keeps refusing while a rebase is unfinished, and the move keeps its rules.

## Acceptance criteria
- [ ] `flai task done T-nnnn -m` performs commit, sync, move to done, narrative log, touches from the diff, `flai check`, and inbox in one call, stops at the first failure with its findings, and answers the result as text and `--json`
- [ ] The same operation is `task_done` over MCP and `task.done` on the host channel, with the same answer
- [ ] `design/conventions/work-management.md`, `git.md`, the template's copies, and the harness prompt send the agent to it at every task transition, and `design/system/flai-cli.md` and the user guide describe it
- [ ] Tests cover the happy path, a refused sync, and a failed check

## Tasks

## Notes

From the epic's log classification: flai move 719 turns, flai check 540, stream sync 486, stream log 417, git commit 405, inbox 336, touches 217, criteria ticks 119.

### Planning

Touches, none declared before. `flai touches suggest S-0269` was seeded with `flai/cmd/stream_sync.go`, `flai/cmd/move.go`, and `flai/internal/mcpserver`, which 96 commits changed:

- `flai/cmd/items.go`, `flai/cmd/task_done.go`, `flai/cmd/task_done_test.go`: layout. The task command is built in `items.go`; the subcommand and its tests are new files.
- `flai/cmd/stream.go`, `stream_sync.go`, `move.go`, `touches.go`, `check.go`: layout. The steps live here as command code, which must become callable from one command.
- `flai/internal/mcpserver/folder.go`, `server.go`: layout. Tools are registered in `folder.go`, and the server's instructions in `server.go` name the task transition.
- `flai/internal/hostapi/writes.go`, `writes_test.go`: co-change (9 and 8 of 96). `task.done` is a write method.
- `flai/internal/harness/harness.go`, `harness_test.go`: design. The story agent's prompt names each step's command.
- `design/conventions/work-management.md`, `git.md`, their `template/root` copies, and `template/CHANGELOG.md`: design (criterion 3).
- `design/adrs`: design. One call in place of ADR-0069's per-task cycle is a decision to record.
- `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`: co-change (30, 36, and 19 of 96) and criterion 3.
- Left out: `docs/operators/settings.md`, `flai/cmd/edit.go`, and `flai/internal/itemedit`, co-changed in 10 to 15% but not reached by the command.

Forecast: flai gave 41m (89 s per unit of size over 21 done large improvement stories, times size 27). I raised it to 55m because the steps live in command code that must become callable, and a refused sync needs a git fixture. S-0217, which composed three commands across the CLI, MCP, and the host channel, took 68m. The delivery is flai's, 2026-10-07T00:46Z, moved by the added 14m.

Cost of delay: 427 USD a week, against flai's 144.74. flai shares E-0017's 900 USD a week (the operator's 6h a cycle, TH-0171) by forecast duration. I shared it by the turns each story removes in the epic's evidence instead. This one removes about 2,900 of about 6,100, the most of any story.
