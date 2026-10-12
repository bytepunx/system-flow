---
id: I-0125
title: A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other
class: efficiency
status: open
count: 4
cost: 9m
first_reported: 2026-10-08T10:21:38Z
last_reported: 2026-10-12T15:29:49Z
updated: 2026-10-12T16:52:49Z
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

## Remediation

Story S-0372 remediates this issue, created from it at 2026-10-12T16:52:49Z.
