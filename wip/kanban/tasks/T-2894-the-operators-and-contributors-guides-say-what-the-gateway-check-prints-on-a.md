---
id: T-2894
type: task
nature: remediation
title: The operators and contributors guides say what the gateway check prints on a 402 and what to do about it
status: ready
parent: S-0685
owner: alex
created: 2026-11-10T08:44:31Z
updated: 2026-11-10T08:49:22Z
transitions:
  - to: ready
    at: 2026-11-10T08:49:22Z
    by: agent-S-0685
stream: S-0685
tags: [gateway, docs]
touches: [docs/operators/index.md, docs/contributors/index.md]
after: [T-2893]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 17
          cache_read: 114072
          cache_write: 3023
          cost: 0.0317
---
# T-2894 The operators and contributors guides say what the gateway check prints on a 402 and what to do about it

## Work

Waits for T-2893, whose failure message it describes.

- In `docs/operators/index.md`, section "The gateway checks on the host's schedule (S-0548)": say that a key out of credit fails the OpenRouter check with one message naming OpenRouter's 402, and that a check with `limit_remaining` at 0 fails before it starts `claude`. Say that topping up the key or raising its limit clears it, and that removing the key from the environment skips the check instead.
- In `docs/contributors/index.md`, section "Gateway checks": add the same in one sentence, linking to the operators guide rather than repeating it.

## Done when

- Both sections name the 402 message and what clears it, and the contributors guide links to the operators guide's section.
- `flai test docs/operators/index.md docs/contributors/index.md` passes the markdown lint.

## Notes
