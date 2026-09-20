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
| Run | `flai dashboard` runs `docker run --detach --rm --name flaiover-<project> --publish <bind>:<port>:3000`, two read-only bind mounts under `/run/secrets` (the login token and the agent credential) with the variables that name them, `--user <uid>:<gid>`, and `ghcr.io/bytepunx/flaiover:<tag>`. No volume: nothing of the project is in the container ([ADR-0031](../adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)) |

## Why

The brief requires a published image the CLI can download and run. Docker avoids asking every user to install Node. The container holds the web server and nothing else: the board is operable from the browser because the server asks `flai serve` on the host for every read and write (ADR-0029), so no repository is mounted ([ADR-0031](../adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md); until S-0077 it was, read-write, at its host path, ADR-0022). It runs as a non-root user with the host's UID passed via `--user`, because the two secret files it is handed are readable by that user alone.

Decision: [ADR 0007](../adrs/0007-sveltekit-spa-with-node-adapter-in-docker.md).
