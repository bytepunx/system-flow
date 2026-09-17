---
title: Docker
updated: 2026-09-15
status: active
---

# Docker

| | |
|-|-|
| Host requirement | Docker Engine 24 or newer on `PATH` (dev host has 29) |
| Image | `ghcr.io/bytepunx/flaiover`, multi-stage build, runtime `node:24-alpine`, non-root user |
| Tags | `latest`, `<semver>`, `<major>`, `sha-<short>` |
| Run | `flai dashboard` runs `docker run --detach --rm --name flaiover-<project> --publish 127.0.0.1:<port>:3000 --volume <repo>:/project --env PROJECT_DIR=/project --user <uid>:<gid> ghcr.io/bytepunx/flaiover:<tag>` |

## Why

The brief requires a published image the CLI can download and run. Docker avoids asking every user to install Node. The mount is read-write so the board is operable from the browser; the container runs as a non-root user with the host's UID passed via `--user` so written files keep sane ownership on Linux.

Decision: [ADR 0007](../adrs/0007-sveltekit-spa-with-node-adapter-in-docker.md).
