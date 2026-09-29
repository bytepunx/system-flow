---
id: S-0155
type: story
nature: improvement
title: Inbox items take users to the story page, not document pages
status: in-progress
parent: E-0013
owner: alex
created: 2026-09-29T05:57:30Z
updated: 2026-09-29T07:02:21Z
transitions:
  - to: ready
    at: 2026-09-29T06:04:02Z
    by: alex
  - to: in-progress
    at: 2026-09-29T07:00:34Z
    by: agent-S-0155
tags: [dashboard]
topics: [client-side]
touches: [flaiover/src, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 160
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 80
      output: 568
      cache_read: 3212716
      cache_write: 96953
      cost: 1.3095
---
# S-0155 Inbox items take users to the story page, not document pages

## Goal

Right now, clicking some parts of an inbox notification for open questions in story threads can take the operator to a documentation page instead of the related story.

Inbox items that require the operator to handle a thread on a story page should always take the operator to the story's page.

## Acceptance criteria
- [ ] Clicking the inbox item where a user needs to respond to a question in a thread always takes the operator to the story page and displays the correct thread.

## Tasks
- T-0550 An inbox thread entry links to its story's page naming the thread
- T-0551 The story page opens on the thread its link names
- T-0552 The dashboard design and user guide say where an inbox thread leads

## Notes
