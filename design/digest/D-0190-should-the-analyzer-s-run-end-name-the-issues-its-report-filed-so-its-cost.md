---
id: D-0190
title: "Should the analyzer's run end name the issues its report filed, so its cost reaches them without an activity_log call: go with option 1"
kind: answer
sources: [TH-0200]
topics: [analysis, planning, orchestration, flai, cli, go, code, all]
paths: [flai/internal/usage, flai/internal/serve, flai/internal/workitem, design/system/metrics.md, flai/internal/issues, flai/internal/metrics, design/system/continuous-improvement.md, design/adrs, flai/cmd/issue.go, flai/cmd/issue_test.go, flai/internal/mcpserver/issues.go, flai/internal/mcpserver/issues_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0058-flai-guard-refuses-a-sub-agent-s-shell-command-whose-heredoc-text-reads-like-a-flai-write.md, design/issues/summary.md, flai/internal/statsread/statsread.go, flai/cmd/activity.go, flai/cmd/activity_test.go, design/system/strategic-agents.md, design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/I-0096-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md, flai/cmd/check.go, flai/cmd/check_test.go]
date: 2026-10-06
status: absorbed
absorbed_by: "design/adrs/0100-an-analyzer-activity-s-usage-is-charged-evenly-to-the-issues-it-names-under.md#decision"
asked_again: []
cited: []
---

# D-0190 Should the analyzer's run end name the issues its report filed, so its cost reaches them without an activity_log call: go with option 1

The operator, alex, answered TH-0200 at 2026-10-06T21:49:44Z.

What was asked:

```text
Criterion 2 charges an analyzer activity to the issues its `activity_log` call names. But the analyzer runs once and ends: flai serve logs its run end with no items (`flai/internal/serve/analyze.go`, `logRunEndSaying(..., nil)`), and neither its prompt nor `analyzer.md` asks it to call `activity_log`. As written, nearly every analyzer run would go to the project total, and no issue would carry a cost.
```

The answer:

```text
go with option 1
```

Source: TH-0200.
