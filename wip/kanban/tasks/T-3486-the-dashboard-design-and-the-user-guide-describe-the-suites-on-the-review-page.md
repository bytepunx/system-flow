---
id: T-3486
type: task
nature: feature
title: The dashboard design and the user guide describe the suites on the review page and the board, the job a state's runs started in the stream, and what board.get, verify.status, and resolve.status give
status: backlog
parent: S-0814
owner: alex
created: 2026-11-16T23:51:41Z
updated: 2026-11-16T23:51:41Z
transitions: []
stream: S-0814
tags: [flaiover, docs]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-3481, T-3483, T-3485]
---
# T-3486 The dashboard design and the user guide describe the suites on the review page and the board, the job a state's runs started in the stream, and what board.get, verify.status, and resolve.status give

## Work

Criterion 4's design and docs. In `design/system/flaiover-dashboard.md`, update the `/review/:id` and `/board` rows under `## Views`, and the `GET /api/board`, `GET /api/items/:id/verify`, and the resolve route's entries under `## API (S-0011)`, with what this story's tasks built: the verification section grouped by suite with the suites the next transitions need, the stream following a job a state's `runs:` started with its skipped steps, the lane header saying what a state runs, and the fields `board.get`, `verify.status`, and `resolve.status` now give. In `docs/users/flaiover.md`, say what the operator sees on the review page and the board. Cite ADR-0312 §4.

Waits for T-3481, T-3483, and T-3485, so it describes what they built.

## Done when

- Both documents describe each of the four changes and name ADR-0312; their `updated` dates are bumped.
- The markdown lint and the `design-routes` tier pass on them, and `flai test` on the changed paths passes.

## Notes

- Drafted by planner-S-0814.
