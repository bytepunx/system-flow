---
id: T-3050
type: task
nature: improvement
title: flai serve opens one thread to the operator on the board when nothing can be made ready, updates it when the list changes, and resolves it once a story can start
status: done
parent: S-0691
owner: alex
created: 2026-11-10T10:35:01Z
updated: 2026-11-11T05:45:32Z
transitions:
  - to: ready
    at: 2026-11-11T05:26:32Z
    by: agent-S-0691
  - to: in-progress
    at: 2026-11-11T05:26:33Z
    by: agent-S-0691
  - to: done
    at: 2026-11-11T05:45:32Z
    by: agent-S-0691
stream: S-0691
tags: [flai, go, serve]
touches: [flai/internal/serve/starved.go, flai/internal/serve/starved_test.go, flai/internal/serve/agents_launch.go]
after: [T-3046]
usage:
  source: log
  seconds: 1139
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 52
      output: 23749
      cache_read: 3532107
      cache_write: 150317
      cost: 2.2086
---
# T-3050 flai serve opens one thread to the operator on the board when nothing can be made ready, updates it when the list changes, and resolves it once a story can start

## Work

Add `flai/internal/serve/starved.go`. Call it from `look()` in `flai/internal/serve/agents_launch.go`, after `noteReadyEmpty`. When the ready column holds no story that can start, by the count S-0673 gives `readyDrained`, and the in-progress limit has room, and `PromotionCandidates` lists no candidate, it does the following:

- With no open thread of its own on `wip/kanban/board.md`, it opens one with `threads.New`, by flai serve's author, as `agents_restart.go`, `forge.go`, and `kluster.go` open theirs. The title says how many lanes stand idle. The text gives T-3046's list, and names `flai serve notify set` when the host has no reminder channel.
- With such a thread open, it adds an entry with `threads.Reply` only when the list's text changed since the thread's last entry.
- Once a story can start again, it resolves the thread with `threads.Resolve`, saying which story.

Keep what it last wrote with the ready-empty state under `.flai-cache`, so that a restart of flai serve neither opens a second thread nor writes the same list again.

It waits on T-3046, whose list it writes.

## Done when

- `starved_test.go` reproduces I-0288. The ready column is empty and the in-progress limit has room. The backlog holds only two drafts, each with a thread awaiting the operator, and stories held by `after` behind them. One look opens one thread naming both drafts and their threads. A second look writes nothing. The test fails when the call in `look()` is removed.
- A change to the list adds one entry. A story that enters ready and can start resolves the thread. A host with no reminder channel gets the `flai serve notify set` line.
- `flai test flai/internal/serve/starved.go flai/internal/serve/starved_test.go flai/internal/serve/agents_launch.go` passes.

## Notes
