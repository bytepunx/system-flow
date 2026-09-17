---
id: S-020
type: story
nature: feature
title: flai upgrade re-integrates the latest template
status: done
parent: E-001
owner: agent
created: 2026-09-15T17:39:42Z
updated: 2026-09-17T03:30:50Z
transitions:
  - to: ready
    at: 2026-09-17T03:19:34Z
    by: agent
  - to: in-progress
    at: 2026-09-17T03:19:34Z
    by: agent
  - to: review
    at: 2026-09-17T03:25:23Z
    by: agent
  - to: done
    at: 2026-09-17T03:30:50Z
    by: alex
tags: [cli, template]
---

# S-020 flai upgrade re-integrates the latest template

## Goal
An existing conforming repository can pick up a newer version of its template without losing project content. `flai upgrade` fetches the template at the configured source and ref, works out which rendered files the project has left untouched, updates those, adds files the template gained, and reports the rest as conflicts for the human to resolve.

## Acceptance criteria
- [x] `flai upgrade` in a conforming repo fetches the template (config or `--template`/`--ref`), compares its `version` with `system-flow.yaml` `template.version`, and is a no-op with a clear message when they match unless `--force`
- [x] Files the template adds are written; files whose rendered content matches what was originally applied are replaced with the new rendering; files the project has modified are reported as conflicts and left alone
- [x] Conflicts offer keep, replace, or show diff in a terminal; non-interactive runs keep the project's version and exit non-zero with the list unless `--keep-all` or `--replace-all`
- [x] `CLAUDE.md` is merged by marker: only the section above `<!-- system-flow:end-of-baseline -->` is replaced and everything below is preserved byte for byte
- [x] Layout renames recorded in `system-flow.yaml` are honoured, so a project with `design` renamed to `architecture` upgrades into the right folders
- [x] `system-flow.yaml` `template.version` and `template.applied` are updated only when the upgrade completes without unresolved conflicts
- [x] `--dry-run` prints every planned add, update, conflict, and skip without writing
- [x] Refuses to run on a dirty git working tree unless `--force`, so the upgrade is reviewable as one diff
- [x] `flai check` passes on the upgraded repo

## Tasks
- T-072 Lock file (ADR-0015) written by flai new; render can collect instead of write
- T-073 upgrade package: classify add, merge by marker, replace unchanged, conflict; apply with policies
- T-074 flai upgrade command: dry-run, dirty-tree guard, interactive keep/replace/diff, keep-all, replace-all
- T-075 Tests, docs, and an upgrade of this repository

## Notes
- Lock file mechanism chosen and recorded as ADR-0015; the alternative (re-rendering the old version) is documented there.
- Detecting "unchanged since applied" needs a record of what was rendered. Proposed: `flai new` and `flai upgrade` write a lock file `system-flow.lock.yaml` next to the manifest with the template version and a sha256 per rendered path. Upgrade compares the project's file against the lock hash: equal means safe to replace, different means conflict, missing means the project deleted it and the file is skipped. This needs an ADR because it adds a file to the repository root and S-005's `flai new` must be extended to write it.
- Alternative considered: re-render the previously applied template version from the cache for a three-way compare. Works only for git sources with a reachable tag per version, not for local paths, so the lock file is preferred.
- Convention files under design/conventions merge like CLAUDE.md: replace above `<!-- system-flow:end-of-baseline -->`, preserve below. Files the template adds are copied whole; a project file with no marker is a conflict.
- Depends on S-005 (rendering) and benefits from S-008 (check). S-003 (publish the template repo) is the first real consumer.
