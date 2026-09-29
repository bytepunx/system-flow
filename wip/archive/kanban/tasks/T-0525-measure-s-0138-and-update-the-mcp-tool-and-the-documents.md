---
id: T-0525
type: task
nature: feature
title: Measure S-0138, and update the MCP tool and the documents
status: done
parent: S-0146
owner: alex
created: 2026-09-29T03:19:38Z
updated: 2026-09-29T03:55:58Z
transitions:
  - to: ready
    at: 2026-09-29T03:19:45Z
    by: agent-S-0146
  - to: in-progress
    at: 2026-09-29T03:55:58Z
    by: agent-S-0146
  - to: done
    at: 2026-09-29T03:55:58Z
    by: agent-S-0146
stream: S-0146
tags: []
touches: [flai/internal/mcpserver, design/system/agent-context.md, design/system/flai-cli.md, docs/users, docs/operators]
---
# T-0525 Measure S-0138, and update the MCP tool and the documents

## Work

Measure `flai prime --story S-0138` and record it in the story notes. Check `--json` and the MCP `prime` tool carry the same pack and update the tool description. Update `design/system/agent-context.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, the reference, and the operator settings index for `prime.budget`.

## Done when

- The S-0138 pack measures under 80 KB.
- `make test`, `make integration`, `make smoke`, lint, and `flai check --strict` pass.

## Notes
