---
id: I-0125
title: A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other
class: efficiency
status: open
count: 1
cost: 3m
first_reported: 2026-10-08T10:21:38Z
last_reported: 2026-10-08T10:21:38Z
updated: 2026-10-08T10:21:38Z
---

# I-0125 A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other

## Description
A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other

## Instances

### 2026-10-08T10:21:38Z
Story: S-0232.
S-0232's agent ended to wait on TH-0324, and flai serve restarted it once for each reply on MS-0021, MS-0022, MS-0037, and others, including replies that only said "agreed, nothing to change". Answering such a reply sets awaiting back to the other story, which wakes that agent in turn. Each restart reprimes the session, reads the conversation, and rewrites the narrative, for no change to either story. A reply could mark the conversation settled, or a message that asks nothing could leave awaiting empty, so that neither agent is started again until a new change arrives.

## Remediation
