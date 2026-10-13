---
id: I-0140
title: "TestActivityWritersTakeTurns fails about one run in three: two writers that both find the activity lock stale can both take it, and an entry is lost"
class: defect
status: closed
count: 1
cost: 15m
first_reported: 2026-10-12T16:22:16Z
last_reported: 2026-10-12T16:22:16Z
updated: 2026-10-13T17:07:24Z
---

# I-0140 TestActivityWritersTakeTurns fails about one run in three: two writers that both find the activity lock stale can both take it, and an entry is lost

## Description
TestActivityWritersTakeTurns fails about one run in three: two writers that both find the activity lock stale can both take it, and an entry is lost

## Instances

### 2026-10-12T16:22:16Z
Story: S-0337.
S-0337's close-out stopped at the integration tier with flai/internal/workitem TestActivityWritersTakeTurns: "an entry was lost: 5 entries, 6 refusals". S-0337 does not change flai/internal/workitem. Run three more times with flai test on the package, it passed twice and failed once, with 4 entries. The cause is in lockActivity (flai/internal/workitem/activity.go), which breaks a lock older than activityLockAge with Stat then Remove. Two writers can both Stat the stale lock. One removes it and creates a fresh lock, then the other removes that fresh lock, believing it stale, and creates its own. Both then hold the lock, and one write of orchestrator.md overwrites the other. The same race can lose a real activity entry or guard refusal when flai serve and flai guard write at once after a crashed flai left a lock. A fix needs breaking a stale lock atomically, for example by renaming it to a unique name and checking the renamed file's age, or an OS file lock.

## Remediation

Story S-0377 remediates this issue, created from it at 2026-10-12T16:53:00Z.
Closed 2026-10-13T17:07:24Z: S-0377 fixed it. lockActivity, the guard's sub-agent record (RecordChange), and flai serve's lockFile now take an operating-system file lock through flai/internal/filelock: flock(2) on unix, LockFileEx on Windows, on a lock file that stays in place. The kernel frees the lock when its holder exits, so no lock is ever broken by its age and two writers cannot both hold it. TestActivityWritersTakeTurns passed 350 runs in a row under -race, where the old code lost an entry. TestLockFreedWhenHolderDies shows that a killed holder's lock is taken at once.
