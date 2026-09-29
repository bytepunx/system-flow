---
id: S-0143
type: story
nature: feature
title: Stories and Tasks track tokens and cost
status: backlog
parent: E-0011
owner: alex
created: 2026-09-29T00:24:52Z
updated: 2026-09-29T00:24:52Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0143 Stories and Tasks track tokens and cost

## Goal

Epics, Stories, and Tasks gain frontmatter to include tokens used and cost or estimated cost of the expenditure. Charts use this to show users the following (all should group by agent model doing the work):

- token rates (epic, task, and story) per hour
- cost or estimated cost (epic, task, and story)
- epic, story, and task completion as a function of time and cost

## Acceptance criteria
- [ ] epics, stories, and tasks track tokens spent
- [ ] epics, stories, and tasks track cost
- [ ] when any item (epic, story, task) moves into done, it triggers a cascade up the hierarchy to update its parent's aggregate front-matter tracking

## Tasks

## Notes
