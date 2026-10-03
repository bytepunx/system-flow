---
id: T-0730
type: task
nature: feature
title: Each issue instance records its story, and an issue can be turned into a story's title, nature, and body
status: done
parent: S-0198
owner: arobson
created: 2026-10-03T01:28:18Z
updated: 2026-10-03T01:34:49Z
transitions:
  - to: ready
    at: 2026-10-03T01:28:52Z
    by: agent-S-0198
  - to: in-progress
    at: 2026-10-03T01:28:52Z
    by: agent-S-0198
  - to: done
    at: 2026-10-03T01:34:49Z
    by: agent-S-0198
stream: S-0198
tags: []
touches: [flai/internal/issues]
usage:
  source: log
  seconds: 357
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 28
      output: 9492
      cache_read: 1622134
      cache_write: 33248
      cost: 0.7123
---
# T-0730 Each issue instance records its story, and an issue can be turned into a story's title, nature, and body

## Work

In `flai/internal/issues`:

- `New` and `Bump` take the story the instance belongs to (an ID such as `S-0198`, or none). The instance names it in a line under its `### <timestamp>` heading, `Story: S-nnnn.`, so readers see it and flai can parse it back. A same-second join keeps one story line per story.
- `Stories(is)` returns the stories an issue's instances name, in order, without repeats. `RecordedBy(list, story)` returns the issues whose instances name a story.
- `ForStory(is)` returns what a story made from the issue needs: a title, the nature from the class, and the body. The nature is `remediation` for `defect` and `blocker`, and `improvement` for `efficiency` and `impression`. The body has a `## Goal` that links the issue's document and gives the issue's `## Remediation` text as the recommended solution, or says to propose one when there is none. Its `## Acceptance criteria` end with closing the issue through `flai issue close I-nnnn --reason`.
- `Links(body, is)` reports whether a work item body names the issue by its ID or by its document's file name. `NoStory(issues, items)` returns the open issues that no open story (`backlog`, `ready`, `in-progress`, `review`) links.
- `SetRemediation(is, storyID)` records the story under `## Remediation`.

Waits for nothing: the command, the check, the MCP tool, and the dashboard all build on it.

## Done when

- Behaviour tests in `flai/internal/issues` cover the story line on new and bump, the same-second join, `Stories`, `ForStory` for each class, `Links`, `NoStory`, and `SetRemediation`.
- `go test -race -short ./internal/issues/` passes.

## Notes
