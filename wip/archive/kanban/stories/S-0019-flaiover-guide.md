---
id: S-0019
type: story
nature: feature
title: flaiover guide
status: done
parent: E-0004
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-24T08:32:17Z
transitions:
  - to: ready
    at: 2026-09-24T07:58:40Z
    by: alex
  - to: in-progress
    at: 2026-09-24T08:28:02Z
    by: agent-S-0019
  - to: review
    at: 2026-09-24T08:31:52Z
    by: agent-S-0019
  - to: done
    at: 2026-09-24T08:32:17Z
    by: alex
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0019 flaiover guide

## Goal
docs/users and docs/operators cover the dashboard views, running it, and the local-only security posture.

## Acceptance criteria
- [x] Every view has a paragraph and the operator page warns against exposure

## Tasks
- T-0406 Users guide: a paragraph for every dashboard view
- T-0407 Operators guide: warn against exposing the dashboard

## Notes
The goal's "local-only" posture predates the dashboard being published on every interface by default; the guides document the posture as it is and warn against exposure (see the narrative's Decisions).
