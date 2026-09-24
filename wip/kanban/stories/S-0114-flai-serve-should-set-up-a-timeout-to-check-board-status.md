---
id: S-0114
type: story
nature: remediation
title: Flai serve should set up a timeout to check board status
status: backlog
parent: E-0008
owner: alex
created: 2026-09-24T08:16:14Z
updated: 2026-09-24T08:16:14Z
transitions: []
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
---
# S-0114 Flai serve should set up a timeout to check board status

## Goal

flai serve needs to perform periodic board state checks to determine if there are stories that are ready and can be moved into in process.

## Acceptance criteria
- [ ] Items that exist in ready and don't get picked up immediately should get picked up when the timer elapses.

## Tasks

## Notes
