---
id: T-0235
type: task
nature: feature
title: The decisions convention, the ADR index text, and the designs say ADRs are made with flai adr new
status: done
parent: S-0060
owner: alex
created: 2026-09-19T09:39:14Z
updated: 2026-09-19T09:47:14Z
transitions:
  - to: ready
    at: 2026-09-19T09:45:52Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T09:45:52Z
    by: system-flow
  - to: done
    at: 2026-09-19T09:47:14Z
    by: system-flow
stream: S-0060
tags: []
---

# T-0235 The decisions convention, the ADR index text, and the designs say ADRs are made with flai adr new

## Work
`design/conventions/decisions.md`, template baseline first and then here above the marker: ADRs are made with `flai adr new`, not by copying the template. The project addition "next is 0015" is removed. `design/adrs/README.md` here and in the template says the same instead of "Copy 0000-template.md". `documentation-standard.md` or wherever the ADR format is described, and `flai-cli.md`.

## Done when
- The convention reads the same in the template and here, and no document tells anyone to copy the template or to look up the next number

## Notes
