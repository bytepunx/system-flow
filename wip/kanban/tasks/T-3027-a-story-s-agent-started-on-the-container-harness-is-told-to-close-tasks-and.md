---
id: T-3027
type: task
nature: remediation
title: A story's agent started on the container harness is told to close tasks and sync through the MCP tools task_done and verify, not flai stream sync in its shell
status: backlog
parent: S-0689
owner: alex
created: 2026-11-10T10:33:06Z
updated: 2026-11-10T10:33:06Z
transitions: []
stream: S-0689
tags: [flai, container]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go]
---
# T-3027 A story's agent started on the container harness is told to close tasks and sync through the MCP tools task_done and verify, not flai stream sync in its shell

## Work

The start prompt `harness.Prompt` builds (`flai/internal/harness/harness.go`, the `fmt.Sprintf` near line 606 and the pre-review steps near lines 968 and 981) is one text for host and container. It tells the agent to close tasks with `flai task done` in the worktree and to run `flai stream sync` before review, and neither works in the container (I-0272). Give `Request` a field saying the agent runs in the container, set by the container adapter (`flai/internal/harness/container.go`, where `pilotStdin` asks the inner adapter for its start) and read by `Prompt`. On that harness the prompt says to close each task with the MCP tool `task_done` and to close the story out with the MCP tool `verify`, both of which sync on the host, and that `flai stream sync`, the shell's `flai task done`, and `git rebase` cannot write the read-only `.git` there. The host's prompt stays as it is.

It shares no file with T-3026 and waits for no task: first layer.

## Done when

- A test renders the prompt for a container request and finds `task_done` and `verify` named for the sync, and no instruction to run `flai stream sync` or the shell's `flai task done`.
- A test renders the prompt for a host request and finds it unchanged.
- `flai test` passes on the changed paths.

## Notes
