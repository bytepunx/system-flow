---
id: T-1718
type: task
nature: feature
title: The integration test asserts the guard's refusal, the permission prompt's answer, the gateway-priced usage, and the card's status from the run
status: cancelled
parent: S-0412
owner: alex
created: 2026-10-16T05:20:39Z
updated: 2026-11-07T21:53:43Z
transitions:
  - to: cancelled
    at: 2026-11-07T21:53:43Z
    by: alex
stream: S-0412
tags: [cli]
touches: [flai/tests/integration/opencode_test.go]
after: [T-1717]
---
# T-1718 The integration test asserts the guard's refusal, the permission prompt's answer, the gateway-priced usage, and the card's status from the run

## Work

- In `flai/tests/integration/opencode_test.go`, after the run T-1717 drives:
  - assert the sub-agent's `flai` write was refused by the guard and the story's own agent's was let through;
  - assert the protected write was held on the story's permission thread and allowed by an owner's answer the test writes, or refused by the deny rule, as S-0408's ADR settled;
  - serve `generation.json` from a fake OpenRouter on `httptest`, as `flai/internal/usage/openrouter_test.go` does, name it as the story's provider, and assert the story's `usage` carries the session's tokens per model and is `priced_by: gateway` with the fixture's cost;
  - assert the story's card shows the session's last step from the OpenCode stream.
- A defect this finds in an earlier E-0023 story's code is fixed here when it is small, or recorded with `flai issue new` and named in the story's Notes.
- Waits for T-1717: it asserts on the run that task drives, in the same file.

## Done when

Each assertion passes in the integration tier with no `opencode`, network, or key.

## Notes
- 2026-11-07T21:53:43Z: moved to cancelled: E-0023 cancelled: OpenCode isn't going to receive first-class support given its aversion to supporting before hooks for tool calls (amongst other things)
