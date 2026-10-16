---
id: T-1732
type: task
nature: feature
title: ADR-0143's own-name rule holds a defined role, its own name the hook's FLAI_AGENT or the name flai serve gives its run
status: backlog
parent: S-0416
owner: alex
created: 2026-10-16T05:44:19Z
updated: 2026-10-16T05:44:19Z
transitions: []
stream: S-0416
tags: [cli]
touches: [flai/internal/guard/authors.go, flai/internal/guard/authors_test.go, flai/internal/guard/shared.go]
after: [T-1729]
---
# T-1732 ADR-0143's own-name rule holds a defined role, its own name the hook's FLAI_AGENT or the name flai serve gives its run

## Work

`heldToSelf` and `selves` in `authors.go` know the planner (`planner-<item>`) and the analyzer (`analyzer`). Extend them to any defined role: its own name is the hook's `FLAI_AGENT`, else `<name>` or `<name>-<item>` from `FLAI_ITEM`, as S-0417 names its runs. `who` in `shared.go` names the role in refusals. Add tests for a fixture role: refuse `--by` naming another, `FLAI_AGENT=` another, and `unset FLAI_SESSION`; let its own name through.

Waits for T-1729, which makes a defined role known to the guard. It runs beside T-1730, since they share no path.

## Done when

`go test ./internal/guard/...` passes. The existing author tests are unchanged.

## Notes
