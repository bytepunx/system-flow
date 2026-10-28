---
id: T-2281
type: task
nature: improvement
title: flaiover-dashboard.md and the flaiover guide say the three agent pages show their settings first, under the host action
status: done
parent: S-0533
owner: alex
created: 2026-10-28T21:14:47Z
updated: 2026-10-28T21:20:52Z
transitions:
  - to: ready
    at: 2026-10-28T21:17:59Z
    by: agent-S-0533
  - to: in-progress
    at: 2026-10-28T21:18:01Z
    by: agent-S-0533
  - to: done
    at: 2026-10-28T21:20:52Z
    by: agent-S-0533
stream: S-0533
tags: [flaiover, docs, settings]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
usage:
  source: log
  seconds: 171
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 7653
      cache_read: 2675555
      cache_write: 71182
      cost: 1.2578
---
# T-2281 flaiover-dashboard.md and the flaiover guide say the three agent pages show their settings first, under the host action

## Work

- In `design/system/flaiover-dashboard.md`, under `## Strategic agents' pages (S-0259, S-0228)`, name the `settings` snippet among those the pages give the panel. Add a Settings row to the sections table after The host action, saying it holds `StrategicSettings.svelte` for the agent's block, or the notice when `settings.get` has no `host.strategic`, and that it comes before Running now. Bump `updated`.
- In `docs/users/flaiover.md`, the `## Planner`, `## Orchestrator`, and `## Analyzer` sections each say the settings come "below that". Say instead that the settings come first, under the host action's state (and Stop and Start on the Orchestrator page), above the run, the orchestrator's last decisions, the activity log, and the runs. Reorder each section's list to match the page.
- Leave `docs/operators/settings.md` as it is: it names the pages but not their order.

It waits for nothing: it shares no path with the code tasks, and it describes the order the story's criteria fix.

## Done when

- [ ] `flaiover-dashboard.md`'s section on the strategic agents' pages lists the settings after the host action and before Running now
- [ ] `docs/users/flaiover.md`'s Planner, Orchestrator, and Analyzer sections describe the settings first
- [ ] The markdown lint passes on both files

## Notes
