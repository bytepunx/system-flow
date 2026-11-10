---
id: T-3017
type: task
nature: remediation
title: Mount the main checkout's checked-write lock file in the close-out's container
status: backlog
parent: S-0679
owner: alex
created: 2026-11-10T10:31:52Z
updated: 2026-11-10T10:31:52Z
transitions: []
stream: S-0679
tags: [flai, container]
touches: [flai/internal/harness/container.go, flai/internal/harness/container_test.go]
---
# T-3017 Mount the main checkout's checked-write lock file in the close-out's container

## Work

The project's default agent runs in a container (`agent.image.layer`), and there the close-out's `flai verify --record-issues` runs inside it (ADR-0274). `ContainerCommand` in `flai/internal/harness/container.go` mounts the main checkout's `wip`, `system-flow.yaml`, and `.flai-cache/verify` (ADR-0253), not `.flai-cache/items-checked.lock`. So inside, `LockCheckedWrites` either fails or locks a file of its own that no checked write on the host holds, and T-3014's confirmation would not exclude a checked write there.

Mount `.flai-cache/items-checked.lock` of the main checkout read-write at its host path, made on the host when absent, as `.flai-cache/verify` is. Check first which purposes `ContainerCommand` serves and mount it for those that record issues, the close-out's verify at least. Check on this host that an `flock` taken in the container excludes one taken on the host across the bind mount (Docker Desktop on WSL here); if it does not, say so on the plan's thread before going on, since T-3014's lock then reaches nothing inside the container.

Add a test in `flai/internal/harness/container_test.go` that the close-out's container command carries the mount and that the file is made when absent.

This task waits for nothing: it is the first layer, beside T-3014, with no path in common.

## Done when

- The close-out's container command mounts the main checkout's `.flai-cache/items-checked.lock` read-write, made on the host when absent, and the test says so.
- A lock taken in the container excludes a checked write on the host, or the plan's thread says why it cannot.
- `flai test` on the changed paths passes.

## Notes
