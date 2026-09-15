---
id: ADR-0011
title: system-flow.yaml marks a project and owns folder names
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0011 system-flow.yaml marks a project and owns folder names

## Context

Tools need to detect a conforming repo, know which template produced it, enumerate its sub-projects, and find its documentation folders even when they were renamed at import.

## Decision

A conforming repo has `system-flow.yaml` at its root with a schema version, project name and key, template source and applied version, a `layout` map of folder names, a `projects` list, and dashboard settings. Tools resolve folders through `layout` only. `flai` owns and rewrites `template.*` and `projects`; humans may edit the rest.

## Consequences

- Renaming `design`, `docs`, or `wip` is supported without touching tooling.
- `flai` commands other than `new` and `import` refuse to run outside a manifest's tree, which prevents accidental writes to unrelated repos.
- `template.version` enables a future `flai upgrade`.

## Alternatives considered

- A hidden `.flai/` directory: hides the marker from humans reading the repo.
- Detecting conformance by folder presence: ambiguous once renaming is allowed.
