---
id: T-2922
type: task
nature: remediation
title: The docs that list the system-flow layer's tools name curl, openssl, GoReleaser, and scripts/agent-image.sh --check
status: done
parent: S-0672
owner: alex
created: 2026-11-10T08:48:23Z
updated: 2026-11-10T10:10:00Z
transitions:
  - to: ready
    at: 2026-11-10T10:09:05Z
    by: agent-S-0672
  - to: in-progress
    at: 2026-11-10T10:09:06Z
    by: agent-S-0672
  - to: done
    at: 2026-11-10T10:10:00Z
    by: agent-S-0672
stream: S-0672
tags: [docs, docker]
touches: [design/tech/docker.md, docs/operators/index.md, docs/users/flai.md, scripts/README.md]
after: [T-2919]
usage:
  source: log
  seconds: 54
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 2913
      cache_read: 1933265
      cache_write: 27725
      cost: 0.6668
---
# T-2922 The docs that list the system-flow layer's tools name curl, openssl, GoReleaser, and scripts/agent-image.sh --check

## Work

Bring each document that lists what this repository's layer carries up to date with T-2917 and T-2919:

- `design/tech/docker.md`, the agent images row: `curl`, `openssl`, and GoReleaser at its pin beside the other tools; `scripts/agent-image.sh` passing `GORELEASER_VERSION`; the image's size as `docker image ls` gives it; and `--check`. Bump `updated`.
- `docs/operators/index.md`, the container harness entry ("This repository's carries Go, golangci-lint, ShellCheck, and pnpm"), and `docs/users/flai.md`, the `agent.image.layer` paragraph: the same list. Name `--check` as how to see that the layer runs the smoke tier. Bump each `updated`.
- `scripts/README.md`, the `agent-image.sh` row: `--check`, and GoReleaser among the pins it passes.

It waits for T-2919, since it documents the flag that task adds and the tools T-2917 adds.

## Done when

- The four files name curl, openssl, GoReleaser, and `--check` as the scripts now behave, and no line in them says the layer carries only Go, golangci-lint, ShellCheck, and pnpm.
- The markdown lint passes on the four files.

## Notes
