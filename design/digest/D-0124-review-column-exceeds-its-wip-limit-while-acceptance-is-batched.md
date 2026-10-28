---
id: D-0124
title: Review column exceeds its WIP limit while acceptance is batched
kind: cause
sources: [I-0007]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/cmd/check.go, flai/cmd/check_stats_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/agents.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/work_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/serve/start_test.go, flai/internal/workitem/boardview.go, flai/internal/workitem/changes_test.go, design/adrs/0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md, design/adrs/README.md, design/conventions/work-management.md, design/issues/I-0007-review-queue-exceeds-limit.md, design/issues/summary.md, design/system/agent-narrative.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/system/workflow.md, docs/operators/index.md, docs/operators/runbooks/migrate.md, docs/users/flai-reference.md, docs/users/flai.md, docs/users/flaiover.md, template/CHANGELOG.md, template/root/design/conventions/work-management.md, template/template.yaml]
date: 2026-10-03
status: absorbed
absorbed_by: "design/adrs/0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md#decision"
asked_again: []
cited: []
---

# D-0124 Review column exceeds its WIP limit while acceptance is batched

I-0007 was closed on 2026-10-03.

The cause:

```text
The agent moves stories to review as they finish; the operator accepts them in batches. With a review limit of 3, the fourth finished story breaches the limit and `flai check --strict` fails, which blocks the commit-at-landing rule until someone accepts.
```

The fix:

```text
S-0243 (ADR-0073): a full review holds the pull, so flai serve, inbox, and wait_for_work start or offer no story while review is at or over its limit, and flai check --strict passes over review over its limit, so it no longer stops another story's close-out. Tests: TestReviewOverItsLimitDoesNotFailStrict, TestNoAgentStartsWhileReviewIsFull, TestWaitForWorkWaitsWhileReviewIsFull.
```

Source: I-0007.
