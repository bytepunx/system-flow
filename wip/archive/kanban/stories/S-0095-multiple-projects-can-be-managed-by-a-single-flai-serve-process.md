---
id: S-0095
type: story
nature: feature
title: Multiple projects can be managed by a single flai serve process
status: done
parent: E-0003
owner: alex
created: 2026-09-23T02:05:47Z
updated: 2026-09-23T03:21:24Z
transitions:
  - to: ready
    at: 2026-09-23T02:06:00Z
    by: alex
  - to: in-progress
    at: 2026-09-23T02:20:36Z
    by: system-flow
  - to: review
    at: 2026-09-23T02:35:25Z
    by: system-flow
  - to: done
    at: 2026-09-23T03:21:24Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# S-0095 Multiple projects can be managed by a single flai serve process

## Goal

When a `flai serve` process fires up, it should support managing and working on multiple repositories.

This means that there should be a way to select a project from the top bar that drives the other screens (overview, board, inbox, etc.).

## Acceptance criteria
- [x] operator can change the active project from the bar
- [x] when an active project changes, the board should get updates from the server

## Tasks
- T-0341 Live changes reach each project's cache and event stream, however late its flai connects
- T-0342 The top bar names the active project; picking another switches every screen in place
- T-0343 Tried with one flai serve, two projects, and a browser; docs say how to add a project

## Notes
