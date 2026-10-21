---
id: T-1732
type: task
nature: feature
title: ADR-0143's own-name rule holds a defined role, its own name the hook's FLAI_AGENT or the name flai serve gives its run
status: backlog
parent: S-0416
owner: alex
created: 2026-10-16T05:44:19Z
updated: 2026-10-21T11:42:48Z
transitions: []
stream: S-0416
tags: [cli]
touches: [flai/internal/guard/authors.go, flai/internal/guard/authors_test.go, flai/internal/guard/shared.go, design/adrs, design/adrs/README.md]
after: [T-1729]
---
# T-1732 ADR-0143's own-name rule holds a defined role, its own name the hook's FLAI_AGENT or the name flai serve gives its run

## Work

`heldToSelf` and `selves` in `authors.go` know the planner (`planner-<item>`) and the analyzer (`analyzer`). Extend them to any defined role: its own name is the hook's `FLAI_AGENT`, else `<name>` or `<name>-<item>` from `FLAI_ITEM`, as S-0417 names its runs. `who` in `shared.go` names the role in refusals. Add tests for a fixture role: refuse `--by` naming another, `FLAI_AGENT=` another, and `unset FLAI_SESSION`; let its own name through.

Record the extension in a short ADR that refines ADR-0143, made with `flai adr new`: a defined role is a served agent held to its own name, and its own name is the hook's `FLAI_AGENT` or the name `flai serve` gives its run. Add it to `design/adrs/README.md`.

Waits for T-1729, which makes a defined role known to the guard. It runs beside T-1730, since they share no path.

## Done when

`go test ./internal/guard/...` passes. The existing author tests are unchanged. The ADR is accepted in the index, and `flai check --strict` passes.

## Notes

- Planner (S-0416, 2026-10-21): the orchestrator approved the ADR on TH-0569. ADR-0143's decision names the agents it holds, and a defined role is not among them.
- `design/adrs` is a folder touch, because the ADR's file name is known only once `flai adr new` makes it. The folder is in `claims.shared`, so it holds no story.
