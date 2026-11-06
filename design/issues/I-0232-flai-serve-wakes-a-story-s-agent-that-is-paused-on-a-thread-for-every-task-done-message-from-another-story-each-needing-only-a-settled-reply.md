---
id: I-0232
title: flai serve wakes a story's agent that is paused on a thread for every task-done message from another story, each needing only a settled reply
class: efficiency
status: closed
count: 2
cost: 4m
first_reported: 2026-10-29T15:42:30Z
last_reported: 2026-10-29T16:27:03Z
updated: 2026-11-06T22:27:28Z
---

# I-0232 flai serve wakes a story's agent that is paused on a thread for every task-done message from another story, each needing only a settled reply

## Description
flai serve wakes a story's agent that is paused on a thread for every task-done message from another story, each needing only a settled reply

## Instances

### 2026-10-29T15:42:30Z
Story: S-0418.
S-0418 is paused on TH-0749 with nothing changed. S-0537 and S-0542 sent five messages between 15:14Z and 15:42Z on MS-0232 and MS-0233: one share request, one claim overlap, and three task-done notices. flai serve restarted S-0418's agent for each notice, and each restart only read the message, replied settled that S-0418 had changed nothing, and ended. A story that has changed none of the paths cannot be broken by the change, so flai could leave the notice unanswered-required (or settle it itself) while the story's branch does not change those paths.

### 2026-10-29T16:27:03Z
Story: S-0418.
At 16:26Z S-0537's T-2314 notice (3832f99) restarted S-0418's agent again while it was still paused on TH-0749 with no change to the five paths named. The only work was to read the notice, check that the strategic-agents.md hunk kept to MS-0232's split, and reply settled.

## Remediation

Story S-0545 remediates this issue, created from it at 2026-10-29T18:11:17Z.
Closed 2026-11-06T22:27:28Z: S-0545: flai task done's notice is headed (notice) (T-2441), and while a run's own question to the operator is unanswered flai serve no longer starts it for a conversation the other story wrote only notices on since the run ended; one plain entry still starts it, and the answer starts it with the notices in its inbox (ADR-0252, T-2443). TestATaskDoneNoticeDoesNotStartAnAgentThatEndedForAThread reproduces it.
