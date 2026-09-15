---
title: Contributors guide
updated: 2026-09-15
status: draft
---

# Contributing to system-flow

This monorepo builds itself with its own conventions. Start with the root `CLAUDE.md` and `wip/agents/index.md`, then the [design index](../../design/system/README.md).

| Part | Folder | Stack |
|------|--------|-------|
| Standard | `design/` | Markdown |
| Template prototype | `template/` | Go text/template files, see [template.md](../../design/system/template.md) |
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

- `flai`: tag `flai/vX.Y.Z`, GoReleaser publishes binaries.
- `flaiover`: tag `flaiover/vX.Y.Z`, the image workflow publishes to GHCR.
- Template: bump `template.yaml` version, add a `CHANGELOG.md` entry, tag the template repository.
