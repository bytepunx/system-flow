---
id: T-1328
type: task
nature: remediation
title: Test that three tasks of one layer closed in turn with flai task done each get a commit of their own, as I-0108 reports
status: backlog
parent: S-0322
owner: alex
created: 2026-10-08T04:32:40Z
updated: 2026-10-08T04:32:40Z
transitions: []
stream: S-0322
tags: [flai]
touches: [flai/internal/taskdone/taskdone_test.go]
---
# T-1328 Test that three tasks of one layer closed in turn with flai task done each get a commit of their own, as I-0108 reports

## Work

I-0108 is the fault S-0312 remediates as I-0104: `flai task done` ran `git add -A`, so closing one task of a layer committed its siblings' files. S-0312's T-1283 changes the commit step and tests it with two open tasks. I-0108's instance is S-0274's layer 4, three tasks in one worktree (T-1115 with T-1117 and T-1118), each of which needed `git commit --amend` or `git reset --soft` to get a commit of its own.

Read the commit step and the tests S-0312 left in `flai/internal/taskdone/taskdone.go` and `flai/internal/taskdone/taskdone_test.go`. When no test there closes three open tasks of one story in turn, add one with the file's `project()` and `story()` helpers: three tasks, each with files changed in the worktree under its touches; close each with `Run` in turn; check that each close makes one commit holding only that task's files, that the tasks still open keep theirs uncommitted and listed in the result, and that no task's touches widen with another's paths. When such a test exists already, change nothing here and name it in the task's close message.

It waits for nothing in this story. The story waits for S-0312, whose commit step this tests.

## Done when

- A test in `flai/internal/taskdone/taskdone_test.go` closes three tasks of one layer in turn and finds one commit per task with only its own files, or an existing test that does so is named.
- `flai test flai/internal/taskdone` passes.

## Notes

Drafted by the planner for S-0322.
