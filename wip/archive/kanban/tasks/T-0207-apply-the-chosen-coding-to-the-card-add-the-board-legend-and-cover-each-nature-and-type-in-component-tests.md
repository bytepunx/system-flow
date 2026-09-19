---
id: T-0207
type: task
nature: feature
title: Apply the chosen coding to the card, add the board legend, and cover each nature and type in component tests
status: done
parent: S-0055
owner: alex
created: 2026-09-19T06:05:58Z
updated: 2026-09-19T06:19:55Z
transitions:
  - to: ready
    at: 2026-09-19T06:17:24Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T06:17:24Z
    by: system-flow
  - to: done
    at: 2026-09-19T06:19:55Z
    by: system-flow
stream: S-0055
tags: []
---

# T-0207 Apply the chosen coding to the card, add the board legend, and cover each nature and type in component tests

## Work
Apply the operator's choice in `BoardCard.svelte` using the token classes, not literals: a map from nature and from type to a class, with an unknown value falling back to today's plain card. Keep the nature, type, and BLOCKED text. Make the blocked state and the dragging state obvious on every tint. Add a legend to the board page that names each colour, built from the same map so it cannot drift from the cards. Component tests: one case per nature and per type, the fallback, blocked, and dragging; a test for the legend.

## Done when
- Every nature and type renders its token class, and the tests fail against main's card
- The legend is on the board and is generated from the same map
- `make flaiover-test` and `make flaiover-build` pass

## Notes
