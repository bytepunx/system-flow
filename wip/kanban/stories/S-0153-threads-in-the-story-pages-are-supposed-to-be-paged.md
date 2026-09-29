---
id: S-0153
type: story
nature: remediation
title: Threads in the story pages are supposed to be paged
status: in-progress
parent: E-0013
owner: alex
created: 2026-09-29T05:52:49Z
updated: 2026-09-29T06:44:20Z
transitions:
  - to: ready
    at: 2026-09-29T05:57:53Z
    by: alex
  - to: in-progress
    at: 2026-09-29T06:44:20Z
    by: agent-S-0153
tags: [dashboard]
topics: [client-side]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0153 Threads in the story pages are supposed to be paged

## Goal

Threads were supposed to only be displayed one at a time on story pages with controls to allow the operator to "page" through them instead of them needing to vertically scroll to the correct one.

## Acceptance criteria
- [ ] Threads are shown one at a time on the story page
- [ ] Operators can page through threads using arrow controls
- [ ] A label next to or in between the arrows identifies the cardinal index (1 based) out of the count of threads

## Tasks

## Notes
