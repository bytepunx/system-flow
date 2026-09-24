---
id: T-0378
type: task
nature: feature
title: ADR, design, and documentation for flai host
status: done
parent: S-0106
owner: alex
created: 2026-09-24T01:30:21Z
updated: 2026-09-24T01:48:58Z
transitions:
  - to: ready
    at: 2026-09-24T01:44:48Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T01:44:48Z
    by: system-flow
  - to: done
    at: 2026-09-24T01:48:58Z
    by: system-flow
stream: S-0106
tags: []
touches: [design/adrs, design/system, docs]
---
# T-0378 ADR, design, and documentation for flai host

## Work

- ADR: one flai host per machine supervises flai serve and each project's MCP server; refines ADR-0034 (serve asks rather than starts).
- `design/system/flai-cli.md`, `docs/users/flai.md`, operator docs: the command, the files, the address, the action.
- Run all three test tiers and `flai check --strict`.

## Done when

- Docs say what the code does; `make test`, `make integration`, `make smoke` pass; check is clean.

## Notes
