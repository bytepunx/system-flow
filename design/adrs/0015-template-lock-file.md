---
id: ADR-0015
title: system-flow.lock.yaml records what the template rendered so upgrade can tell edits from baseline
status: accepted
date: 2026-09-17
supersedes: []
superseded_by: []
---

# ADR-0015 system-flow.lock.yaml records what the template rendered so upgrade can tell edits from baseline

## Context

`flai upgrade` must bring a project to a newer template without destroying project work. For each rendered path it needs to know whether the project has changed the file since the template wrote it. A local template directory has no history to diff against, a git template's old version may be unreachable, and re-rendering the old version needs the old variables. Only a record made at render time answers the question reliably.

## Decision

`flai new` and `flai upgrade` write `system-flow.lock.yaml` next to the manifest. It records the template repo, ref, and version that were applied and, for every path the template rendered, the sha256 of the content as written. On upgrade, a project file whose current hash equals the lock's hash is unchanged since it was applied and is replaced; a file that differs is a conflict the operator resolves; a file the template newly adds is added; files carrying the baseline marker (`CLAUDE.md`, the conventions) are merged, template above the marker and project below, regardless of hash. The lock is committed. A project without a lock treats every differing file as a conflict; `flai upgrade --relock` writes the lock at the current version without changing files.

## Consequences

- One more generated file at the root, small and diff-friendly, owned by `flai`.
- Upgrades are safe by default: nothing the project edited is overwritten without an explicit choice.
- Renamed layout folders are handled because the lock stores project-relative paths after the rename.
- `flai check` warns when the lock is missing on a project whose manifest names a template version, so hand-assembled projects get relocked.

## Alternatives considered

- Three-way merge against the previously applied template version: needs the old template and the old variables; fails for local sources and moved refs.
- Treating every difference as a conflict: safe but noisy; every upgrade would ask about every file the template ever changed.
- Storing hashes inside `system-flow.yaml`: mixes a human-edited manifest with a generated table that changes every upgrade.
