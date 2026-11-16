---
id: D-0646
title: "S-0742 cannot be enqueued: its branch is 51 commits behind main, and the close-out stops at sync: done"
kind: answer
sources: [TH-1227]
topics: [cli, orchestration, planning, template, flai, go, code, all]
paths: [flai/internal/serve/agents.go, flai/internal/serve/agents_launch.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_activity.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/serve/stop.go, flai/internal/serve/agents_test.go, flai/internal/serve/restart_test.go, flai/internal/serve/readyempty.go, flai/internal/serve/readyempty_test.go, flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/role.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/usage.go, flai/internal/serve/confirm.go, flai/internal/serve/remind.go, flai/internal/serve/remind_test.go, flai/internal/serve/commit.go, flai/internal/serve/claudecheck.go, flai/internal/serve/components.go, flai/internal/serve/flow.go, flai/internal/serve/kinds.go, flai/internal/serve/kinds_test.go, flai/internal/serve/serve.go, flai/internal/serve/role_events.go, flai/internal/serve/imagejob.go, design/system/workflow.md]
date: 2026-11-16
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0646 S-0742 cannot be enqueued: its branch is 51 commits behind main, and the close-out stops at sync: done

Merged into D-0645: when a story in review cannot be enqueued because its branch is behind main and its close-out stops at `sync`, the operator brings the branch up to main, and the orchestrator asks for the close-out again and enqueues the story once it passes. Here the operator answered "done" on TH-1227, having synced S-0742. Read D-0645 for the whole fact.
