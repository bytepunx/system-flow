---
id: T-1751
type: task
nature: feature
title: The design and the user and operator docs say how each trigger starts a role
status: in-progress
parent: S-0418
owner: alex
created: 2026-10-16T05:45:58Z
updated: 2026-10-29T17:20:52Z
transitions:
  - to: ready
    at: 2026-10-29T17:20:50Z
    by: agent-S-0418
  - to: in-progress
    at: 2026-10-29T17:20:52Z
    by: agent-S-0418
stream: S-0418
tags: [cli, template]
touches: [design/system/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md, template/CHANGELOG.md, design/system/flai-cli.md]
after: [T-2340]
---
# T-1751 The design and the user and operator docs say how each trigger starts a role

## Work

- In `design/system/strategic-agents.md`, under `## A role a project defines (S-0417)` only (the split agreed with S-0537 on MS-0232), describe how an event and a schedule start a role, the coalescing, and what each run records.
- In `docs/users/flai.md` and `flai-reference.md`, say how a role's trigger starts it.
- In `docs/operators/index.md`, say how to give a role a trigger; in `docs/operators/settings.md`, add `roles.<name>.schedule`.
- Name the triggers in `template/CHANGELOG.md`, under the next version after the last entry on main (D-0410).
- Add to `design/system/flai-cli.md` what the serve loop now starts.

Last layer: it describes what T-1748 wires.

## Done when

The markdown lint and `flai check --strict` pass.

## Notes
