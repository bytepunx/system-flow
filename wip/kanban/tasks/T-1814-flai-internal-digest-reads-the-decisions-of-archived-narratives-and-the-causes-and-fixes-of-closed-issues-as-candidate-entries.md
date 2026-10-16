---
id: T-1814
type: task
nature: feature
title: flai/internal/digest reads the decisions of archived narratives and the causes and fixes of closed issues as candidate entries
status: backlog
parent: S-0429
owner: alex
created: 2026-10-16T06:03:54Z
updated: 2026-10-16T06:03:54Z
transitions: []
stream: S-0429
tags: [cli]
touches: [flai/internal/digest/sources.go, flai/internal/digest/sources_test.go]
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
