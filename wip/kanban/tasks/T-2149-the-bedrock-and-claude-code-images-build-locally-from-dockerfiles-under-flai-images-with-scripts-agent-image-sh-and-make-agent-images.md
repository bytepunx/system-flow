---
id: T-2149
type: task
nature: feature
title: The bedrock and claude-code images build locally from Dockerfiles under flai/images with scripts/agent-image.sh and make agent-images
status: backlog
parent: S-0505
owner: alex
created: 2026-10-26T04:47:03Z
updated: 2026-10-26T04:47:03Z
transitions: []
stream: S-0505
tags: [cli, docker]
touches: [flai/images/bedrock/Dockerfile, flai/images/claude-code/Dockerfile, scripts/agent-image.sh, Makefile, scripts/README.md, ".dockerignore"]
---
# T-2149 The bedrock and claude-code images build locally from Dockerfiles under flai/images with scripts/agent-image.sh and make agent-images

## Work

Write `flai/images/bedrock/Dockerfile`, built from the repository root as `flaiover/Dockerfile` is: a `golang:1.26-alpine` stage builds `flai` from the tree, and the runtime stage is a small Linux base (Alpine or Debian slim, whichever Claude Code's install needs on the next layer) with `git`, `ca-certificates`, `tini`, and the `flai` binary on `PATH`, running as a non-root user whose uid `--user` overrides. Write `flai/images/claude-code/Dockerfile` with `ARG BASE=flai-agent-bedrock:local`, adding Node 24 and Claude Code installed as Anthropic documents it, with no login state and no key baked in. Write `scripts/agent-image.sh`, POSIX `sh` under `set -eu`, that builds the two in order and tags them `flai-agent-bedrock:local` and `flai-agent-claude-code:local`, and add `agent-images` to the `Makefile` calling it, with a line in `scripts/README.md`. Extend `.dockerignore` if the build context needs it. Add nothing under `.github/workflows`. This task waits for nothing.

## Done when

`make agent-images` builds both images on this host from a clean tree, `docker run --rm flai-agent-bedrock:local flai version` and `docker run --rm flai-agent-claude-code:local claude --version` both answer, and no workflow names either image.

## Notes

The images carry no release stamp and are not signed; E-0029 builds locally and leaves the published image to a later story.
