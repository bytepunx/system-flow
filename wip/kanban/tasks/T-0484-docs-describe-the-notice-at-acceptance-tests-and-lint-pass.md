---
id: T-0484
type: task
nature: feature
title: "Docs describe the notice at acceptance; tests and lint pass"
status: done
parent: S-0132
owner: alex
created: 2026-09-26T18:17:27Z
updated: 2026-09-26T18:25:57Z
transitions:
  - to: ready
    at: 2026-09-26T18:17:41Z
    by: agent-S-0132
  - to: in-progress
    at: 2026-09-26T18:23:06Z
    by: agent-S-0132
  - to: done
    at: 2026-09-26T18:25:57Z
    by: agent-S-0132
stream: S-0132
tags: []
touches: [design/system, docs]
---
# T-0484 Docs describe the notice at acceptance; tests and lint pass

## Work

Describe the notice in `design/system/workflow.md`, `design/system/flai-cli.md`, and `docs/users/flai.md`. Run `make test`, lint, and `flai check --strict`.

## Done when

- The three documents describe it; `make test`, lint, and `flai check --strict` pass.

## Notes
