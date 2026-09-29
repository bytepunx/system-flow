---
id: T-0528
type: task
nature: feature
title: Document doc_search and heading, point the pack's briefs at them, and pass every tier
status: done
parent: S-0147
owner: alex
created: 2026-09-29T05:00:19Z
updated: 2026-09-29T05:14:59Z
transitions:
  - to: ready
    at: 2026-09-29T05:00:26Z
    by: agent-S-0147
  - to: in-progress
    at: 2026-09-29T05:06:24Z
    by: agent-S-0147
  - to: done
    at: 2026-09-29T05:14:59Z
    by: agent-S-0147
stream: S-0147
tags: []
touches: [docs/users, design/system/flai-cli.md, flai/internal/context/context.go]
---
# T-0528 Document doc_search and heading, point the pack's briefs at them, and pass every tier

## Work

- `docs/users/flai.md`, `docs/users/flai-reference.md` (`make flai-reference`), and `design/system/flai-cli.md` document `doc_search`, `flai doc search`, and `heading`.
- The pack's briefs, decisions, and catalog say how to fetch one section.
- Run `make test`, `make integration`, `make smoke`, lint, and `flai check --strict`.

## Done when

- All three tiers, lint, and `flai check --strict` pass, and every criterion of S-0147 is ticked.

## Notes
