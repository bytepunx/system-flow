---
id: I-0087
title: One in-progress story holds every ready story by overlap, through folder-wide touches and docs/users/flai.md, so the board runs one story at a time under a limit of three
class: efficiency
status: closed
count: 3
cost: 1h44m
first_reported: 2026-10-06T10:32:27Z
last_reported: 2026-10-06T11:15:23Z
updated: 2026-10-06T19:22:47Z
---

# I-0087 One in-progress story holds every ready story by overlap, through folder-wide touches and docs/users/flai.md, so the board runs one story at a time under a limit of three

## Description

A ready story is held while its `touches` overlap an in-progress story's, so that two agents do not edit the same file (ADR-0046). The in-progress limit is three. On 2026-10-06 the board still ran one story at a time, because every ready story overlapped the one in progress.

Two things made the overlap total. Stories claim whole folders, such as `flai/internal/mcpserver`, and a folder claim holds every file in it. And one document, `docs/users/flai.md`, is changed by almost every flai story. A story in review keeps its claim, so the hold lasts until the operator accepts it.

## Instances

### 2026-10-06T10:32:27Z
Story: S-0283.
While S-0283 was in progress and in review, 09:57Z to 10:32Z, all nine ready stories were held by overlap with it and two of three in-progress slots stayed empty: seven by a claim on the folder flai/internal/mcpserver, which holds permission.go, two by docs/users/flai.md, one by permission.go itself. The cost is the time the two empty slots could have been working.

Story: S-0220.
While S-0220 was in progress and in review, 05:58Z to 09:57Z, all ten ready stories were held by overlap with it, through flai/internal/harness, flai/internal/serve, flai/internal/mcpserver, and flaiover/src/routes. Nothing else ran for four hours, 2h45m of them with S-0220 waiting in review and no agent running at all.

### 2026-10-06T11:15:23Z
Story: S-0285.
While S-0285 was in progress and in review, 10:32Z to 11:14Z, all eight ready stories were held by overlap with it, through flai/internal/harness, flai/internal/mcpserver, and docs/users. Its touches were folder-wide because they were a first guess written with the story, before any task named a file. At 11:15Z S-0221 started and holds all seven that remain.

## Remediation

Not designed yet. Directions to weigh:

- The planner narrows a folder claim to the files a story will change, where its tasks already name them.
- A document every story edits is split by section or by command, or is generated, so that two stories change different files.
- Whether a story in review still needs to hold others is a question for ADR-0046: its branch is finished and synced, and a later story syncs onto main when it is accepted.

Story S-0295 remediates this issue, created from it at 2026-10-06T11:44:50Z.
Closed 2026-10-06T19:22:47Z: S-0295, ADR-0096: only a story in progress holds a ready story (one in review holds nothing); a story's folder touch is narrowed in its claim to the files its tasks name; an overlap wholly inside the manifest's claims.shared patterns holds nothing. flai/internal/workitem/hold_i0087_test.go rebuilds the board of 2026-10-06 and shows the ready stories clear, and each rule's revert holding one again.
