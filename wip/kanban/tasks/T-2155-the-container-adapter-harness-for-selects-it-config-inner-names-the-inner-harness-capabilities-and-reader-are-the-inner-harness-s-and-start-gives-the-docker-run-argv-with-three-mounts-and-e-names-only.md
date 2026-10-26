---
id: T-2155
type: task
nature: feature
title: "The container adapter: harness.For selects it, config.inner names the inner harness, Capabilities and Reader are the inner harness's, and Start gives the docker run argv with three mounts and -e names only"
status: backlog
parent: S-0508
owner: alex
created: 2026-10-26T04:48:02Z
updated: 2026-10-26T04:48:02Z
transitions: []
stream: S-0508
tags: [cli, docker]
touches: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
---
# T-2155 The container adapter: harness.For selects it, config.inner names the inner harness, Capabilities and Reader are the inner harness's, and Start gives the docker run argv with three mounts and -e names only

## Work

Write `flai/internal/harness/container.go`: an `Adapter` named `container`, registered in `Adapters` in `harness.go` and listed by `Names()`. It reads `config.inner` from the agent, `claude-code` when unset, and refuses an inner harness `Adapters` lacks or that is `container` itself. `Capabilities()` returns the inner adapter's; `Reader()` the inner adapter's, since pilot relays the inner harness's log unchanged. `Start(req, host)` returns `Harness: "container"`, `Argv` of `docker run --rm -i --name flai-agent-<story or role>-<session> --user <uid>:<gid> --network <n> --add-host host.docker.internal:host-gateway`, three `--mount` or `-v` entries (the story's worktree read-write at its host path, `.flai-cache/agents/<S-nnnn>` at the same path, the login path read-only) and no other, `-e NAME` for each of `FLAI_AGENT`, `FLAI_SESSION`, `FLAI_STORY`, `FLAI_STARTED_BY`, `FLAI_MCP_URL`, `FLAI_MCP_TOKEN`, and each provider key's `key_env` (names only), the image from `host.Args` or the harness settings for the inner harness, and `Stdin` the JSON start request pilot reads, with `Inner`, `MCP`, `Cache`, and `ConfigDir` set. `Start.Env` and `EnvFrom` carry what the docker client must have in its own environment for `-e NAME` to pass. Test in `container_test.go`: the argv holds no key, no token, no `docker.sock`, and no mount of the home, `~/.flai`, or the main checkout; an unknown inner harness is refused; the readers and capabilities are the inner adapter's. This task waits for nothing.

## Done when

`flai test flai/internal/harness` passes with the new tests, and `harness.For` returns the container adapter for `harness: container`.

## Notes

`Start` may need a `Stdin` field; add it to `harness.Start` here, left empty by the other adapters, and let the serve task feed it to the process.
