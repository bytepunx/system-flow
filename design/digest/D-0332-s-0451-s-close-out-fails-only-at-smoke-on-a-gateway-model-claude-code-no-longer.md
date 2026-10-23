---
id: D-0332
title: "S-0451's close-out fails only at smoke, on a gateway model Claude Code no longer recognises (I-0170): recommend setting FLAI_TEST_LITELLM_HAIKU on the host, then I run it again: what version of haiku…"
kind: recommendation
sources: [TH-0529]
topics: [flai, cli, go, code, all]
paths: [flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go, design/issues/I-0169-message-share-refuses-a-split-between-two-stories-in-progress-so-an-agreed-split-cannot-lift-the-overlap-that-refuses-the-other-story-s-flai-task-new.md, design/issues/summary.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0170-the-litellm-gateway-smoke-check-s-default-free-model-is-refused-by-claude-code-as-unrecognized-model-so-every-close-out-on-this-host-fails-smoke.md, docs/contributors/index.md, docs/operators/index.md, flai/internal/serve/gateway_test.go, scripts/litellm.sh]
date: 2026-10-20
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0332 S-0451's close-out fails only at smoke, on a gateway model Claude Code no longer recognises (I-0170): recommend setting FLAI_TEST_LITELLM_HAIKU on the host, then I run it again: what version of haiku…

The operator, alex, answered TH-0529 at 2026-10-20T22:15:15Z.

What was asked:

```text
S-0451 is done apart from the close-out's smoke tier. Both tasks are closed and both criteria are ticked.
```

The recommendation it overruled:

```text
**Recommendation:** option (1). Set `FLAI_TEST_LITELLM_HAIKU` on the host to a model your LiteLLM proxy serves and Claude Code 2.1.296 accepts, or fix the proxy's model list. Then reply here, and S-0451's agent runs the close-out again.

**Why this needs you.** It is a change to the host's environment, which no agent makes.

**Why not (2).**
- ADR-0093 lets me accept only when "the verifier's run passed at the story branch's head". A close-out that fails at smoke leaves the story's verify failing, so I could not accept S-0451, however small its diff.
- It would wait in review for you instead, and every story after it would hit the same wall. S-0450 is in progress now and changes `flai/internal/mcpserver`, so it selects smoke too.

**On (3).** A story from I-0170 that changes `litellmHaikuDefault`, or makes the gateway smoke check skip when the model is unrecognised rather than fail, is still worth doing so the next Claude Code upgrade cannot break every close-out. It is not a reason to hold S-0451 behind it once the host setting is fixed.

If you would rather not touch the host now, answer (2) and accept S-0451 yourself: its two criteria and its narrow diff need no smoke run.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md § Decision
```

The operator's answer:

```text
what version of haiku will claude code 2.1.296 accept?
```

Source: TH-0529.
