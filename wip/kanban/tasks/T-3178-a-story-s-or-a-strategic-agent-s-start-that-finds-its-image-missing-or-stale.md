---
id: T-3178
type: task
nature: feature
title: A story's or a strategic agent's start that finds its image missing or stale waits on an image job, recorded on its run, in place of refusing
status: backlog
parent: S-0751
owner: alex
created: 2026-11-12T16:01:44Z
updated: 2026-11-12T16:02:29Z
transitions: []
stream: S-0751
tags: [cli]
touches: [flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents.go]
after: [T-3177]
---
# T-3178 A story's or a strategic agent's start that finds its image missing or stale waits on an image job, recorded on its run, in place of refusing

## Work

In `containerAllowed` and `containerImage` (`flai/internal/serve/container.go`, lines 108 to 143), when S-0735's status says an image is missing or stale and it is not the operator's own:

- request an image job from T-3177's runner, in place of refusing;
- record the run as waiting for its image, with the job's id, in `agents.json` (`AgentRun` in `agents.go`, lines 206 to 216);
- go on with the start in `startOf` (`agents_start.go`, line 398) when the job passes;
- fail the run with the job's cause, step, and log file when it fails (`agents_start.go`, around line 179).

An image named by the host's `--image` still refuses as today. A strategic agent waits for the layer, or starts on the harness image, as S-0733's ADR decides.

Tests in `container_test.go` cover:

- a missing image built and the start gone on;
- a stale image after an upgrade, whose label names another flai version;
- a failed job failing the run with its log;
- an own image refused;
- a strategic start.

`TestAContainerStartIsRefusedWhatItLacks` changes for the missing image alone.

It waits for T-3177, whose runner it calls.

## Done when

- `go test ./internal/serve/ -run Container` passes, under `-race` too.
- No path in `container.go` names `make agent-images`.

## Notes
