---
id: T-2782
type: task
nature: improvement
title: Run the operator's checks for a container story in the story's container
status: backlog
parent: S-0628
owner: alex
created: 2026-11-09T21:59:00Z
updated: 2026-11-09T21:59:00Z
transitions: []
stream: S-0628
tags: [flai, serve, security]
touches: [flai/internal/serve/checks.go, flai/internal/serve/checks_test.go]
after: [T-2775, T-2777]
---
# T-2782 Run the operator's checks for a container story in the story's container

## Work

`RunChecksWith` in `flai/internal/serve/checks.go` runs each named check (`ResolveChecks`, S-0082) with `runOneCheck`. That calls `exec.Command` in the story's worktree on the host. Its callers are `flai checks run` (`flai/cmd/checks.go`) and the kluster pass after a deploy (`flai/internal/serve/kluster.go`). For a story whose agent's harness is `container`, run each check's argv through T-2775's argv instead, in the image T-2777 exposes, with the run's environment given as `-e NAME` and its values in the docker client's environment. Keep the process group, the cancel file, and the time limit. Cancelling a check must also `docker stop` its container by name, as a stop does for an agent (ADR-0058). Pass the harness in from the callers, or read it from the story, without changing what a host-harness story runs.

Waits for T-2775, whose argv it runs, and for T-2777, whose image choice it reuses.

## Done when

- A test in `flai/internal/serve/checks_test.go` shows that a container story's check runs as `docker run` with the story's mounts, and that no check program runs on the host.
- A test shows a cancelled check stops its container by name.
- A host-harness story's checks run as before.
- `flai test` on the two paths passes.

## Notes
