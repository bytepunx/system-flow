---
id: I-0301
title: flai serve starts a story's agent for a ready story the board holds no-touches, and story_start then refuses it
class: efficiency
status: closed
count: 1
cost: 5m
first_reported: 2026-11-07T08:44:09Z
last_reported: 2026-11-07T08:44:09Z
updated: 2026-11-09T23:20:32Z
---

# I-0301 flai serve starts a story's agent for a ready story the board holds no-touches, and story_start then refuses it

## Description
flai serve starts a story's agent for a ready story the board holds no-touches, and story_start then refuses it

## Instances

### 2026-11-07T08:44:09Z
Story: S-0568.
Story: S-0568. flai serve started agent-S-0568 "because it entered ready", but S-0568 declared no touches, so story_start refused it as held (no-touches) by S-0562, S-0564, and S-0584 in progress. The agent read the issue and the guard to find the paths, declared the story's touches with flai touches, and started it; the planner, asked by the operator at the same moment, then added its own. S-0586 and S-0587 entered ready the same morning with no touches and were held the same way. flai serve could hold its start until the story's touches are declared or the planner has run, or story_start could tell the agent that declaring touches clears the hold.

## Remediation

Story S-0633 remediates this issue, created from it at 2026-11-07T23:41:06Z.
Closed 2026-11-09T23:20:32Z: Fixed by S-0633. flai serve's look started agent-S-0568 because it read the planner's task T-2550 in the seconds between itemnew.Create writing it and the check removing it, and so judged S-0568's claim by that task's touches. Now a checked creation holds a checked-write lock from its write until it keeps or removes the item. The look reads its items under that lock, with the folders looked at again (ADR-0272, T-2786; TestALookWaitsForACheckedCreationAndStartsNoStoryHeldOnceItIsRemoved reproduces it). story_start's refusal of a story held no-touches by its own empty claim now tells the agent to declare its touches and call story_start again (T-2787).
