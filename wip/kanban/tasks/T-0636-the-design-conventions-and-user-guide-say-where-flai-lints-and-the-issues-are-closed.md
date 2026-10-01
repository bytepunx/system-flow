---
id: T-0636
type: task
nature: remediation
title: The design, conventions, and user guide say where flai lints, and the issues are closed
status: ready
parent: S-0179
owner: arobson
created: 2026-10-01T08:23:14Z
updated: 2026-10-01T08:23:32Z
transitions:
  - to: ready
    at: 2026-10-01T08:23:32Z
    by: claude-opus-5-5
stream: S-0179
tags: []
touches: [design/system/flai-cli.md, design/adrs, design/conventions, template/, docs/users, design/issues]
---
# T-0636 The design, conventions, and user guide say where flai lints, and the issues are closed

## Work

Record the decision to embed a subset of markdownlint in an ADR; describe the lint in `design/system/flai-cli.md`, the tooling convention (template first), and `docs/users/flai.md`; close I-0027 and I-0043 with what fixed them; confirm `scripts/lint-md.sh` passes on main.

## Done when

- [ ] The ADR, design, conventions, and user guide describe where flai lints what it writes
- [ ] I-0027 and I-0043 are closed
- [ ] `make test`, `make integration`, `make smoke`, and `flai check --strict` pass

## Notes

Part of S-0179.
