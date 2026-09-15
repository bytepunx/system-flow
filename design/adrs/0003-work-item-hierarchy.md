---
id: ADR-0003
title: Epics, stories, tasks, each with a nature
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0003 Epics, stories, tasks, each with a nature

## Context

Work must be traceable from a specified outcome to the pieces of work that deliver it, sized so agents can complete units in a session, and classified so metrics can show what kind of work the system spends its time on.

## Decision

Three levels: epics span multiple stories, stories are incremental deliverables, tasks are pieces of work required by a story. Every item carries a `nature` from a closed list: `feature`, `improvement`, `remediation`, `research`, `experiment`. IDs are `E-nnn`, `S-nnn`, `T-nnn` with per-type sequences. Front matter and body headings are fixed and specified in `design/system/work-hierarchy.md`.

## Consequences

- Metrics group by type and nature out of the box.
- Adding a nature or a level is an ADR-level change because the dashboard and `flai check` depend on the closed lists.
- No sub-tasks. If a task needs splitting, it becomes a story.

## Alternatives considered

- Two levels (story, task): epics are needed for burn-up scope and for human-level direction.
- Free-text natures: not chartable.
- Project-prefixed IDs (`SF-S-001`): unnecessary while items live in one repo; the manifest reserves a `key` for a future opt-in.
