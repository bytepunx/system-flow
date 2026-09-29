---
id: S-0155
type: story
nature: improvement
title: Inbox items take users to the story page, not document pages
status: backlog
parent: E-0013
owner: alex
created: 2026-09-29T05:57:30Z
updated: 2026-09-29T05:57:30Z
transitions: []
tags: [dashboard]
topics: [client-side]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0155 Inbox items take users to the story page, not document pages

## Goal

Right now, clicking some parts of an inbox notification for open questions in story threads can take the operator to a documentation page instead of the related story.

Inbox items that require the operator to handle a thread on a story page should always take the operator to the story's page.

## Acceptance criteria
- [ ] Clicking the inbox item where a user needs to respond to a question in a thread always takes the operator to the story page and displays the correct thread.

## Tasks

## Notes
