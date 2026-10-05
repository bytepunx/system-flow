---
id: S-0269
type: story
nature: improvement
title: "One command closes a task: flai task done commits, syncs, moves, logs, widens touches, checks, and answers the inbox"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:28Z
updated: 2026-10-05T01:35:28Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
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
