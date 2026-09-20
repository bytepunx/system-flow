---
id: T-0259
type: task
nature: feature
title: "flai serve: one process per user, the projects registered with it, started detached, and its state readable"
status: backlog
parent: S-0072
owner: alex
created: 2026-09-20T07:28:28Z
updated: 2026-09-20T07:28:28Z
transitions: []
stream: S-0072
tags: []
---
# T-0259 flai serve: one process per user, the projects registered with it, started detached, and its state readable

## Work
The command, its registry of projects and its state under flai's home (with an override so tests never touch the real one), foreground and detached start without root, stop, and a state file that says since when it runs, which projects it serves, and which dashboards it is connected to.

## Done when
- Command tests for register, unregister, state, and a second start that finds the first
- Nothing is written under the real flai home by any test

## Notes
