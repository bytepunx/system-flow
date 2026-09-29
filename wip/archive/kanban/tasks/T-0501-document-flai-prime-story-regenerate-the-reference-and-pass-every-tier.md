---
id: T-0501
type: task
nature: feature
title: Document flai prime --story, regenerate the reference, and pass every tier
status: done
parent: S-0136
owner: alex
created: 2026-09-28T22:57:55Z
updated: 2026-09-28T23:04:06Z
transitions:
  - to: ready
    at: 2026-09-28T22:57:59Z
    by: agent-S-0136
  - to: in-progress
    at: 2026-09-28T23:02:00Z
    by: agent-S-0136
  - to: done
    at: 2026-09-28T23:04:06Z
    by: agent-S-0136
stream: S-0136
tags: []
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
---
# T-0501 Document flai prime --story, regenerate the reference, and pass every tier

## Work

- Describe `--story` in `design/system/flai-cli.md` and `docs/users/flai.md`; drop the "not yet built" notes this story makes untrue.
- Regenerate the command reference (`make flai-reference`).
- Run `make test`, `make integration`, `make smoke`, lint, `make lint-md`, and `flai check --strict`.

## Done when

- Docs describe the flag, the reference is current, all three tiers and lint pass, and `flai check --strict` is clean.

## Notes
