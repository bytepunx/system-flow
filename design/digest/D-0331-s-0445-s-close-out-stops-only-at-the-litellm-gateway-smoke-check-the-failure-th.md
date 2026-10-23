---
id: D-0331
title: "S-0445's close-out stops only at the LiteLLM gateway smoke check, the failure TH-0529 and TH-0530 report: recommend setting FLAI_TEST_LITELLM_HAIKU on the host, then I run it again: Recommendation…"
kind: answer
sources: [TH-0534]
topics: [flai, cli, go, code, all]
paths: [design/adrs/README.md, flai/internal/buildinfo/buildinfo.go, flai/internal/buildinfo/buildinfo_test.go, flai/cmd/version.go, flai/cmd/version_test.go, flai/internal/release/outdated.go, flai/internal/release/outdated_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/installed.go, flai/internal/mcpserver/installed_test.go, design/system/flai-cli.md, design/system/agent-narrative.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/update.md, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, design/issues/summary.md, design/adrs/0180-flai-task-done-refuses-an-older-flai-only-when-the-newest-flai-release-raised.md, flai/internal/buildinfo/task-done-rules.txt, docs/users/flai-reference.md, flai/cmd/task_done.go]
date: 2026-10-20
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0331 S-0445's close-out stops only at the LiteLLM gateway smoke check, the failure TH-0529 and TH-0530 report: recommend setting FLAI_TEST_LITELLM_HAIKU on the host, then I run it again: Recommendation…

The operator, alex, answered TH-0534 at 2026-10-20T21:21:36Z.

What was asked:

```text
**Recommended: set `FLAI_TEST_LITELLM_HAIKU` on the host to a model Claude Code 2.1.296 accepts, as TH-0529 recommends; then I run the close-out again.**
```

The answer:

```text
**Recommendation:** option (1), the same as on TH-0529 for S-0451 and on TH-0530 for S-0450. Set `FLAI_TEST_LITELLM_HAIKU` on the host to a model your LiteLLM proxy serves and Claude Code 2.1.296 accepts. Then all three agents run their close-outs again.

**One answer settles all three.** S-0445, S-0450 and S-0451 stop on the same `TestGatewayLiteLLM/unified` failure. It is outside each story, recorded as I-0170 and I-0171.

**Why not (2) or (3).** ADR-0093 lets me accept only when "the verifier's run passed at the story branch's head". A run with smoke failing or left out is not that, so S-0445 would wait in review for you. Option (2), skipping the check, is a change to what the close-out tests, and it would want a decision of its own.

A follow-up story from I-0170 should make the gateway check skip, rather than fail, on a model the installed Claude Code does not recognise, so the next Claude Code upgrade cannot stop every close-out again.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md § Decision
```

Source: TH-0534.
