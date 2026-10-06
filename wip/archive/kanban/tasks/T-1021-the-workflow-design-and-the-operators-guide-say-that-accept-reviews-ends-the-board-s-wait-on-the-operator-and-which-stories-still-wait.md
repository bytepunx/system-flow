---
id: T-1021
type: task
nature: remediation
title: The workflow design and the operators' guide say that accept_reviews ends the board's wait on the operator, and which stories still wait
status: done
parent: S-0296
owner: alex
created: 2026-10-06T11:50:11Z
updated: 2026-10-06T20:15:41Z
transitions:
  - to: ready
    at: 2026-10-06T20:11:19Z
    by: agent-S-0296
  - to: in-progress
    at: 2026-10-06T20:11:20Z
    by: agent-S-0296
  - to: done
    at: 2026-10-06T20:15:41Z
    by: agent-S-0296
stream: S-0296
tags: [docs]
touches: [design/system/workflow.md, docs/operators/settings.md, docs/operators/index.md]
usage:
  source: log
  seconds: 261
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 30
      output: 9664
      cache_read: 1486910
      cache_write: 57004
      cost: 0.8556
---
# T-1021 The workflow design and the operators' guide say that accept_reviews ends the board's wait on the operator, and which stories still wait

## Work

Say in `design/system/workflow.md`, where it describes acceptance and the review column, what still waits on a story's acceptance now that a story in review holds nothing by overlap (ADR-0096): a ready story that names it in `after:`, and every ready story while review is at its limit. While the operator alone accepts, that wait is as long as the operator is away (I-0088). Say that `orchestrate` with `orchestration.permissions.accept_reviews` on lets the orchestrator accept it under ADR-0093's conditions, so that the wait remains only for a story the orchestrator cannot vouch for and one the operator keeps for review. S-0286's rule for a story that changes a path Claude Code protects is not built: leave it to S-0286 rather than describe it.

In `docs/operators/settings.md` and `docs/operators/index.md`, tell the operator how to turn it on and what it changes for an idle board. It changes no path T-1020 changes and needs nothing it builds, so the two run together.

## Done when

- The design and both guides say it, linking I-0088's cause and S-0221's conditions rather than repeating them.
- The markdown lint and `flai check --strict` pass.

## Notes
