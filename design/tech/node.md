---
title: Node.js and pnpm
updated: 2026-09-15
status: active
---

# Node.js and pnpm

| | |
|-|-|
| Node version | 24 LTS, pinned in `flaiover/.nvmrc` and `engines`, and the Dockerfile base image `node:24-alpine`; the dev host runs 25 and has no corepack, so pnpm is installed by `scripts/install-tools.sh` into `.flai-cache/pnpm` and all Node caches live under `.flai-cache` |
| pnpm version | 10, pinned via `packageManager` in `flaiover/package.json` and enabled by corepack |
| Used in | `flaiover` |

## Why

- Node 24 is the active LTS line through 2027. Host machines may run newer (the dev host has 25) but builds and the image use the LTS to keep the container boring.
- pnpm for a strict `node_modules`, fast installs, and a single lockfile format that the Docker build can cache.
