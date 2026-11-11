---
id: T-3003
type: task
nature: improvement
title: merge.queue reads the heads once and runs git status only when a story in review passed verify, at the same time as the heads
status: done
parent: S-0698
owner: alex
created: 2026-11-10T10:31:13Z
updated: 2026-11-11T08:42:44Z
transitions:
  - to: ready
    at: 2026-11-11T08:42:38Z
    by: agent-S-0698
  - to: in-progress
    at: 2026-11-11T08:42:38Z
    by: agent-S-0698
  - to: done
    at: 2026-11-11T08:42:43Z
    by: agent-S-0698
stream: S-0698
tags: [cli]
touches: [flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/server_test.go]
after: [T-3000]
usage:
  source: log
  seconds: 5
  estimated: true
  priced_by: estimate
  models: []
---
# T-3003 merge.queue reads the heads once and runs git status only when a story in review passed verify, at the same time as the heads

## Work

Write the test first in `flai/internal/mergequeue/eligible_test.go`. It reproduces I-0333: a runner that counts git processes, and five stories in review whose last `flai verify` passed. On `main` it counts five `git rev-parse` calls for the heads, plus `git status` and the two main-head reads.

Then change `ComputeAt` in `flai/internal/mergequeue/eligible.go`:

- Read the verify reports of the stories in review first. They are files, so this runs no git.
- When no report passed, run no git for the heads or the dirt. No story has a head that dirt could hold back.
- Otherwise run T-3000's heads read and `mainDirt`'s `git status --porcelain` at the same time. `verified` takes each story's head from the heads read, and `mainDirt` and `filling` take main's head from it, in place of `mainHead`'s two rev-parse calls.
- `mainHead` in `land.go` stays as it is for the landing.

Every reason, its kind and text, and the window stay as they are now.

It waits for T-3000, whose heads read it calls.

## Done when

- The new test shows two git processes, the heads read and `git status`, with five stories in review that passed verify. It shows the same two with one such story.
- With stories in review and no passing verify, the test shows no git process.
- The existing tests of `eligible_test.go`, `batch_test.go`, and `land_test.go` pass unchanged. This covers the branch-unread, verify-stale, dirty-main, and window cases.
- `flai test flai/internal/mergequeue` passes.

## Notes
