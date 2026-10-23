---
id: D-0194
title: A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand
kind: cause
sources: [I-0093]
topics: [cli, conventions, flai, go, template, code, all]
paths: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, flai/internal/hostapi/writes.go, docs/users/flai-reference.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md, design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai.md, docs/operators/settings.md, design/adrs, design/adrs/README.md, ".claude/settings.json", template/root/.claude/settings.json, design/issues/I-0093-a-story-s-sub-agent-that-writes-under-claude-blocks-its-layer-for-thirty-minutes-on-a-permission-thread-nobody-answers-and-the-start-prompt-does-not-warn-the-agent-beforehand.md, design/issues/summary.md, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, docs/operators/index.md]
date: 2026-10-06
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0194 A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand

I-0093 was closed on 2026-10-06.

The cause:

```text
A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand
```

The fix:

```text
S-0299: flai guard refuses a story's sub-agent an Edit, MultiEdit, Write, or NotebookEdit of a file in a .claude/ folder at once while auto-approve is off (T-1035), the settings now run the guard before a story session's file writes (T-1053), and the story agent's start prompt says before any layer that such a write is its own, made after the layer or handed over with cp commands on a thread (T-1036); ADR-0102.
```

Source: I-0093.
