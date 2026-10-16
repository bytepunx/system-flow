---
id: T-1620
type: task
nature: feature
title: flai serve starts a story's agent when flai sent the story back with a thread on a failing check or a requested change
status: backlog
parent: S-0394
owner: alex
created: 2026-10-16T04:37:02Z
updated: 2026-10-16T04:37:02Z
transitions: []
stream: S-0394
tags: [cli]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
after: [T-1617]
---
# T-1620 flai serve starts a story's agent when flai sent the story back with a thread on a failing check or a requested change

## Work

- In `flai/internal/serve/agents.go`, beside `answered`, start the story's agent again, in its last session, when T-1617 has opened a thread on the story and moved it to in-progress. The prompt names the thread, and the agent's first `inbox` lists it.
- The start does not count toward `agent.auto_restarts`: it answers a thread, it does not retry a failure.
- Respect the in-progress limit as a restart does. A story sent back from review takes the place it had.
- Test it: a sent-back story gets its agent once, and a second look does not start a second.
- Waits for T-1617: it reacts to that task's thread and move.

## Done when

A story T-1617 sent back gets its agent started once with the thread in its prompt, and the serve tests pass.

## Notes
