---
id: T-0498
type: task
nature: feature
title: Document topics on items for builders and users, and run all three tiers
status: done
parent: S-0135
owner: alex
created: 2026-09-28T22:41:24Z
updated: 2026-09-28T22:54:26Z
transitions:
  - to: ready
    at: 2026-09-28T22:41:40Z
    by: agent-S-0135
  - to: in-progress
    at: 2026-09-28T22:50:25Z
    by: agent-S-0135
  - to: done
    at: 2026-09-28T22:54:26Z
    by: agent-S-0135
stream: S-0135
tags: []
touches: [design/system, docs/users, docs/operators/settings.md]
---

# T-0498 Document topics on items for builders and users, and run all three tiers

## Work

Document the key in `design/system/work-hierarchy.md`, the commands in `docs/users/flai.md` and the reference, the dashboard in `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md`, and the MCP tools in `design/system/flai-cli.md` where they are listed. Note that an older flai refuses an item with `topics`. Run `make test`, `make integration`, `make smoke`, lint, and `flai check --strict`.

## Done when

All three tiers and lint pass, `flai check --strict` is clean, and the docs say how to set topics.

## Notes
