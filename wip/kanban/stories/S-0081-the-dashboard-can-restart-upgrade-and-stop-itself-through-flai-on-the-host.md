---
id: S-0081
type: story
nature: feature
title: The dashboard can restart, upgrade, and stop itself through flai on the host
status: in-progress
parent: E-0003
owner: alex
created: 2026-09-20T07:26:52Z
updated: 2026-09-22T21:13:07Z
transitions:
  - to: ready
    at: 2026-09-20T13:01:46Z
    by: alex
  - to: in-progress
    at: 2026-09-22T21:12:17Z
    by: system-flow
tags: [cli, dashboard]
touches: [flai/cmd, flaiover/src, design/system]
---
# S-0081 The dashboard can restart, upgrade, and stop itself through flai on the host

## Goal
A host action for managing the dashboard from its own page: restart, pull and start a newer image, stop. Off until the operator enables it.

## Acceptance criteria
- [ ] Restart, upgrade to the image tag the configuration names, and stop are host actions, enabled by name, journalled, and run by flai with the same code as `flai dashboard`
- [ ] The page says what version runs, whether a newer one exists, and what the action will do; after a restart or upgrade the browser reconnects by itself and says what changed
- [ ] An upgrade that fails leaves the previous container running, and the page says why

## Tasks
- T-0322 flai dashboard gets restart, check, and upgrade subcommands, upgrade a safe blue-green swap
- T-0323 A dashboard host action gates restart, check, and upgrade over the channel, journalled like push and agent
- T-0324 A host page shows the dashboard's own version and update state, with restart/upgrade/stop actions
- T-0325 Tried end to end: restart, a real upgrade swap, a failed upgrade, and stop, against real Docker

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Later. Depends on the host-actions groundwork in the push story.
