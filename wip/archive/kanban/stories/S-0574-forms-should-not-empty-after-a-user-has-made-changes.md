---
id: S-0574
type: story
nature: remediation
title: Forms should not empty after a user has made changes
status: done
owner: alex
created: 2026-11-06T17:55:11Z
updated: 2026-11-06T18:14:24Z
transitions:
  - to: ready
    at: 2026-11-06T17:55:18Z
    by: alex
  - to: in-progress
    at: 2026-11-06T17:55:24Z
    by: agent-S-0574
  - to: review
    at: 2026-11-06T18:13:02Z
    by: agent-S-0574
  - to: done
    at: 2026-11-06T18:14:24Z
    by: orchestrator
tags: [dashboard]
touches: [flaiover/src, design/system/flaiover-dashboard.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1069
  priced_by: harness
  turns:
    - day: 2026-11-06
      ceremony: 2
      hand_edits: 1
      work: 47
  models:
    - model: claude-haiku-5-5
      input: 68
      output: 31797
      cache_read: 2877122
      cache_write: 138325
      cost: 0.1928
    - model: claude-opus-5-5
      input: 170
      output: 40767
      cache_read: 17827639
      cache_write: 280281
      cost: 6.4995
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: alex
    at: 2026-11-06T17:55:11Z
---
# S-0574 Forms should not empty after a user has made changes

## Goal

Once a user begins entering data into form fields, they should not revert on a push from the server (or whatever other possible cause is underlying this behavior).

## Acceptance criteria
- [x] Controls do not revert to a previous value once the user begins making changes

## Tasks
- T-2463 A form helper keeps the fields the user changed when the data it was filled from is read again
- T-2464 A strategic agent's settings keep what the operator typed when a push or a reconnect reads them again
- T-2465 A save in one section of the host settings leaves unsaved edits in the other sections as typed

## Notes

### Accepted by the orchestrator

- Verified: b1b1711f26e8a9282e37e41c9dfc8e3e0647fd04
- At: 2026-11-06T18:14:24Z

Verdict: pass
- 1: flaiover/src/lib/form.ts, flaiover/src/lib/form.test.ts, flaiover/src/lib/components/StrategicSettings.svelte, flaiover/src/lib/components/StrategicSettings.svelte.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts
