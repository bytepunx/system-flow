---
id: I-0076
title: "flai check finds `wip.overlap` outside the story at close-out"
class: efficiency
status: closed
count: 16
first_reported: 2026-10-05T03:24:33Z
last_reported: 2026-10-07T08:30:24Z
updated: 2026-10-07T09:42:24Z
---

# I-0076 flai check finds `wip.overlap` outside the story at close-out

## Description
flai check finds `wip.overlap` outside the story at close-out

## Instances

### 2026-10-05T03:24:33Z
Story: S-0260.
flai check found outside the story:
`wip/kanban/stories/S-0253-an-acceptance-merge-committed-conflict-markers-to-a-design-document-on-main-and-nothing-caught-it.md`: S-0253 touches docs/users/flai.md, which S-0260 (in progress) also touches as docs/users/flai.md

### 2026-10-05T03:23:20Z
Story: S-0253.
flai check found outside the story:
`wip/kanban/stories/S-0253-an-acceptance-merge-committed-conflict-markers-to-a-design-document-on-main-and-nothing-caught-it.md`: S-0253 touches docs/users/flai.md, which S-0260 (in progress) also touches as docs/users/flai.md

### 2026-10-05T04:26:15Z
Story: S-0262.
flai check found outside the story:
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches design/issues, which S-0262 (in progress) also touches as design/issues
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches design/system/flai-cli.md, which S-0262 (in progress) also touches as design/system/flai-cli.md
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches docs/users/flai.md, which S-0262 (in progress) also touches as docs/users/flai.md

### 2026-10-05T04:30:02Z
Story: S-0262.
flai check found outside the story:
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches design/issues, which S-0262 (in progress) also touches as design/issues
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches design/system/flai-cli.md, which S-0262 (in progress) also touches as design/system/flai-cli.md
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches docs/users/flai.md, which S-0262 (in progress) also touches as docs/users/flai.md
`wip/kanban/stories/S-0262-flai-s-wip-markdown-lint-does-not-flag-a-space-inside-a-code-span-so-a-task-body-reached-main-and-failed-a-story-s-close-out.md`: S-0262 touches design/issues, which T-0877 (in progress) also touches as design/issues
`wip/kanban/stories/S-0262-flai-s-wip-markdown-lint-does-not-flag-a-space-inside-a-code-span-so-a-task-body-reached-main-and-failed-a-story-s-close-out.md`: S-0262 touches design/system/flai-cli.md, which T-0877 (in progress) also touches as design/system/flai-cli.md
`wip/kanban/stories/S-0262-flai-s-wip-markdown-lint-does-not-flag-a-space-inside-a-code-span-so-a-task-body-reached-main-and-failed-a-story-s-close-out.md`: S-0262 touches docs/users/flai.md, which T-0877 (in progress) also touches as docs/users/flai.md

### 2026-10-05T04:38:00Z
Story: S-0257.
flai check found outside the story:
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches design/issues, which S-0262 (in progress) also touches as design/issues
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches design/system/flai-cli.md, which S-0262 (in progress) also touches as design/system/flai-cli.md
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches docs/users/flai.md, which S-0262 (in progress) also touches as docs/users/flai.md

### 2026-10-06T20:18:53Z
Story: S-0296.
flai check found outside the story:
`wip/kanban/stories/S-0223-the-analyzer-runs-on-demand-or-on-a-schedule-and-writes-a-report-under-design-analysis.md`: S-0223 touches docs/operators/settings.md, which S-0296 (in progress) also touches as docs/operators/settings.md
`wip/kanban/stories/S-0223-the-analyzer-runs-on-demand-or-on-a-schedule-and-writes-a-report-under-design-analysis.md`: S-0223 touches flai/internal/serve, which S-0296 (in progress) also touches as flai/internal/serve/review_wait_test.go

### 2026-10-06T22:10:04Z
Story: S-0227.
flai check found outside the story:
`wip/kanban/stories/S-0227-the-analyzer-s-cost-is-recorded-on-the-issues-it-filed-and-the-stories-made-from-them.md`: S-0227 touches design/system/strategic-agents.md, which S-0229 (in progress) also touches as design/system/strategic-agents.md

### 2026-10-06T22:53:33Z
Story: S-0264.
flai check found outside the story:
`wip/kanban/stories/S-0299-a-story-s-sub-agent-that-writes-under-claude-blocks-its-layer-for-thirty-minutes-on-a-permission-thread-nobody-answers-and-the-start-prompt-does-not-warn-the-agent-beforehand.md`: S-0299 touches flai/cmd/guard.go, which S-0301 (in progress) also touches as flai/cmd
`wip/kanban/stories/S-0299-a-story-s-sub-agent-that-writes-under-claude-blocks-its-layer-for-thirty-minutes-on-a-permission-thread-nobody-answers-and-the-start-prompt-does-not-warn-the-agent-beforehand.md`: S-0299 touches flai/cmd/guard_test.go, which S-0301 (in progress) also touches as flai/cmd

