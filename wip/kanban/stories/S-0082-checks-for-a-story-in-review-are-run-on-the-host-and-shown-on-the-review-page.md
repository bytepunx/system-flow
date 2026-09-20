---
id: S-0082
type: story
nature: feature
title: Checks for a story in review are run on the host and shown on the review page
status: ready
parent: E-0003
owner: alex
created: 2026-09-20T07:26:52Z
updated: 2026-09-20T13:01:56Z
transitions:
  - to: ready
    at: 2026-09-20T13:01:56Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, design/system]
---
# S-0082 Checks for a story in review are run on the host and shown on the review page

## Goal
A host action that runs the project's checks for a story in review, in the story's worktree, and shows the result on the review page, so that acceptance is decided with the evidence at hand.

## Acceptance criteria
- [ ] The commands are named in the project's manifest or the host's configuration by the operator, as argument lists, never sent by the dashboard; the action is off until enabled
- [ ] A run streams its output to the review page, can be cancelled (terminate, then kill), has a time limit, and its outcome and duration are kept with the story until it is accepted
- [ ] One run per story at a time; what runs where is journalled

## Tasks

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Later. Depends on the host-actions groundwork in the push story.
