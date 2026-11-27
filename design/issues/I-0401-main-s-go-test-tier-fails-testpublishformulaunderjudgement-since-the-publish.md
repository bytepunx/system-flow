---
id: I-0401
title: Main's go-test tier fails TestPublishFormulaUnderJudgement since the publish formula was removed from design/workflows/formulas ahead of S-0837
class: defect
status: closed
count: 1
cost: 15m
first_reported: 2026-11-17T08:45:07Z
last_reported: 2026-11-17T08:45:07Z
updated: 2026-11-27T22:55:29Z
---

# I-0401 Main's go-test tier fails TestPublishFormulaUnderJudgement since the publish formula was removed from design/workflows/formulas ahead of S-0837

## Description
Main's go-test tier fails TestPublishFormulaUnderJudgement since the publish formula was removed from design/workflows/formulas ahead of S-0837

## Instances

### 2026-11-17T08:45:07Z
Story: S-0827.
S-0827's flai test on flai/internal/verify/select_test.go selects ./cmd as an importer, and go-test fails TestPublishFormulaUnderJudgement (flai/cmd/publish_judgement_test.go:72): open ../../design/workflows/formulas/publish.md: no such file or directory. Main's 4052fc603 removed that file as a stopgap for S-0837, which is still a backlog draft; flai/cmd/orchestrate_permissions_test.go and flai/internal/mcpserver/orchestrate_test.go name the file too. Every story whose change selects ./cmd fails its close-out until the file is back or S-0837 lands.

## Remediation
Closed 2026-11-27T22:55:29Z: No longer reproduces on main at 716793afe: go test ./cmd -run TestPublishFormulaUnderJudgement passes, and the integration tier passed in the close-outs of S-0839 and S-0811 on 2026-11-27
