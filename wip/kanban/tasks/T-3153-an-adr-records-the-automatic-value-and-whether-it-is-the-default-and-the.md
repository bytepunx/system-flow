---
id: T-3153
type: task
nature: improvement
title: An ADR records the automatic value and whether it is the default, and the operator docs and flai-cli.md say what mcp.bind holds
status: done
parent: S-0732
owner: alex
created: 2026-11-12T15:55:17Z
updated: 2026-11-12T18:28:40Z
transitions:
  - to: ready
    at: 2026-11-12T18:23:06Z
    by: agent-S-0732
  - to: in-progress
    at: 2026-11-12T18:23:06Z
    by: agent-S-0732
  - to: done
    at: 2026-11-12T18:28:40Z
    by: agent-S-0732
stream: S-0732
tags: [cli, docs]
touches: [design/adrs/README.md, docs/operators/settings.md, docs/operators/index.md, docs/users/flai.md, design/system/flai-cli.md, design/tech/docker.md, design/adrs/0293-mcp-bind-holds-one-address-or-several-or-auto-the-default-which-flai-host.md, docs/operators/runbooks/agent-container.md, docs/operators/runbooks/install.md]
after: [T-3148, T-3150]
usage:
  source: log
  seconds: 334
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 83
      output: 33441
      cache_read: 4943432
      cache_write: 181229
      cost: 2.7303
---
# T-3153 An ADR records the automatic value and whether it is the default, and the operator docs and flai-cli.md say what mcp.bind holds

## Work

Write an ADR with `flai adr new`, refining ADR-0269. It records that `mcp.bind` holds a list or `auto`, what `auto` resolves to on each engine, and whether `auto` is the default when `mcp.bind` is unset, as the plan's thread settles. It also records what a running server does when the engine or its gateway changes. Rewrite the `mcp.bind` row of `docs/operators/settings.md`, and the MCP over HTTP section and the container harness paragraph of `docs/operators/index.md`, so they name the list, `auto`, and the refusal's new first advice. Do the same in the `flai mcp http` row and the container harness's "What it does not do" paragraph of `design/system/flai-cli.md`, and in the `flai serve mcp` commands of `docs/users/flai.md`.

It waits for T-3148 and T-3150, whose behaviour it describes as built.

## Done when

- The ADR is in `design/adrs/` and its index, and says whether `auto` is the default.
- The four documents name the list, `auto`, and the refusal's first advice, with no sentence left saying `mcp.bind` holds one address.
- `TestSettingsIndex` passes.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes
