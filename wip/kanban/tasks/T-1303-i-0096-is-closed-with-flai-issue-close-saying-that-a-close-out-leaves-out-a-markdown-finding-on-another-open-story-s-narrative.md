---
id: T-1303
type: task
nature: improvement
title: I-0096 is closed with flai issue close, saying that a close-out leaves out a markdown finding on another open story's narrative
status: done
parent: S-0318
owner: alex
created: 2026-10-08T00:10:51Z
updated: 2026-10-08T04:45:28Z
transitions:
  - to: ready
    at: 2026-10-08T04:45:22Z
    by: agent-S-0318
  - to: in-progress
    at: 2026-10-08T04:45:23Z
    by: agent-S-0318
  - to: done
    at: 2026-10-08T04:45:28Z
    by: agent-S-0318
stream: S-0318
tags: [flai]
touches: [design/issues/I-0096-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-1301]
usage:
  source: log
  seconds: 5
  estimated: true
  models: []
---
# T-1303 I-0096 is closed with flai issue close, saying that a close-out leaves out a markdown finding on another open story's narrative

## Work

Close I-0096 with `flai issue close I-0096 --reason`, in the story's worktree. It waits for T-1301, because the issue closes only once the fix is built and tested. It shares no path with T-1302, so the two can run together.

The reason names:

- the cause: S-0229's agent hand-edited its narrative's `## Decisions` with a code span that broke MD038, and S-0227's close-out recorded it, though only S-0229 could fix it;
- the fix: a check scoped to a story leaves out a markdown finding on another open story's narrative, so a close-out records none, and that story's own close-out still finds it;
- the ADR T-1300 wrote, and the tests T-1301 added.

## Done when

- I-0096 is closed, with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes
