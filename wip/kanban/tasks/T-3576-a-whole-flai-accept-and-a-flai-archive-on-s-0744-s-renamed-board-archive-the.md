---
id: T-3576
type: task
nature: remediation
title: A whole flai accept and a flai archive on S-0744's renamed board archive the shipped story
status: backlog
parent: S-0842
owner: alex
created: 2026-11-17T12:29:45Z
updated: 2026-11-17T12:29:45Z
transitions: []
stream: S-0842
tags: [flai, workflow, test]
touches: [flai/cmd/accept_workflow_test.go, flai/cmd/archive_test.go]
after: [T-3572]
---
# T-3576 A whole flai accept and a flai archive on S-0744's renamed board archive the shipped story

## Work

Reproduce I-0377 at the command line. In `flai/cmd/accept_workflow_test.go`, on the board `renamedProject` (`flai/cmd/move_workflow_test.go`) builds, whose accepted state is `shipped`, run a whole `flai accept` of a story in review: it walks the story to `shipped`, archives it and its tasks with the narrative, and commits. Today's test checks only the refusal and the dry run. In `flai/cmd/archive_test.go`, `flai archive` with an ID and without archives a shipped story there. Each test fails at the base of the story and passes after T-3572. Waits for T-3572, whose `Closed` and `PlanArchive` it tests.

## Done when

- The whole acceptance on the renamed board archives the story; run against the code before T-3572, it fails with the refusal I-0377 quotes.
- `flai archive` on the renamed board archives a shipped story, with an ID and without.
- `flai test` on the changed paths passes.

## Notes
