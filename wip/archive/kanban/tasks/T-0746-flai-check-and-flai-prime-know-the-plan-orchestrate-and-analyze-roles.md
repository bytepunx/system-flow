---
id: T-0746
type: task
nature: improvement
title: flai check and flai prime know the plan, orchestrate, and analyze roles
status: done
parent: S-0207
owner: alex
created: 2026-10-03T07:11:12Z
updated: 2026-10-03T07:20:59Z
transitions:
  - to: ready
    at: 2026-10-03T07:11:43Z
    by: agent-S-0207
  - to: in-progress
    at: 2026-10-03T07:11:43Z
    by: agent-S-0207
  - to: done
    at: 2026-10-03T07:20:59Z
    by: agent-S-0207
stream: S-0207
tags: [flai]
touches: [flai/internal/conventions, flai/internal/context, flai/cmd/prime.go, flai/internal/mcpserver, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 556
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 112
      output: 33791
      cache_read: 6592682
      cache_write: 134115
      cost: 2.8393
---
# T-0746 flai check and flai prime know the plan, orchestrate, and analyze roles

## Work

- Replace the nouns `orchestrator`, `planner`, `analyzer` in `flai/internal/conventions` with the roles `plan`, `orchestrate`, `analyze`, which `flai check` accepts (`conventions.roles`), and list them as the strategic roles.
- Build a strategic agent's pack in `flai/internal/context`: the conventions whose roles are empty or list the role, with the sections the pack's topics leave out taken out; the role's topic (`planning`, `orchestration`, `analysis`) and, for the planner, the item's topics; for the planner, what its epic or story names, whole, and the design its topics and links select, ranked and catalogued as a story's pack is; for the orchestrator and the analyzer, briefs of the design their topic selects.
- `flai prime --role plan --epic E-nnnn|--story S-nnnn`, `--role orchestrate`, and `--role analyze` print the packs, and the MCP tool `prime` takes `epic` and the new roles; regenerate `docs/users/flai-reference.md`.
- Waits for nothing: the first layer.

## Done when

- `flai prime --role plan --epic E-0016`, `--role plan --story S-0207`, `--role orchestrate`, and `--role analyze` print packs; `--role orchestrate --story` and `--role plan` without an item are refused with a message naming what to give.
- Tests in `flai/internal/context` and `flai/internal/conventions` cover each role's pack and the roles check accepts, and pass.

## Notes
