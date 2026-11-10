---
id: T-2976
type: task
nature: improvement
title: flai serve readies, names, and stops a defined role's container by the role's name
status: backlog
parent: S-0705
owner: alex
created: 2026-11-10T10:29:09Z
updated: 2026-11-10T10:29:09Z
transitions: []
stream: S-0705
tags: [flai, serve]
touches: [flai/internal/serve/container.go, flai/internal/serve/container_test.go]
after: [T-2973]
---
# T-2976 flai serve readies, names, and stops a defined role's container by the role's name

## Work

Second layer, beside the pilot task, with no path in common. It waits for T-2973, whose exported container subject it calls.

- In `flai/internal/serve/container.go`, `containerReady` makes the cache folder and the hooks for the subject T-2973's function gives, so a defined role's is `.flai-cache/agents/<name>`, not its verb's.
- `AgentRun.containerSubject` gives `RoleName` for a defined role's run (`defined()`), so that `stopContainer` stops `flai-agent-<name>-<session>`, the name the adapter gave it, and `containerSaid` logs it by that name.
- Tests in `flai/internal/serve/container_test.go`, with docker faked: a defined role on the container harness starts, its cache folder is made under its name, and a stop runs `docker stop` on the container the start named.

## Done when

- A defined role's run on the container harness starts and stops through `flai serve`, by the role's name.
- `flai test` on the changed paths passes.

## Notes
