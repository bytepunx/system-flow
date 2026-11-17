---
id: S-0823
type: story
nature: remediation
title: A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back
status: backlog
owner: alex
created: 2026-11-17T07:13:19Z
updated: 2026-11-17T07:13:19Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h
    by: flai
    at: 2026-11-17T07:13:19Z
---
# S-0823 A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back

## Goal

This story remediates [I-0397](../../../design/issues/I-0397-a-story-the-operator-moves-from-review-to-in-progress-by-hand-gets-no-agent-and.md), "A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back". The issue recommends this solution:

A move of a story from a state of kind `waiting` back to a working state by the operator's hand, from the dashboard or `flai move`, with a reason, is a send back: flai opens the send back's thread on the story with the reason as its text, as `sendback.Thread` does for `flai accept` and the rebase and resolve job, so that flai serve starts the story's agent on it. Where the operator gives no reason, the dashboard asks for one, or the board shows the story as in progress with no agent and no thread awaiting it, with what to do. Say which in the story, and bring the workflow design and the user guide in step.

## Acceptance criteria
- [ ] The cause I-0397 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0397 is closed with `flai issue close I-0397 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0397. time_lost_per_cycle 1h carried over from I-0397's Impact section. Its Impact time_lost_per_cycle was taken rather than the 2h30m derived from its cost and count.
