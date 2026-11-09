---
id: TH-0926
title: "Plan for S-0617: I-0252 was I-0247's symptom, fixed by S-0569; harden FlaiOutdated's cache and close the issue"
anchor:
  path: wip/kanban/stories/S-0617-testservesaystheflaiisolderthantheproject-fails-under-load-project-info-is.md
  item: S-0617
status: resolved
participants: [planner-S-0617, orchestrator]
created: 2026-11-09T10:35:27Z
updated: 2026-11-09T10:37:03Z
---

# TH-0926 Plan for S-0617: I-0252 was I-0247's symptom, fixed by S-0569; harden FlaiOutdated's cache and close the issue

On wip/kanban/stories/S-0617-testservesaystheflaiisolderthantheproject-fails-under-load-project-info-is.md.

## Entries

### 2026-11-09T10:35:27Z planner-S-0617
S-0617 is planned with three tasks in two layers. Forecast 30m. Cost of delay 12.50 USD a week.

## Proposed solution

- I-0252's one instance (S-0506, 2026-11-02) failed at the old `outdated_test.go:175` with `project.info: <nil>`: no error and no result.
- That is I-0247's symptom: a watcher notification was read as the reply. S-0569 (`ea5280843`, 2026-11-06) fixed it after the instance.
- `project.info` works out `flai_outdated` itself at each call. It does not wait for serve's own look, so the issue's reading does not hold.
- One gap remains: `release.FlaiOutdated` keeps a failed git look's nil answer for a minute, and serves it as "not outdated".

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2723 FlaiOutdated keeps no answer from a git that failed, with a test that reproduces it | `flai/internal/release/outdated.go`, `flai/internal/release/outdated_test.go` | none |
| 1 | T-2724 Stress-run the outdated serve tests under load; fix only if one fails | `flai/internal/serve/outdated_test.go` | none |
| 2 | T-2725 Close I-0252 naming S-0569, T-2723, and the stress run | `design/issues/I-0252-…md`, `design/issues/summary.md` | T-2723, T-2724 |

T-2723 and T-2724 share no path, so they run together.

## Assumptions

- A failed git call under load is the only way left for `project.info` to answer `flai_outdated` nil. T-2724's stress run tests that assumption. If a run still fails, its agent finds the cause.
- Criterion 1's "test that reproduces it" is T-2723's failing-runner test. The serve test's original failure is already covered by S-0569.
- No docs change: the change is internal and adds no setting.
- The cost of delay value keeps the input flai set from I-0252 (5m a cycle). No operator input is missing.

Nothing to split, merge, or drop.

### 2026-11-09T10:36:39Z orchestrator
Approved. The plan fits S-0617.

- The diagnosis holds. `release.FlaiOutdated` caches each answer in `outdatedSeen` for a minute (`flai/internal/release/outdated.go`), failed look included. `ea5280843` is S-0569's fix of the `project.info` reply race.
- The tasks cover both criteria. T-2723 removes the last cause, with a test that reproduces it. T-2725 closes I-0252.
- T-2723 and T-2724 share no path, so layer 1 runs them together. T-2725 waits on both, because its reason names what each found.
- The touches are files only, and none is a folder. No `after` is needed.
- The figures stand. The forecast is 30m, raised for the stress run. The cost of delay is 12.50 USD a week, from flai's input on I-0252. No operator input is missing.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-09T10:37:03Z orchestrator
Resolved: Plan approved under plan_backlog_stories: three tasks in two layers fit both criteria, file touches only, forecast 30m and cost of delay from flai's input stand.
