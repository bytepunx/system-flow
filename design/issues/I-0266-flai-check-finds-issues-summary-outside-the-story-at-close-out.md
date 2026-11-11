---
id: I-0266
title: "flai check finds `issues.summary` outside the story at close-out"
class: efficiency
status: closed
count: 1
first_reported: 2026-11-07T00:11:02Z
last_reported: 2026-11-07T00:11:02Z
updated: 2026-11-11T09:18:16Z
---

# I-0266 flai check finds `issues.summary` outside the story at close-out

## Description
flai check finds `issues.summary` outside the story at close-out

## Instances

### 2026-11-07T00:11:02Z
Story: S-0576.
flai check found outside the story:
`design/issues/summary.md`: summary.md differs from what its issue files render; run flai check --fix or flai issue summary, which rewrite it

## Remediation

Story S-0702 remediates this issue, created from it at 2026-11-10T08:55:37Z.
Closed 2026-11-11T09:18:16Z: S-0702 (ADR-0286): once a rebase that wrote the issues summary or the ADR index again at a stop finishes, flai stream sync (and so flai task done, the close-out, and flai accept) and the merge queue's batch write every generated file again at the branch's head and commit what differs, so a later commit replayed cleanly over a file written at a stop no longer leaves the head's summary stale. TestSyncWritesTheGeneratedFilesAgainAtTheHeadAfterAStopOnThem and TestTheBatchWritesTheGeneratedFilesAgainAtTheRebasedHeadAfterAStopOnThem reproduce it.
