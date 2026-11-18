---
id: I-0377
title: workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it
class: defect
status: closed
count: 1
cost: 15m
first_reported: 2026-11-15T01:36:52Z
last_reported: 2026-11-15T01:36:52Z
updated: 2026-11-18T08:08:32Z
---

# I-0377 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it

## Description
workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it

## Instances

### 2026-11-15T01:36:52Z
Story: S-0744.
S-0744 (T-3227): on a workflow whose accepted state is shipped, flai accept walks the story to shipped and then PlanArchive (flai/internal/workitem/archive.go:73) refuses it: "S-0002 is shipped; only done or cancelled items can be archived", since Item.Closed() compares the status with the names done and cancelled. No open story of E-0034 changes Item.Closed or PlanArchive (S-0739 reads kinds in follow, holds, order, and plan; S-0738 the move rules). The cmd test of S-0744 checks the refusal and the dry run on the renamed board, not a whole acceptance. Fix: give Closed and PlanArchive the repository's workflow and read the accepted and cancelled kinds.

## Remediation

Story S-0842 remediates this issue, created from it at 2026-11-17T08:32:53Z.
Closed 2026-11-18T08:08:32Z: S-0842 fixed it. T-3572: Repo.List, Repo.Get, and Repo.ChildrenOf give each item the repository's workflow, so Item.Closed, CompletedAt, and PlanArchive read the accepted and cancelled kinds, and PlanArchive's refusal names the workflow's states. T-3576: TestAcceptArchivesAStoryInTheAcceptedStateOfTheWorkflow and TestArchiveTakesAStoryInTheAcceptedStateOfTheWorkflow in flai/cmd run a whole flai accept and flai archive on a board whose accepted state is shipped. Both fail on main's workitem with this issue's refusal and pass with the fix.
