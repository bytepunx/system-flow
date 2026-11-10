---
title: Node.js and pnpm
updated: 2026-11-10
status: active
topics: [sveltekit]
---

# Node.js and pnpm

| | |
|-|-|
| Node version | 24 LTS, pinned in `flaiover/.nvmrc` and `engines`, and the Dockerfile base image `node:24.22.2-alpine3.24`, by digest ([docker.md](docker.md)); the dev host runs 25 and has no corepack, so pnpm is installed by `scripts/install-tools.sh` into `.flai-cache/pnpm` and all Node caches live under `.flai-cache`. The same script installs markdownlint-cli2 into `.flai-cache/markdownlint` and links it from `bin/`, so `scripts/lint-md.sh` starts it without `npx` ([ci.md](ci.md)) |
| pnpm version | 10, pinned via `packageManager` in `flaiover/package.json` and enabled by corepack |
| pnpm store | `.flai-cache/pnpm-store` of the main checkout, shared by every worktree, set by `scripts/env.sh` as `npm_config_store_dir`. Add a dependency with `scripts/flaiover-pnpm.sh add <package>`, not a bare `pnpm add`, which uses the default store and may be another pnpm on the host's `PATH`; `scripts/flaiover-install.sh` purges an install from another store and makes it again from this one (I-0177) |
| Used in | `flaiover` |

## Why

- Node 24 is the active LTS line through 2027. Host machines may run newer (the dev host has 25) but builds and the image use the LTS to keep the container boring.
- pnpm for a strict `node_modules`, fast installs, and one lockfile the Docker build can cache. The Docker build and CI install from `flaiover/pnpm-lock.yaml` alone, with `--frozen-lockfile` (`flaiover/Dockerfile`, `.github/workflows/flaiover.yml`). It is the one lockfile: S-0630 removed the `package-lock.json` files at the root and in `flaiover/`, which nothing read (I-0297), and `scripts/ci-pins.sh` fails a second one. The caret ranges in `flaiover/package.json` stay: the frozen lockfile is the pin.
