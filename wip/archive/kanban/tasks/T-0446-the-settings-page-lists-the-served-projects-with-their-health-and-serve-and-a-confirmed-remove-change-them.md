---
id: T-0446
type: task
nature: improvement
title: The settings page lists the served projects with their health, and Serve and a confirmed Remove change them
status: done
parent: S-0122
owner: alex
created: 2026-09-26T06:05:08Z
updated: 2026-09-26T06:12:43Z
transitions:
  - to: ready
    at: 2026-09-26T06:05:19Z
    by: agent-S-0119
  - to: in-progress
    at: 2026-09-26T06:08:41Z
    by: agent-S-0119
  - to: done
    at: 2026-09-26T06:12:43Z
    by: agent-S-0119
stream: S-0122
tags: []
touches: [flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/settings.ts, flaiover/src/routes/api/settings]
---
# T-0446 The settings page lists the served projects with their health, and Serve and a confirmed Remove change them

## Work

A Projects section on `SettingsPanel.svelte`: each served project with key, name, root, and connected since or the last error or why it is unavailable; each unserved project below an import folder with why and a Serve button; each registered project with a Remove button behind a confirmation that says none of its files is touched and how to add it back. Both go through `/api/settings` as the new kinds `serve` and `unserve`, and flai's refusal shows in the section. After Serve the page refreshes the switcher until the project is connected, as the board's import does, so it appears without a reload.

## Done when

- component tests cover the list, both actions, the confirmation, and a refusal
- the route test covers the two new kinds
- `make flaiover-test` passes

## Notes
