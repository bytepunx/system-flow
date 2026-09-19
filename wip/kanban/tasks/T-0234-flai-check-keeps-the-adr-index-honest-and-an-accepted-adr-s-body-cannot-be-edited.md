---
id: T-0234
type: task
nature: feature
title: flai check keeps the ADR index honest, and an accepted ADR's body cannot be edited
status: done
parent: S-0060
owner: alex
created: 2026-09-19T09:39:14Z
updated: 2026-09-19T09:45:52Z
transitions:
  - to: ready
    at: 2026-09-19T09:45:51Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T09:45:52Z
    by: system-flow
  - to: done
    at: 2026-09-19T09:45:52Z
    by: system-flow
stream: S-0060
tags: []
---

# T-0234 flai check keeps the ADR index honest, and an accepted ADR's body cannot be edited

## Work
`flai check`: a warning `adr.index` when an ADR file has no row in `design/adrs/README.md` and when a row links to a file that is not there; the template file is exempt. `flai doc show` and `flai doc save` refuse the body of an ADR whose status is `accepted` (mode `none` with the reason: an accepted ADR is immutable; write a new one that supersedes it), which settles the question S-0040 left open; a proposed ADR stays editable. Tests for both.

## Done when
- The rule and the refusal are tested, and this repository passes `flai check --strict` with them

## Notes
