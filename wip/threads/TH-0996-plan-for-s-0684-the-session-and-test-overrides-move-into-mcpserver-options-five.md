---
id: TH-0996
title: "Plan for S-0684: the session and test overrides move into mcpserver.Options, five tasks in four layers, 2h"
anchor:
  path: wip/kanban/stories/S-0684-flai-internal-mcpserver-s-tests-still-take-136-to-162-s-under-race-since-the.md
  item: S-0684
status: resolved
participants: [planner-S-0684, agent-S-0684, orchestrator]
created: 2026-11-10T08:49:33Z
updated: 2026-11-10T08:54:04Z
---

# TH-0996 Plan for S-0684: the session and test overrides move into mcpserver.Options, five tasks in four layers, 2h

On wip/kanban/stories/S-0684-flai-internal-mcpserver-s-tests-still-take-136-to-162-s-under-race-since-the.md.

## Entries

### 2026-11-10T08:49:33Z planner-S-0684
The plan takes I-0234's first remedy, which D-0443 also names. The session's story, role, session, reasons, and started-by go into `mcpserver.Options`, which `flai mcp` fills from the environment. The package's tests pass them instead of setting the process, so the serial phase of about 70 s goes. Priming faster under race is left out, because T-2420 measured no gain.

Tasks and layers:

| Layer | Task | What |
|-------|------|------|
| 1 | T-2904 | Adds `Options` session fields and a `SessionFromEnv` helper, used by `flai mcp` over stdio and HTTP. `newServer`, `servedSession`, and `recordingStory` read the fields. Updates the cmd and serve tests that set the variables. Adds `session_test.go` |
| 2 | T-2905 (after T-2904, which also changes `server.go`) | Adds `Options` fields for what the tests now swap: `permissionPoll`, `flaiExecutable`, `buildinfo.Version`, `FLAI_RELEASES_API` and its tokens, and the child's `FLAI_CONFIG` |
| 2 | T-2906 (after T-2904) | Makes the session tests, 12 test files, top-level and parallel. Their serial parents go. Git's identity is set once in `TestMain` |
| 3 | T-2910 (after T-2905 and T-2906) | Makes the permission, subproject, versions, and outdated tests parallel through the new fields. Removes `together` and `inEnv` |
| 4 | T-2915 (after T-2906 and T-2910) | Adds `process_test.go`, which fails on any test that sets the process or is not parallel. Measures `-race` three times with the load, closes I-0234 |

Figures:

- Forecast: raised from flai's 56m to 2h, because of about 90 tests to convert and repeated `-race` runs of 2 to 3 minutes each. Delivery is 2026-11-10T15:00Z.
- Cost of delay: 1.79 USD a week, `flai cod`'s value from the operator's 43 s input, unchanged.
- `after`: none.

Assumptions:

- `flai mcp` behaves exactly as before. Over HTTP the helper reads the listening process's environment, as `New` does today. Per-session story and role over HTTP is out of scope.
- The issue's count of 95 is stale since S-0548. Today 36 serial tests and 6 serial parents of 54 subtests remain.
- The two prime-on-this-repository tests stay as they are. The package's `-race` time then falls to roughly theirs, about 55 to 70 s. T-2915 measures it and claims no figure in advance.
- No user or design documentation changes, because no command's behaviour changes.

I propose no split, merge, or drop. Nothing here waits on the operator. The story stays a draft for the operator to finalize.

### 2026-11-10T08:51:10Z agent-S-0684
Plan accepted as drafted: five tasks in four layers. One change in how they are checked. The hand `go test -race` runs in T-2906's and T-2910's done-when are met with `flai test` on the package. T-2915's timed runs stay, as the measurement the story needs. Both are recorded in the narrative's Decisions. Starting T-2904.

### 2026-11-10T08:52:10Z planner-S-0684
Second plan, asked for after the 08:48:56Z touches edit logged as by system-flow. The plan stands.

What the edit did: it added `flai/internal/mcpserver/task.go` to the story's touches. The first plan's Notes already name that file, and T-2904 touches it. Nothing is missing.

State now: alex finalized the story at 08:50Z, and agent-S-0684 has it in progress. T-2904 is in progress. T-2905, T-2906, T-2910, and T-2915 are in the backlog.

Revisited, unchanged:

| Layer | Task | after |
|-------|------|-------|
| 1 | T-2904 | none |
| 2 | T-2905, T-2906 | T-2904 |
| 3 | T-2910 | T-2905, T-2906 |
| 4 | T-2915 | T-2906, T-2910 |

Figures checked again:

- `flai touches suggest`: its top candidate is still `folder.go`. It stays out, because it reads no session variable.
- `flai forecast`: 58m from size 34. The 2h stands for the reasons in the Notes.
- `flai cod`: 1.79 USD a week. It stands.

Changed: only the story's `### Planning` notes. They now record the touches edit, the second plan's figures, and the delivery that flai's replan set. The first plan's 15:00Z delivery was stale.

I propose no new task, split, merge, or drop.

### 2026-11-10T08:53:28Z orchestrator
Approved as written: T-2904, then T-2905 and T-2906, then T-2910, then T-2915; 2h, 1.79 USD a week.

- The plan takes I-0234's first remedy, as D-0443 names it. The session moves into `mcpserver.Options`, so the tests stop setting the process and the serial phase goes. That is criterion 1, and T-2915's `process_test.go` keeps the cause from coming back.
- Leaving out faster priming under race is right: T-2420 measured no gain.
- `flai mcp` keeps its behaviour, so no user or design document changes.
- The layers follow what each task reads, and T-2915 measures the `-race` time rather than claiming a figure in advance.

### 2026-11-10T08:54:04Z orchestrator
Resolved: Plan approved as written by the orchestrator under plan_backlog_stories.
