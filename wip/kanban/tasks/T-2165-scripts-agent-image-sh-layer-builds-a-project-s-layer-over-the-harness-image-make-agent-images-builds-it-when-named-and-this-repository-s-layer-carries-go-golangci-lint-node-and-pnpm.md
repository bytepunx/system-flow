---
id: T-2165
type: task
nature: feature
title: scripts/agent-image.sh --layer builds a project's layer over the harness image, make agent-images builds it when named, and this repository's layer carries Go, golangci-lint, Node, and pnpm
status: backlog
parent: S-0510
owner: alex
created: 2026-10-26T04:49:45Z
updated: 2026-10-26T04:49:45Z
transitions: []
stream: S-0510
tags: [cli, docker]
touches: [scripts/agent-image.sh, Makefile, flai/images/system-flow/Dockerfile, system-flow.yaml, scripts/README.md]
---
# T-2165 scripts/agent-image.sh --layer builds a project's layer over the harness image, make agent-images builds it when named, and this repository's layer carries Go, golangci-lint, Node, and pnpm

## Work

Extend `scripts/agent-image.sh` with `--layer <Dockerfile> --inner <harness>`: it builds the Dockerfile with `--build-arg BASE=flai-agent-<inner>:local` and tags the result `flai-agent-<project key>-<inner>:local`, the project key read from `system-flow.yaml`. Make `make agent-images` build the two base images and then, when `agent.image.layer` is set in the manifest, the layer for `claude-code`. Write `flai/images/system-flow/Dockerfile` with `ARG BASE`, adding Go at the version `flai/go.mod` names, golangci-lint v2 at the version `scripts/install-tools.sh` installs, Node 24, and pnpm, so that every tier the manifest's `tests` select runs; name it in `system-flow.yaml` as `agent.image.layer`, and describe the flag in `scripts/README.md`. This task waits for nothing: it changes the script and the manifest file, not the manifest reader.

## Done when

`make agent-images` builds `flai-agent-system-flow-claude-code:local` on this host, and `docker run --rm flai-agent-system-flow-claude-code:local sh -c 'go version && golangci-lint version && node --version && pnpm --version'` answers for each.

## Notes

The image is large; record its size in the story's notes so the operator knows the cost of the layer.
