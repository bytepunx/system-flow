---
id: T-1054
type: task
nature: improvement
title: The workitem package replaces a narrative's Current state and Next steps and leaves the rest of it alone
status: in-progress
parent: S-0271
owner: alex
created: 2026-10-06T22:49:37Z
updated: 2026-10-07T06:47:07Z
transitions:
  - to: ready
    at: 2026-10-07T06:47:06Z
    by: agent-S-0271
  - to: in-progress
    at: 2026-10-07T06:47:07Z
    by: agent-S-0271
stream: S-0271
tags: [cli]
touches: [flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go]
---
# T-1054 The workitem package replaces a narrative's Current state and Next steps and leaves the rest of it alone

## Work

This is the core that the command, the MCP tool, and the host method share. Today agents rewrite the two sections with Edit or Write.

- In `flai/internal/workitem/narrative.go`, add a function on `Repo`, such as `SetStreamState(storyID, current, next string, opt StreamOptions)`.
  - It replaces the body of `## Current state` and of `## Next steps` with the text given, up to the next second-level heading.
  - It leaves every other byte of the narrative as it was, and appends nothing to `## Log`.
  - Either text may be left out, which leaves that section as it is. Giving neither is an error.
- Refuse a story that is not `in-progress` or `review`, and a story with no narrative. Each refusal names the story's state and what to do.
- Refuse text that would bring a finding of the project's markdown lint, through `Repo.LintGuard`, as `LogStream` does.
- Regenerate `wip/agents/index.md` as `LogStream` does, if the index shows the narrative's state.
- Export the section parser, if one is needed, so that `flai check` can reuse it for the sections' shape.

## Done when

- [ ] Tests in `narrative_test.go` show the following:
  - Both sections are replaced, and every other section stays byte for byte the same.
  - Leaving one text out keeps that section.
  - A story in `ready` or `done` is refused, and so is a story with no narrative.
  - Text that fails the lint is refused, and the file is left unchanged.
- [ ] `scripts/flai-test.sh` passes.

## Notes

Drafted by the planner for S-0271.
