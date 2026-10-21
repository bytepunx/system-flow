---
id: I-0176
title: A usage write to a story overwrites its move to in-progress made meanwhile, so a started story reads as ready again
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-10-21T07:48:30Z
last_reported: 2026-10-21T07:48:30Z
updated: 2026-10-21T12:48:00Z
---

# I-0176 A usage write to a story overwrites its move to in-progress made meanwhile, so a started story reads as ready again

## Description
A usage write to a story overwrites its move to in-progress made meanwhile, so a started story reads as ready again

## Instances

### 2026-10-21T07:48:30Z
Story: S-0456.
story_start moved S-0456 to in-progress at about 07:28:40Z on 2026-10-21 and answered status in-progress. The story's file on main, committed by S-0457's acceptance (d8fb3eab), holds status ready with its one ready transition and updated 07:33:37Z, a write that changed only its usage (seconds 9639 to 9911, tokens, cost): a read-modify-write of usage from a copy read before the move dropped the in-progress transition. T-1944 was worked and closed while the board showed S-0456 ready; task_done's inbox listed it among the ready stories. Moved it to in-progress again at the time found. S-0463 shows the same: its agent closed T-1977 at 07:36:28Z and moved S-0463 to in-progress again at 07:37:05Z.

## Remediation

Story S-0480 remediates this issue, created from it at 2026-10-21T11:59:52Z.
Closed 2026-10-21T12:48:00Z: S-0479 (ADR-0184), into which the operator merged S-0480 on TH-0577: Repo.Save takes one items lock, reads the file again, and merges a copy read before another writer's save by front-matter key and body section, so a usage write from a copy read before a move keeps the move. TestSaveKeepsAMoveSavedSinceTheCopyWasRead reproduces the stale write.
