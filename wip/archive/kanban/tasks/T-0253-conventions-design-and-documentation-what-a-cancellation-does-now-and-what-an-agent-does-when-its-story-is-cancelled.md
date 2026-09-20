---
id: T-0253
type: task
nature: improvement
title: "Conventions, design, and documentation: what a cancellation does now, and what an agent does when its story is cancelled"
status: done
parent: S-0070
owner: alex
created: 2026-09-20T06:16:44Z
updated: 2026-09-20T06:33:48Z
transitions:
  - to: ready
    at: 2026-09-20T06:30:42Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T06:30:43Z
    by: system-flow
  - to: done
    at: 2026-09-20T06:33:48Z
    by: system-flow
stream: S-0070
tags: []
---
# T-0253 Conventions, design, and documentation: what a cancellation does now, and what an agent does when its story is cancelled

## Work
The baseline rule for agents in `template/root/design/conventions/work-management.md` first, then the same text here. `design/system/workflow.md`, `work-hierarchy.md`, `flai-cli.md`, `flaiover-dashboard.md`, and `docs/operators`. A test that `flai stats` counts cascaded items as cancelled at the time of the cascade, with `metrics.md` unchanged.

## Done when
- `make flai-test`, `make flaiover-test`, and `make flaiover-build` pass
- `flai check --strict` passes

## Notes
