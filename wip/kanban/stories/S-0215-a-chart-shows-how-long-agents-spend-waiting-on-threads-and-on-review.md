---
id: S-0215
type: story
nature: feature
title: A chart shows how long agents spend waiting on threads and on review
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:15Z
updated: 2026-10-02T11:54:42Z
transitions: []
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0205]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0215 A chart shows how long agents spend waiting on threads and on review

## Goal

Agents wait for the operator: on threads they open and on stories in review. That time is the operator's bottleneck and the orchestrator's reason to exist. The designer asked for a chart of it on 2026-10-02.

## Acceptance criteria
- [ ] `/charts/agent-waiting`: stacked bar per week of the window of hours agents waited, split into thread waits (opened to answered, from thread timestamps while the story was in progress) and review waits (review to done, from transitions), with the mean wait per story as a line
- [ ] A table under it lists the longest waits in the window with the story, the thread, and who was awaited
- [ ] Spans the window, reads `/api/stats`, matches `flai stats --json`; listed under Flow; design and user guide describe it; tests cover the mapping

## Tasks

## Notes
