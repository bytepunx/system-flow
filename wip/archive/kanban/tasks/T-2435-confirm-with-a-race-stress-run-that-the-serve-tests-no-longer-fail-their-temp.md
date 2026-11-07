---
id: T-2435
type: task
nature: remediation
title: Confirm with a race stress run that the serve tests no longer fail their temp folder's cleanup, and fix the writer if one still does
status: done
parent: S-0543
owner: alex
created: 2026-11-02T18:12:56Z
updated: 2026-11-07T10:20:20Z
transitions:
  - to: ready
    at: 2026-11-07T10:03:51Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:03:55Z
    by: system-flow
  - to: done
    at: 2026-11-07T10:20:20Z
    by: agent-S-0543
stream: S-0543
tags: [flai, serve, test]
touches: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go]
usage:
  source: log
  seconds: 985
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 25
      output: 4710
      cache_read: 2274346
      cache_write: 96227
      cost: 1.319
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

2026-11-07, on story/S-0543 at 8cb8d4c5e, on a host at load average 54 to 71 on 24 CPUs:

- `GOMAXPROCS=2 go test -race -count=20 -run '^(TestAChangedFileReachesTheDashboardAsANotification|TestServeStartsADefinedRoleOnAnEventAndOnASchedule)$' ./internal/serve/` passed 40 of 40 runs, 20 of each test, in 25.4 s.
- `GOMAXPROCS=2 go test -race -count=3 -skip '^TestGateway' ./internal/serve/` passed 3 of 3 whole-package runs in 398.5 s. No "directory not empty" and no data race.
- A first whole-package run without `-skip '^TestGateway'` failed only TestGatewayLiteLLM, at its settings check. That is a separate defect, recorded as I-0303.

No code changed.
