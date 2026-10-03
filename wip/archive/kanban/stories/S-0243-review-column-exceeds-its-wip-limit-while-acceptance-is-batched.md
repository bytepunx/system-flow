---
id: S-0243
type: story
nature: improvement
title: Review column exceeds its WIP limit while acceptance is batched
status: done
owner: arobson
created: 2026-10-03T01:46:35Z
updated: 2026-10-03T03:31:01Z
transitions:
  - to: ready
    at: 2026-10-03T01:52:11Z
    by: alex
  - to: in-progress
    at: 2026-10-03T03:04:59Z
    by: agent-S-0243
  - to: review
    at: 2026-10-03T03:30:07Z
    by: agent-S-0243
  - to: done
    at: 2026-10-03T03:31:01Z
    by: alex
tags: [flai]
touches: [flai/cmd/check.go, flai/cmd/check_stats_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/agents.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/work_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/serve/start_test.go, flai/internal/workitem/boardview.go, flai/internal/workitem/changes_test.go, design/adrs/0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md, design/adrs/README.md, design/conventions/work-management.md, design/issues/I-0007-review-queue-exceeds-limit.md, design/issues/summary.md, design/system/agent-narrative.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/system/workflow.md, docs/operators/index.md, docs/operators/runbooks/migrate.md, docs/users/flai-reference.md, docs/users/flai.md, docs/users/flaiover.md, template/CHANGELOG.md, template/root/design/conventions/work-management.md, template/template.yaml]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2800
  models:
    - model: claude-haiku-4-5-20251001
      input: 226
      output: 6213
      cache_read: 1262015
      cache_write: 71252
      cost: 0.2466
    - model: claude-opus-5-5
      input: 350
      output: 104078
      cache_read: 13281357
      cache_write: 378917
      cost: 7.0865
    - model: claude-sonnet-5-5
      input: 12
      output: 3133
      cache_read: 100544
      cache_write: 30263
      cost: 0.1271
---
# S-0243 Review column exceeds its WIP limit while acceptance is batched

## Goal

This story remediates [I-0007](../../../design/issues/I-0007-review-queue-exceeds-limit.md), "Review column exceeds its WIP limit while acceptance is batched". The issue recommends this solution:

Either raise the review limit in `wip/kanban/board.md` to match the operator's acceptance cadence, or accept stories before pulling the next one. Operator's call; see the S-0026 open questions.

## Acceptance criteria
- [x] The cause I-0007 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0007 is closed with `flai issue close I-0007 --reason` saying what fixed it

## Tasks
- T-0736 flai check --strict passes over the review column over its limit
- T-0737 The pull respects the review limit
- T-0738 An ADR, the design, and the docs say a full review holds the pull, and I-0007 is closed

## Notes