### 2026-10-06T23:21:51Z
Story: S-0261.
flai check found outside the story:
`wip/kanban/stories/S-0261-the-mcp-prime-pack-for-a-story-is-larger-than-claude-code-s-tool-result-limit-so-the-agent-reads-it-back-from-a-saved-file.md`: S-0261 touches docs/operators/settings.md, which S-0301 (in progress) also touches as docs/operators/settings.md
`wip/kanban/stories/S-0261-the-mcp-prime-pack-for-a-story-is-larger-than-claude-code-s-tool-result-limit-so-the-agent-reads-it-back-from-a-saved-file.md`: S-0261 touches flai/internal/mcpserver, which S-0300 (in progress) also touches as flai/internal/mcpserver/plan.go

### 2026-10-07T00:07:32Z
Story: S-0251.
flai check found outside the story:
`wip/kanban/stories/S-0228-the-workflow-menu-has-orchestrator-and-analyzer-pages-showing-their-status-activity-log-and-runs.md`: S-0228 touches design/system/dashboard-host-channel.md, which S-0273 (in progress) also touches as design/system/dashboard-host-channel.md
`wip/kanban/stories/S-0228-the-workflow-menu-has-orchestrator-and-analyzer-pages-showing-their-status-activity-log-and-runs.md`: S-0228 touches flai/internal/hostapi, which S-0273 (in progress) also touches as flai/internal/hostapi/hostapi.go
`wip/kanban/stories/S-0228-the-workflow-menu-has-orchestrator-and-analyzer-pages-showing-their-status-activity-log-and-runs.md`: S-0228 touches flai/internal/hostapi, which S-0273 (in progress) also touches as flai/internal/hostapi/writes.go
`wip/kanban/stories/S-0228-the-workflow-menu-has-orchestrator-and-analyzer-pages-showing-their-status-activity-log-and-runs.md`: S-0228 touches flai/internal/hostapi, which S-0273 (in progress) also touches as flai/internal/hostapi/writes_test.go
`wip/kanban/stories/S-0273-flai-test-runs-the-project-s-test-and-lint-tiers-for-a-path-or-a-package-and-answers-pass-or-the-first-failures-as-findings.md`: S-0273 touches flai/internal/hostapi/hostapi.go, which T-0917 (in progress) also touches as flai/internal/hostapi
`wip/kanban/stories/S-0273-flai-test-runs-the-project-s-test-and-lint-tiers-for-a-path-or-a-package-and-answers-pass-or-the-first-failures-as-findings.md`: S-0273 touches flai/internal/hostapi/writes.go, which T-0917 (in progress) also touches as flai/internal/hostapi
`wip/kanban/stories/S-0273-flai-test-runs-the-project-s-test-and-lint-tiers-for-a-path-or-a-package-and-answers-pass-or-the-first-failures-as-findings.md`: S-0273 touches flai/internal/hostapi/writes_test.go, which T-0917 (in progress) also touches as flai/internal/hostapi

### 2026-10-07T00:38:11Z
Story: S-0273.
flai check found outside the story:
`wip/kanban/stories/S-0228-the-workflow-menu-has-orchestrator-and-analyzer-pages-showing-their-status-activity-log-and-runs.md`: S-0228 touches design/system/dashboard-host-channel.md, which S-0273 (in progress) also touches as design/system/dashboard-host-channel.md
`wip/kanban/stories/S-0228-the-workflow-menu-has-orchestrator-and-analyzer-pages-showing-their-status-activity-log-and-runs.md`: S-0228 touches flai/internal/hostapi, which S-0273 (in progress) also touches as flai/internal/hostapi/hostapi.go
`wip/kanban/stories/S-0228-the-workflow-menu-has-orchestrator-and-analyzer-pages-showing-their-status-activity-log-and-runs.md`: S-0228 touches flai/internal/hostapi, which S-0273 (in progress) also touches as flai/internal/hostapi/writes.go
`wip/kanban/stories/S-0228-the-workflow-menu-has-orchestrator-and-analyzer-pages-showing-their-status-activity-log-and-runs.md`: S-0228 touches flai/internal/hostapi, which S-0273 (in progress) also touches as flai/internal/hostapi/writes_test.go
`wip/kanban/stories/S-0228-the-workflow-menu-has-orchestrator-and-analyzer-pages-showing-their-status-activity-log-and-runs.md`: S-0228 touches flaiover/src/lib/server/agent.ts, which S-0273 (in progress) also touches as flaiover/src/lib/server/agent.ts
`wip/kanban/stories/S-0273-flai-test-runs-the-project-s-test-and-lint-tiers-for-a-path-or-a-package-and-answers-pass-or-the-first-failures-as-findings.md`: S-0273 touches design/system/dashboard-host-channel.md, which T-0938 (in progress) also touches as design/system/dashboard-host-channel.md

