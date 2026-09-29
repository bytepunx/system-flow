---
id: S-0167
type: story
nature: feature
title: Right click menu for the board
status: backlog
parent: E-0013
owner: alex
created: 2026-09-29T23:30:41Z
updated: 2026-09-29T23:30:41Z
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
# S-0167 Right click menu for the board

## Goal

A right click menu for lanes on the board gives the operator convenient access to the following actions:

- create item
- move stories forward
- move stories backward
- change WIP limit

## Acceptance criteria
- [ ] when taking the operator to the create screen, the lane the menu was created for is the starting lane for the card (only works for backlog, ready, and in-progress, all other lanes get changed to backlog as the default)
- [ ] moving stories forward should only be an option for the backlog lane
- [ ] moving stories backward should only be an option in ready, in progress, review, and cancelled (stories cannot be moved from done). cancelled stories move back to backlog
- [ ] changing the wip limit should change the correct file(s) needed to make the change permanent

## Tasks

## Notes
