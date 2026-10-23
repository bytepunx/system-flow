---
id: D-0198
title: "S-0299: accept ADR-0102, flai guard refuses a story's sub-agent a .claude/ write while auto-approve is off: accept"
kind: answer
sources: [TH-0211]
topics: [cli, conventions, flai, go, template, code, all]
paths: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, flai/internal/hostapi/writes.go, docs/users/flai-reference.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md, design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai.md, docs/operators/settings.md, design/adrs, design/adrs/README.md, ".claude/settings.json", template/root/.claude/settings.json, design/issues/I-0093-a-story-s-sub-agent-that-writes-under-claude-blocks-its-layer-for-thirty-minutes-on-a-permission-thread-nobody-answers-and-the-start-prompt-does-not-warn-the-agent-beforehand.md, design/issues/summary.md, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, docs/operators/index.md]
date: 2026-10-06
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0198 S-0299: accept ADR-0102, flai guard refuses a story's sub-agent a .claude/ write while auto-approve is off: accept

The operator, alex, answered TH-0211 at 2026-10-06T22:53:59Z.

What was asked:

```text
ADR-0102 (proposed, on story/S-0299: `design/adrs/0102-while-auto-approve-is-off-flai-guard-refuses-a-story-s-sub-agent-a-write-to-a.md`) refines ADR-0060 and ADR-0086. While `auto-approve` is off, `flai guard` refuses a story's sub-agent an `Edit`, `MultiEdit`, `Write`, or `NotebookEdit` of a file in a `.claude/` folder at once, so that the write never holds a layer on a `permission_prompt` thread (I-0093). The story's agent makes such writes itself, after the layer.
```

The answer:

```text
accept
```

Source: TH-0211.
