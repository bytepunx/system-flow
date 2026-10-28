---
id: D-0289
title: "TestActivityWritersTakeTurns fails about one run in three: two writers that both find the activity lock stale can both take it, and an entry is lost"
kind: cause
sources: [I-0140]
topics: [flai, cli, go, code, all]
paths: [flai/internal/filelock/filelock.go, flai/internal/filelock/filelock_unix.go, flai/internal/filelock/filelock_windows.go, flai/internal/filelock/filelock_test.go, flai/go.mod, design/tech/go-libraries.md, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/internal/serve/restart.go, flai/internal/serve/restart_test.go, design/system/strategic-agents.md, design/issues/I-0140-testactivitywriterstaketurns-fails-about-one-run-in-three-two-writers-that-both-find-the-activity-lock-stale-can-both-take-it-and-an-entry-is-lost.md, design/issues/summary.md, flai/cmd/guard_test.go, flai/cmd/promote_test.go]
date: 2026-10-13
status: absorbed
absorbed_by: "design/tech/go-libraries.md#go-libraries-used-by-flai"
asked_again: []
cited: []
---

# D-0289 TestActivityWritersTakeTurns fails about one run in three: two writers that both find the activity lock stale can both take it, and an entry is lost

I-0140 was closed on 2026-10-13.

The cause:

```text
TestActivityWritersTakeTurns fails about one run in three: two writers that both find the activity lock stale can both take it, and an entry is lost
```

The fix:

```text
S-0377 fixed it. lockActivity, the guard's sub-agent record (RecordChange), and flai serve's lockFile now take an operating-system file lock through flai/internal/filelock: flock(2) on unix, LockFileEx on Windows, on a lock file that stays in place. The kernel frees the lock when its holder exits, so no lock is ever broken by its age and two writers cannot both hold it. TestActivityWritersTakeTurns passed 350 runs in a row under -race, where the old code lost an entry. TestLockFreedWhenHolderDies shows that a killed holder's lock is taken at once.
```

Source: I-0140.
