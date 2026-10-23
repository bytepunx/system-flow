---
id: TH-0622
title: "Plan for S-0488: narrow the test's leak check to the registered entry's URL, then close I-0197"
anchor:
  path: wip/kanban/stories/S-0488-testservesettingsareshowntothedashboard-in-flai-cmd-fails-now-and-then-in-the-close-out-s-go-test-tier-on-a-branch-that-does-not-touch-serve.md
  item: S-0488
status: resolved
participants: [planner-S-0488, orchestrator]
created: 2026-10-23T18:11:17Z
updated: 2026-10-23T18:11:27Z
---

# TH-0622 Plan for S-0488: narrow the test's leak check to the registered entry's URL, then close I-0197

On wip/kanban/stories/S-0488-testservesettingsareshowntothedashboard-in-flai-cmd-fails-now-and-then-in-the-close-out-s-go-test-tier-on-a-branch-that-does-not-touch-serve.md.

## Entries

### 2026-10-23T18:11:17Z planner-S-0488
Cause, confirmed in `flai/cmd/serve_actions_test.go` line 980: the test refuses any `4242` in the settings JSON, which carries two `t.TempDir()` paths (the project root and the import root the test adds); Go's random temp-directory suffix can contain `4242`, so the test fails on an unlucky run and passes on the next. `hostSettings` is right to show those paths, so only the test changes.

Tasks, two layers:

1. T-2087 (layer 1): the check refuses `127.0.0.1:4242`, the entry's URL as registered, instead of the bare port, keeping the `"token"` and `/p/harbour/k` checks; the import root goes under a folder named with `4242`, so every run carries the string the old check refused. That case is the reproducing test: it fails under the old check on every run and passes under the new one.
2. T-2088 (layer 2, after T-2087): `flai issue close I-0197 --reason` naming that fix, which regenerates `design/issues/summary.md`; then the criteria are ticked.

Enrichment: touches are the test file and the two issue files; the co-change list (`serve_actions.go`, `hostapi/writes.go`, and the serve, config, and settings sources) is left out since no source changes. Forecast raised from flai's 5m to 15m for the go-test tier and the close-out around two tasks; delivery 2026-10-24T12:12Z. Cost of delay value 17.50 USD a week, straight from the inputs flai set.

Assumptions: no ADR, design, or docs change, since no behaviour changes; no story in `after`. Nothing to split, merge, or drop.

### 2026-10-23T18:11:24Z orchestrator
Approved: the plan fits S-0488.

- **Cause.** Confirmed in `serve_actions_test.go`: a `t.TempDir()` suffix can contain `4242`. The fix narrows the check to `127.0.0.1:4242`. It adds a case whose import root carries `4242` on every run, so the test proves the fix.
- **Tasks.** Two, in two layers: T-2087 the test, then T-2088 to close I-0197.
- **Touches.** Three files.
- **Forecast.** 15m, raised from flai's 5m for the go-test tier and the close-out. The story's agent waits on both in its session, so that is agent time.
- **Cost of delay.** 17.50 USD a week, from `flai cod`.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T18:11:27Z orchestrator
Resolved: Plan approved under plan_backlog_stories.
