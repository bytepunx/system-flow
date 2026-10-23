---
id: T-2074
type: task
nature: remediation
title: An ADR settles that a checked edit of a work item refuses only a finding the edit introduces inside its story's scope, and leaves out what a check scoped to the story leaves out
status: backlog
parent: S-0486
owner: alex
created: 2026-10-23T17:50:12Z
updated: 2026-10-23T17:50:12Z
transitions: []
stream: S-0486
tags: [flai]
touches: [design/adrs, design/adrs/README.md]
---
# T-2074 An ADR settles that a checked edit of a work item refuses only a finding the edit introduces inside its story's scope, and leaves out what a check scoped to the story leaves out

## Work

`flai criteria tick`, the MCP tool `criteria_tick`, `flai edit`, `item_edit`, and the dashboard's edits all save through `itemedit.Apply` (`flai/internal/itemedit/itemedit.go`), which runs `check.Run` on the whole repository before the write and again after it, and refuses the edit on any finding the second run has that the first did not, "flai check has N finding(s) with this change". That diff is not scoped: a finding another writer brings between the two runs refuses the edit, as I-0182's two instances show, a `narrative.state` on the narrative of a story that had just started and a `markdown.MD012` on an untracked task file of another story. A check scoped to a story (`check.ScopeToStories`, `flai/internal/check/scope.go`) leaves the first out (ADR-0125) and marks the second a note outside the story (ADR-0085), and `flai task done`, `flai verify`, and `flai check --story` run it so.

Write the ADR with `flai adr new`, refining ADR-0085 and ADR-0089, and read ADR-0122, ADR-0123, ADR-0125, ADR-0133, ADR-0136, and ADR-0137 first, whose rules it reuses. The decision to record, as the plan proposes it on the story's plan thread:

1. A checked edit of a story or a task scopes the second run to the item's story, the story itself or the task's parent, with `check.ScopeToStories` and the paths the edit wrote as the story's changed paths. A new finding the scope leaves out never refuses the edit; a new finding the scope marks outside is returned as a warning, not a refusal; a new finding inside the scope refuses it as now.
2. A checked edit of an epic, which has no story, refuses a new finding only on a path the edit wrote, and returns the others as warnings.
3. A finding on the item's own file at level error refuses the edit whether new or not, as now; a `wip.overlap` stays a warning (ADR-0019).
4. `flai item new` and the other checked writes (`itemnew`, `docedit`, `adr`) keep their unscoped diff: no instance asks for more, and the ADR says so under alternatives.

Amend the proposal where the code says otherwise, and name in the ADR what the code task tests: the two instances of I-0182. This task is the plan's first layer: the code task waits for the number and the words this ADR fixes.

## Done when

`design/adrs/NNNN-<slug>.md` exists with status `proposed`, `refines: [ADR-0085, ADR-0089]`, its row in `design/adrs/README.md`, and sections Context, Decision, Consequences, and Alternatives considered; its Context names I-0182's two instances; `flai check --strict` reports nothing on it.

## Notes

The folder touch `design/adrs` stands because `flai adr new` allocates the number when it runs; `flai task done` widens the touches to the file.
