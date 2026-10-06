---
id: S-0269
type: story
nature: improvement
title: "One command closes a task: flai task done commits, syncs, moves, logs, widens touches, checks, and answers the inbox"
status: ready
parent: E-0017
owner: alex
created: 2026-10-05T01:35:28Z
updated: 2026-10-06T23:33:22Z
transitions:
  - to: ready
    at: 2026-10-06T22:47:57Z
    by: alex
tags: [cli, mcp]
topics: [automation, mcp, hostapi, conventions, git]
touches: [flai/cmd/items.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/move.go, flai/cmd/touches.go, flai/cmd/check.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/cursor.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/work-management.md, design/conventions/git.md, template/root/design/conventions/work-management.md, template/root/design/conventions/git.md, template/CHANGELOG.md, design/adrs, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 427
  by: planner-E-0017
  at: 2026-10-06T11:36:15Z
forecast:
  duration: 80m
  delivery: 2026-10-07T01:37:00Z
  basis: "Its own forecast of 1h20m; 5th in the pull order with an in-progress limit of 3, behind S-0300, S-0302, S-0273, S-0303 and S-0228."
  by: flai
  at: 2026-10-06T23:33:22Z
finalized:
  by: alex
  at: 2026-10-06T22:47:53Z
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
- T-1080 An ADR and flai-cli.md set flai task done's steps, stop rule, answer, and its MCP and host-channel names
- T-1085 The story-branch sync is a storygit function that cmd's stream sync calls and that answers a structured result
- T-1092 An agent's inbox is computed by a package that the MCP server and the CLI both call
- T-1098 A taskdone package closes a task in the ADR's order and stops at the first step that fails, with its findings
- T-1102 flai task done T-nnnn -m closes a task from the story's worktree and answers the result as text and --json
- T-1107 The MCP tool task_done closes a task with the same answer as flai task done --json
- T-1113 The host channel's task.done runs flai task done --json and answers its result
- T-1116 The conventions, their template copies, and the harness prompt send the agent to flai task done at every task transition
- T-1119 The user guide describes flai task done, task_done, and task.done, and the reference is regenerated

## Notes

From the epic's log classification: flai move 719 turns, flai check 540, stream sync 486, stream log 417, git commit 405, inbox 336, touches 217, criteria ticks 119.

### Planning

Planned by planner-S-0269 on 2026-10-06. The plan thread on S-0269 gives the tasks and their layers.

Tasks, in five layers:

- Layer 1: T-1080 (the ADR and `flai-cli.md`), T-1085 (the sync moves into storygit), T-1092 (the agent inbox moves into a package).
- Layer 2: T-1098 (the `taskdone` package and its three tests).
- Layer 3: T-1102 (the CLI) and T-1107 (the MCP tool).
- Layer 4: T-1113 (`task.done` on the host channel) and T-1116 (the conventions, the template, and the harness prompt).
- Layer 5: T-1119 (the user guide and the reference).

Touches. I kept all 23 declared ones (the epic planner's, from `flai touches suggest`, the layout, the criteria, and co-change) and added ten for files the code layout showed. The explorer read the code:

- Layout: `flai/cmd/stream_sync_test.go`, `flai/internal/storygit/sync.go`, `sync_test.go`. The sync is in methods of cmd's `app` (`newStreamSyncCmd` in `stream.go`; `checkSync`, `trialMerge`, and `reportConflicts` in `stream_sync.go`). The MCP server cannot import cmd, so the sync moves to storygit (T-1085).
- Layout: `flai/internal/inbox/inbox.go`, `inbox_test.go`, `flai/internal/mcpserver/cursor.go`. The agent's inbox and its cursor are built only in `mcpserver` (`server.go`, `cursor.go`), and there is no `flai inbox` command, so the CLI could not answer one (T-1092).
- Layout: `flai/internal/taskdone/taskdone.go`, `taskdone_test.go`. This is one composition that all three surfaces call, as `criteria_tick` calls `itemedit` (T-1098).
- Layout: `flai/internal/mcpserver/task.go`, `task_test.go`. The new tool's file, following `criteria.go` (T-1107).
- Declared and unchanged: `flai/internal/hostapi/writes.go` holds the host write specs that build flai arguments (`item.move`, `item.criteria`), so `task.done` is a spec there (T-1113).
- Kept as a possible change: `flai/cmd/move.go`, `touches.go`, `check.go`. They change only where a rule `taskdone` needs sits in a RunE.
- Left out, as the epic planner did: `docs/operators/settings.md`, `design/system/flaiover-dashboard.md`, and `docs/operators/index.md`, co-changed in 12 to 17% of commits but not reached. `design/system/dashboard-host-channel.md` does not list the host methods; `flai-cli.md` does.
- Folder touch kept: `design/adrs`, since T-1080's ADR takes its number when it is written.
- Topic added: `git`, which T-1085 and T-1098 reach through the story-branch sync and the commit.

Forecast: 80m, against flai's 38m (83 s per unit of size over 25 done large improvement stories, times size 27). I raised it because the plan has nine tasks in five layers. Two of them move code out of cmd and mcpserver into packages before the composition can be written. The tests need git fixtures for a refused sync. S-0217, which composed three commands across the CLI, MCP, and the host channel with less to move, took 68m. The delivery is flai's, 2026-10-07T00:22Z, moved on by the 42m added.

Cost of delay: 427 USD a week stands, against flai's 129.58. flai shares E-0017's 900 USD a week (the operator's 6h a cycle, TH-0171) by forecast duration. The epic's planner shared it by the turns each story removes, as the epic's evidence counts them. This story removes about 2,900 of about 6,100, the most of any story. That share still holds, and the figure is unchanged. The values of E-0017's stories now sum to 959, a little over 900, because S-0293 was valued after the split.
