---
id: T-0211
type: task
nature: feature
title: "Living design and CLI documentation for the pull order and flai order; the project convention that said the board cannot reorder"
status: done
parent: S-0057
owner: alex
created: 2026-09-19T06:41:25Z
updated: 2026-09-19T06:48:26Z
transitions:
  - to: ready
    at: 2026-09-19T06:47:45Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T06:47:45Z
    by: system-flow
  - to: done
    at: 2026-09-19T06:48:26Z
    by: system-flow
stream: S-0057
tags: []
---

# T-0211 Living design and CLI documentation for the pull order and flai order; the project convention that said the board cannot reorder

## Work
`design/system/workflow.md`: what the pull order is, which states have one, how unplaced stories sort, what `flai move` does to it, and `flai order`. `design/system/flai-cli.md` and `docs/users/flai.md`: the command. The project addition in `design/conventions/work-management.md` that says the board gives the operator no way to change the order until S-0057 is rewritten for what is now true; it is a project rule below the marker, not baseline, so it is edited here. `template/` only if its copies of the design or docs describe the order.

## Done when
- The design, the CLI reference, and the user docs describe the order and the command
- The project convention no longer says the order cannot be changed from the board

## Notes
