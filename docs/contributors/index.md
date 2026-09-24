---
title: Contributors guide
updated: 2026-09-24
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
3. Keep the narrative's current state and next steps true as you go.
4. Meet the definition of done in [workflow.md](../../design/system/workflow.md), move the story to `review`, open a pull request that names the story.

## Changing a decision

Write a new ADR in `design/adrs` that supersedes the old one, then update `design/system` and `design/tech` in the same pull request.

## Releasing

Releases are per component and computed from what was accepted, never tagged by hand; the rules are in [git.md](../../design/conventions/git.md). `flai release --pending`, the board's Publish action, or `flai push --pending` bumps each component that has accepted work since its last release:

- `flai`: tag `flai/vX.Y.Z`, GoReleaser publishes binaries.
- `flaiover`: tag `flaiover/vX.Y.Z`, the image workflow publishes to GHCR.
- Template: the version in `template/template.yaml` and an entry in `template/CHANGELOG.md`, then published to its own repository.

## The template

The template is developed in `./template` and published to [bytepunx/system-flow-template](https://github.com/bytepunx/system-flow-template), which is flai's default source. This repository keeps `template.repo: ./template` in `system-flow.yaml` so it develops against the working copy. How to edit, test, version, and publish it, here or in a fork: the [template guide](template.md).
