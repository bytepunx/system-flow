---
id: S-021
type: story
nature: feature
title: flai template push publishes template changes to a remote
status: backlog
parent: E-001
owner: agent
created: 2026-09-15T17:39:42Z
updated: 2026-09-15T17:39:42Z
transitions: []
tags: [cli, template, devops]
---

# S-021 flai template push publishes template changes to a remote

## Goal
A repository that develops a template locally, as this monorepo does under `./template`, can publish the template to its git remote in one command so other projects can consume it with `flai new` and `flai upgrade`. The command assumes the user already has push permission; any git failure is surfaced as the git error and a non-zero exit.

## Acceptance criteria
- [ ] `flai template push [dir]` takes the template directory from the argument, else from config `template.repo` when it is a local path, and errors clearly otherwise
- [ ] The remote and branch come from `--remote` and `--ref`, defaulting to a new `publish` section in `template.yaml` (`publish.repo`, `publish.ref`); missing both is an error
- [ ] Clones the remote at the branch into the cache (creating the branch from the default branch when it does not exist), replaces its contents with the local template excluding `.git`, commits with a message that includes the template `version`, and pushes
- [ ] No changes produces "nothing to push" and exit 0
- [ ] Git failures (authentication, non-fast-forward, missing remote) are printed verbatim from git and exit non-zero; no retries and no force push unless `--force`
- [ ] `--tag` also creates and pushes an annotated tag `v<version>`, refusing if the tag already exists on the remote
- [ ] `--dry-run` prints the file-level change summary and the commit message without cloning into a temporary working copy being left behind or pushing
- [ ] `--json` reports remote, ref, commit, files changed, and tag
- [ ] The project's own `system-flow.yaml` is not modified; switching a project from the local path to the remote stays an explicit `flai template use`

## Tasks

## Notes
- S-003 published the template once by hand with `git subtree split` and a tag (see docs/contributors); this story replaces that procedure. Keep the subtree approach or a clone-and-replace, but preserve history where practical.
- Implementation goes through `execx.Runner` and real `git` (ADR-0010). Tests use a local bare repository as the remote, including a non-fast-forward case to prove the error path.
- `template.yaml` gains `publish: {repo, ref}`; `design/system/template.md` and the prototype manifest are updated in this story.
