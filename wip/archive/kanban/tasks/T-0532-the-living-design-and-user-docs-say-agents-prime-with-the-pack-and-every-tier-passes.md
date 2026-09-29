---
id: T-0532
type: task
nature: feature
title: The living design and user docs say agents prime with the pack, and every tier passes
status: done
parent: S-0148
owner: alex
created: 2026-09-29T05:24:54Z
updated: 2026-09-29T05:29:45Z
transitions:
  - to: ready
    at: 2026-09-29T05:25:01Z
    by: agent-S-0148
  - to: in-progress
    at: 2026-09-29T05:26:52Z
    by: agent-S-0148
  - to: done
    at: 2026-09-29T05:29:45Z
    by: agent-S-0148
stream: S-0148
tags: []
touches: [design/system/conventions.md, design/system/agent-context.md, docs/users]
---
# T-0532 The living design and user docs say agents prime with the pack, and every tier passes

## Work

`design/system/conventions.md` "Priming", `design/system/agent-context.md`, and `docs/users/flai.md` (prime section and the MCP `prime` row) say agents with a story prime with the budgeted pack and fetch what it briefs. Run `make test`, `make integration`, `make smoke`, lint, and `flai check --strict`.

## Done when

- [x] The three documents say agents prime with the pack; nothing says they still use `--cat` with a story.
- [x] All three test tiers, lint, `make smoke`, and `flai check --strict` pass.

## Notes
