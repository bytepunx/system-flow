---
id: S-0095
type: story
nature: feature
title: Multiple projects can be managed by a single flai serve process
status: backlog
parent: E-0003
owner: alex
created: 2026-09-23T02:05:47Z
updated: 2026-09-23T02:05:47Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# S-0095 Multiple projects can be managed by a single flai serve process

## Goal

When a `flai serve` process fires up, it should support managing and working on multiple repositories.

This means that there should be a way to select a project from the top bar that drives the other screens (overview, board, inbox, etc.).

## Acceptance criteria
- [ ] operator can change the active project from the bar
- [ ] when an active project changes, the board should get updates from the server

## Tasks

## Notes
