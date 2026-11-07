---
id: I-0293
title: Design documents describe a /conventions page, a /streams page, a mounted-repository readiness check, a flai OTEL span, a dashboard scan of wip/kanban, and a single lockfile that are not there
class: impression
status: open
count: 1
first_reported: 2026-11-07T08:17:36Z
last_reported: 2026-11-07T08:17:36Z
updated: 2026-11-07T23:40:21Z
usage:
  strategic:
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
---

# I-0293 Design documents describe a /conventions page, a /streams page, a mounted-repository readiness check, a flai OTEL span, a dashboard scan of wip/kanban, and a single lockfile that are not there

## Description
Design documents describe a /conventions page, a /streams page, a mounted-repository readiness check, a flai OTEL span, a dashboard scan of wip/kanban, and a single lockfile that are not there

## Instances

### 2026-11-07T08:17:36Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: each sentence was true when written; documentation.md asks that behaviour and its document change in one commit.

## Impact
flaiover-dashboard.md lines 126 and 129 vs flaiover/src/routes (no conventions or streams); telemetry.md Project additions vs flaiover/src/routes/_ready/+server.ts:17-30 and flai/go.mod (no go.opentelemetry.io); metrics.md line 605 vs flaiover/src/lib/server/stats.ts:1-2; code-quality.md line 85 vs flai/tests/integration/; node.md line 20 vs tracked flaiover/package-lock.json; dashboard-host-channel.md line 10 'Nothing here is built' vs S-0082 done

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0626 remediates this issue, created from it at 2026-11-07T23:40:21Z.
