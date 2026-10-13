---
id: I-0125
title: A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other
class: efficiency
status: open
count: 8
cost: 6m
first_reported: 2026-10-08T10:21:38Z
last_reported: 2026-10-13T19:15:11Z
updated: 2026-10-13T19:15:11Z
---

# I-0125 A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other

## Description
A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other

## Instances

### 2026-10-08T10:21:38Z
Story: S-0232.
S-0232's agent ended to wait on TH-0324, and flai serve restarted it once for each reply on MS-0021, MS-0022, MS-0037, and others, including replies that only said "agreed, nothing to change". Answering such a reply sets awaiting back to the other story, which wakes that agent in turn. Each restart reprimes the session, reads the conversation, and rewrites the narrative, for no change to either story. A reply could mark the conversation settled, or a message that asks nothing could leave awaiting empty, so that neither agent is started again until a new change arrives.

### 2026-10-08T10:43:24Z
Story: S-0232.
MS-0041 with S-0347: S-0232's agent answered that there was nothing to change and asked for no reply. S-0347's agent replied anyway, restating that, which woke S-0232's agent again. S-0232 left the conversation unanswered, even though it is listed as awaiting S-0232, so that S-0347's agent is not woken in turn.

### 2026-10-12T14:41:55Z
Story: S-0232.
MS-0053 from S-0362, like MS-0021 and MS-0031 before it, notified S-0232 that a task changed `design/issues/summary.md`. That file is generated and regenerated at every sync (ADR-0098), so the notice cannot break anything. Each such notice still restarts the claiming story's agent. A notice for a path that flai regenerates at sync could be left out.

### 2026-10-12T15:29:49Z
Story: S-0337.
S-0337 was blocked for four days (TH-0390) with no branch content, and flai serve restarted its agent 12 times between 2026-10-08 and 2026-10-12 for messages on MS-0045, MS-0046, MS-0049, and MS-0054. Most were "changed paths" notices from another story's task done that could not affect a story with no commits. Others were the other agent's acknowledgement of a settled split, or a "Closed: ... was accepted" entry, which awaits no one. Asking in the conversation not to be told about further commits on the same files changed nothing, because the next task done sent another notice. A notice to a story that has no commits on its branch, or one that is blocked on a thread, could be held until the story can act on it.

### 2026-10-12T22:07:22Z
Story: S-0356.
MS-0101 (S-0356 and S-0239, on a share with an agreed split) woke S-0356's agent three times while S-0356 was otherwise idle on TH-0429: for 242d70b's ack, d989ac0, and 7a64631. Each change sat wholly inside S-0239's side of the split, so each needed only a "settled" reply, and each reply in turn awaited S-0239. A share's split could stop the notice for paths inside the other story's side.

### 2026-10-13T19:12:42Z
Story: S-0380.
MS-0109 (S-0378 and S-0380) was settled in its first reply. S-0378's agreement then woke S-0380's agent a second time, while S-0380 waited only on the operator's answer to TH-0445. S-0380 had to reply again ("settled, no reply needed"), which wakes S-0378 in turn.

### 2026-10-13T19:12:53Z
Story: S-0378.
MS-0109 between S-0378 and S-0380: S-0378 agreed nothing changes, S-0380 replied "Settled ... no reply is needed", and the conversation still awaited S-0378, so flai serve started S-0378's agent again to read a reply that needed no answer.

### 2026-10-13T19:15:11Z
Story: S-0378.
MS-0110: S-0357 agreed to the shared split on docs/operators/settings.md, which left the settled conversation awaiting S-0378; its closing reply in turn awaits S-0357 and woke S-0378 from waiting on TH-0436.

## Remediation

Story S-0372 remediates this issue, created from it at 2026-10-12T16:52:49Z.
