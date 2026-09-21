---
id: S-0092
type: story
nature: improvement
title: New stories should not require an epic
status: backlog
parent: E-0003
owner: alex
created: 2026-09-21T04:12:16Z
updated: 2026-09-21T04:12:16Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# S-0092 New stories should not require an epic

## Goal

When creating a new story, an epic should *not* be required. Not every story belongs to an active epic, sometimes they alter code that belonged to a delivered epic. Creating epics as placeholder containers is an anti-pattern.

## Acceptance criteria
- [ ] The user should be allowed to choose "No epic" from the parent epic dropdown in the New epic or story creation screen.
- [ ] The associated flai operation should allow for stories without an epic to be created.

## Tasks

## Notes
