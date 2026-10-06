---
id: T-1021
type: task
nature: remediation
title: The workflow design and the operators' guide say that accept_reviews ends the board's wait on the operator, and which stories still wait
status: backlog
parent: S-0296
owner: alex
created: 2026-10-06T11:50:11Z
updated: 2026-10-06T11:50:11Z
transitions: []
stream: S-0296
tags: [docs]
touches: [design/system/workflow.md, docs/operators/settings.md, docs/operators/index.md]
---
# T-1021 The workflow design and the operators' guide say that accept_reviews ends the board's wait on the operator, and which stories still wait

## Work

Say in `design/system/workflow.md`, where it describes acceptance and the review column, that a story in review keeps its claim and, while the operator alone accepts, holds every ready story that overlaps it. Say that `orchestrate` with `orchestration.permissions.accept_reviews` on lets the orchestrator accept it under S-0221's conditions, so that the wait remains only for a story that changes a path Claude Code protects (S-0286) and one the operator keeps for review.

In `docs/operators/settings.md` and `docs/operators/index.md`, tell the operator how to turn it on and what it changes for an idle board. It changes no path T-1020 changes and needs nothing it builds, so the two run together.

## Done when

- The design and both guides say it, linking I-0088's cause and S-0221's conditions rather than repeating them.
- The markdown lint and `flai check --strict` pass.

## Notes
