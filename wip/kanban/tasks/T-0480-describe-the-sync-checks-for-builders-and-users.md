---
id: T-0480
type: task
nature: feature
title: Describe the sync checks for builders and users
status: in-progress
parent: S-0131
owner: alex
created: 2026-09-26T18:05:56Z
updated: 2026-09-26T18:10:50Z
transitions:
  - to: ready
    at: 2026-09-26T18:06:14Z
    by: agent-S-0131
  - to: in-progress
    at: 2026-09-26T18:10:50Z
    by: agent-S-0131
stream: S-0131
tags: []
touches: [design/system, docs, design/tech]
---
# T-0480 Describe the sync checks for builders and users

## Work

- `design/system/workflow.md` and `design/system/flai-cli.md` describe the trial merge, the conflict thread, and the drift check at sync.
- `docs/users/flai.md` describes the new output of `flai stream sync` and its `--json`.
- `design/tech` records git 2.38 as the minimum for the trial merge, if no minimum covers it already.
- `make test`, `make integration`, lint, and `flai check --strict` pass.

## Done when

- The three documents describe it and every tier and lint pass.

## Notes
