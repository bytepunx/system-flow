---
id: T-2919
type: task
nature: remediation
title: scripts/agent-image.sh --check runs the smoke tier inside the built layer and fails on an image that lacks a tool it runs
status: done
parent: S-0672
owner: alex
created: 2026-11-10T08:48:14Z
updated: 2026-11-10T10:08:53Z
transitions:
  - to: ready
    at: 2026-11-10T10:00:39Z
    by: agent-S-0672
  - to: in-progress
    at: 2026-11-10T10:00:40Z
    by: agent-S-0672
  - to: done
    at: 2026-11-10T10:08:53Z
    by: agent-S-0672
stream: S-0672
tags: [docker, scripts]
touches: [scripts/agent-image.sh, flai/images/system-flow/Dockerfile]
after: [T-2917]
usage:
  source: log
  seconds: 493
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 6013
      cache_read: 3990031
      cache_write: 57221
      cost: 1.3761
---
# T-2919 scripts/agent-image.sh --check runs the smoke tier inside the built layer and fails on an image that lacks a tool it runs

## Work

Add `--check` to `scripts/agent-image.sh`: the test that reproduces I-0281. It runs, in the manifest's layer image already built, `flai-agent-sf-claude-code:local`:

1. a `command -v` of each tool the test tiers' scripts call: `go`, `gofmt`, `golangci-lint`, `shellcheck`, `node`, `pnpm`, `git`, `curl`, `openssl`, `tar`, and `goreleaser`. It names each one missing.
2. `flai test smoke` on the tree, mounted as the close-out's container mounts a story's worktree (ADR-0274, ADR-0246), with the entrypoint overridden and no login, key, or MCP variable. It exits as the tier does.

It fails, naming the image and the step, when the manifest names no layer or the image is not built.

It waits for T-2917, which writes the same script and adds the tools this checks.

## Done when

- `scripts/agent-image.sh --check` passes on the layer T-2917 builds.
- On a layer built from main's Dockerfile, without curl, it fails, and its output says "curl is required". The run's last lines are written in this task's `## Notes`.
- ShellCheck passes on `scripts/agent-image.sh`.

## Notes
