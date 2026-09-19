---
id: T-0183
type: task
nature: remediation
title: Conventions, server instructions, and documentation for agents that end their turn
status: done
parent: S-0058
owner: alex
created: 2026-09-19T02:46:21Z
updated: 2026-09-19T02:56:16Z
transitions:
  - to: ready
    at: 2026-09-19T02:55:02Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T02:55:02Z
    by: system-flow
  - to: done
    at: 2026-09-19T02:56:16Z
    by: system-flow
stream: S-0058
tags: []
touches: [design/conventions, template, docs]
---

# T-0183 Conventions, server instructions, and documentation for agents that end their turn

## Work
Template baseline first, then this repository: in `work-management.md`, `session-start.md`, and the `CLAUDE.md` template's priming steps, say that `inbox` lists ready work and the designer's changes, that a ready story found there is pulled without waiting to be told when nothing is in progress, that an agent which ends its turn between messages calls `inbox` at the start of every turn, and that one which stays running holds `wait_for_events`. Do not ask the operator to place stories in the pull order: they cannot from the board until S-0057; pull what they name or what they moved to ready. Update `docs/users/flai.md` (the MCP tools table) and the project additions here.

## Done when
- Template baseline and this repository's copies are identical above the marker
- `docs/users/flai.md` lists the tools as built
- `make smoke` passes

## Notes
