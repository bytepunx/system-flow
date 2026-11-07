---
id: T-2581
type: task
nature: remediation
title: The Harnesses settings page edits this project's command and harness programs, and changes every project only when the operator chooses everywhere
status: done
parent: S-0587
owner: alex
created: 2026-11-07T10:12:25Z
updated: 2026-11-07T11:41:12Z
transitions:
  - to: ready
    at: 2026-11-07T10:57:41Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:57:45Z
    by: system-flow
  - to: done
    at: 2026-11-07T11:24:06Z
    by: agent-S-0587
stream: S-0587
tags: [dashboard]
touches: [flaiover/src/lib/settings.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/writes.test.ts]
after: [T-2580]
usage:
  source: log
  seconds: 1581
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 94
      output: 27707
      cache_read: 18077507
      cache_write: 156338
      cost: 5.1229
---
# T-2581 The Harnesses settings page edits this project's command and harness programs, and changes every project only when the operator chooses everywhere

## Work

It waits for T-2580, which makes `settings.agent` and `settings.harness` per project and adds the request's "everywhere" choice that this page sends.

- In `flaiover/src/lib/settings.ts`, take `agent` and `harness` out of `HOSTWIDE`. They are now the project's, with a host-wide fallback. `provider` stays host-wide.
- In `flaiover/src/lib/components/SettingsPanel.svelte` (`/settings/harnesses`):
  - Show whether each value is this project's or the host-wide fallback.
  - Save to this project by default.
  - Offer an explicit "every project" choice that needs `settings` on everywhere.

## Done when

- `SettingsPanel.svelte.test.ts` shows a save on the Harnesses page sent for this project alone, and the "every project" choice sent host-wide.
- `flai test` on the changed paths passes the `flaiover-lint`, `vitest`, and `svelte-check` tiers.

## Notes
