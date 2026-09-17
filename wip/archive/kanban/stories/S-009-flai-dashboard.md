---
id: S-009
type: story
nature: feature
title: flai dashboard
status: done
parent: E-002
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-17T03:58:09Z
transitions:
  - to: ready
    at: 2026-09-17T03:50:04Z
    by: agent
  - to: in-progress
    at: 2026-09-17T03:50:04Z
    by: agent
  - to: review
    at: 2026-09-17T03:54:28Z
    by: agent
  - to: done
    at: 2026-09-17T03:58:09Z
    by: alex
tags: []
---

# S-009 flai dashboard

## Goal
flai dashboard pulls and runs the flaiover image against the current repo and opens the browser.

## Acceptance criteria
- [x] Pulls image if missing, honours config image, tag, and port
- [x] Mounts the repo read-write at /project with the host UID
- [x] flai dashboard stop stops the container
- [x] Clear error when docker is missing

## Tasks
- T-079 dashboard command: run, stop, status, logs via docker with image, tag, port precedence
- T-080 Tests with a recording runner; manual smoke with a stand-in image
- T-081 Docs for users and operators

## Notes
- Verified with a stand-in image; the real image arrives in S-015 and must run as an arbitrary UID.
