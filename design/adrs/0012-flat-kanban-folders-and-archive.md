---
id: ADR-0012
title: Kanban items are flat per type, hierarchy by parent key, archive on completion
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0012 Kanban items are flat per type, hierarchy by parent key, archive on completion

## Context

Items could be nested on disk to mirror the hierarchy, or kept flat with the hierarchy in front matter. The board must stay small while history remains available for metrics.

## Decision

`wip/kanban` has three flat folders, `epics`, `stories`, `tasks`, with files named `<ID>-<slug>.md`. Hierarchy is expressed by the `parent` key. `flai archive` moves done and cancelled items, and finished narratives, into `wip/archive` with the same layout. Metrics scan both.

## Consequences

- Re-parenting is a one-line front matter edit, not a file move.
- Tools list a folder to get all items of a type without walking a tree.
- Human browsing of the hierarchy on disk is weaker; the dashboard and `flai board` provide the tree view.

## Alternatives considered

- Nested `E-0001/S-0001/T-0001.md`: pretty on disk, painful to re-parent and to glob.
- Never archiving: the board folder grows unbounded and every listing gets slower.
