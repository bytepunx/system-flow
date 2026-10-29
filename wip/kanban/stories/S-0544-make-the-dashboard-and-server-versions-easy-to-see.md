---
id: S-0544
type: story
nature: improvement
title: Make the dashboard and server versions easy to see
status: backlog
owner: alex
created: 2026-10-29T17:39:08Z
updated: 2026-10-29T17:39:08Z
transitions: []
tags: []
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0544 Make the dashboard and server versions easy to see

## Goal

Right now, the host page doesn't show the dashboard version and you tend to only see that new versions are available after a push. It should be very easy to see what versions are running.

## Acceptance criteria
- [ ] Put a label on the nav rail under the license that shows the server and dashboard versions (newline separated)
- [ ] Show the versions of the current dashboard and the available dashboard version when the user checks, not just that a new version is available

## Tasks

## Notes
