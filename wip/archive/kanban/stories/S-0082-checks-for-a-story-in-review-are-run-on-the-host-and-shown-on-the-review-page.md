---
id: S-0082
type: story
nature: feature
title: Checks for a story in review are run on the host and shown on the review page
status: done
parent: E-0003
owner: alex
created: 2026-09-20T07:26:52Z
updated: 2026-09-22T23:52:35Z
transitions:
  - to: ready
    at: 2026-09-20T13:01:56Z
    by: alex
  - to: in-progress
    at: 2026-09-22T22:34:59Z
    by: system-flow
  - to: review
    at: 2026-09-22T23:45:54Z
    by: system-flow
  - to: done
    at: 2026-09-22T23:52:35Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, design/system]
---
# S-0082 Checks for a story in review are run on the host and shown on the review page

## Goal
A host action that runs the project's checks for a story in review, in the story's worktree, and shows the result on the review page, so that acceptance is decided with the evidence at hand.

## Acceptance criteria
- [x] The commands are named in the project's manifest or the host's configuration by the operator, as argument lists, never sent by the dashboard; the action is off until enabled
- [x] A run streams its output to the review page, can be cancelled (terminate, then kill), has a time limit, and its outcome and duration are kept with the story until it is accepted
- [x] One run per story at a time; what runs where is journalled

## Tasks
- T-0326 Named check commands are configured in the manifest or the host's config, never sent by the dashboard
- T-0327 flai checks run/status/cancel/tail: sequential, in the story's worktree, one at a time, cancellable by killing the whole process group
- T-0328 A checks host action gates run/cancel over the channel, journalled like the others
- T-0329 The review page shows check state and, once enabled, Run and Cancel
- T-0330 Tried end to end: a real run, a real cancel that kills a spawned grandchild, a real timeout

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Later. Depends on the host-actions groundwork in the push story.
