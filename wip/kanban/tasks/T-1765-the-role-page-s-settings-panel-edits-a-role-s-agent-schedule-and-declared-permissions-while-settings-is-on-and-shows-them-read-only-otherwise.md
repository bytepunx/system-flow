---
id: T-1765
type: task
nature: feature
title: The role page's settings panel edits a role's agent, schedule, and declared permissions while settings is on and shows them read-only otherwise
status: done
parent: S-0420
owner: alex
created: 2026-10-16T05:46:50Z
updated: 2026-11-02T01:04:54Z
transitions:
  - to: ready
    at: 2026-11-02T00:59:35Z
    by: agent-S-0420
  - to: in-progress
    at: 2026-11-02T00:59:39Z
    by: agent-S-0420
  - to: done
    at: 2026-11-02T01:04:54Z
    by: agent-S-0420
stream: S-0420
tags: [dashboard]
touches: [flaiover/src/lib/components/StrategicSettings.svelte, flaiover/src/lib/components/StrategicSettings.svelte.test.ts, flaiover/src/lib/settings.ts, flaiover/src/lib/components/RolePage.svelte, "flaiover/src/routes/workflow/[role]/+page.svelte"]
after: [T-1763, T-1764]
usage:
  source: log
  seconds: 315
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 58
      output: 358
      cache_read: 2821017
      cache_write: 115918
      cost: 1.03
---
# T-1765 The role page's settings panel edits a role's agent, schedule, and declared permissions while settings is on and shows them read-only otherwise

## Work

- Generalize `StrategicSettings.svelte` from the three fixed panels to a role's panel, built from the keys the settings methods list for it.
- Keep the planner's, the orchestrator's, and the analyzer's panels exactly as they are, `currency` and `orchestration.agent` left out as today.
- Mount the panel on `RolePage`.

Waits for the settings methods task and for `RolePage`.

## Done when

`StrategicSettings.svelte.test.ts` passes, its existing cases unchanged, with a case for a fixture role.

## Notes
