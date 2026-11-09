---
id: T-2775
type: task
nature: improvement
title: Give the container adapter a docker run argv for one command in a story's container with its mounts
status: backlog
parent: S-0628
owner: alex
created: 2026-11-09T21:58:16Z
updated: 2026-11-09T21:59:10Z
transitions: []
stream: S-0628
tags: [flai, harness]
touches: [flai/internal/harness/container.go, flai/internal/harness/container_test.go]
---
# T-2775 Give the container adapter a docker run argv for one command in a story's container with its mounts

## Work

`flai/internal/harness/container.go` has no function that runs an arbitrary command in a story's container. `Start` and `containerStart` are tied to pilot's start request on standard input. Add one, such as `ContainerCommand(host, root, story, image, argv)`. It returns `docker run --rm --name flai-closeout-<story>-<id> --user <uid>:<gid> --network <network>`, built from `dockerRun`, `dockerUser`, `mountArgs`, and `createMounts`. It carries the story's mounts from `storygit.GitMounts` and `storygit.MainMounts`: the worktree read-write as `-w`, the main checkout's `.git` read-only with the registration read-write, `wip/` and `system-flow.yaml` read-only, and `.flai-cache/verify` read-write. It sets `--entrypoint` to the command's program, and the `safe.directory` variables `safeMain` gives. It gives no login mount, no MCP variables, no provider key, and no `.flai-cache/agents/<story>` folder. It makes the same refusals `Start` makes for a missing registration, a missing main checkout, or a path holding a colon.

Needs no earlier task. T-2777's close-out job and T-2782's checks run their commands through it.

## Done when

- A test in `flai/internal/harness/container_test.go` pins the argv: the six story mounts in order with their modes, the entrypoint, `-w`, and no login, home, `~/.flai`, Docker socket, or key.
- A test pins each refusal `Start` shares.
- `flai test flai/internal/harness/container.go flai/internal/harness/container_test.go` passes.

## Notes
