---
id: ADR-0004
title: One state machine, transitions recorded in front matter, blocked is a flag
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0004 One state machine, transitions recorded in front matter, blocked is a flag

## Context

Cycle time, lead time, time in state, and cumulative flow all need the full history of state changes, not just the current state. Blocked work needs to be measured without hiding which stage it is stuck in.

## Decision

All item types share the states `backlog`, `ready`, `in-progress`, `review`, `done`, `cancelled`. Every state change is appended to a `transitions` list in front matter with `to`, `at`, and `by`. `status` is a denormalised copy of the last transition and is validated against it. Blocked is a list of timestamped intervals with reasons, not a state. Derived timestamps (`started`, `completed`) are computed, never stored.

## Consequences

- Front matter grows with every transition, which is acceptable for items that see five to ten transitions in their life.
- History survives without a database and is reviewable in git.
- Tools must compute, not read, `started` and `completed`, so both `flai` and `flaiover` implement `design/system/metrics.md` and are fixture-tested against each other.

## Alternatives considered

- Deriving history from git log: fragile, depends on commit granularity, and unavailable in uncommitted work.
- A `blocked` state: hides the column the item was in and distorts time-in-state.
- Storing `started` and `completed` directly: two sources of truth that drift.
