---
id: T-2435
type: task
nature: remediation
title: Confirm with a race stress run that the serve tests no longer fail their temp folder's cleanup, and fix the writer if one still does
status: backlog
parent: S-0543
owner: alex
created: 2026-11-02T18:12:56Z
updated: 2026-11-02T18:12:56Z
transitions: []
stream: S-0543
tags: [flai, serve, test]
touches: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go]
---
# T-2435 Confirm with a race stress run that the serve tests no longer fail their temp folder's cleanup, and fix the writer if one still does

## Work

I-0233 was closed on 2026-10-31 as fixed by S-0418's fe589e207, where `Run` waits for each served project's looks before it returns, and S-0548's c6c0310db, where the test waits for the watcher's first look. This task confirms that fix on the story's base. It waits for no other task.

Run `GOMAXPROCS=2 go test -race -count=20 -run 'TestAChangedFileReachesTheDashboardAsANotification|TestServeStartsADefinedRoleOnAnEventAndOnASchedule' ./internal/serve/` from `flai/`. Then run the whole package once with `-race`.

If every run passes, change no code. Write the commands and their pass counts in the story's Notes.

If a run fails its cleanup with "directory not empty", find what still writes under the temp folder after `Run` returns. Look at the inbox cursors, the reason queues, the journal, and the launcher. Make `Run`'s halt in `serve.go` wait for that writer, or make the test wait for it, whichever owns it. Add a test in `serve_test.go` that fails without the wait. Record the new instance with `flai issue bump I-0233` before the fix.

## Done when

- The stress run passes 20 of 20, and the whole serve package passes under `-race`, with the counts in the story's Notes.
- Or a failing run's writer is waited for, a test that fails without the wait passes with it, and I-0233 has the new instance and is closed again with the reason naming the fix.

## Notes
