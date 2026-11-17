---
id: T-3468
type: task
nature: feature
title: A flaiover test of the built flai reads the suites project's board and review, with the integration column and the suites each step needs and ran
status: backlog
parent: S-0815
owner: alex
created: 2026-11-16T23:50:25Z
updated: 2026-11-16T23:50:25Z
transitions: []
stream: S-0815
tags: [flaiover, dashboard]
touches: [flaiover/src/lib/server/suites.test.ts]
after: [T-3463]
---
# T-3468 A flaiover test of the built flai reads the suites project's board and review, with the integration column and the suites each step needs and ran

## Work

It waits for T-3463 because it builds its project from that task's fixture, `flai/cmd/testdata/e2e-suites-board.md`.

Write `suites.test.ts`, importing `testing.ts` and skipped with `describe.skipIf(!haveFlai)` as `inbox.test.ts` is. Make a project with that fixture as its workflow file, move a story to `review` and on to `integration` with `shell`, and read `board.get` and the review data with `flaiAsk`. Check that the board has the `integration` column after `review` with its kind, and that the review data names the suites the move to review and the acceptance need, the suites the `integration` state runs, and those that ran, as S-0814 gives them.

## Done when

The test passes in the `flaiover-flai` tier against the built flai, and fails when the board or the review data lacks the `integration` state or its suites.

## Notes
