---
id: T-1290
type: task
nature: remediation
title: The Dashboard area shows Check for updates while the dashboard host action is off
status: done
parent: S-0314
owner: alex
created: 2026-10-07T23:41:52Z
updated: 2026-10-07T23:44:48Z
transitions:
  - to: ready
    at: 2026-10-07T23:43:46Z
    by: agent-S-0314
  - to: in-progress
    at: 2026-10-07T23:43:46Z
    by: agent-S-0314
  - to: done
    at: 2026-10-07T23:44:48Z
    by: agent-S-0314
stream: S-0314
tags: [dashboard]
touches: [flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, design/system/flaiover-dashboard.md]
usage:
  source: log
  seconds: 62
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 8
      output: 84
      cache_read: 413367
      cache_write: 3479
      cost: 0.1838
---
# T-1290 The Dashboard area shows Check for updates while the dashboard host action is off

## Work

In `flaiover/src/lib/components/HostPanel.svelte`, move the Check for updates button (`data-testid="host-panel-check"`) out of the `{#if status.dashboard_enabled}` block, so it shows whatever the action says, beside Versions, as `docs/users/flaiover.md` § Host promises. `dashboard.check` is a read the host API does not gate, and the route runs it with `repo().run`, so nothing else changes. Leave Restart, Upgrade, and Stop gated. While the action is off, the check's "Not running; Restart starts it." points at a button that is hidden: word that case so it does not name Restart while the action is off. Update the header comment's "Restart / Upgrade / Stop, gated" sentence if it no longer reads true.

In `HostPanel.svelte.test.ts`, change the action-off test, which expects only `['Versions']` today, to expect Check for updates and Versions, and add a test that clicks Check for updates with the action off and sees the result (`host-panel-check-result`) and a `POST` with `{ action: 'check' }`. That test fails before the change and passes after it: it reproduces I-0115.

Waits for nothing: first layer, beside the design task, with which it shares no path.

## Done when

- With `dashboard_enabled` false, the Dashboard area shows Check for updates and Versions, and neither Restart, Upgrade, nor Stop
- A check with the action off shows its result and names no hidden button
- The new test fails on the old gate and passes now
- `flai test flaiover/src/lib/components/HostPanel.svelte flaiover/src/lib/components/HostPanel.svelte.test.ts` passes

## Notes
