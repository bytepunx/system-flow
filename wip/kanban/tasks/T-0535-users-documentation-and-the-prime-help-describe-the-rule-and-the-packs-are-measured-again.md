---
id: T-0535
type: task
nature: improvement
title: Users' documentation and the prime help describe the rule, and the packs are measured again
status: in-progress
parent: S-0149
owner: alex
created: 2026-09-29T05:25:26Z
updated: 2026-09-29T05:30:28Z
transitions:
  - to: ready
    at: 2026-09-29T05:25:48Z
    by: agent-S-0149
  - to: in-progress
    at: 2026-09-29T05:30:28Z
    by: agent-S-0149
stream: S-0149
tags: []
touches: [docs/users, flai/cmd/prime.go]
---
# T-0535 Users' documentation and the prime help describe the rule, and the packs are measured again

## Work

- Update the `flai prime --story` help, the MCP `prime` row and the Prime a session section of the users' flai guide, and the generated reference.
- Measure the packs of the stories TH-0032 listed (S-0138, S-0141, S-0146, S-0147) with the change, and record the numbers in the story notes and in `agent-context.md`.

## Done when

- The docs say what the pack does, the reference is regenerated, and the measurements are in the notes.
- `make test`, `make integration`, `make smoke`, and `flai check --strict` pass.

## Notes
