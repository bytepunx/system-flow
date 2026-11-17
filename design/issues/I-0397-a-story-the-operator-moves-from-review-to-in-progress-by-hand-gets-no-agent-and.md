---
id: I-0397
title: A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back
class: defect
status: closed
count: 1
cost: 2h30m
first_reported: 2026-11-17T07:13:19Z
last_reported: 2026-11-17T07:13:19Z
updated: 2026-11-17T10:52:17Z
---

# I-0397 A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back

## Description
A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back

## Instances

### 2026-11-17T07:13:19Z
The operator moved S-0817 from review to in-progress from the dashboard at 2026-11-17T04:29:20Z with a reason. .flai-cache/events.jsonl holds no agent event for S-0817 after its agent ended at 03:23:45Z, and no thread was opened on it: the send back's thread, which flai serve starts a story's agent on, is opened only by flai accept's send back and by a failing rebase and resolve job. The story stayed in in-progress, counted against the WIP limit, with no agent, until the operator asked at 06:54Z.

## Impact
S-0817 held an in-progress slot for 2h25m with no agent; any send back by hand does the same until the operator notices

- time_lost_per_cycle: 1h

## Remediation

A move of a story from a state of kind `waiting` back to a working state by the operator's hand, from the dashboard or `flai move`, with a reason, is a send back: flai opens the send back's thread on the story with the reason as its text, as `sendback.Thread` does for `flai accept` and the rebase and resolve job, so that flai serve starts the story's agent on it. Where the operator gives no reason, the dashboard asks for one, or the board shows the story as in progress with no agent and no thread awaiting it, with what to do. Say which in the story, and bring the workflow design and the user guide in step.

Story S-0823 remediates this issue, created from it at 2026-11-17T07:13:19Z.
Closed 2026-11-17T10:52:17Z: S-0823: flai move along the send-back transition by anyone but flai and the story's own agent, from the command line or the dashboard, opens the send back's thread by flai with the move's reason, and flai serve starts the story's agent on it, told who sent it back by hand, and resolves the thread once the story is back in review.
