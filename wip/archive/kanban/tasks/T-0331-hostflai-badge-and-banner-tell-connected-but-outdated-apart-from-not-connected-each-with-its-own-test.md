---
id: T-0331
type: task
nature: feature
title: HostFlai badge and banner tell connected-but-outdated apart from not-connected, each with its own test
status: done
parent: S-0086
owner: alex
created: 2026-09-22T23:46:31Z
updated: 2026-09-22T23:50:57Z
transitions:
  - to: ready
    at: 2026-09-22T23:46:42Z
    by: system-flow
  - to: in-progress
    at: 2026-09-22T23:46:44Z
    by: system-flow
  - to: done
    at: 2026-09-22T23:50:57Z
    by: system-flow
stream: S-0086
tags: []
---

# T-0331 HostFlai badge and banner tell connected-but-outdated apart from not-connected, each with its own test

## Work
`HostFlaiBanner.svelte` and `HostFlai.svelte` both branch on `hostFlai.usable` (`connected && !error`) to choose their message, so a flai that is genuinely connected but missing required methods (`status.connected: true`, `status.error` set, from `/api/agent`'s own already-correct payload) reads identically to no flai ever having connected at all. Branch on `status.connected` itself instead, keeping the existing `!status.configured` branch as is:
- Banner: `!status.configured` (unchanged) → `status.connected` with an error → a new paragraph using `status.error` verbatim (it already carries the version, what is missing, and the fix commands) → the existing not-connected paragraph only when `!status.connected`.
- Badge (`HostFlai.svelte`): a third `data-host-flai` state (e.g. `"outdated"`) for `status.connected && status.error`, its own label (not "not connected"), `title` still `status.error`.

## Done when
The four acceptance criteria: the connected-but-outdated case reads its own message with version/missing/fix commands, never "not connected"; the true not-connected and never-configured cases are unchanged; the header badge tells all three apart; a test for each of the three states (plus the existing usable/connected-and-fine case). `pnpm run check`/`test:unit`/`lint` clean.

## Notes
`status.error` for this case is already built server-side in `src/routes/api/agent/+server.ts` and needs no change.
