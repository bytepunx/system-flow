---
id: T-3262
type: task
nature: feature
title: A runbook says how to change the board, and the guides and design say where the workflow file is and how flai reads it
status: backlog
parent: S-0755
owner: alex
created: 2026-11-12T16:10:29Z
updated: 2026-11-12T16:10:29Z
transitions: []
stream: S-0755
tags: [cli, template]
touches: [docs/operators/runbooks/change-the-board.md, docs/operators/runbooks/README.md, docs/operators/index.md, docs/users/flai.md, docs/users/conventions.md, design/system/workflow.md, design/system/template.md, design/workflows/README.md, template/root/design/workflows/README.md]
after: [T-3260]
---
# T-3262 A runbook says how to change the board, and the guides and design say where the workflow file is and how flai reads it

## Work

Waits for T-3260: it documents the file as the template ships it and as upgrade keeps it. Write `docs/operators/runbooks/change-the-board.md`: adding a state, a rule, and a timer; renaming or removing a state with items in it, as the ADR says; what `flai check` reports; and that the change goes through a story the operator accepts. List it in the runbooks' README and `docs/operators/index.md`. Say in `docs/users/flai.md`, `docs/users/conventions.md`, `design/system/workflow.md`, `design/system/template.md`, and both `design/workflows/README.md` files where the file is and how flai reads it.

## Done when

The runbook exists and is listed, and each document names the file and how flai reads it.

## Notes
