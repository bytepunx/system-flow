---
id: T-1251
type: task
nature: feature
title: The dashboard's connection list says per project whether its flai is signed or was refused
status: done
parent: S-0237
owner: alex
created: 2026-10-07T23:00:52Z
updated: 2026-10-12T17:35:40Z
transitions:
  - to: ready
    at: 2026-10-12T17:21:46Z
    by: agent-S-0237
  - to: in-progress
    at: 2026-10-12T17:21:47Z
    by: agent-S-0237
  - to: done
    at: 2026-10-12T17:35:40Z
    by: agent-S-0237
stream: S-0237
tags: [dashboard]
touches: [flaiover/src/routes/api/projects/+server.ts, flaiover/src/routes/api/projects/projects.test.ts, flaiover/src/lib/components/ProjectSwitcher.svelte, flaiover/src/lib/components/ProjectSwitcher.svelte.test.ts, flaiover/src/lib/settings.ts]
after: [T-1248]
usage:
  source: log
  seconds: 833
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 80
      output: 401
      cache_read: 3552163
      cache_write: 96599
      cost: 1.6533
---
# T-1251 The dashboard's connection list says per project whether its flai is signed or was refused

## Work

The second half of criterion 3. It waits for T-1248, which marks each connection signed and keeps each refusal. It runs beside the banner task, with no path in common.

- `flaiover/src/routes/api/projects/+server.ts`: `ProjectGlance` gains whether the project's flai is signed, and the refusal when the hub refused it. A refused project is listed even though the registry never adopted it.
- `flaiover/src/lib/components/ProjectSwitcher.svelte`: each project's label says "signed" or "unsigned" when connected. A refused project says "(refused: unsigned flai)" and appears in the "not connected: why?" list with the reason.
- `flaiover/src/lib/settings.ts`: the Settings page's project state shows the same, for a project flai reports as not connected because the dashboard refused it.
- Tests: `projects.test.ts` for a signed, an unsigned, and a refused project; `ProjectSwitcher.svelte.test.ts` for their labels.

## Done when

- `/api/projects` and the project switcher say per project whether its flai is signed, and name a refused project with the reason.
- The tests pass with `flai test` on the changed paths.

## Notes

Drafted by the planner. If the Settings panel needs a change of its own (`SettingsPanel.svelte`), the story's agent widens the touches.
