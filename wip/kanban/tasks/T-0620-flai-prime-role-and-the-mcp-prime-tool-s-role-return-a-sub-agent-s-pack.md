---
id: T-0620
type: task
nature: improvement
title: flai prime --role and the MCP prime tool's role return a sub-agent's pack
status: done
parent: S-0175
owner: arobson
created: 2026-10-01T07:53:01Z
updated: 2026-10-01T08:03:41Z
transitions:
  - to: ready
    at: 2026-10-01T07:53:25Z
    by: agent-S-0175
  - to: in-progress
    at: 2026-10-01T07:54:42Z
    by: agent-S-0175
  - to: done
    at: 2026-10-01T08:03:41Z
    by: agent-S-0175
stream: S-0175
tags: []
touches: [flai/internal/context, flai/internal/conventions, flai/cmd/prime.go, flai/internal/mcpserver]
usage:
  source: log
  seconds: 539
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 92
      output: 463
      cache_read: 8447770
      cache_write: 110215
      cost: 3.8384
---
# T-0620 flai prime --role and the MCP prime tool's role return a sub-agent's pack

## Work

- Conventions gain `roles` in front matter: the sub-agent roles that read the file.
- `context.ForRole`: the conventions whose roles include the role, filtered by the story's topics; the story's goal and acceptance criteria; then briefs of what the story names and what its topics and one link step select, in order, while the budget has room. Default budget half the project's.
- `flai prime --story S-nnnn --role explore|verify`, and `role` on the MCP `prime` tool.
- Behaviour tests for the pack, the flag, and the tool.

## Done when

- `flai prime --story S-0175 --role explore` and `--role verify` print packs under half the budget with the right conventions, and `make test` and lint pass.

## Notes
