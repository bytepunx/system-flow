---
id: TH-0529
title: "S-0451's close-out fails only at smoke, on a gateway model Claude Code no longer recognises (I-0170): recommend setting FLAI_TEST_LITELLM_HAIKU on the host, then I run it again"
anchor:
  path: wip/kanban/stories/S-0451-message-share-refuses-a-split-between-two-stories-in-progress-so-an-agreed-split-cannot-lift-the-overlap-that-refuses-the-other-story-s-flai-task-new.md
  item: S-0451
status: answered
participants: [agent-S-0451, orchestrator, alex]
created: 2026-10-20T19:17:18Z
updated: 2026-10-20T22:20:02Z
---

# TH-0529 S-0451's close-out fails only at smoke, on a gateway model Claude Code no longer recognises (I-0170): recommend setting FLAI_TEST_LITELLM_HAIKU on the host, then I run it again

On wip/kanban/stories/S-0451-message-share-refuses-a-split-between-two-stories-in-progress-so-an-agreed-split-cannot-lift-the-overlap-that-refuses-the-other-story-s-flai-task-new.md.

## Entries

### 2026-10-20T19:17:18Z agent-S-0451
S-0451 is done apart from the close-out's smoke tier. Both tasks are closed and both criteria are ticked.

The close-out passed rebase, sync, narrative, check, gofmt, vet, golangci-lint, markdown, and integration. It failed only at smoke, in `TestGatewayLiteLLM/unified`. The host's Claude Code 2.1.296 answers `[claude-code:unrecognized_model]` for the check's default model, `openrouter-messages/nvidia/nemotron-3-super-120b-a12b:free` (`litellmHaikuDefault` in `flai/internal/serve/gateway_test.go`). S-0451 changes only `flai/internal/itemnew` and two docs, so it does not cause this. It will fail every close-out on this host. I recorded it as I-0170.

Which do you want?

1. (Recommended) Set `FLAI_TEST_LITELLM_HAIKU` on the host to a model the proxy serves and Claude Code accepts, or fix the proxy's model list. Then reply here and I run the close-out again.
2. Let S-0451 go to review with this smoke failure recorded in I-0170. I will move it without a clean close-out.
3. Something else, such as a story from I-0170 that changes `litellmHaikuDefault`, which S-0451 waits for.

### 2026-10-20T19:17:32Z orchestrator (recommendation)
**Recommendation:** option (1). Set `FLAI_TEST_LITELLM_HAIKU` on the host to a model your LiteLLM proxy serves and Claude Code 2.1.296 accepts, or fix the proxy's model list. Then reply here, and S-0451's agent runs the close-out again.

**Why this needs you.** It is a change to the host's environment, which no agent makes.

**Why not (2).**
- ADR-0093 lets me accept only when "the verifier's run passed at the story branch's head". A close-out that fails at smoke leaves the story's verify failing, so I could not accept S-0451, however small its diff.
- It would wait in review for you instead, and every story after it would hit the same wall. S-0450 is in progress now and changes `flai/internal/mcpserver`, so it selects smoke too.

**On (3).** A story from I-0170 that changes `litellmHaikuDefault`, or makes the gateway smoke check skip when the model is unrecognised rather than fail, is still worth doing so the next Claude Code upgrade cannot break every close-out. It is not a reason to hold S-0451 behind it once the host setting is fixed.

If you would rather not touch the host now, answer (2) and accept S-0451 yourself: its two criteria and its narrow diff need no smoke run.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md § Decision

### 2026-10-20T22:15:15Z alex
what version of haiku will claude code 2.1.296 accept?

### 2026-10-20T22:18:58Z agent-S-0451
Claude Code 2.1.296 accepts any model name the proxy serves. My first diagnosis was wrong. `unrecognized_model` is only a warning Claude Code prints for a model it has no metadata for, and it prints it on passing runs too.

The real cause is upstream. OpenRouter answers 404 ("Function ... Not found for account") for `nvidia/nemotron-3-super-120b-a12b:free` through the proxy. The free model is retired for this account. I corrected I-0170 to say so (6a855883 on story/S-0451).

What works through the proxy now, at one token each:

| Model (`openrouter-messages/...`) | Result |
|---|---|
| `nvidia/nemotron-3-super-120b-a12b:free` (the current default) | 404 |
| `nvidia/nemotron-3-ultra-550b-a55b:free` | answers |
| `nvidia/nemotron-3.5-super-vl:free` | answers |
| `minimax/minimax-m3.1-flash:free` | answers |
| `poolside/laguna-s-2.1:free` | answers |
| `google/gemma-4-31b-it:free` | rate-limited |

`TestGatewayLiteLLM` passes with `FLAI_TEST_LITELLM_HAIKU=openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free`. It ran in 143s at $0 of key spend. The Write passed and the guard refused the sub-agent's `item_edit`.

Which do you want?

1. (Recommended) I add a task to S-0451 that changes `litellmHaikuDefault` in `flai/internal/serve/gateway_test.go` to `openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free`. Then I run the close-out again and move S-0451 to review. It is a one-line test change, it unblocks every story's close-out without touching the host, and I-0170 stays open for a lasting fix: a skip or fallback when the proxy answers 404.
2. You set `FLAI_TEST_LITELLM_HAIKU` to that model in the environment `flai serve` and its agents start from. Then I run the close-out again with no code change.

### 2026-10-20T22:20:02Z alex
1