### 2026-10-07T01:04:51Z
Story: S-0272.
flai check found outside the story:
`wip/kanban/stories/S-0272-an-agent-with-an-open-question-ends-instead-of-waiting-flai-serve-restarts-it-on-the-answer-and-wait-for-events-keeps-a-timeout-only-for-an-agent-with-work-in-hand.md`: S-0272 touches flai/internal/harness/harness.go, which S-0286 (in progress) also touches as flai/internal/harness
`wip/kanban/stories/S-0272-an-agent-with-an-open-question-ends-instead-of-waiting-flai-serve-restarts-it-on-the-answer-and-wait-for-events-keeps-a-timeout-only-for-an-agent-with-work-in-hand.md`: S-0272 touches flai/internal/harness/harness_test.go, which S-0286 (in progress) also touches as flai/internal/harness
`wip/kanban/stories/S-0272-an-agent-with-an-open-question-ends-instead-of-waiting-flai-serve-restarts-it-on-the-answer-and-wait-for-events-keeps-a-timeout-only-for-an-agent-with-work-in-hand.md`: S-0272 touches flai/internal/serve/restart_test.go, which S-0286 (in progress) also touches as flai/internal/serve

### 2026-10-07T02:00:55Z
Story: S-0307.
flai check found outside the story:
`wip/kanban/stories/S-0269-one-command-closes-a-task-flai-task-done-commits-syncs-moves-logs-widens-touches-checks-and-answers-the-inbox.md`: S-0269 touches flai/internal/harness/harness.go, which S-0294 (in progress) also touches as flai/internal/harness/harness.go

### 2026-10-07T03:10:53Z
Story: S-0277.
flai check found outside the story:
`wip/kanban/stories/S-0269-one-command-closes-a-task-flai-task-done-commits-syncs-moves-logs-widens-touches-checks-and-answers-the-inbox.md`: S-0269 touches design/system/workflow.md, which S-0277 (in progress) also touches as design/system/workflow.md

### 2026-10-07T08:25:30Z
Story: S-0213.
flai check found outside the story:
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches design/system/flaiover-dashboard.md, which S-0214 (in progress) also touches as design/system/flaiover-dashboard.md
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches design/system/metrics.md, which S-0214 (in progress) also touches as design/system/metrics.md
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches design/system/metrics.md, which T-0924 (in progress) also touches as design/system/metrics.md
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches docs/users/flaiover.md, which S-0214 (in progress) also touches as docs/users/flaiover.md
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flai/cmd/check_stats_test.go, which S-0214 (in progress) also touches as flai/cmd/check_stats_test.go
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flai/cmd/stats.go, which S-0214 (in progress) also touches as flai/cmd/stats.go
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flai/internal/metrics, which S-0214 (in progress) also touches as flai/internal/metrics
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flaiover/src/lib/charts, which S-0214 (in progress) also touches as flaiover/src/lib/charts
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flaiover/src/lib/sitemenu.ts, which S-0214 (in progress) also touches as flaiover/src/lib/sitemenu.ts
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flaiover/src/lib/viz, which S-0214 (in progress) also touches as flaiover/src/lib/viz
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flaiover/src/routes/charts, which S-0214 (in progress) also touches as flaiover/src/routes/charts

### 2026-10-07T08:30:24Z
Story: S-0213.
flai check found outside the story:
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches design/system/flaiover-dashboard.md, which S-0214 (in progress) also touches as design/system/flaiover-dashboard.md
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches design/system/metrics.md, which S-0214 (in progress) also touches as design/system/metrics.md
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches docs/users/flaiover.md, which S-0214 (in progress) also touches as docs/users/flaiover.md
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flai/cmd/check_stats_test.go, which S-0214 (in progress) also touches as flai/cmd/check_stats_test.go
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flai/cmd/stats.go, which S-0214 (in progress) also touches as flai/cmd/stats.go
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flai/internal/metrics, which S-0214 (in progress) also touches as flai/internal/metrics
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flai/internal/metrics, which T-0929 (in progress) also touches as flai/internal/metrics
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flaiover/src/lib/charts, which S-0214 (in progress) also touches as flaiover/src/lib/charts
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flaiover/src/lib/sitemenu.ts, which S-0214 (in progress) also touches as flaiover/src/lib/sitemenu.ts
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flaiover/src/lib/viz, which S-0214 (in progress) also touches as flaiover/src/lib/viz
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flaiover/src/lib/viz, which T-0936 (in progress) also touches as flaiover/src/lib/viz
`wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md`: S-0213 touches flaiover/src/routes/charts, which S-0214 (in progress) also touches as flaiover/src/routes/charts

## Remediation

Story S-0279 remediates this issue, created from it at 2026-10-05T04:40:48Z.
Closed 2026-10-07T09:42:24Z: S-0279 fixed it as ADR-0115 decides: a close-out's scoped flai check keeps only a wip.overlap that names its story, as a note, and --record-issues records no wip.overlap (T-1145); wip.overlap compares the claims of the stories in progress, as the pull hold does, one finding per pair of stories (T-1143).
