---
id: E-0018
type: epic
nature: feature
title: Improve Inter-Agent Coordination Via Inboxes
status: in-progress
owner: alex
created: 2026-10-07T20:05:16Z
updated: 2026-10-08T04:31:41Z
transitions:
  - to: ready
    at: 2026-10-07T20:05:17Z
    by: alex
  - to: in-progress
    at: 2026-10-07T20:24:33Z
    by: agent-S-0330
tags: []
usage:
  source: sum
  seconds: 15007
  estimated: true
  turns:
    - day: 2026-10-07
      ceremony: 10
      test_runs: 2
      hand_edits: 11
      work: 293
    - day: 2026-10-08
      ceremony: 1
      test_runs: 4
      hand_edits: 4
      work: 161
  models:
    - model: claude-opus-5-5
      input: 2578
      output: 831287
      cache_read: 153153530
      cache_write: 5069427
      cost: 78.4183
  strategic:
    - kind: planner
      seconds: 1011
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 442
          output: 9689
          cache_read: 2347064
          cache_write: 68686
          cost: 0.3695
        - model: claude-opus-5-5
          input: 728
          output: 107451
          cache_read: 70599021
          cache_write: 339947
          cost: 35.112
    - kind: orchestrator
      seconds: 9130
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 727
          output: 12160
          cache_read: 103858419
          cache_write: 269744
          cost: 25.6568
        - model: claude-sonnet-5-5
          input: 99
          output: 571
          cache_read: 1605467
          cache_write: 347469
          cost: 1.6834
cost_of_delay:
  inputs:
    penalty_per_week: 1000
    by: alex
    at: 2026-10-07T20:05:16Z
  value: 1000
  by: planner-E-0018
  at: 2026-10-07T20:13:59Z
---
# E-0018 Improve Inter-Agent Coordination Via Inboxes

## Outcome

Right now, agents often get "held" based on the predicted overlapping touches (expected modifications to the same files or folders).

Additionally - agents that do run in parallel can run into issues when it's time to merge their work in and introduce a number of conflicting changes when require one or more story agents to resolve the conflicts at the end before moving from review to done.

Introduce inter-agent inboxes (not an overlapping concept with the current user inbox) so that agents are able to exchange messages in order to coordinate changes and reduce the amount of late arriving conflicts or simplify the conflict resolution process.

The ultimate goal here is to allow for greater system-wide throughput without incurring more work at the end of each story's process to resolve conflicts.

## Stories
- S-0330 flai message sends a message from one open story's agent to another's, kept apart from the operator's threads
- S-0331 A story's agent sends and answers messages through MCP, inbox lists those to its story, and wait_for_events wakes on one
- S-0332 flai tells two stories' agents of a trial-merge conflict or a grown overlap with a message between them, and asks the operator only when they do not agree
- S-0333 Closing a task messages every open story whose claim covers a path the task changed, before either story is accepted
- S-0334 A story held on overlap is asked about by message, and starts when the holding story's agent shares the paths or narrows its claim
- S-0335 A story's agent that waits only on another agent's reply ends, and flai serve starts it again when the reply or a new message to its story comes
- S-0336 The dashboard shows the messages between stories' agents on each story's page and in a Messages view
- S-0337 flai stats reports the conversations between agents, the conflicts found at sync and at acceptance, and the hold time shares saved, per week
- S-0338 A held card says when its holding story's agent was asked about the hold, and a story started on a share names the paths it shares

## Notes

Possible ways inboxes could be used by agents - determining opportunities for clearer boundaries and separation of concerns so that there is less shared/overlapping code. This may add some additional duration to the work but should save total time in process by making the codebase easier to evolve independently.
- 2026-10-07T20:24:33Z: moved to in-progress: follows S-0330, which moved to in-progress
