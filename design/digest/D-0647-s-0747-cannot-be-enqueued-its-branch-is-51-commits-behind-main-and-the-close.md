---
id: D-0647
title: "S-0747 cannot be enqueued: its branch is 51 commits behind main, and the close-out stops at sync: done"
kind: answer
sources: [TH-1228]
topics: [cli, dashboard, template, flai, go, code, all]
paths: [flai/internal/hostapi/writes_items.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/writes_describe.go, flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flai/internal/hostapi/components.go, flai/internal/hostapi/people.go, flai/internal/hostapi/people_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/flow.go, flai/internal/hostapi/components_test.go, flai/internal/hostapi/hostapi_test.go, design/system/flaiover-dashboard.md]
date: 2026-11-16
status: absorbed
absorbed_by: design/adrs/0310-a-close-out-whose-first-flai-verify-stops-at-sync-syncs-and-verifies-again.md
asked_again: []
cited: []
---

# D-0647 S-0747 cannot be enqueued: its branch is 51 commits behind main, and the close-out stops at sync: done

Merged into D-0645: when a story in review cannot be enqueued because its branch is behind main and its close-out stops at `sync`, the operator brings the branch up to main, and the orchestrator asks for the close-out again and enqueues the story once it passes. Here the operator answered "done" on TH-1228, having synced S-0747. Read D-0645 for the whole fact.
