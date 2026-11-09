---
id: T-2623
type: task
nature: improvement
title: The activity agent.status builds reads the conversations once, for every live run and for held, with the same answer as before
status: done
parent: S-0603
owner: alex
created: 2026-11-07T23:47:49Z
updated: 2026-11-08T02:04:26Z
transitions:
  - to: ready
    at: 2026-11-08T01:48:39Z
    by: system-flow
  - to: in-progress
    at: 2026-11-08T01:48:43Z
    by: system-flow
  - to: done
    at: 2026-11-08T02:04:26Z
    by: agent-S-0603
stream: S-0603
tags: [cli]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/messages/messages.go, flai/internal/workitem/hold.go, flai/internal/workitem/share.go]
after: [T-2621]
usage:
  source: log
  seconds: 943
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 40
      output: 6515
      cache_read: 7551617
      cache_write: 116192
      cost: 2.3488
---
# T-2623 The activity agent.status builds reads the conversations once, for every live run and for held, with the same answer as before

## Work

Criterion 2. `serve.Activity` in `flai/internal/serve/agents.go` reads every conversation under `wip/messages` once for each live run, through `awaitingOther` and `messages.AwaitingOther`, which lists them all (`flai/internal/messages/messages.go` `For`). `held` reads them again through `Repo.Holds` and `Repo.Conversations` (`flai/internal/workitem/hold.go`, `flai/internal/workitem/share.go`).

- Read the conversations once per `Activity` call, before `runActivity` and `held`.
- In `runActivity`, give each live run its awaited conversations from that one read: a form of `messages.AwaitingOther` that takes the list read, with the same filters (closed, awaiting, the other story in progress or in review).
- In `held`, build the holds with the shares and asks of that same read, through a form of `Repo.Holds` that takes them, in place of `Repo.Conversations` reading the files again. `Repo.Holds` keeps its signature for its other callers.
- Change nothing in what `Activity` answers.

Waits for T-2621: both change `flai/internal/serve/agents.go` and `flai/internal/serve/agents_test.go`, so they cannot run in one layer.

## Done when

- A test in `flai/internal/serve/agents_test.go` with three live runs, one waiting on a conversation, one on a thread, and one working, and a ready story held on overlap with a share in force on another, gives the same `Activity` answer as before the change, and shows the conversation files read once for the call.
- `flai test` on the changed paths passes.

## Notes
