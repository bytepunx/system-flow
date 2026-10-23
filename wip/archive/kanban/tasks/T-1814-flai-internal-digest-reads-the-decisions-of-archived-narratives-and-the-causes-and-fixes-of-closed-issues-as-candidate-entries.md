---
id: T-1814
type: task
nature: feature
title: flai/internal/digest reads the decisions of archived narratives and the causes and fixes of closed issues as candidate entries
status: done
parent: S-0429
owner: alex
created: 2026-10-16T06:03:54Z
updated: 2026-10-23T17:47:09Z
transitions:
  - to: ready
    at: 2026-10-23T17:37:50Z
    by: agent-S-0429
  - to: in-progress
    at: 2026-10-23T17:37:54Z
    by: agent-S-0429
  - to: done
    at: 2026-10-23T17:47:09Z
    by: agent-S-0429
stream: S-0429
tags: [cli]
touches: [flai/internal/digest/sources.go, flai/internal/digest/sources_test.go]
usage:
  source: log
  seconds: 555
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 37
      output: 16951
      cache_read: 2778768
      cache_write: 89990
      cost: 1.4358
---
# T-1814 flai/internal/digest reads the decisions of archived narratives and the causes and fixes of closed issues as candidate entries

## Work

Add `sources.go`:

- `Decisions(root, since)`: each bullet under `## Decisions` of each narrative in `wip/archive/agents/` (and `wip/agents/` for a story done but not yet archived) whose story was completed on or after `since`, as a `decision` candidate with the story as its source, the story's topics and touches, and its completion date. Read the section with `workitem.NarrativeSection`; a narrative whose path is in the archive is read from there.
- `Causes(root, since)`: each issue under `design/issues/` closed on or after `since` with a cause and a fix in its body, as a `cause` candidate with the issue as its source, the topics and paths of the story that closed it, and its close date.

First layer: runs beside the answer rule's task, which touches other files.

## Done when

`sources_test.go` covers a narrative with decisions, one without the section, an archived and an unarchived narrative, a closed issue with a cause and fix, one without, and an open issue left out; `flai test flai/internal/digest` passes.

## Notes

If `workitem.NarrativeSection` cannot read an archived narrative's path, extend it there and add `flai/internal/workitem/narrative.go` to this task's touches.

Planner-S-0429, 2026-10-23: `workitem.NarrativeSection(text, name)` takes the narrative's text, so an archived narrative is read with `workitem.ReadNarrative` from `wip/archive/agents/S-nnnn.md` and the section taken from its `Body`; `narrative.go` needs no change. A story's completion date is `Item.CompletedAt()` on the archived story under `wip/archive/kanban/stories/`. An issue file (`flai/internal/issues`) has no cause or fix heading: its sections are `## Description`, `## Instances`, and `## Remediation`, and `issues.Close` appends `Closed <ts>: <reason>` to the body. Read the cause from `## Description`, the fix from the `Closed` line's reason, with `## Remediation` when the reason names none, the close date from that line's timestamp, and the story that closed it from the `S-nnnn` the reason names. A rewording of `## Work` to say so is proposed on the plan's thread. On 2026-10-23, 168 of the 190 issues are closed and 442 narratives are archived.
