---
id: T-0186
type: task
nature: feature
title: flai stream diff, and progress events from flai accept
status: done
parent: S-0041
owner: alex
created: 2026-09-19T04:08:29Z
updated: 2026-09-19T04:11:03Z
transitions:
  - to: ready
    at: 2026-09-19T04:08:53Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T04:08:53Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:11:03Z
    by: system-flow
stream: S-0041
tags: []
touches: [flai/cmd]
---

# T-0186 flai stream diff, and progress events from flai accept

## Work
Add `flai stream diff <story-id> [--json]`: the story branch against its merge base with the main branch, as files (path, status, additions, deletions) each with its unified hunks, read with git in the main checkout so it works whether or not the worktree is readable; binary files are listed without hunks; a file's patch over a size limit is cut and marked `truncated`, and the whole result is capped, so a large branch cannot swamp the browser. A story with no branch is an error that says so. Make `flai accept` log one info event per step as it completes (branch merged, moved to done, archived, committed, tagged, pushed or not pushed), with stable field names, so a client can show progress; the final JSON result is unchanged. Tests with real git: added, modified, deleted, and renamed files, a binary file, truncation, no branch; the accept events in order.

## Done when
- The new tests pass with `-race`
- The existing accept tests pass unchanged

## Notes
