---
id: T-1696
type: task
nature: feature
title: flai guard holds an OpenCode permission request through ask.Asker.Hold within the ADR's bound and answers allow or deny with the owner's word
status: cancelled
parent: S-0408
owner: alex
created: 2026-10-16T05:16:27Z
updated: 2026-11-07T22:13:46Z
transitions:
  - to: ready
    at: 2026-11-07T21:50:51Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T21:50:52Z
    by: system-flow
  - to: cancelled
    at: 2026-11-07T21:53:43Z
    by: alex
stream: S-0408
tags: [cli]
touches: [flai/cmd/guard.go, flai/cmd/guard_test.go, flai/internal/guard/call.go, flai/internal/guard/call_test.go, flai/internal/guard/opencode.go, flai/internal/guard/opencode_test.go, docs/users/flai-reference.md]
after: [T-1695]
usage:
  source: log
  seconds: 171
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 71
      output: 25664
      cache_read: 10241218
      cache_write: 420809
      cost: 5.047
---
# T-1696 flai guard holds an OpenCode permission request through ask.Asker.Hold within the ADR's bound and answers allow or deny with the owner's word

## Work

- Under hold-and-ask: have the `opencode` reader in `flai/internal/guard/opencode.go` read a `permission.asked` input into a `Call`, and `flai guard` in `flai/cmd/guard.go` hold an `Ask` verdict through `guard.Hold` and `ask.Asker` with the ADR's bound, then print allow or deny in the form the plugin reads. Reuse the story's open permission thread for the same tool, path, and input, as ADR-0124 says.
- Under `--auto` with deny rules: have the reader turn a refused protected write into a refusal that names the thread the agent should open, and leave `Hold` as it is.
- Test a held call answered allow, answered deny, and unanswered past the bound, against a temporary repository.
- Waits for T-1695: the ADR chooses which of the two this task builds.

## Done when

`flai guard --harness opencode` behaves as the ADR says for a protected write, and its tests pass.

## Notes

- If this task gives `flai guard` a flag, add `docs/users/flai-reference.md` to its touches with `flai touches` and regenerate it with `make flai-reference`, as the orchestrator approved on TH-0490.
- 2026-11-07T21:53:43Z: moved to cancelled: E-0023 cancelled: OpenCode isn't going to receive first-class support given its aversion to supporting before hooks for tool calls (amongst other things)
