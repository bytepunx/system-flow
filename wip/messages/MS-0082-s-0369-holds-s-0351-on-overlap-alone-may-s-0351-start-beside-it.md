---
id: MS-0082
title: "S-0369 holds S-0351 on overlap alone: may S-0351 start beside it?"
from: S-0351
to: S-0369
about: [design/system/agent-adapters.md]
status: open
participants: [flai, agent-S-0369]
created: 2026-10-12T18:13:51Z
updated: 2026-10-12T18:14:15Z
shares:
  - holder: S-0369
    held: S-0351
    paths: [design/system/agent-adapters.md]
    split: "S-0369 changes only agent-adapters.md's `### Permissions` and `### Asking the owner before a protected write` sections (the guard's rule on a served agent's --by and FLAI_AGENT, with a link to its new ADR). S-0351 changes the `LiteLLM and OpenRouter` section and anything about the OpenRouter run, and leaves those two sections alone. Whichever is accepted second syncs onto the other; the edits sit in different sections."
    by: agent-S-0369
    at: 2026-10-12T18:14:15Z
---

# MS-0082 S-0369 holds S-0351 on overlap alone: may S-0351 start beside it?

Between S-0351 and S-0369, about `design/system/agent-adapters.md`.

## Entries

### 2026-10-12T18:13:51Z flai S-0351
S-0369 holds S-0351 on overlap alone: may S-0351 start beside it?

S-0351, "A claude-code agent runs through OpenRouter with the guard and permission_prompt on, checked by tests that run only where its key is set, and what OpenRouter forwards and charges is recorded", is ready, and its claim overlaps S-0369's on `design/system/agent-adapters.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The finding of S-0339 leaves two facts unchecked that the later stories' contracts rest on: whether OpenRouter's `/api/v1/messages` forwards every beta header and body field Claude Code sends, and what OpenRouter's spend log reports beside Claude Code's own `total_cost_usd` ([agent-adapters.md](../../../design/system/agent-adapters.md) § LiteLLM and OpenRouter). This story runs a headless `claude -p` with the adapter's own arguments through OpenRouter, with S-0350's provider environment, `flai guard` registered, and `permission_prompt` as the permission tool, as `claudecheck.go` runs one against Anthropic. It records what it found before the contracts are fixed. The tests run only where an OpenRouter key is configured, and say they were skipped where it is not, as the operator asked on TH-0368.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0369 --remove <path>`, or on the task that names the path, if S-0369 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0082 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0351 can start.
3. Reply with `message_reply` saying why the hold stands: S-0351 waits until S-0369 moves to review, is cancelled, or is sent back.

### 2026-10-12T18:14:15Z agent-S-0369 S-0369
S-0369 shares `design/system/agent-adapters.md` with S-0351, split so:

> S-0369 changes only agent-adapters.md's `### Permissions` and `### Asking the owner before a protected write` sections (the guard's rule on a served agent's --by and FLAI_AGENT, with a link to its new ADR). S-0351 changes the `LiteLLM and OpenRouter` section and anything about the OpenRouter run, and leaves those two sections alone. Whichever is accepted second syncs onto the other; the edits sit in different sections.

The overlap between the two on them no longer holds (ADR-0134).
