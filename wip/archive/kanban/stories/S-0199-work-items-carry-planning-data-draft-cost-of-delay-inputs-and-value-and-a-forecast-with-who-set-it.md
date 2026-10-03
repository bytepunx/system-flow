---
id: S-0199
type: story
nature: feature
title: "Work items carry planning data: draft, cost of delay inputs and value, and a forecast with who set it"
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:11Z
updated: 2026-10-03T07:06:10Z
transitions:
  - to: ready
    at: 2026-10-03T05:33:48Z
    by: alex
  - to: in-progress
    at: 2026-10-03T05:34:17Z
    by: agent-S-0199
  - to: review
    at: 2026-10-03T07:00:09Z
    by: agent-S-0199
  - to: done
    at: 2026-10-03T07:06:10Z
    by: alex
tags: [flai, dashboard]
touches: [flai/internal/workitem, flai/internal/itemedit, flai/internal/itemnew, flai/internal/check, flai/internal/mcpserver, flai/internal/hostapi, flai/internal/manifest, flai/internal/metrics/cascade_test.go, flai/cmd, flaiover/src, docs/operators/settings.md, design/adrs, design/system/work-hierarchy.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, template/template.yaml, docs/users/flai-reference.md, design/issues/summary.md, design/system/flaiover-dashboard.md, design/issues/I-0061-flai-check-splits-an-item-s-front-matter-errors-on-so-a-message-that-holds-one-becomes-two-findings-the-second-on-line-1.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5186
  models:
    - model: claude-haiku-4-5-20251001
      input: 386
      output: 9063
      cache_read: 2585583
      cache_write: 87703
      cost: 0.4139
    - model: claude-opus-5-5
      input: 602
      output: 222159
      cache_read: 35920872
      cache_write: 906189
      cost: 16.9748
    - model: claude-sonnet-5
      input: 136
      output: 38939
      cache_read: 4014068
      cache_write: 320478
      cost: 1.9937
---
# S-0199 Work items carry planning data: draft, cost of delay inputs and value, and a forecast with who set it

## Goal

E-0016's planner, orchestrator, and analyzer need data the work items do not carry: whether a story is a draft an agent wrote, what delaying an item costs, and a forecast of how long it will take and when it lands. Nothing records who set a value, so a human's estimate cannot be told from a planner's. `estimate` exists (a Go duration) and nothing writes or reads it.

## Acceptance criteria
- [x] Stories carry `draft: true` when an agent wrote them (the planner, S-0198's story from an issue); `flai story new --draft`, MCP `item_new` `draft`, and `flai edit --draft/--no-draft` set and clear it; `flai move` and `item_move` refuse to move a draft story to `ready` with "finalize it first", unless the caller is allowed to finalize (the orchestrator's permission, or `--yes` for the operator)
- [x] Epics and stories carry `cost_of_delay`: `inputs` the operator supplies (`revenue_per_week`, `penalty_per_week`, `time_lost_per_cycle` as a duration, each optional, in the project's currency), `value` per week the planner derives, `by` and `at`; the manifest can set the hourly rate the planner uses for time lost (`planning.hour_rate`) and the cycle length (`planning.cycle`); a story without its own inputs inherits its epic's value share as the planner apportions it
- [x] Stories carry `forecast`: `duration` (agent wall clock), `delivery` (a timestamp), `basis` (what it was computed from, in a sentence), `by`, `at`; `estimate` stays as the human's figure, and the metrics compare both to the actual
- [x] Every field is set through flai (`flai edit`, MCP `item_edit`, hostapi), never by a command that owns state; `flai check` validates shapes, currencies, and durations, and warns on a `draft` story in `ready` or later
- [x] `front-matter-fields.txt` lists the new keys and the release raises `flai.minimum` (S-0181), so an older host flai stops before reading the project
- [x] An ADR records the schema; `design/system/work-hierarchy.md`, the template's conventions where they describe items, and the user guide describe the fields
- [x] Tests cover each field's set, clear, validation, the refused move of a draft, and an older-flai read

## Tasks
- T-0739 Items carry draft, cost_of_delay, and forecast, the manifest carries planning, and a move refuses a draft story to ready
- T-0740 flai check validates the manifest's planning and warns on a draft story in ready or later
- T-0741 flai edit, flai story new --draft, and flai move --yes set and clear the planning fields
- T-0742 MCP item_new, item_edit, item_move and hostapi item.new, item.edit, item.move carry the planning fields
- T-0743 The dashboard's item type carries the planning fields and the story page shows cost of delay and forecast
- T-0744 An ADR records the planning fields, and the design, the template's conventions, and the user guide describe them

## Notes

Decided by the designer on 2026-10-02 for E-0016. The orchestrator's permission to finalize drafts is defined in the orchestrator settings story.

- ADR-0074 records the schema. A story may carry a cost of delay `value` with no `inputs`, its share of the epic's value. The planner apportions that share (S-0210); this story gives it a place.
- No agent finalizes a draft yet: MCP `item_move` refuses a draft to ready, and `item_edit` refuses `draft: false`. `workitem.MoveOptions.Finalize` is where S-0218's permission plugs in. The operator finalizes with `flai move <story> ready --yes` or `flai edit --no-draft`, and the dashboard will through hostapi `item.move` `finalize` (S-0201).
- `planning.currency` (ISO 4217, default USD) is the currency the criterion calls the project's; `flai check` reports a bad one as `manifest.planning`.
- `flai.minimum` rises when the flai release carrying this story is published: `release.RaiseMinimum` sees `front-matter-fields.txt` changed. The host's flai must be upgraded then.
- The metrics comparing `estimate` and `forecast` with the actual (`estimate_error`, `forecast_error`, `delivery_error`) are S-0205's, as the designer chose on TH-0080. This story records the split between the two figures in ADR-0074 and gives S-0205 the fields.
