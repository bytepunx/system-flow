---
id: T-0476
type: task
nature: feature
title: "Describe after: for builders and users"
status: done
parent: S-0130
owner: alex
created: 2026-09-26T17:55:03Z
updated: 2026-09-26T18:04:48Z
transitions:
  - to: ready
    at: 2026-09-26T17:55:13Z
    by: agent-S-0130
  - to: in-progress
    at: 2026-09-26T18:02:25Z
    by: agent-S-0130
  - to: done
    at: 2026-09-26T18:04:48Z
    by: agent-S-0130
stream: S-0130
tags: []
touches: [design/system, docs, flai/cmd]
---
# T-0476 Describe after: for builders and users

## Work

- `design/system/work-hierarchy.md`, `workflow.md`, and `flai-cli.md` describe `after:`, the `after` hold, and its check findings.
- `docs/users/flai.md` documents `flai edit --after` / `--clear-after` and the hold; regenerate the command reference if one is generated.
- The MCP tool descriptions for `board`, `inbox`, and `wait_for_work` say a story is held by overlap or by `after`.

## Done when

- The documents say what was built; `make lint-md` finds nothing in the story's files; `flai check --strict` passes.

## Notes
