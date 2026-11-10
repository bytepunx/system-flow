---
id: TH-0966
title: "Plan for S-0647: the dev server tests check the pid flai started, then I-0311 closes"
anchor:
  path: wip/kanban/stories/S-0647-teststartdevserverstopsvitewhenitscontextends-fails-under-load-the-fake-vite.md
  item: S-0647
status: resolved
participants: [planner-S-0647, orchestrator]
created: 2026-11-10T06:24:22Z
updated: 2026-11-10T06:26:10Z
---

# TH-0966 Plan for S-0647: the dev server tests check the pid flai started, then I-0311 closes

On wip/kanban/stories/S-0647-teststartdevserverstopsvitewhenitscontextends-fails-under-load-the-fake-vite.md.

## Entries

### 2026-11-10T06:24:22Z planner-S-0647
The plan for S-0647 has two tasks in two layers.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2825 The dev server tests check the pid flai started, and a slow-starting fake vite reproduces the missing pid | `flai/internal/components/devserver_test.go`, `flai/internal/components/devserver.go` | none |
| 2 | T-2826 Close I-0311 with what fixed it | `design/issues/I-0311-…md`, `design/issues/summary.md` | T-2825 |

Cause, from the code: the fake vite is the Go test binary, started again by a script that `exec`s it. It writes its pid only once that binary runs. Under heavy load it does not run within the context test's 300 ms, so flai stops it before it writes the pid. flai's own behaviour is correct. The test loses the pid it checks.

Figures:

- Forecast 40m, delivery 2026-11-10T08:02Z. `flai forecast` gave 5m; the time is raised for the integration close-out and for repeated test runs.
- Cost of delay 127.50 USD a week, as `flai cod` gives it from flai's 51m input.

Assumptions:

- Product code changes only through a test seam in `devserver.go`, beside `pickPort`, that hands the tests `cmd.Process.Pid`. The pid is the same because the script `exec`s the binary. If the story's agent finds the seam wrong, the fallback is a test-only fix: the script writes `$$` before its `exec`.
- A start delay in the fake reproduces the failure without load, so criterion 1's test fits.
- `TestStartDevServerStopsAViteThatDoesNotAnswerInTime`, at 500 ms, has the same window and is fixed in T-2825 too.
- No docs or changelog change, since no user-facing behaviour changes.
- No `after` on another story.

I would not split, merge, or drop anything.

### 2026-11-10T06:25:19Z orchestrator
Approved as planned. The plan fits S-0647.

- T-2825 reproduces I-0311 with a start delay before it fixes anything, which meets criterion 1's "a test that reproduces it". T-2826 closes I-0311 for criterion 2, after T-2825, because its reason names the fix.
- The touches name files only, and each task's touches sit within the story's.
- The fallback, the script writing `$$` before its `exec`, is a fair second way. The story's agent records why in Decisions if it takes it.
- Figures: 40m over flai forecast's 5m has its reason, the integration tier's close-out under load. The 127.50 USD a week value is `flai cod`'s, from the 51m input flai set from I-0311.

S-0638 also touches the I-0311 file. That is an occurrence it recorded, not a fix. A trial merge settles any overlap.

I finalize S-0647 next, under `finalize_drafts`.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T06:26:10Z orchestrator
Resolved: Plan approved under plan_backlog_stories: two tasks fit S-0647's criteria, touches, and flai's figures.
