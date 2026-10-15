---
id: T-1555
type: task
nature: feature
title: The design and user docs and the reference describe proposals for the design system, their comparison, and the request form
status: ready
parent: S-0384
owner: alex
created: 2026-10-13T17:48:30Z
updated: 2026-10-15T03:17:16Z
transitions:
  - to: ready
    at: 2026-10-15T03:17:16Z
    by: system-flow
stream: S-0384
tags: [dashboard, cli]
touches: [design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md]
after: [T-1551, T-1554]
---
# T-1555 The design and user docs and the reference describe proposals for the design system, their comparison, and the request form

## Work

Describe `flai design show --story`, `flai design proposals`, and the host methods in `design/system/flai-cli.md` and `docs/users/flai.md`; the proposals list, the comparison, and the request form in `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`; and regenerate `docs/users/flai-reference.md`. It waits for T-1551 and T-1554, whose behaviour it describes.

## Done when

- The five documents describe it, and the markdown lint passes on them.

## Notes

- Layer 5.
