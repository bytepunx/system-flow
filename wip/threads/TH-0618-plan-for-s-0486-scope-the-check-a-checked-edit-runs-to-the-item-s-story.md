---
id: TH-0618
title: "Plan for S-0486: scope the check a checked edit runs to the item's story"
anchor:
  path: wip/kanban/stories/S-0486-criteria-tick-is-refused-on-another-open-story-s-narrative-state-finding-which-the-story-s-own-scoped-check-leaves-out.md
  item: S-0486
status: answered
participants: [planner-S-0486, orchestrator]
created: 2026-10-23T17:51:54Z
updated: 2026-10-23T17:52:20Z
---

# TH-0618 Plan for S-0486: scope the check a checked edit runs to the item's story

On wip/kanban/stories/S-0486-criteria-tick-is-refused-on-another-open-story-s-narrative-state-finding-which-the-story-s-own-scoped-check-leaves-out.md.

## Entries

### 2026-10-23T17:51:54Z planner-S-0486
Plan summary from planner-S-0486.

**Cause found.** `flai criteria tick`, `criteria_tick`, `flai edit`, `item_edit`, and the dashboard's edits save through `itemedit.Apply` (`flai/internal/itemedit/itemedit.go`), which runs `check.Run` on the whole repository before and after the write and refuses the edit on any finding the second run has that the first did not. The diff is not scoped: a finding another writer brings between the two runs refuses the edit. Both I-0182 instances are that: a `narrative.state` on the narrative of S-0471, just started, and a `markdown.MD012` on T-2042, an untracked task file another agent was writing. A check scoped to a story (`check.ScopeToStories`, `flai/internal/check/scope.go`) leaves the first out (ADR-0125) and marks the second a note outside the story (ADR-0085); `flai task done`, `flai verify`, and `flai check --story` run it so, and `Apply` does not.

**Proposed fix.** `Apply` scopes its second run to the edited item's story, the story itself or a task's parent, with `ScopeToStories` and the paths the edit wrote as the changed paths: a new finding the scope leaves out refuses nothing, a new finding outside the scope is returned as a warning, a new finding inside, or an error on the item's own file, refuses as now. An epic, which has no story, is refused only on a new finding on a path the edit wrote. `itemnew`, `docedit`, and `adr` keep their unscoped diff: no instance asks for more. An ADR records it, refining ADR-0085 and ADR-0089.

**Tasks and layers.**

1. T-2074: the ADR (`design/adrs`, `design/adrs/README.md`).
2. T-2075, after T-2074: the change in `itemedit.Apply` with tests that reproduce both instances (`flai/internal/itemedit/itemedit.go`, `criteria_test.go`, `flai/internal/mcpserver/criteria_test.go`).
3. T-2076, after T-2075: the design and the user docs (`design/system/flai-cli.md` rows `flai edit` and `flai criteria`, `flai/cmd/edit.go` and `flai/cmd/criteria.go` long help, `docs/users/flai-reference.md` regenerated, `docs/users/flai.md`).
4. T-2077, after T-2075 and T-2076: `flai issue close I-0182` and the criteria ticks.

One task per layer; nothing runs together.

**Figures.** Touches: 12, every one on a task but the folder `design/adrs`, kept for the ADR whose number `flai adr new` allocates. Forecast: flai's 22m at size 14 raised to 2h, delivery 2026-10-24T13:00Z, for an ADR that awaits the operator, tests that must bring a finding between two check runs in a writer `flai edit` and the dashboard share, and docs in three places with a generated reference. Cost of delay: flai's 10.00 USD a week from the input flai set from I-0182 (4m per 168h cycle), unchanged.

**Assumptions.**

- The ADR is worth a task: ADR-0122, -0123, -0125, -0133, -0136, and -0137 each recorded one left-out rule for the scoped check, and this changes what every checked item edit refuses.
- `ScopeToStories` can take the edit's written paths as `changed` without a new helper; if the code says otherwise, T-2075 adds one in `flai/internal/check/scope.go`, which `flai task done` widens the touches to.
- `flai/internal/itemedit/itemedit.go`, `check/scope.go`, `design/system/flai-cli.md`, and the two user docs are touches of S-0468 (in review) and `flai-cli.md` and the user docs of S-0429 (in progress). S-0468 holds nothing in review; S-0429 holds this story in ready until it is accepted, unless its agent shares the doc paths. T-2076 changes only the `flai edit` and `flai criteria` rows, which neither story changes.
- No template change: the fix is in flai alone, so `template/CHANGELOG.md` is not a touch.

Nothing to split, merge, or drop.

### 2026-10-23T17:52:20Z orchestrator
The plan's shape is approved. I will resolve this thread once two things are settled.

**Approved:**
- **Cause.** `itemedit.Apply` diffs two unscoped `check.Run` calls, so a finding another writer brings in between refuses the edit. That explains both instances of I-0182.
- **Fix.** Scope the second run with `ScopeToStories`, and record it in an ADR that refines ADR-0085 and ADR-0089.
- **Tasks.** Four tasks in four layers, each criterion covered.
- **Touches.** 12 files. The one folder, `design/adrs`, sits inside `claims.shared`.

**1. Write the figures.** `flai promote --drafts` still reports S-0486 lacks a forecast duration, a forecast delivery, and a cost of delay value. Please write them with `flai forecast` and `flai cod`.

**2. Revisit the 2h forecast.** `forecast.duration` is agent time: the seconds per unit of size times the size (`design/system/strategic-agents.md` § Forecast). An ADR that awaits the operator adds waiting, not agent time, so it does not justify a raise. Raising flai's 22m for the two-writer test and the docs in three places is a fair reason. Please state the raise on those grounds alone, or say why 2h still stands without the wait.

Source: design/system/strategic-agents.md § Forecast
