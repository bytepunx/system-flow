---
id: S-0852
type: story
nature: improvement
title: The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right
status: backlog
owner: alex
created: 2026-11-27T22:55:19Z
updated: 2026-11-27T23:57:00Z
transitions: []
tags: [flai, template]
topics: [orchestration, review, git]
touches: [design/adrs, design/adrs/README.md, flai/internal/manifest/manifest_roles.go, flai/internal/manifest/settings.go, flai/internal/manifest/manifest_test.go, flai/cmd/accept.go, flai/cmd/accept_undone_test.go, flai/cmd/accept_queue_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/mcpserver/accept.go, flai/internal/mcpserver/accept_test.go, flai/internal/guard/guard_rules.go, flai/internal/guard/orchestrate_test.go, flai/internal/mcpserver/queue.go, flai/internal/mcpserver/queue_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/operators/index.md, template/CHANGELOG.md, design/issues/I-0415-the-orchestrator-s-accept-takes-no-keep-removal-so-a-story-that-rewrites-a-line.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: analyzer
      seconds: 85
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 55
          output: 5532
          cache_read: 391176
          cache_write: 29513
          cost: 0.9652
        - model: claude-opus-5-5
          input: 29
          output: 4664
          cache_read: 1264948
          cache_write: 28114
          cost: 0.487
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 148h
    by: flai
    at: 2026-11-27T22:55:19Z
  value: 22200
  by: planner-S-0852
  at: 2026-11-27T23:57:00Z
forecast:
  duration: 1h13m
  delivery: 2026-11-28T06:55:00Z
  basis: "flai forecast: median 109 s per unit of size over 187 done improvement stories on claude-opus-5-5 in the large band, times size 40 (2 criteria, 38 touches), 12th in the pull order with an in-progress limit of 3"
  by: planner-S-0852
  at: 2026-11-27T23:57:00Z
---
# S-0852 The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right

## Goal

This story remediates [I-0415](../../../design/issues/I-0415-the-orchestrator-s-accept-takes-no-keep-removal-so-a-story-that-rewrites-a-line.md), "The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0415 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0415 is closed with `flai issue close I-0415 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0415: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0415. time_lost_per_cycle 148h carried over from I-0415's Impact section. Its Impact time_lost_per_cycle was taken rather than the 148h20m derived from its cost and count.

Strategic usage carried over from I-0415: analyzer 1m25s, 1.4522 USD.

### Planning

- `design/adrs` is kept as a folder touch: the story writes a new ADR refining ADR-0155, ADR-0093, and ADR-0278, whose number and slug are taken when it is written.
