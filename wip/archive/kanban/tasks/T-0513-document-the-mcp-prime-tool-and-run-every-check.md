---
id: T-0513
type: task
nature: feature
title: Document the MCP prime tool and run every check
status: done
parent: S-0138
owner: alex
created: 2026-09-29T01:08:52Z
updated: 2026-09-29T03:15:07Z
transitions:
  - to: ready
    at: 2026-09-29T01:08:56Z
    by: agent-S-0138
  - to: in-progress
    at: 2026-09-29T01:13:31Z
    by: agent-S-0138
  - to: done
    at: 2026-09-29T03:15:07Z
    by: agent-S-0138
stream: S-0138
tags: [cli]
touches: [docs/users]
---
# T-0513 Document the MCP prime tool and run every check

## Work

Document the `prime` MCP tool in `docs/users/flai.md` (and wherever the tool list appears). Run `make test`, `make integration`, `make smoke`, lint, and `flai check --strict`.

## Done when

- The docs name the tool; all three tiers, lint, and `flai check --strict` pass.

## Notes
