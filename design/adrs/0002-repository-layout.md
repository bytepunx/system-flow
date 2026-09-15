---
id: ADR-0002
title: "Three top-level documentation folders: design, docs, wip"
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0002 Three top-level documentation folders: design, docs, wip

## Context

A monorepo mixes internal design, outward-facing documentation, and transient work state. Readers and agents need to know at a glance which is which, and tooling needs stable places to look.

## Decision

A conforming monorepo has `design/` for internal documentation with subfolders `adrs`, `system`, and `tech`; `docs/` for outward-facing documentation with one subfolder per audience; and `wip/` for work in process with subfolders `agents`, `kanban`, and `archive`. Code sub-projects sit at the repo root, one folder each. The three folder names are defaults and may be renamed at import; the chosen names are recorded in the project manifest (ADR-0011).

## Consequences

- Sub-project design lives in `design/system/<project>.md`, not inside the sub-project. This keeps one place to look but means a sub-project is not self-contained if extracted.
- Renaming support means every tool resolves folder names through the manifest.
- The full layout is specified in `design/system/repository-layout.md`.

## Alternatives considered

- A single `docs/` with internal and external mixed: readers cannot tell what is safe to publish.
- Per-sub-project `design/` folders: fragments the living design and makes cross-project ADRs homeless.
- `projects/` wrapper for sub-projects: one more level for no benefit in the common case.
