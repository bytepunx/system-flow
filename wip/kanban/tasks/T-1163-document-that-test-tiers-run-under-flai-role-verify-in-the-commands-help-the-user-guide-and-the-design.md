---
id: T-1163
type: task
nature: feature
title: "Document that test tiers run under FLAI_ROLE=verify, in the commands' help, the user guide, and the design"
status: done
parent: S-0311
owner: alex
created: 2026-10-07T14:29:39Z
updated: 2026-10-07T14:42:35Z
transitions:
  - to: ready
    at: 2026-10-07T14:32:45Z
    by: agent-S-0311
  - to: in-progress
    at: 2026-10-07T14:32:45Z
    by: agent-S-0311
  - to: done
    at: 2026-10-07T14:42:35Z
    by: agent-S-0311
stream: S-0311
tags: [cli, docs]
touches: [flai/cmd/verify.go, flai/cmd/test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/system/strategic-agents.md]
after: [T-1161]
usage:
  source: log
  seconds: 590
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 35
      output: 9885
      cache_read: 1473896
      cache_write: 62372
      cost: 0.8946
---
# T-1163 Document that test tiers run under FLAI_ROLE=verify, in the commands' help, the user guide, and the design

## Work

Say what T-1161 built, where each reader looks:

- The `Long` help of `flai verify` (`flai/cmd/verify.go`) and `flai test` (`flai/cmd/test.go`): each tier runs with `FLAI_ROLE=verify`, beside `CLOSE_OUT_STORY`, whatever role ran the command. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- `docs/users/flai.md`, under Run the tests for what changed and Verify a story before review: the same, in a sentence each.
- `design/system/flai-cli.md`, the rows of `flai test` and `flai verify` in Commands: the same, citing S-0311.
- `design/system/strategic-agents.md`, under The orchestrator › Accepting a story (S-0221), step 1: the orchestrator's `flai verify` runs its tiers under `verify`, so a story whose tests run flai writes is not refused by the orchestrator's own checks (TH-0260). Bump `updated` on each design and docs file changed.

This task waits for T-1161, so that it describes what was built. It shares no file with T-1162, so the two run together.

## Done when

- [ ] The help of `flai verify` and `flai test`, `docs/users/flai-reference.md`, `docs/users/flai.md`, `design/system/flai-cli.md`, and `design/system/strategic-agents.md` say that tiers run under `FLAI_ROLE=verify`.
- [ ] `flai test` on the changed paths passes, the markdown lint among its tiers.

## Notes
