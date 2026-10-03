---
id: S-0200
type: story
nature: improvement
title: "An epic follows its stories: to ready and in-progress with the first, to review and done with the last"
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:11Z
updated: 2026-10-03T05:33:53Z
transitions:
  - to: ready
    at: 2026-10-03T05:33:53Z
    by: alex
tags: [flai, dashboard]
touches: [flai/internal/workitem, flai/internal/itemedit, flai/cmd/move.go, flai/cmd/accept.go, flai/internal/hostapi, flai/internal/mcpserver, design/system/workflow.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0200 An epic follows its stories: to ready and in-progress with the first, to review and done with the last

## Goal

An epic's status is moved by hand and drifts from its stories'. The designer decided on 2026-10-02 that an epic follows its children: it moves to `ready` when its first story is moved to ready or created in ready, to `in-progress` when its first story starts, to `review` when its last open story moves to review, and to `done` when its last story is accepted. This also gives the planner and orchestrator epics whose state means something.

## Acceptance criteria
- [ ] When a story moves from backlog to ready, or is created in ready, and its epic is in backlog, the epic moves to ready in the same write, with a transition `by` the same actor and a reason naming the story
- [ ] When a story moves to in-progress and its epic is in ready or backlog, the epic moves to in-progress
- [ ] When the last story of an epic that is not done or cancelled moves to review, the epic moves to review; when the last one is accepted (done), the epic moves to done through the same acceptance flow (`flai accept`, `flai move done`, the dashboard), archiving it as an epic done by hand is archived today
- [ ] A story moved back (review to in-progress, ready to backlog) moves the epic back only when no other story holds it where it is; cancelled stories do not count
- [ ] The epic's move is reported in the command's output, the MCP result, and the inbox as a change on the epic
- [ ] `flai check` warns about an epic whose status lags its stories under these rules, so epics from before this change are found
- [ ] `design/system/workflow.md`, `work-hierarchy.md`, and the user guide describe it; an ADR records the rule, refining ADR-0004
- [ ] Tests cover each transition, the back moves, and acceptance of the last story

## Tasks

## Notes
