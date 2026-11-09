---
id: I-0282
title: A formula's instance declares no touches, so the story that started it holds it until that story leaves in-progress
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-11-07T08:04:18Z
last_reported: 2026-11-07T08:04:18Z
updated: 2026-11-09T17:33:37Z
---

# I-0282 A formula's instance declares no touches, so the story that started it holds it until that story leaves in-progress

## Description
A formula's instance declares no touches, so the story that started it holds it until that story leaves in-progress

## Instances

### 2026-11-07T08:04:18Z
Story: S-0564.
S-0564's agent ran flai formula run analysis-cadence --ready --autocommit; the instance S-0584 entered ready with no touches and was held (no-touches) by S-0562 and S-0564 in progress, which clears only when they move to review. S-0564 waits on S-0584's acceptance, so neither could move. S-0585, S-0562's publish instance, was held the same way. Cleared by hand with flai touches S-0584 --add on its own story file. flai formula run could give the instance a touch of its own, or the no-touches hold could pass a story whose branch a formula says changes nothing.

## Remediation

Story S-0620 remediates this issue, created from it at 2026-11-07T23:39:44Z.
Closed 2026-11-09T17:33:37Z: S-0620: flai formula run now gives the instance story its own story file as its touch before --ready moves it (formulaRun, flai/cmd/formula.go, commit e1aeff6), so its claim is not empty and no story in progress holds it no-touches. flai import --formula and the orchestrator's publish go through the same path. TestFormulaRunInstanceIsNotHeldByAStoryInProgress reproduces the hold and fails without the fix.
