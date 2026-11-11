---
id: T-3121
type: task
nature: remediation
title: wait_for_events, in one project and across a folder of projects, answers past a task file that holds only PLACEHOLDER, with tests
status: done
parent: S-0717
owner: alex
created: 2026-11-11T09:07:28Z
updated: 2026-11-11T22:49:15Z
transitions:
  - to: ready
    at: 2026-11-11T22:32:57Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T22:32:57Z
    by: system-flow
  - to: done
    at: 2026-11-11T22:49:14Z
    by: agent-S-0717
stream: S-0717
tags: [flai]
touches: [flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder_test.go]
usage:
  source: log
  seconds: 977
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 45
      output: 8538
      cache_read: 3289122
      cache_write: 81162
      cost: 1.4122
---
# T-3121 wait_for_events, in one project and across a folder of projects, answers past a task file that holds only PLACEHOLDER, with tests

## Work

`wait_for_events` reads the items through `Repo.List`: `server.catchUp` in `mcpserver/cursor.go` calls `listItems`. When a server serves a folder of projects, `folderWaits.catchUp` in `mcpserver/folder.go` calls it for each project and returns the first error for every project. It does so both before the hold and inside the hold's callback. Today one file with no front matter fails the whole call.

S-0716 makes `Repo.List` skip a file that does not parse (T-3107). Its T-3109 pins `wait_for_events` and `inbox` for one project in `mcpserver/inbox_test.go`. This task adds what that leaves out. Each test writes an 11-byte `PLACEHOLDER` file into `wip/kanban/tasks/`, as I-0345's instance saw:

- One project, `server_test.go`: a `wait_for_events` call already behind the cursor, and one held while the file appears and another item changes. Both answer with the other item's change and no error.
- A folder of two projects, `folder_test.go`: the bad file in one project. The call answers with the other project's events, and with the first project's other changes, and no error, both before the hold and while held.

If a test still fails once S-0716 is on main, change `catchUp` in `cursor.go` or `folder.go` so that one project's unreadable item does not fail the others. Change neither otherwise. Look at `TestAHeldWaitFailsOnATruncatedStory` in `folder_test.go` before adding a test beside it. It may pin the failure that this task removes, and then it changes with it.

This task waits for no task of this story. It builds on S-0716's loader change, which the story's `after` names, and it touches no file that T-3120 touches.

## Done when

- The tests above pass with a `PLACEHOLDER` task file present.
- `flai test` on the changed files passes.

## Notes
