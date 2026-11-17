---
id: T-3470
type: task
nature: feature
title: The board runbook gives the suites workflow file as a worked example of verify, runs, and verify_at
status: backlog
parent: S-0815
owner: alex
created: 2026-11-16T23:50:30Z
updated: 2026-11-16T23:50:30Z
transitions: []
stream: S-0815
tags: [template]
touches: [docs/operators/runbooks/change-the-board.md, docs/users/conventions.md]
after: [T-3463]
---
# T-3470 The board runbook gives the suites workflow file as a worked example of verify, runs, and verify_at

## Work

It waits for T-3463 so that the example is the file the test proves, `flai/cmd/testdata/e2e-suites-board.md`, and not a copy that can drift from it.

In `docs/operators/runbooks/change-the-board.md`, add a section after `## Add a rule to a move`, such as `## Choose what each step verifies`. Give the fixture's front matter as the worked example: review needs the fast suites, an `integration` state runs the slow suites and the cluster's deploy and checks, and publish runs the slow suites. Say what each key does, that a cluster step is skipped where no cluster is configured, and how to see the suites each step ran. In `docs/users/conventions.md` § States, add one line that links to the section, so a reader of the users' guide finds it.

## Done when

The runbook's example parses as the fixture does, the markdown lint passes, and `docs/users/conventions.md` links to it.

## Notes
