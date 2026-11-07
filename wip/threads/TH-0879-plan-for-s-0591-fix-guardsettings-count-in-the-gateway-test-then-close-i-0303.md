---
id: TH-0879
title: "Plan for S-0591: fix guardSettings' count in the gateway test, then close I-0303"
anchor:
  path: wip/kanban/stories/S-0591-testgatewaylitellm-fails-its-setup-check-the-template-s-settings-run-flai-guard.md
  item: S-0591
status: resolved
participants: [planner-S-0591, orchestrator]
created: 2026-11-07T21:48:34Z
updated: 2026-11-07T21:48:57Z
---

# TH-0879 Plan for S-0591: fix guardSettings' count in the gateway test, then close I-0303

On wip/kanban/stories/S-0591-testgatewaylitellm-fails-its-setup-check-the-template-s-settings-run-flai-guard.md.

## Entries

### 2026-11-07T21:48:34Z planner-S-0591
The plan for S-0591 has two tasks in two layers.

Cause: `guardSettings` in `flai/internal/serve/gateway_test.go` counts every `flai guard` substring in `template/root/.claude/settings.json`. S-0518 (ece3f9e69) added the message `flai guard did not run (exit %s)` to the two refusing hooks, so it counts 4 invocations plus 2 messages: 6. The settings are right; the test is wrong. `TestGatewayOpenRouter` shares `checkGateway`, so it fails the same way.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2592: count and rewrite only the hooks' invocations, with a test that `-short` runs on the template's real settings | `flai/internal/serve/gateway_test.go` | none |
| 2 | T-2593: close I-0303 with a reason naming the fix | the I-0303 file and `design/issues/summary.md` | T-2592 |

Figures:

- Forecast 12m, flai's figure kept. It comes from 16 done remediation stories at size 5.
- Cost of delay 12.50 USD a week, from the 5m input. Kept.
- No `after` on another story.

Assumptions:

- The fix goes in the test, not the settings. The message is the guard's own error text, which `flai/cmd/guard_test.go` checks. Changing it would also need a template changelog line and a release.
- The story's agent does not run the paid `gateway` tier. The new `-short` test reproduces the defect without a model call. A real gateway run stays the operator's, through `flai test gateway` or its schedule.

Nothing to split, merge, or drop.

### 2026-11-07T21:48:46Z orchestrator
Approved. The two tasks cover both criteria. T-2592 fixes the test and adds a `-short` test on the template's real settings (criterion 1). T-2593 closes I-0303 (criterion 2).

- The touches match what each task writes.
- T-2593 waits for T-2592 so that the close reason can name the fix.
- The 12m forecast and the 12.50 USD a week value are flai's figures, kept.

Both assumptions stand: the fix goes in the test, and the paid `gateway` tier stays the operator's to run. I am finalizing S-0591 next.

### 2026-11-07T21:48:57Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: tasks, touches, and figures fit S-0591.
