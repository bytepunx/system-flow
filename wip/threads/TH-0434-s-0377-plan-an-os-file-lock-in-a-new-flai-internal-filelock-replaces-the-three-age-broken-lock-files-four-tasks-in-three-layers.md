---
id: TH-0434
title: "S-0377 plan: an OS file lock in a new flai/internal/filelock replaces the three age-broken lock files, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0377-testactivitywriterstaketurns-fails-about-one-run-in-three-two-writers-that-both-find-the-activity-lock-stale-can-both-take-it-and-an-entry-is-lost.md
  item: S-0377
status: resolved
participants: [planner-S-0377, orchestrator]
created: 2026-10-13T16:48:46Z
updated: 2026-10-13T16:48:59Z
---

# TH-0434 S-0377 plan: an OS file lock in a new flai/internal/filelock replaces the three age-broken lock files, four tasks in three layers

On wip/kanban/stories/S-0377-testactivitywriterstaketurns-fails-about-one-run-in-three-two-writers-that-both-find-the-activity-lock-stale-can-both-take-it-and-an-entry-is-lost.md.

## Entries

### 2026-10-13T16:48:46Z planner-S-0377
**Plan:** every flai writer that today makes a lock file with `O_EXCL` and breaks it by age takes an operating-system file lock instead. That is `flock` on unix and `LockFileEx` on Windows, in a new package `flai/internal/filelock`. The kernel frees the lock when its holder dies, so no lock is ever stale and two writers can never both hold it.

## Remedy, from I-0140's instance

- The cause is the age break in `lockActivity`: `os.Stat`, then `os.Remove`. Two writers both see the stale lock. One removes it and creates a fresh one, and the other then removes that fresh one.
- I did not propose a rename-to-a-unique-name break. It narrows the window, but a third writer that creates a lock between the rename and the age check still splits the lock.
- `lockRecord` in `flai/internal/guard/subagents.go` and `lockFile` in `flai/internal/serve/restart.go` break their locks the same way, so the plan fixes them too.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1510 `filelock` package: unix and Windows files, tests including a child process killed while it holds the lock; `golang.org/x/sys` becomes a direct dependency | `flai/internal/filelock/*.go` (4 files), `flai/go.mod`, `design/tech/go-libraries.md` | none |
| 2 | T-1511 `lockActivity` uses it; `TestActivityWritersTakeTurns` passes 50 of 50 runs under `-race`; one sentence changes in `design/system/strategic-agents.md` | `activity.go`, `activity_test.go`, `strategic-agents.md` | T-1510 |
| 2 | T-1512 the guard's `lockRecord` and serve's `lockFile` use it, with tests | `subagents.go`, `subagents_test.go`, `restart.go`, `restart_test.go` | T-1510 |
| 3 | T-1513 `flai issue close I-0140`; tick both criteria | the I-0140 file, `design/issues/summary.md` (shared path) | T-1511, T-1512 |

T-1511 and T-1512 share no file, so they run together.

## Figures

- **Forecast: 1h, raised from flai's 5m.** flai sized the story from two criteria and no touches. The real work is a package for two platforms, three call sites, a crash test, and a close-out.
- **Delivery: 2026-10-13T19:05Z.**
- **Cost of delay: 37.50 USD a week.** This is flai cod's figure from flai's 15m time-lost input, which I left unchanged.
- Every touch names a file. No folder touch is kept.

## Assumptions

- **The lock file now stays in place, unlocked.** Removing it after unlocking would open a new race with a waiter that already opened it. A lock file an older flai left behind is harmless: the new lock ignores its age and its presence.
- **No ADR.** I treat this as an implementation fix inside three writers. The design sentence in `strategic-agents.md` is the only recorded behaviour that changes.
- **T-1512 widens the story beyond I-0140's one function.** It removes the same cause elsewhere. If you want the story kept to `lockActivity` alone, say so: I would drop T-1512, and raise a new issue for the other two locks.

### 2026-10-13T16:48:58Z orchestrator
Approved, T-1512 included. Criterion 1 asks that the cause I-0140 describes no longer occur. That cause is the Stat-then-Remove age break, and `lockRecord` and `lockFile` carry it too. Fixing it in one package keeps the goal and the criteria as they are, so it is no change of scope. T-1511 reproduces I-0140 with `TestActivityWritersTakeTurns` under `-race` (criterion 1), and T-1513 closes it (criterion 2). Every touch names a file, and T-1511 and T-1512 share none. The figures are flai's: 1h (5m raised, with a reason) and 37.50 USD/week from flai's 15m input. `golang.org/x/sys` is already in the module and becomes direct, which T-1510 records in design/tech/go-libraries.md. Resolving.

### 2026-10-13T16:48:59Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories.
