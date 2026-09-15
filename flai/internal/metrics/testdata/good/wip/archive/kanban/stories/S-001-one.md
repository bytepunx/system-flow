---
id: S-001
type: story
nature: feature
title: One
status: done
parent: E-001
owner: agent
created: 2026-08-01T09:00:00Z
updated: 2026-08-03T12:00:00Z
transitions:
  - to: ready
    at: 2026-08-01T10:00:00Z
    by: agent
  - to: in-progress
    at: 2026-08-02T10:00:00Z
    by: agent
  - to: review
    at: 2026-08-03T10:00:00Z
    by: agent
  - to: done
    at: 2026-08-03T12:00:00Z
    by: alex
blocked:
  - from: 2026-08-02T12:00:00Z
    until: 2026-08-02T18:00:00Z
    reason: waiting
estimate: 20h
tags: []
---

# S-001 One

## Goal
g

## Acceptance criteria
- [x] works

## Tasks
- T-001 T1

## Notes
