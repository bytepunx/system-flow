---
id: S-0192
type: story
nature: feature
title: Add a way to create a sibling story
status: done
parent: E-0015
owner: alex
created: 2026-10-01T11:12:56Z
updated: 2026-10-02T16:35:51Z
transitions:
  - to: ready
    at: 2026-10-01T11:14:37Z
    by: alex
  - to: in-progress
    at: 2026-10-02T16:27:41Z
    by: agent-S-0192
  - to: review
    at: 2026-10-02T16:35:04Z
    by: agent-S-0192
  - to: done
    at: 2026-10-02T16:35:51Z
    by: alex
tags: [dashboard]
topics: [client-side-story-page]
touches: [flaiover/src, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 465
  models:
    - model: claude-opus-5-5
      input: 104
      output: 24433
      cache_read: 3039680
      cache_write: 142182
      cost: 2.07
    - model: claude-sonnet-5-5
      input: 10
      output: 2515
      cache_read: 114807
      cache_write: 41110
      cost: 0.1509
---
# S-0192 Add a way to create a sibling story

## Goal

Add a `New sibling` next to the `New story` button on the story page so that it's easy to create multiple stories for the same epic.

## Acceptance criteria
- [x] A `New sibling` button appears next to `New story` and when clicked, it navigates to the new story page with the same parent epic as the original story's.

## Tasks
- T-0703 A story's page offers New story on its own and New sibling under its open epic

## Notes
