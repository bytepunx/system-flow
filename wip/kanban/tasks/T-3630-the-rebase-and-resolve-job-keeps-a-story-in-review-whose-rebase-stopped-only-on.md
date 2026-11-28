---
id: T-3630
type: task
nature: improvement
title: The rebase and resolve job keeps a story in review whose rebase stopped only on two appends to one Go test file
status: backlog
parent: S-0857
owner: alex
created: 2026-11-27T23:58:25Z
updated: 2026-11-27T23:58:25Z
transitions: []
stream: S-0857
tags: [go, serve]
touches: [flai/internal/serve/resolvejob.go, flai/internal/serve/resolvejob_test.go]
after: [T-3628]
---
# T-3630 The rebase and resolve job keeps a story in review whose rebase stopped only on two appends to one Go test file

## Work

Reproduce I-0409 through the acceptance's rebase and resolve job: a test beside `TestARebaseAndResolveWhoseRebaseConflictsAbortsItAndSendsTheStoryBack` where the story's branch and main each appended a test function after the same function in one `_test.go` file. It waits for T-3628, since the job settles the stop only once the sync it runs takes the Go merge. Add the Go merge to the kinds `conflictText` names as settled, so the send-back thread stays true.

## Done when

- The new test passes: the job's rebase ends clean, the close-out runs, and the story stays in review with both functions in the file.
- The existing send-back test still passes, and `conflictText` names the Go merge among what the rebase settles.
- `flai test flai/internal/serve` passes.

## Notes
