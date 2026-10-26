---
id: T-2163
type: task
nature: feature
title: The container adapter mounts a strategic role's checkout and cache folder as the ADR says, and pilot's request carries Role, Item, and Focus so the inner harness starts with --agent
status: backlog
parent: S-0509
owner: alex
created: 2026-10-26T04:49:31Z
updated: 2026-10-26T04:49:31Z
transitions: []
stream: S-0509
tags: [cli, docker]
touches: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/pilot/run.go, flai/internal/pilot/run_test.go]
after: [T-2158]
---
# T-2163 The container adapter mounts a strategic role's checkout and cache folder as the ADR says, and pilot's request carries Role, Item, and Focus so the inner harness starts with --agent

## Work

In `container.go`, when `Request.Role` is set, replace the worktree mount with the main checkout mounted as the ADR of T-2158 settles (read-write at its host path, or read-only), the cache folder `.flai-cache/agents/<role>` in place of the story's, the container's working directory the checkout, and the name `flai-agent-<role>-<session>`; `FLAI_STORY` is not passed for a role. Keep `Role`, `Item`, and `Focus` in the start request on standard input, and in `pilot/run.go` make the inner harness's start carry them so the claude-code adapter's `Start` adds `--agent <role>` and the role's prompt as on the host (`adapters.go`). Test in `container_test.go` the argv and mounts for each of the four roles and in `run_test.go` that a role's request reaches the inner start with `--agent`. Waits for T-2158, which settles the mounts.

## Done when

`flai test flai/internal/harness flai/internal/pilot` passes with a test per role on the mounts and on `--agent`.

## Notes

The orchestrator runs for as long as its host action is on; nothing here differs for it but the stop by name, which the next task covers.
