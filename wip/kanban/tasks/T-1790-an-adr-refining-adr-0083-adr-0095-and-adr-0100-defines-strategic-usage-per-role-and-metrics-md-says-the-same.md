---
id: T-1790
type: task
nature: feature
title: An ADR refining ADR-0083, ADR-0095, and ADR-0100 defines strategic usage per role, and metrics.md says the same
status: backlog
parent: S-0425
owner: alex
created: 2026-10-16T05:49:19Z
updated: 2026-10-23T15:34:23Z
transitions: []
stream: S-0425
tags: [cli]
touches: [design/adrs, design/adrs/README.md, design/system/metrics.md, design/system/work-hierarchy.md]
---
# T-1790 An ADR refining ADR-0083, ADR-0095, and ADR-0100 defines strategic usage per role, and metrics.md says the same

## Work

Write the ADR with `flai adr new`, refining ADR-0083, ADR-0095, and ADR-0100:

- a strategic kind is a role a project defines;
- each role's file names its charge rule, `item`, `named`, or `issues`, the three built-in roles keeping the rules they have;
- `flai stats` reports per role.

Rewrite `design/system/metrics.md` § Strategic usage and § Strategic agents to match. `CLAUDE.md` asks that `metrics.md` change only with an ADR.

First layer.

## Done when

The ADR is accepted and listed in `design/adrs/README.md`. `flai check --strict` passes.

## Notes

Touches the shared `design/adrs` folder for the ADR's new file.
