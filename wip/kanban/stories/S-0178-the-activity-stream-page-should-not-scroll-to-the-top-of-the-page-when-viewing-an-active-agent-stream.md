---
id: S-0178
type: story
nature: remediation
title: The activity stream page should not scroll to the top of the page when viewing an active agent stream
status: ready
owner: alex
created: 2026-10-01T07:38:11Z
updated: 2026-10-02T16:14:45Z
transitions:
  - to: ready
    at: 2026-10-01T07:39:43Z
    by: alex
tags: [dashboard]
topics: [client-side-activity]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0178 The activity stream page should not scroll to the top of the page when viewing an active agent stream

## Goal

Prevent the activity page from scrolling automatically (this causes the operator to have to constantly scroll back to the activity window they'd been watching).

## Acceptance criteria
- [ ] When viewing an agent's activity stream, the page does not scroll automatically
- [ ] If the user scrolls away from the end of the stream, the page should not scroll back down automatically until the operator returns to the bottom again

## Tasks

## Notes
