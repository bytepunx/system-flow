---
id: S-0199
type: story
nature: feature
title: "Work items carry planning data: draft, cost of delay inputs and value, and a forecast with who set it"
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:11Z
updated: 2026-10-02T11:54:11Z
transitions: []
tags: [flai, dashboard]
touches: [flai/internal/workitem, flai/internal/itemedit, flai/internal/itemnew, flai/internal/check, flai/internal/mcpserver, flai/internal/hostapi, flaiover/src, design/system/work-hierarchy.md, template/]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0199 Work items carry planning data: draft, cost of delay inputs and value, and a forecast with who set it

## Goal

E-0016's planner, orchestrator, and analyzer need data the work items do not carry: whether a story is a draft an agent wrote, what delaying an item costs, and a forecast of how long it will take and when it lands. Nothing records who set a value, so a human's estimate cannot be told from a planner's. `estimate` exists (a Go duration) and nothing writes or reads it.

## Acceptance criteria
- [ ] Stories carry `draft: true` when an agent wrote them (the planner, S-0198's story from an issue); `flai story new --draft`, MCP `item_new` `draft`, and `flai edit --draft/--no-draft` set and clear it; `flai move` and `item_move` refuse to move a draft story to `ready` with "finalize it first", unless the caller is allowed to finalize (the orchestrator's permission, or `--yes` for the operator)
- [ ] Epics and stories carry `cost_of_delay`: `inputs` the operator supplies (`revenue_per_week`, `penalty_per_week`, `time_lost_per_cycle` as a duration, each optional, in the project's currency), `value` per week the planner derives, `by` and `at`; the manifest can set the hourly rate the planner uses for time lost (`planning.hour_rate`) and the cycle length (`planning.cycle`); a story without its own inputs inherits its epic's value share as the planner apportions it
- [ ] Stories carry `forecast`: `duration` (agent wall clock), `delivery` (a timestamp), `basis` (what it was computed from, in a sentence), `by`, `at`; `estimate` stays as the human's figure, and the metrics compare both to the actual
- [ ] Every field is set through flai (`flai edit`, MCP `item_edit`, hostapi), never by a command that owns state; `flai check` validates shapes, currencies, and durations, and warns on a `draft` story in `ready` or later
- [ ] `front-matter-fields.txt` lists the new keys and the release raises `flai.minimum` (S-0181), so an older host flai stops before reading the project
- [ ] An ADR records the schema; `design/system/work-hierarchy.md`, the template's conventions where they describe items, and the user guide describe the fields
- [ ] Tests cover each field's set, clear, validation, the refused move of a draft, and an older-flai read

## Tasks

## Notes

Decided by the designer on 2026-10-02 for E-0016. The orchestrator's permission to finalize drafts is defined in the orchestrator settings story.
