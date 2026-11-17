---
id: I-0400
title: TestPublishFormulaUnderJudgement in flai/cmd reads design/workflows/formulas/publish.md, which main removed as S-0837's stopgap, so every flai test that selects flai/cmd fails on main
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-11-17T08:37:17Z
last_reported: 2026-11-17T08:37:17Z
updated: 2026-11-17T09:04:54Z
---

# I-0400 TestPublishFormulaUnderJudgement in flai/cmd reads design/workflows/formulas/publish.md, which main removed as S-0837's stopgap, so every flai test that selects flai/cmd fails on main

## Description
TestPublishFormulaUnderJudgement in flai/cmd reads design/workflows/formulas/publish.md, which main removed as S-0837's stopgap, so every flai test that selects flai/cmd fails on main

## Instances

### 2026-11-17T08:37:17Z
Story: S-0752.
S-0752's flai test on flai/internal/metrics selected ./cmd as an importer and failed at flai/cmd/publish_judgement_test.go:72 (open ../../design/workflows/formulas/publish.md: no such file or directory). Commit 4052fc603 removed the file on main; S-0837, in the backlog, retires the formula. Not S-0752's change.

## Remediation
Closed 2026-11-17T09:04:54Z: Duplicate of I-0401. S-0827's T-3557 fixed it: TestPublishFormulaUnderJudgement reads the publish formula the template ships.
