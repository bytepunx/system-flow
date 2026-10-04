---
id: T-0825
type: task
nature: improvement
title: The settings page shows the planning triggers
status: done
parent: S-0211
owner: alex
created: 2026-10-04T20:42:11Z
updated: 2026-10-04T21:02:24Z
transitions:
  - to: ready
    at: 2026-10-04T20:42:48Z
    by: agent-S-0211
  - to: in-progress
    at: 2026-10-04T20:55:56Z
    by: agent-S-0211
  - to: done
    at: 2026-10-04T21:02:24Z
    by: agent-S-0211
stream: S-0211
tags: [flai, dashboard]
touches: [flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flaiover/src/lib/settings.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts]
after: [T-0823]
usage:
  source: log
  seconds: 388
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 49
      output: 18949
      cache_read: 3043210
      cache_write: 76912
      cost: 1.4491
---
# T-0825 The settings page shows the planning triggers

## Work

`settings.get`'s host view gains `planning`: whether the `plan` action is on here, so edits replan; `replan` with its default; `schedule` as written; and its next run in UTC when set. The dashboard's Settings page shows them in a Planning section, read-only, saying they are set in `system-flow.yaml` (ADR-0084). It waits for T-0823, whose keys it reads; it shares no path with T-0824 and runs beside it.

## Done when

- [ ] `settings.get` carries the planning triggers, with a test
- [ ] The Settings page shows them, with a component test
- [ ] flaiover's `npm run check` and tests for the panel pass

## Notes
