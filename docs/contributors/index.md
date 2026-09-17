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

## Publishing the template

The template is developed in `./template` and published to [bytepunx/system-flow-template](https://github.com/bytepunx/system-flow-template), which is flai's default source. Until `flai template push` (S-021) automates it, sync by hand after a story that changes the template lands:

```bash
git subtree split --prefix=template -b template-main
git push git@github.com:bytepunx/system-flow-template.git template-main:main
git tag -a template/vX.Y.Z template-main -m "system-flow-template vX.Y.Z"
git push git@github.com:bytepunx/system-flow-template.git refs/tags/template/vX.Y.Z:refs/tags/vX.Y.Z
git branch -D template-main
```

`template/template.yaml` carries the version and `template/CHANGELOG.md` the entry; both are bumped at acceptance per the release rule. The subtree split preserves the template's history without rewriting this repository. This repository keeps `template.repo: ./template` in `system-flow.yaml` so it develops against the working copy; other projects use the published repo and a tag.
