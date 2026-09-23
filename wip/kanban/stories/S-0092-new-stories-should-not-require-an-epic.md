---
id: S-0092
type: story
nature: improvement
title: New stories should not require an epic
status: review
parent: E-0003
owner: alex
created: 2026-09-21T04:12:16Z
updated: 2026-09-23T00:37:36Z
transitions:
  - to: ready
    at: 2026-09-22T22:39:04Z
    by: alex
  - to: in-progress
    at: 2026-09-23T00:37:34Z
    by: system-flow
  - to: review
    at: 2026-09-23T00:37:36Z
    by: system-flow
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd, flai/internal/workitem, flai/internal/hostapi]
---
# S-0092 New stories should not require an epic

## Goal

When creating a new story, an epic should *not* be required. Not every story belongs to an active epic, sometimes they alter code that belonged to a delivered epic. Creating epics as placeholder containers is an anti-pattern.

## Acceptance criteria
- [x] The user should be allowed to choose "No epic" from the parent epic dropdown in the New epic or story creation screen.
- [x] The associated flai operation should allow for stories without an epic to be created.

## Tasks
- T-0335 flai story new and the dashboard's new-item form allow a story with no epic

## Notes
