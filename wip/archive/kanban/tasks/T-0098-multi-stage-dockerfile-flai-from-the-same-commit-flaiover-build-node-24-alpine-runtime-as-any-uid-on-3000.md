---
id: T-0098
type: task
nature: feature
title: "Multi-stage Dockerfile: flai from the same commit, flaiover build, node:24-alpine runtime as any UID on 3000"
status: done
parent: S-0015
owner: alex
created: 2026-09-17T05:29:33Z
updated: 2026-09-17T05:33:21Z
transitions:
  - to: ready
    at: 2026-09-17T05:33:20Z
    by: agent
  - to: in-progress
    at: 2026-09-17T05:33:20Z
    by: agent
  - to: done
    at: 2026-09-17T05:33:21Z
    by: agent
stream: S-0015
tags: [dashboard, docker]
---

# T-0098 Multi-stage Dockerfile: flai from the same commit, flaiover build, node:24-alpine runtime as any UID on 3000

## Work
flaiover/Dockerfile with three stages: golang:1.26-alpine builds flai (CGO off, version from build args), node:24-alpine builds flaiover with pnpm pinned by packageManager and prunes to production dependencies, node:24-alpine runtime copies build/ and node_modules and the flai binary to /usr/local/bin, sets PORT=3000, HOST=0.0.0.0, PROJECT_DIR=/project, FLAI_BIN, HOME=/tmp, runs as an unprivileged numeric user by default but works under any --user; .dockerignore keeps the context small.

## Done when
docker build from the repo root succeeds; the container serves /api/manifest for a mounted repo as UID 1000.

## Notes
