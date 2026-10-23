---
id: D-0334
title: "S-0450's close-out stops only at the LiteLLM gateway smoke check, which Claude Code 2.1.296 fails outside this story: recommend a close-out with the gateway check skipped: Recommendation: set…"
kind: answer
sources: [TH-0530]
topics: [mcp, inbox, flai, cli, go, template, conventions, code, all]
paths: [design/adrs/README.md, design/system/agent-narrative.md, design/system/flai-cli.md, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/storystart/start.go, flai/internal/storystart/start_test.go, flai/internal/mcpserver/story_start.go, flai/internal/mcpserver/story_start_test.go, flai/cmd/story_start.go, flai/cmd/story_start_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/session-start.md, template/root/design/conventions/session-start.md, design/issues/I-0167-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md, design/issues/summary.md, design/adrs/0179-task-done-and-story-start-carry-an-inbox-that-lists-the-story-s-threads-and.md, design/issues/I-0171-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md]
date: 2026-10-20
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0334 S-0450's close-out stops only at the LiteLLM gateway smoke check, which Claude Code 2.1.296 fails outside this story: recommend a close-out with the gateway check skipped: Recommendation: set…

The operator, alex, answered TH-0530 at 2026-10-20T22:28:30Z.

What was asked:

```text
S-0450 is built and both criteria are ticked. Its close-out passed every step and tier except smoke, where `TestGatewayLiteLLM/unified` fails.
```

The answer:

```text
**Recommendation:** set `FLAI_TEST_LITELLM_HAIKU` on the host to a model your LiteLLM proxy serves and Claude Code 2.1.296 accepts, then run S-0450's close-out again. This is the same as I recommend on TH-0529 (S-0451) and TH-0534 (S-0445). One host fix lets all three close out cleanly.

**Why not a close-out with the gateway check skipped.**
- ADR-0093 lets me accept only when "the verifier's run passed at the story branch's head". Skipping a check the manifest selects is not a passing run, unless you decide that skipping is allowed.
- That would be a decision about what the close-out tests, recorded nowhere yet, so it is yours.
- If you prefer it, say so here, and the follow-up from I-0170 and I-0171 can make the gateway check skip by itself on a model the installed Claude Code does not recognise.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md § Decision
```

Source: TH-0530.
