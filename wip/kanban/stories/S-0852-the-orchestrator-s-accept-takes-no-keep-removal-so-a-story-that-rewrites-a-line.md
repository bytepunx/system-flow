---
id: S-0852
type: story
nature: improvement
title: The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right
status: in-progress
owner: alex
created: 2026-11-27T22:55:19Z
updated: 2026-11-28T01:17:12Z
transitions:
  - to: ready
    at: 2026-11-28T00:02:04Z
    by: claude-concierge
  - to: in-progress
    at: 2026-11-28T01:17:12Z
    by: agent-S-0852
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
    - kind: planner
      seconds: 12
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 17
          output: 4073
          cache_read: 924285
          cache_write: 26770
          cost: 0.4594
    - kind: orchestrator
      seconds: 14
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 12
          cache_read: 309097
          cache_write: 6502
          cost: 0.1005
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
  delivery: 2026-11-28T03:07:00Z
  basis: "Its own forecast of 1h13m; 2nd in the pull order with an in-progress limit of 3, behind S-0835, S-0837, S-0844 and S-0814."
  by: flai
  at: 2026-11-28T00:26:08Z
finalized:
  by: alex
  at: 2026-11-28T00:01:52Z
---
# S-0852 The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right

## Goal

This story remediates [I-0415](../../../design/issues/I-0415-the-orchestrator-s-accept-takes-no-keep-removal-so-a-story-that-rewrites-a-line.md), "The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0415 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0415 is closed with `flai issue close I-0415 --reason` saying what fixed it

## Tasks
- T-3618 Propose in an ADR refining ADR-0155 how the orchestrator keeps a recorded removal of main's lines, behind a permission off by default
- T-3619 The manifest reads orchestration.permissions.keep_removals, off by default and needing accept_reviews, and the settings index lists it
- T-3620 flai accept --by orchestrator keeps a recorded removal with --keep-removal and --reason while keep_removals is on and the narrative's Decisions name the path
- T-3621 The MCP tool accept takes keep_removals, each a path and a reason, and passes them to flai accept --by orchestrator on the host
- T-3622 queue_add records the orchestrator's keep_removals with the entry, the queue leaves a kept path out of undone_main, and its landing passes the keeps
- T-3624 The orchestrator's role file and the strategic-agents convention say when it keeps a removal and how it logs it
- T-3626 The design, the user and operator guides, and the template changelog describe the orchestrator's keep of a removal
- T-3629 Close I-0415 as fixed by the orchestrator's keep of a removal the story's narrative records as a rewrite

## Notes

The analyzer found I-0415: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0415. time_lost_per_cycle 148h carried over from I-0415's Impact section. Its Impact time_lost_per_cycle was taken rather than the 148h20m derived from its cost and count.

Strategic usage carried over from I-0415: analyzer 1m25s, 1.4522 USD.

### Planning

The proposal the goal asks for is the planner's recommendation, written into T-3618: a new permission, `orchestration.permissions.keep_removals`, off by default and needing `accept_reviews`. With it on, the orchestrator keeps a recorded removal through `flai accept --by orchestrator --keep-removal`, the MCP tool `accept`, and `queue_add`. It may keep only a path the story's narrative names under `## Decisions`, with a reason. The story's agent confirms or changes the proposal in that ADR, and the plan's thread asks the operator first.

Touches. The story declared none, and `flai touches suggest` needs a path to start from, so it ran from `flai/cmd/accept.go`, `flai/internal/mcpserver/accept.go`, `flai/internal/mergequeue/eligible.go`, and `flai/internal/manifest/manifest_roles.go`.

- Design: ADR-0155's refusal of a keep with `--by orchestrator` in `flai/cmd/accept.go`. ADR-0278's tool `accept` in `flai/internal/mcpserver/accept.go`. The queue's `undone_main` in `flai/internal/mergequeue/eligible.go`, its entry in `state.go`, and its enqueue in `flai/internal/mcpserver/queue.go` and `flai/cmd/queue.go`. The permission in `flai/internal/manifest/manifest_roles.go` and `settings.go`. The orchestrator's accept rule in `flai/internal/guard/guard_rules.go`.
- Co-change: `flai/internal/preview/accept.go` (25%), `docs/users/flai-reference.md` (21%), `design/system/flai-cli.md` (18%), `docs/users/flai.md` (18%), `docs/operators/settings.md` (9%), `docs/operators/index.md` (8%), `design/issues/summary.md` (8%), `design/system/workflow.md` (7%), `flai/internal/mergequeue/eligible_test.go` and `flai/cmd/accept_queue_test.go`.
- Layout: the tests beside each changed file. The orchestrator's role file and its template and baseline copies. `design/conventions/strategic-agents.md` and its template copy. `design/system/strategic-agents.md`, `design/system/project-manifest.md`, `template/CHANGELOG.md`, `design/adrs/README.md`, and I-0415's file.
- Left out of the co-change list: `flai/cmd/branch.go`, `flai/internal/preview/orchestrator.go`, `flai/internal/hostapi/writes_test.go`, and the flaiover files. The dashboard's `accept.run` stays the operator's, and the keep reaches flai only through `flai accept` and the queue.
- `design/adrs` is kept as a folder touch: the story writes a new ADR refining ADR-0155, ADR-0093, and ADR-0278, whose number is taken when it is written. T-3618 names 0322 as a guess.

Forecast: 1h13m, delivery 2026-11-28T06:55:00Z, as `flai forecast` gives it. Its basis is 187 done improvement stories in the large band at size 40. Eight tasks over three code paths fit that size, and no measure argues for another figure. It leaves out any wait for the operator's answer on the plan's thread, which is not agent time.

Cost of delay: 22200 USD a week, as `flai cod` works it out from the operator's input of 148h lost per 168h cycle at 150 USD an hour. Kept as given. The 148h is S-0811's one wait in review, so it is the cost of one instance per cycle. The input is the operator's to change.
