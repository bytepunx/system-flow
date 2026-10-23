---
id: D-0377
title: A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself
kind: cause
sources: [I-0150]
topics: [orchestration, flai, cli, go, template, conventions, code, all]
paths: [flai/cmd/accept.go, flai/cmd/accept_batch_dirty.go, flai/cmd/accept_own_thread_test.go, flai/cmd/accept_conflict_test.go, flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/threads/threads.go, flai/internal/threads/threads_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, design/adrs/README.md, design/system/strategic-agents.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md, design/issues/summary.md, design/adrs/0189-the-orchestrator-s-acceptance-lands-on-a-dirty-main-as-the-merge-queue-s-does.md, flai/cmd/accept_orchestrator_dirty_test.go, docs/users/flai-reference.md]
date: 2026-10-23
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0377 A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself

I-0150 was closed on 2026-10-23.

The cause:

```text
A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself
```

The fix:

```text
S-0467, by ADR-0189, fixes each of the three causes.
```

Source: I-0150.
