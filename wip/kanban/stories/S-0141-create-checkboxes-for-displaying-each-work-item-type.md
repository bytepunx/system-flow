---
id: S-0141
type: story
nature: improvement
title: Create checkboxes for displaying each work item type
status: backlog
parent: E-0003
owner: alex
created: 2026-09-27T04:06:55Z
updated: 2026-09-27T04:06:55Z
transitions: []
tags: [dashboard]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0141 Create checkboxes for displaying each work item type

## Goal

Include a checkbox for each work item type in the board. Use local browser storage to remember the user's last configuration choices so it remains the same between page refreshes/navigation.

## Acceptance criteria
- [ ] a set of checkboxes for tasks, stories, and epics are visible
- [ ] toggling the checkbox for a particular work type toggles its visibility
- [ ] stories are on by default until a user changes the toggle

## Tasks

## Notes
