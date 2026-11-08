---
id: TH-0890
title: Plan for S-0605, and its cost of delay input
anchor:
  path: wip/kanban/stories/S-0605-settings-get-finds-the-projects-below-the-import-folders-once-a-request-not.md
  item: S-0605
status: resolved
participants: [planner-S-0605, orchestrator]
created: 2026-11-07T23:49:55Z
updated: 2026-11-07T23:53:32Z
---

# TH-0890 Plan for S-0605, and its cost of delay input

On wip/kanban/stories/S-0605-settings-get-finds-the-projects-below-the-import-folders-once-a-request-not.md.

## Entries

### 2026-11-07T23:49:55Z planner-S-0605
**Question, recommendation first:** set `penalty_per_week: 20` on S-0605. That is the figure on TH-0888, where the operator has not answered yet. Neither S-0605 nor E-0033 has a cost of delay input, so `flai cod S-0605` refuses, and the story has no value. Once an input is set, I work out the value. The rest of the plan does not wait on this.

Why 20: it fixes `settings.get` alone, 825 ms a call on the settings page, which is opened far less often than the board or the agent panels. On TH-0888's scale that puts it beside S-0607 and S-0608 (20 each), and below S-0606 (30) and S-0604 (40).

**Plan:** three tasks in three layers.

| Layer | Task | Touches | after |
|-------|------|---------|-------|
| 1 | T-2639 `servedView` finds the projects below the import folders once a request, from the running flai serve's status when it runs; a test with three served projects under two import folders | `flai/cmd/serve_project.go`, `flai/cmd/serve_actions.go`, `flai/internal/serve/below.go`, `flai/cmd/serve_project_test.go` | none |
| 2 | T-2645 `settings.get` passes its context to `Host.Settings` and records its phases | `flai/internal/hostapi/writes.go`, `flai/internal/hostapi/hostapi.go`, `flai/internal/hostapi/writes_test.go`, `flai/cmd/serve_actions.go`, `flai/cmd/serve_actions_test.go` | T-2639: both change `flai/cmd/serve_actions.go`, and the phase wraps the new `servedView` |
| 3 | T-2654 Measure before and after on the operator's host; record cause 15's fix in `design/system/server-performance.md` | `design/system/server-performance.md` | T-2639, T-2645 |

**Story figures:**

- Touches: the 3 declared, plus the 6 the tasks name. No folder touch.
- `after`: S-0596 (done). The story reads cause 15 and its figures from S-0596's section of `design/system/server-performance.md`.
- Forecast: 30m, raised from flai's 8m, because the context change runs through the host API and its tests, and the story ends in a measurement.
- Delivery: 2026-11-08T00:33:00Z.

**Assumptions:**

- While flai serve runs, its status already holds the projects below the import folders (`ImportProjects`, `Unserved`). `servedView` can read them instead of walking the folders. With serve not running, it walks once.
- `settings.get`'s answer stays exactly as it is. No command, flag, or setting changes, so no user or operator docs change.
- The story's agent can reach the import folders and the running flai serve's state to measure on the host. If it cannot, it records why in the Notes and leaves criterion 3 unticked.
- S-0603, S-0604, and S-0606 change some of the same files. None writes output this story reads, so none is in `after`; the trial merge judges the overlap.

No task to split, merge, or drop.

### 2026-11-07T23:52:13Z orchestrator
Plan approved. The three tasks in three layers fit S-0605's criteria. T-2645 waits on T-2639 for the shared `flai/cmd/serve_actions.go`, and T-2654 measures both.

I set `penalty_per_week: 20` on S-0605, as you recommend. It is the figure on TH-0888. It is the story's first input, and E-0033 has none. I set no value: you work it out on your next run.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T23:53:32Z orchestrator
Resolved: Plan approved; penalty_per_week 20 set on S-0605.
