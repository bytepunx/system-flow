---
title: Docker
updated: 2026-09-18
status: active
---

# Docker

| | |
|-|-|
| Host requirement | Docker Engine 24 or newer on `PATH` (dev host has 29) |
| Image | `ghcr.io/bytepunx/flaiover`, three-stage build from the repo root (`flaiover/Dockerfile`): `golang:1.26-alpine` builds flai, `node:24-alpine` builds flaiover and prunes to production deps, `node:24-alpine` runtime with `tini`, `git`, the OpenSSH client (about 0.7 MB, for pushing an acceptance with a key the operator gives the container, [ADR 0026](../adrs/0026-the-dashboard-may-push-with-a-key-the-operator-gives-it.md), S-0062), the flai binary, `USER 65532` by default and any `--user` supported; about 670 MB (S-0015) |
| Tags | `latest` on main, `<semver>` and `<major>` from `flaiover/v*` tags, `sha-<short>`; built by `release-flaiover.yml` with buildx and GHA cache |
| Run | `flai dashboard` runs `docker run --detach --rm --name flaiover-<project> --publish 127.0.0.1:<port>:3000 --volume <repo>:<repo> --env PROJECT_DIR=<repo> --user <uid>:<gid> ghcr.io/bytepunx/flaiover:<tag>` |

## Why

The brief requires a published image the CLI can download and run. Docker avoids asking every user to install Node. The mount is read-write so the board is operable from the browser; the container runs as a non-root user with the host's UID passed via `--user` so written files keep sane ownership on Linux. The repository is mounted at its host path, not a fixed one, because git links story worktrees with absolute paths and the container runs git for acceptance ([ADR 0022](../adrs/0022-repository-mounted-at-its-host-path.md)).

Decision: [ADR 0007](../adrs/0007-sveltekit-spa-with-node-adapter-in-docker.md).
