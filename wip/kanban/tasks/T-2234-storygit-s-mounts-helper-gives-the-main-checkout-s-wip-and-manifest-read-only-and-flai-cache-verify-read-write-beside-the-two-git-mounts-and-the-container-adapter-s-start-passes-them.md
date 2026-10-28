---
id: T-2234
type: task
nature: remediation
title: storygit's mounts helper gives the main checkout's wip and manifest read-only and .flai-cache/verify read-write beside the two git mounts, and the container adapter's Start passes them
status: backlog
parent: S-0520
owner: alex
created: 2026-10-28T13:57:41Z
updated: 2026-10-28T13:57:41Z
transitions: []
stream: S-0520
tags: [cli, docker]
touches: [flai/internal/storygit/mounts.go, flai/internal/storygit/mounts_test.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go]
after: [T-2232]
---
# T-2234 storygit's mounts helper gives the main checkout's wip and manifest read-only and .flai-cache/verify read-write beside the two git mounts, and the container adapter's Start passes them

## Work

Waits for T-2232, whose ADR settles the mounts. S-0517's T-2219 writes `storygit.GitMounts` in `flai/internal/storygit/mounts.go`, which reads a story worktree's `.git` file and `commondir` and gives the main checkout's `.git` read-only and the story's registration read-write. Beside it, give the three mounts this story adds, from the main checkout's path the same `.git` file names: `<main>/wip` read-only, `<main>/system-flow.yaml` read-only, and `<main>/.flai-cache/verify` read-write, each at its host path, with the mode on each; create `<main>/.flai-cache/verify` on the host when it is absent, with the caller's uid, since Docker would create a missing bind-mount source as root and the container's `--user` could then not write the report. Whether the three join `GitMounts`'s list or a sibling function gives them is the story's agent's choice; the adapter and `scripts/agent-container.sh` (T-2235) take the same list. Test it with the real-git fixture `mounts_test.go` has from S-0517: the five mounts, their modes, and the folder created. In `flai/internal/harness/container.go`, S-0517's T-2221 makes `Start` take the git mounts; make it take these too, so `docker run` carries eight mounts in all: the worktree, the cache folder, the login, the two git mounts, and the three here. Extend the argv test in `container_test.go`: `wip` and `system-flow.yaml` are `:ro`, `.flai-cache/verify` is read-write, nothing else of the main checkout and nothing under `.flai-cache` but the story's agent folder, its worktree, and `verify` is mounted, and no socket.

## Done when

`go test ./internal/storygit/ ./internal/harness/` in `flai/` passes with the new cases; the argv test fails when a mount of the main checkout outside the ADR's list is added; `gofmt`, `go vet`, and golangci-lint pass on the changed files.

## Notes

`mounts.go` and `container.go` are S-0517's and S-0508's files, which is why the story waits on S-0517 (and through it S-0508) with `after`. `scripts/agent-container.sh` gets the same mounts in T-2235, which shares no path with this task.
