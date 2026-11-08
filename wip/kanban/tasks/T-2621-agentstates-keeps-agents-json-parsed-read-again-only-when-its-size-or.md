---
id: T-2621
type: task
nature: improvement
title: AgentStates keeps agents.json parsed, read again only when its size or modification time changed, and holds agentsMu only to read or replace it
status: backlog
parent: S-0603
owner: alex
created: 2026-11-07T23:47:30Z
updated: 2026-11-07T23:47:30Z
transitions: []
stream: S-0603
tags: [cli]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
---
# T-2621 AgentStates keeps agents.json parsed, read again only when its size or modification time changed, and holds agentsMu only to read or replace it

## Work

Criterion 1. In `flai/internal/serve/agents.go`, `Dir.AgentStates` reads and parses the host's whole `agents.json`, 663 KB, under the host-wide `agentsMu` on every call, so the 13 projects' `agent.status` calls queue on it.

- Keep the parsed states per `agents.json` path, with the file's size and modification time when it was read.
- `AgentStates` stats the file. It parses again only when the size or the modification time differs from what is kept, or nothing is kept. Parse outside `agentsMu`; hold the lock only to read what is kept or to replace it.
- `updateAgent` replaces what is kept with the states it wrote, and the stat after its write, so the next read does not parse again. `flai serve agent restart` writes from a process of its own (S-0116): its write changes the stat, and the next read parses again.
- Callers must not change the kept states through the map or the runs it returns. Return a copy, or show that no caller writes to them, and say which in the narrative's `## Decisions`.

Waits for nothing: the first layer.

## Done when

- A test in `flai/internal/serve/agents_test.go` shows that two reads with the file unchanged parse it once, that a write by `updateAgent` is read back without a parse, and that a write from outside, with a new size or modification time, is read on the next call.
- A test shows that concurrent `AgentStates` calls on an unchanged file do not wait on one another's parse.
- `flai test` on the changed paths passes.

## Notes
