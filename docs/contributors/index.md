---
title: Contributors guide
updated: 2026-10-07
status: draft
---

# Contributing to system-flow

This monorepo builds itself with its own conventions. Start with the root `CLAUDE.md` and `wip/agents/index.md`, then the [design index](../../design/system/README.md).

| Part | Folder | Stack |
|------|--------|-------|
| Standard | `design/` | Markdown |
| Template | `template/` | Go text/template files, see the [template guide](template.md) |
| CLI | `flai/` | Go 1.26, see [go.md](../../design/tech/go.md) |
| Dashboard | `flaiover/` | SvelteKit 2, Svelte 5, Tailwind 4, see [sveltekit.md](../../design/tech/sveltekit.md) |

## Working a story

1. Pull a `ready` story from `wip/kanban/board.md` to `in-progress`.
2. Open its narrative under `wip/agents/`.
3. Keep the narrative's current state and next steps true as you go, with `flai stream state`.
4. Meet the definition of done in [workflow.md](../../design/system/workflow.md), move the story to `review`, open a pull request that names the story.

## Tools and environment

Every script in `scripts/` sources `scripts/env.sh`, which sets this environment:

| What | Value |
|------|-------|
| `PATH` | The tree's own `bin/`, then the main checkout's `bin/` (from a story worktree under `.flai-cache/worktrees/`), then the main checkout's pnpm, then `~/go/bin` |
| `FLAI_CONFIG` | `.flai-cache/config.json` in the main checkout, unless already set |
| `FLAI_CACHE_DIR`, npm and pnpm caches | Under the main checkout's `.flai-cache/` |

Run `make install-tools` once, from any tree: it installs the pinned golangci-lint v2 and GoReleaser into the main checkout's `bin/`, so every story worktree lints with them. `make flai-test` names the golangci-lint it runs and stops if it is not v2.

An agent `flai serve` starts gets flai serve's environment without `FLAI_HOST_URL`, `FLAI_HOST_TOKEN`, and `FLAI_CONFIG`, plus `FLAI_AGENT`, `FLAI_STORY`, `FLAI_SESSION`, and `FLAI_STARTED_BY`. A `flai serve` it runs to try a change, with a `--config` of its own, does not reach the operator's host. `scripts/flaiover-dev.sh` gives the dev server the dashboard token file that `flai dashboard token --json` names.

## Changing a decision

Write a new ADR in `design/adrs` that supersedes the old one, then update `design/system` and `design/tech` in the same pull request.

## Releasing

Releases are per component and computed from what was accepted, never tagged by hand; the rules are in [git.md](../../design/conventions/git.md). Publishing, `flai release --pending` after a fetch or the board's Publish action, bumps each component that has accepted work since its last release, and is the one way accepted work reaches the remote ([ADR-0067](../../design/adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md)):

- `flai`: tag `flai/vX.Y.Z`, GoReleaser publishes binaries.
- `flaiover`: tag `flaiover/vX.Y.Z`, the image workflow publishes to GHCR.
- Template: the version in `template/template.yaml` and an entry in `template/CHANGELOG.md`, then published to its own repository.

## The template

The template is developed in `./template` and published to [bytepunx/system-flow-template](https://github.com/bytepunx/system-flow-template), which is flai's default source. This repository keeps `template.repo: ./template` in `system-flow.yaml` so it develops against the working copy. How to edit, test, version, and publish it, here or in a fork: the [template guide](template.md).
