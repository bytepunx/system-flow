---
id: T-3026
type: task
nature: remediation
title: flai stream sync and the shell's flai task done refuse at once in a story's agent container, naming the MCP tools that sync on the host
status: backlog
parent: S-0689
owner: alex
created: 2026-11-10T10:32:53Z
updated: 2026-11-10T10:32:53Z
transitions: []
stream: S-0689
tags: [flai, git, container]
touches: [flai/cmd/verify.go, flai/cmd/stream.go, flai/cmd/stream_sync_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, docs/users/flai-reference.md]
---
# T-3026 flai stream sync and the shell's flai task done refuse at once in a story's agent container, naming the MCP tools that sync on the host

## Work

ADR-0246 says `flai stream sync` refuses in a story's container, but it refuses only through `Repo.RequireMain`, which keys on `MainAbsent`. ADR-0253 mounts the main checkout's `system-flow.yaml`, so `MainAbsent` is empty there and the sync runs against a read-only `.git` (I-0272). Add a story-container signal beside `(*app).orchestratorInRoleContainer` in `flai/cmd/verify.go`: the config path `FLAI_CONFIG` names lies under the main checkout's `.flai-cache/agents` (`harness.ContainerCache`), with `FLAI_STORY` set. Reuse `pathWithin`. Generalise the orchestrator's helper only if both read the same way.

Before any git call, `flai stream sync` (`flai/cmd/stream.go`) and the shell's `flai task done` (`flai/cmd/task_done.go`) refuse there with one error. The error says the main checkout's `.git` is read-only in the container, so no rebase or trial merge can be written. It names the MCP tool `task_done`, which closes a task and syncs on the host, and the MCP tool `verify`, whose close-out syncs on the host (ADR-0274). Add a line to `flai stream sync`'s help saying it refuses in a story's container, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

This is the first layer: it waits for no task.

## Done when

- In a test whose `FLAI_CONFIG` lies under `.flai-cache/agents/<S-nnnn>` with `FLAI_STORY` set, `flai stream sync` exits non-zero before any rebase, with an error that names `task_done` and `verify`, and the branch and HEAD are as they were.
- `flai task done` refuses the same way there.
- The same commands outside a container, and an orchestrator's `flai accept` in a role's container, behave as before.
- `TestSyncRefusesAWorktreeWhoseMainCheckoutIsAbsent` still passes.
- `flai test` passes on the changed paths, `TestReferenceIsCurrent` included.

## Notes
