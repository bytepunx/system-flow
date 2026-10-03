---
id: T-0745
type: task
nature: improvement
title: An epic follows its story's move in the same write
status: done
parent: S-0200
owner: alex
created: 2026-10-03T07:11:09Z
updated: 2026-10-03T07:19:02Z
transitions:
  - to: ready
    at: 2026-10-03T07:11:52Z
    by: agent-S-0200
  - to: in-progress
    at: 2026-10-03T07:11:52Z
    by: agent-S-0200
  - to: done
    at: 2026-10-03T07:19:02Z
    by: agent-S-0200
stream: S-0200
tags: []
touches: [flai/internal/workitem]
usage:
  source: log
  seconds: 430
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 51
      output: 18137
      cache_read: 3139020
      cache_write: 70338
      cost: 1.4206
    - model: claude-sonnet-5
      input: 54
      output: 16609
      cache_read: 1920718
      cache_write: 103164
      cost: 0.8083
---
# T-0745 An epic follows its story's move in the same write

## Work

Add the rule in `flai/internal/workitem`. `EpicFollows` computes the state an epic's stories put it in, cancelled stories left out: done when all are done, review when all are in review or done, in-progress when any has started, ready when any is ready, else backlog. `Follow` walks the epic one allowed transition at a time: forward when the story moved forward or was cancelled and its stories put the epic further on, back when the story moved back or came back from cancelled and no other story holds the epic where it is. It never takes an epic into done except in acceptance, and leaves an epic that is done, cancelled, or archived alone. Each step has the story's actor and time and a reason naming the story. `TransitionAll` applies it before anything is written and reports it as `MoveResult.Followed`. `Changes` names the story an epic's move followed, found as `CancelledWith` finds its cause.

It waits for nothing: the other tasks build on its API.

## Done when

- [x] `TransitionAll` moves the epic to ready and in-progress with its first story, to review with its last open one, and back only when no other story holds it
- [x] Behavior tests cover each transition, the back moves, cancelled stories, and an epic that is done, cancelled, or archived
- [x] `go test ./internal/workitem` passes

## Notes
