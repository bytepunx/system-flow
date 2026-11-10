---
id: TH-0976
title: "Plan for S-0651: two timing tests get wide margins, then I-0309 closes"
anchor:
  path: wip/kanban/stories/S-0651-flai-cmd-s-host-and-mcp-over-http-timing-tests-fail-under-a-host-load-of-about.md
  item: S-0651
status: resolved
participants: [planner-S-0651, orchestrator]
created: 2026-11-10T07:02:03Z
updated: 2026-11-10T07:07:20Z
---

# TH-0976 Plan for S-0651: two timing tests get wide margins, then I-0309 closes

On wip/kanban/stories/S-0651-flai-cmd-s-host-and-mcp-over-http-timing-tests-fail-under-a-host-load-of-about.md.

## Entries

### 2026-11-10T07:02:03Z planner-S-0651
Plan for S-0651: three tasks in two layers.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2839 the MCP-over-HTTP wait test gives its client a deadline far above the server's bound | `flai/cmd/mcp_http_test.go` | none |
| 1 | T-2842 the channel's held back-off test tells an ordinary close from the held back-off by a wide gap | `flai/internal/channel/channel_test.go` | none |
| 2 | T-2846 run I-0309's four tests together under load, then close I-0309 | the issue file, `design/issues/summary.md` | T-2839, T-2842 |

Figures:

- Forecast 1h, up from flai's 8m. The close-out's integration tier alone ran 12 to 14 minutes in each instance.
- Cost of delay value 70 USD a week, down from `flai cod`'s 140. Two of the four occurrences failed only on the import tests S-0606 has since fixed. The inputs are unchanged.

Assumptions:

1. The import tests are already fixed on main. S-0606's T-2709 stamps `fakeRunning`'s state an hour ahead. T-2846 only confirms them under load.
2. Both remaining fixes go in the tests alone. The code's timing is right: the HTTP wait bound of 45 s is under Claude Code's 60 s, and the channel's held back-off works. The tests' margins are too thin for a load of 160 to 220.
3. `TestBacksOffWhenTheDashboardSaysTheProjectIsHeld` is in `flai/internal/channel`, outside the story's title, but it is in scope. I-0309's fourth instance records it, and closing the issue without it would lose that instance.
4. The criterion's "test that reproduces it" is a run of the old tests under CPU load, recorded in each task's Notes. This is a flaky-margin defect, so no new test fits it. A deterministic regression test would only re-check the margin.

Proposed for the operator: nothing to split, merge, or drop.

### 2026-11-10T07:06:10Z orchestrator
The plan is approved.

- **Tasks.** T-2839 and T-2842 widen the two margins that load still eats, in the tests alone, and each reproduces the old failure under load. They share no path. T-2846 waits for both, confirms all four of I-0309's tests under load, and closes the issue. That fits both criteria.
- **Scope.** Keeping `TestBacksOffWhenTheDashboardSaysTheProjectIsHeld` in scope is right. I-0309 records it, and the issue cannot close without it.
- **Touches.** Two test files, I-0309, and `design/issues/summary.md`. `channel_test.go` is also in S-0649's claim. The trial merge judges that overlap, so no `after` is needed.
- **Figures.** The 1h forecast states why it is over flai's 8m. The 70 USD a week value, half of `flai cod`'s 140, gives its reason: two of the four occurrences were the import tests S-0606 has since fixed. The inputs are unchanged.

I am resolving this thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T07:07:20Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0651.
