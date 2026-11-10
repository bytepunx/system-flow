---
id: T-2973
type: task
nature: improvement
title: The container adapter carries a defined role's name and inputs to pilot and names its container and cache by the role's name
status: backlog
parent: S-0705
owner: alex
created: 2026-11-10T10:28:54Z
updated: 2026-11-10T10:28:54Z
transitions: []
stream: S-0705
tags: [flai, harness]
touches: [flai/internal/harness/harness.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go]
---
# T-2973 The container adapter carries a defined role's name and inputs to pilot and names its container and cache by the role's name

## Work

First layer: it waits for no task. The other code tasks read the request shape it writes.

- Write the failing test first, in `flai/internal/harness/container_test.go`: a `container` start of a request with `RoleFile` set (a defined role, its verb, an item, and inputs) is refused today by `containerSubject` with "the container harness cannot carry to pilot". That test reproduces I-0258.
- In `flai/internal/harness/harness.go`, give `Request` a field pilot reads, such as `RoleName string` with `json:"role_name,omitempty"`, the name of the defined role's file, set when `RoleFile` is. Keep `RoleFile` `json:"-"`: pilot reads the file from the mounted checkout, as it reads a built-in one (ADR-0247 §3). Let `Inputs` reach pilot in the start request, since a defined role's prompt names them and the checkout does not hold them.
- In `flai/internal/harness/container.go`:
  - `containerSubject` lets a request with `RoleFile` through, with its role's verb, and still refuses any other role that is not one of `conventions.StrategicRoles`.
  - Export one function that gives a run's container subject, the story's ID, a defined role's name, or the built-in role's verb, so that `flai serve` names and stops the container by the same subject. Use it in `roleStart` for `ContainerCache` and `ContainerName`, so the cache folder is `.flai-cache/agents/<name>` and the container `flai-agent-<name>-<session>`.
  - `pilotRequest` keeps `RoleName` and `Inputs` and drops `RoleFile`. `pilotStdin`'s check of the inner harness's start must still see the role file, so it refuses what pilot would refuse; give it the request before the file is dropped.
- Test: a defined role's container start mounts what a built-in role's does, in that order; names the container and cache by the role's name; and writes a start request on stdin with `role_name`, the role's verb, its item, and its inputs, and no role file. A role that is neither built in nor defined is still refused.

## Done when

- The new test fails on main and passes with the change.
- `flai test` on the changed paths passes.
- The container adapter starts a defined role, with its name and inputs in pilot's start request.

## Notes
