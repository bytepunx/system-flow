---
id: S-0486
type: story
nature: remediation
title: criteria_tick is refused on another open story's narrative.state finding, which the story's own scoped check leaves out
status: backlog
owner: alex
created: 2026-10-23T17:43:07Z
updated: 2026-10-23T17:56:28Z
transitions: []
tags: [flai]
topics: [cli]
touches: [design/adrs, design/adrs/README.md, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/criteria_test.go, flai/internal/mcpserver/criteria_test.go, flai/cmd/edit.go, flai/cmd/criteria.go, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0182-criteria-tick-is-refused-on-another-open-story-s-narrative-state-finding-which-the-story-s-own-scoped-check-leaves-out.md, design/issues/summary.md]
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
    - kind: orchestrator
      seconds: 512
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 210
          output: 3714
          cache_read: 35996176
          cache_write: 42148
          cost: 7.9266
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4m
    by: flai
    at: 2026-10-23T17:43:07Z
  value: 10
  by: planner-S-0486
  at: 2026-10-23T17:52:54Z
forecast:
  duration: 2h
  delivery: 2026-10-24T13:14:00Z
  basis: "Its own forecast of 2h; 27th in the pull order with an in-progress limit of 5, behind S-0429, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0430, S-0431, S-0432, S-0433, S-0434 and S-0435."
  by: flai
  at: 2026-10-23T17:56:28Z
---
# S-0486 criteria_tick is refused on another open story's narrative.state finding, which the story's own scoped check leaves out

## Goal

This story remediates [I-0182](../../../design/issues/I-0182-criteria-tick-is-refused-on-another-open-story-s-narrative-state-finding-which-the-story-s-own-scoped-check-leaves-out.md), "criteria_tick is refused on another open story's narrative.state finding, which the story's own scoped check leaves out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0182 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0182 is closed with `flai issue close I-0182 --reason` saying what fixed it

## Tasks
- T-2074 An ADR settles that a checked edit of a work item refuses only a finding the edit introduces inside its story's scope, and leaves out what a check scoped to the story leaves out
- T-2075 itemedit.Apply refuses a new finding only inside the edited item's story's scope, so criteria_tick is no longer refused on another story's narrative.state or an untracked task file, with tests that reproduce I-0182's two instances
- T-2076 The design and the user docs say a checked edit refuses only what it introduces inside its story's scope, and warns on the rest
- T-2077 I-0182 is closed with the reason naming the scoped check a checked edit now applies

## Notes

Cost of delay inputs set by flai from I-0182. time_lost_per_cycle 4m: 2m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-21T12:14:14Z, 2.2 days before this story; under one cycle counts as one).

### Planning

- **Plan thread.** TH-0618 names the cause, the proposed fix, the tasks, their layers, and the assumptions: `itemedit.Apply` runs an unscoped before-and-after `flai check` and refuses on any finding that appears between the two runs, so another writer's finding refuses a tick; the fix scopes the second run to the edited item's story with `check.ScopeToStories`, as `flai task done` and `flai verify` do; T-2074's ADR confirms or amends it.
- **Layers.** T-2074 (the ADR) first; T-2075 (the change in `Apply` with tests) after it; T-2076 (design and user docs) after T-2075; T-2077 (closing I-0182 and ticking the criteria) last. One task per layer.
- **Touches.** 12, one folder among them, every one on a task.
  - **Declared:** none; the story was made from I-0182 with no touches, and `flai touches suggest S-0486` refused for want of a path to start from.
  - **From co-change** (`flai touches suggest S-0486 flai/internal/itemedit/itemedit.go flai/internal/check/scope.go flai/internal/mcpserver/criteria.go`, 29 commits): `docs/users/flai-reference.md` (38%), `flai/cmd/edit.go` (34%), and `docs/users/flai.md` (14%) kept; `flai/internal/check/scope_test.go`, `flai/cmd/check.go`, `check_test.go`, `flai/cmd/items.go`, and the rest left out, since the fix reads the scope and does not change it.
  - **From the design:** `design/adrs/README.md` for the ADR's index row; `design/system/flai-cli.md`, whose Commands table rows for `flai edit` and `flai criteria` describe the checked edit; `design/issues/I-0182-….md` and `design/issues/summary.md`, which `flai issue close` rewrites.
  - **From the code layout:** `flai/internal/itemedit/itemedit.go`, where `Apply` runs the two checks; `flai/internal/itemedit/criteria_test.go` and `flai/internal/mcpserver/criteria_test.go`, the tests of the tick on both ways in; `flai/cmd/criteria.go`, whose long help `make flai-reference` renders. `flai/internal/mcpserver/criteria.go` and `flai/internal/check/scope.go` were read and left out: neither changes under the proposed fix.
  - **Folder touch kept:** `design/adrs`, for the ADR T-2074 adds, whose number `flai adr new` allocates when it runs; `flai task done` widens the touches to the file. While the story is in progress it holds a ready story that touches an ADR.
  - **Overlaps:** `itemedit.go`, `flai-cli.md`, `flai.md`, `flai-reference.md`, and `design/adrs/README.md` are touches of S-0468 (in review, which holds nothing) and the three documents of S-0429 (in progress, which holds this story in ready until it is accepted or shares them).
- **Forecast.** flai's 22m at size 14 (2 criteria, 12 touches), 91 s a unit over 35 done remediation stories on claude-opus-5-5 in the large band, raised to 2h: an ADR that waits for the operator's agreement, tests that must bring a finding between two check runs in a writer `flai edit`, `item_edit`, and the dashboard share, and docs in three places with a generated reference are work the touch count does not weigh. Delivery flai's 2026-10-24T11:20Z, 27th in the pull order behind S-0429 and the E-0016 and digest stories, moved to 2026-10-24T13:00Z by the raise.
- **Cost of delay.** 10.00 USD a week, flai's figure from the input flai set from I-0182: 4m lost per 168h cycle at 150 USD an hour, 1.00 cycles a week. Unchanged: two instances in two days, each a retry of a few seconds to a minute once the other file was finished.
- **Tags and topics.** Tag `flai`, the component every touch is in; topic `cli`, which the command docs reach.
