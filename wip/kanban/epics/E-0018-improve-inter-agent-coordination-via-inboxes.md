---
id: E-0018
type: epic
nature: feature
title: Improve Inter-Agent Coordination Via Inboxes
status: backlog
owner: alex
created: 2026-10-07T20:05:16Z
updated: 2026-10-07T20:05:16Z
transitions: []
tags: []
cost_of_delay:
  inputs:
    penalty_per_week: 1000
    by: alex
    at: 2026-10-07T20:05:16Z
---
# E-0018 Improve Inter-Agent Coordination Via Inboxes

## Outcome

Right now, agents often get "held" based on the predicted overlapping touches (expected modifications to the same files or folders).

Additionally - agents that do run in parallel can run into issues when it's time to merge their work in and introduce a number of conflicting changes when require one or more story agents to resolve the conflicts at the end before moving from review to done.

Introduce inter-agent inboxes (not an overlapping concept with the current user inbox) so that agents are able to exchange messages in order to coordinate changes and reduce the amount of late arriving conflicts or simplify the conflict resolution process.

The ultimate goal here is to allow for greater system-wide throughput without incurring more work at the end of each story's process to resolve conflicts.

## Stories

## Notes

Possible ways inboxes could be used by agents - determining opportunities for clearer boundaries and separation of concerns so that there is less shared/overlapping code. This may add some additional duration to the work but should save total time in process by making the codebase easier to evolve independently.
