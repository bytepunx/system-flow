---
id: T-1751
type: task
nature: feature
title: The design and the user and operator docs say how each trigger starts a role
status: backlog
parent: S-0418
owner: alex
created: 2026-10-16T05:45:58Z
updated: 2026-10-16T05:50:22Z
transitions: []
stream: S-0418
tags: [cli, template]
touches: [design/system/strategic-agents.md, docs/users/flai.md, docs/operators/index.md, template/CHANGELOG.md, design/system/flai-cli.md]
after: [T-1748]
---
# T-1751 The design and the user and operator docs say how each trigger starts a role

## Work

- In `design/system/strategic-agents.md`, describe the events, the schedule, and `continuous`, what each run records, and the hold.
- Describe `flai role stop` and `start` in `docs/users/flai.md`.
- Add how to give a role a trigger to `docs/operators/index.md`.
- Name the triggers in `template/CHANGELOG.md`.

Last layer.

## Done when

The markdown lint and `flai check --strict` pass.

## Notes
