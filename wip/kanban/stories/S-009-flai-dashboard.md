---
id: S-009
type: story
nature: feature
title: flai dashboard
status: backlog
parent: E-002
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-15T16:22:05Z
transitions: []
tags: []
---

# S-009 flai dashboard

## Goal
flai dashboard pulls and runs the flaiover image against the current repo and opens the browser.

## Acceptance criteria
- [ ] Pulls image if missing, honours config image, tag, and port
- [ ] Mounts the repo read-write at /project with the host UID
- [ ] flai dashboard stop stops the container
- [ ] Clear error when docker is missing

## Tasks


## Notes
