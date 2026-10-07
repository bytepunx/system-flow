---
id: S-0314
type: story
nature: remediation
title: The flaiover guide says Check for updates always works, but the Dashboard area hides it while the dashboard host action is off
status: done
owner: alex
created: 2026-10-07T18:59:44Z
updated: 2026-10-07T23:47:54Z
transitions:
  - to: ready
    at: 2026-10-07T23:42:28Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T23:43:18Z
    by: agent-S-0314
  - to: review
    at: 2026-10-07T23:47:16Z
    by: agent-S-0314
  - to: done
    at: 2026-10-07T23:47:54Z
    by: orchestrator
tags: [dashboard]
touches: [flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, design/system/flaiover-dashboard.md, design/issues/I-0115-the-flaiover-guide-says-check-for-updates-always-works-but-the-dashboard-area-hides-it-while-the-dashboard-host-action-is-off.md, design/issues/summary.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 246
  turns:
    - day: 2026-10-07
      ceremony: 2
      hand_edits: 1
      work: 29
  models:
    - model: claude-opus-5-5
      input: 66
      output: 11199
      cache_read: 3262863
      cache_write: 127672
      cost: 1.8982
  strategic:
    - kind: orchestrator
      seconds: 129
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 32
          output: 462
          cache_read: 6360982
          cache_write: 18111
          cost: 1.5717
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-07T18:59:44Z
  value: 12.5
  by: planner-S-0314
  at: 2026-10-07T23:41:35Z
forecast:
  duration: 12m
  delivery: 2026-10-08T06:25:00Z
  basis: "flai forecast: median 101 s per unit of size over 16 done remediation stories on claude-opus-5-5 in the medium band, times size 7 (2 criteria, 5 touches); 21st in the pull order with an in-progress limit of 3."
  by: planner-S-0314
  at: 2026-10-07T23:41:35Z
finalized:
  by: orchestrator
  at: 2026-10-07T23:42:23Z
---
# S-0314 The flaiover guide says Check for updates always works, but the Dashboard area hides it while the dashboard host action is off

## Goal

This story remediates [I-0115](../../../design/issues/I-0115-the-flaiover-guide-says-check-for-updates-always-works-but-the-dashboard-area-hides-it-while-the-dashboard-host-action-is-off.md), "The flaiover guide says Check for updates always works, but the Dashboard area hides it while the dashboard host action is off". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0115 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0115 is closed with `flai issue close I-0115 --reason` saying what fixed it

## Tasks
- T-1290 The Dashboard area shows Check for updates while the dashboard host action is off
- T-1291 The flaiover design says Check for updates shows whatever the dashboard host action says
- T-1292 I-0115 is closed with the reason that the Dashboard area now shows Check for updates with the action off

## Notes

Cost of delay inputs set by flai from I-0115. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T15:11:31Z, 0.2 days before this story; under one cycle counts as one).

### Planning

Proposed fix: show Check for updates in the Dashboard area whatever the `dashboard` host action says, as Versions already is, so the code matches `docs/users/flaiover.md` § Host. `dashboard.check` is a read the host API does not gate (`flai/internal/hostapi/writes.go`, `read(...)`; `design/system/dashboard-host-channel.md` item 10), and `POST /api/dashboard` with `check` runs it with `repo().run`, not a write, so no host or route change is needed. The guide stays as it is; the design's sentence that the page shows "once enabled, the four actions" is corrected instead.

Touches, all files, no folder touch:

| Touch | From |
|-------|------|
| `flaiover/src/lib/components/HostPanel.svelte` | layout: the `{#if status.dashboard_enabled}` around the button, named in I-0115 |
| `flaiover/src/lib/components/HostPanel.svelte.test.ts` | co-change (3 of 4 commits with HostPanel.svelte); its "action off" test expects only Versions today |
| `design/system/flaiover-dashboard.md` | co-change (2 of 4) and design: § Managing the dashboard from its own page says "once enabled, the four actions" |
| `design/issues/I-0115-…-is-off.md` | criterion 2: `flai issue close` writes it |
| `design/issues/summary.md` | criterion 2: `flai issue close` regenerates it |

Left out: `HostProcesses.svelte` and its test (co-change 2 of 4), whose Check for upgrade is already outside its `host_enabled` gate; `docs/users/flaiover.md`, which already says what the fix makes true; `flaiover/src/routes/api/dashboard/+server.ts`, which already runs `check` ungated.

Forecast 12m, delivery 2026-10-08T06:25:00Z: `flai forecast` as it stands, rerun after the touches were set (it gave 5m on 0 touches). Kept: the change is one gate, one test, one design sentence, and a mechanical close, which the 16 medium-band remediations it rests on fit.

Cost of delay value 12.50 USD a week: `flai cod` from the inputs flai set from I-0115 (5m lost per 168h cycle at 150 USD an hour). Kept: one occurrence, and the defect misleads but blocks nothing.

### Accepted by the orchestrator

- Verified: 8a90ac5f53e8d101e8f74944efe48951407a3828
- At: 2026-10-07T23:47:54Z

Verdict: accept; both criteria met (verifier at 8a90ac5f53e8d101e8f74944efe48951407a3828; flai verify passed every step at that commit). The I-0104 bump on the branch is the close-out's own record.

- 1: flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, design/system/flaiover-dashboard.md
- 2: design/issues/I-0115-the-flaiover-guide-says-check-for-updates-always-works-but-the-dashboard-area-hides-it-while-the-dashboard-host-action-is-off.md, design/issues/summary.md
