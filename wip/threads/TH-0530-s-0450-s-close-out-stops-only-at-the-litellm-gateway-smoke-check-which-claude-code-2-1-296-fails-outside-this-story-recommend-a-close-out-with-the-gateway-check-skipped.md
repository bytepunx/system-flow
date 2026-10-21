---
id: TH-0530
title: "S-0450's close-out stops only at the LiteLLM gateway smoke check, which Claude Code 2.1.296 fails outside this story: recommend a close-out with the gateway check skipped"
anchor:
  path: wip/kanban/stories/S-0450-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md
  item: S-0450
status: resolved
participants: [agent-S-0450, orchestrator, alex]
created: 2026-10-20T20:07:50Z
updated: 2026-10-21T00:29:28Z
---

# TH-0530 S-0450's close-out stops only at the LiteLLM gateway smoke check, which Claude Code 2.1.296 fails outside this story: recommend a close-out with the gateway check skipped

On wip/kanban/stories/S-0450-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md.

## Entries

### 2026-10-20T20:07:50Z agent-S-0450
S-0450 is built and both criteria are ticked. Its close-out passed every step and tier except smoke, where `TestGatewayLiteLLM/unified` fails.

- **Cause:** Claude Code 2.1.296 refuses the default haiku model, `openrouter-messages/nvidia/nemotron-3-super-120b-a12b:free`, with `unrecognized_model`, and the session ends before any model call. The model is set in `flai/internal/serve/gateway_test.go` and `scripts/litellm.sh`.
- **Not this story's:** S-0450 changes nothing under `flai/internal/serve`. The test failed the same way when run again alone. It is recorded as I-0171, which is on S-0450's review page.
- **Effect:** every close-out on this host stops here while `LITELLM_BASE_URL` and `LITELLM_API_KEY` are set.

Which do you want?

1. **Recommended:** I run the close-out again with `LITELLM_BASE_URL` and `LITELLM_API_KEY` unset for that one run. The gateway check then skips, as it does on a host with no gateway, and S-0450 goes to review with I-0171 for you to choose at acceptance.
2. S-0450 waits in progress until another story fixes I-0171.
3. I change the default model in S-0450, if you name one that Claude Code 2.1.296 accepts through the proxy.

### 2026-10-20T21:16:45Z orchestrator (recommendation)
**Recommendation:** set `FLAI_TEST_LITELLM_HAIKU` on the host to a model your LiteLLM proxy serves and Claude Code 2.1.296 accepts, then run S-0450's close-out again. This is the same as I recommend on TH-0529 (S-0451) and TH-0534 (S-0445). One host fix lets all three close out cleanly.

**Why not a close-out with the gateway check skipped.**
- ADR-0093 lets me accept only when "the verifier's run passed at the story branch's head". Skipping a check the manifest selects is not a passing run, unless you decide that skipping is allowed.
- That would be a decision about what the close-out tests, recorded nowhere yet, so it is yours.
- If you prefer it, say so here, and the follow-up from I-0170 and I-0171 can make the gateway check skip by itself on a model the installed Claude Code does not recognise.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md § Decision

### 2026-10-20T22:28:30Z alex
Confirmed the recommendation of 2026-10-20T21:16:45Z orchestrator.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md § Decision

### 2026-10-20T22:40:44Z agent-S-0450
S-0450's close-out passed every step, smoke included, and S-0450 is moving to review.

- **How it ran:** I set `FLAI_TEST_LITELLM_HAIKU=openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free` for that one run only. It is the model S-0451 found works through the proxy at no cost (TH-0529).
- **Host still to fix:** I changed no host file. `.flai-cache/litellm/litellm.env` still names the retired `nvidia/nemotron-3-super-120b-a12b:free`, so other close-outs will fail the same way until you change or remove that line.
- **I-0171 corrected:** my diagnosis was wrong. `unrecognized_model` is only a warning. The real cause is OpenRouter's 404 for the retired free model. The correction is committed on story/S-0450. It records the same cause as I-0170, so at acceptance you may prefer to keep I-0170 alone.

### 2026-10-21T00:29:28Z alex
Resolved: S-0450 was accepted
