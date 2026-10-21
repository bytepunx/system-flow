---
id: TH-0534
title: "S-0445's close-out stops only at the LiteLLM gateway smoke check, the failure TH-0529 and TH-0530 report: recommend setting FLAI_TEST_LITELLM_HAIKU on the host, then I run it again"
anchor:
  path: wip/kanban/stories/S-0445-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md
  item: S-0445
status: resolved
participants: [agent-S-0445, orchestrator, alex]
created: 2026-10-20T21:16:28Z
updated: 2026-10-21T02:14:44Z
---

# TH-0534 S-0445's close-out stops only at the LiteLLM gateway smoke check, the failure TH-0529 and TH-0530 report: recommend setting FLAI_TEST_LITELLM_HAIKU on the host, then I run it again

On wip/kanban/stories/S-0445-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md.

## Entries

### 2026-10-20T21:16:28Z agent-S-0445
**Recommended: set `FLAI_TEST_LITELLM_HAIKU` on the host to a model Claude Code 2.1.296 accepts, as TH-0529 recommends; then I run the close-out again.**

S-0445 is otherwise ready for review. Every task is done, both criteria are ticked, and I-0157 is closed.

`scripts/close-out.sh S-0445` passed every step but smoke:

- rebase, sync, narrative, check, gofmt, vet, golangci-lint, markdown: passed
- integration: passed (5m30s), including the new tests that reproduce I-0157
- smoke: failed in `TestGatewayLiteLLM/unified`

The failure is outside S-0445. Claude Code 2.1.296 refuses the gateway's default model as `unrecognized_model`:

```text
[claude-code:unrecognized_model] {"model":"openrouter-messages/nvidia/nemotron-3-super-120b-a12b:free","query_source":"sdk"}
```

S-0451 (TH-0529, I-0170) and S-0450 (TH-0530, I-0171) stopped on the same failure. I did not record a third issue, because both of those live on unmerged branches.

Answer one of:

1. **Set the variable** (recommended). I run the close-out again and move S-0445 to review when it passes.
2. **Skip the gateway check for this close-out**, as TH-0530 recommends. Say how, and I move S-0445 to review on that run.
3. **Move it to review as it is**, with the smoke failure stated in the review.

### 2026-10-20T21:16:41Z orchestrator (recommendation)
**Recommendation:** option (1), the same as on TH-0529 for S-0451 and on TH-0530 for S-0450. Set `FLAI_TEST_LITELLM_HAIKU` on the host to a model your LiteLLM proxy serves and Claude Code 2.1.296 accepts. Then all three agents run their close-outs again.

**One answer settles all three.** S-0445, S-0450 and S-0451 stop on the same `TestGatewayLiteLLM/unified` failure. It is outside each story, recorded as I-0170 and I-0171.

**Why not (2) or (3).** ADR-0093 lets me accept only when "the verifier's run passed at the story branch's head". A run with smoke failing or left out is not that, so S-0445 would wait in review for you. Option (2), skipping the check, is a change to what the close-out tests, and it would want a decision of its own.

A follow-up story from I-0170 should make the gateway check skip, rather than fail, on a model the installed Claude Code does not recognise, so the next Claude Code upgrade cannot stop every close-out again.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md § Decision

### 2026-10-20T21:21:36Z alex
Confirmed the recommendation of 2026-10-20T21:16:41Z orchestrator.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md § Decision

### 2026-10-20T21:22:18Z agent-S-0445
Noted: option 1. The value has not changed yet, so a close-out run now would fail the same way:

- `.flai-cache/litellm/litellm.env`, which `scripts/env.sh` loads for the gateway tests, still holds `FLAI_TEST_LITELLM_HAIKU=openrouter-messages/nvidia/nemotron-3-super-120b-a12b:free`. That is the model Claude Code 2.1.296 refuses. The file was last written 2026-10-13T12:30.
- My session's environment has no `FLAI_TEST_LITELLM_HAIKU`, so nothing overrides the file.

**Next step for you:** in `.flai-cache/litellm/litellm.env`, change `FLAI_TEST_LITELLM_HAIKU` to a model your LiteLLM proxy serves and Claude Code accepts, then reply here. I will run `scripts/close-out.sh S-0445` again and move S-0445 to review when it passes.

### 2026-10-20T22:21:12Z agent-S-0445
One more thing to know. S-0451's T-1928 is changing the gateway check's built-in default model to nemotron-3-ultra, because OpenRouter retired nemotron-3-super. That fix alone will not clear this failure. `.flai-cache/litellm/litellm.env` sets `FLAI_TEST_LITELLM_HAIKU` to the retired nemotron-3-super, and that setting overrides the built-in default. Change that line, or remove it so the new default applies once S-0451 lands, then reply here. I will run the close-out again.

### 2026-10-21T00:51:20Z agent-S-0445
Correction to my last two entries: you do not need to change `.flai-cache/litellm/litellm.env`.

- This session's environment already sets `LITELLM_BASE_URL`, so `scripts/env.sh` never loads that file.
- The smoke check therefore used its built-in default model, nemotron-3-super. OpenRouter has retired it, as I-0170's correction found.
- S-0451's T-1928 changes that default to `openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free`, and S-0451's close-out passed smoke with it.

Under option 1, which you confirmed, I am running S-0445's close-out now with `FLAI_TEST_LITELLM_HAIKU` set to that model for this run only. Nothing on the host is changed. If smoke passes, I move S-0445 to review. Once S-0451 is accepted, the default covers every close-out and the variable is not needed.

### 2026-10-21T02:14:44Z alex
Resolved: S-0445 was accepted
