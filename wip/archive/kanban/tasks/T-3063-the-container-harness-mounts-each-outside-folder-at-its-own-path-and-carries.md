---
id: T-3063
type: task
nature: improvement
title: The container harness mounts each outside folder at its own path and carries them to pilot, refusing what mountable refuses
status: done
parent: S-0701
owner: alex
created: 2026-11-10T10:35:37Z
updated: 2026-11-12T00:21:11Z
transitions:
  - to: ready
    at: 2026-11-12T00:05:46Z
    by: agent-S-0701
  - to: in-progress
    at: 2026-11-12T00:05:46Z
    by: agent-S-0701
  - to: done
    at: 2026-11-12T00:21:11Z
    by: agent-S-0701
stream: S-0701
tags: [harness, docker]
touches: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/pilot/request.go, flai/internal/pilot/request_test.go]
after: [T-3059]
usage:
  source: log
  seconds: 925
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 41
      output: 12133
      cache_read: 2323658
      cache_write: 104408
      cost: 1.3121
---
# T-3063 The container harness mounts each outside folder at its own path and carries them to pilot, refusing what mountable refuses

## Work

In `flai/internal/harness/container.go`, have `container.Start` add a read-write `-v <folder>:<folder>` for each outside folder in the start request, for a story's agent and for a strategic agent. Check each with `mountable`, extended so that it refuses an outside folder that is home, `~/.flai`, the main checkout, or the Docker socket, or that holds one of them. Keep the folders in `pilotRequest`, so that the inner claude-code adapter passes them as `--add-dir` inside the container. In `flai/internal/pilot/request.go`, make sure pilot reads and keeps them, and refuses a folder that is not absolute.

This task waits for T-3059, which adds the folders to the start request and the adapter's `--add-dir`.

## Done when

- Tests in `flai/internal/harness/container_test.go` show the `docker run` mounts for two outside folders, and a refusal for one that holds the main checkout.
- A test in `flai/internal/pilot/request_test.go` shows the folders reach the inner adapter's arguments as `--add-dir`.
- `flai test` passes on the changed paths.

## Notes
