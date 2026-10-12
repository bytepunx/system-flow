---
id: I-0126
title: flai task done commits the changed paths of another open task of the same layer
class: defect
status: closed
count: 1
cost: 3m
first_reported: 2026-10-08T10:36:24Z
last_reported: 2026-10-08T10:36:24Z
updated: 2026-10-12T14:41:37Z
---

# I-0126 flai task done commits the changed paths of another open task of the same layer

## Description
flai task done commits the changed paths of another open task of the same layer

## Instances

### 2026-10-08T10:36:24Z
Story: S-0347.
S-0347's T-1377 and T-1379 ran together in one layer with disjoint touches. Closing T-1377 with the MCP tool task_done (installed flai 1.39.13) committed T-1379's three changed paths in T-1377's commit and widened T-1377's touches to them, though ADR-0128 says task done leaves a path only another open task covers. T-1379 then closed with nothing to commit, and T-1377's touches were narrowed back by hand with flai touches --remove.

## Remediation

Story S-0362 remediates this issue, created from it at 2026-10-12T14:12:29Z.
Closed 2026-10-12T14:41:37Z: Fixed by S-0322, which implemented ADR-0128 and was published in flai 1.39.14 at 2026-10-08T09:30Z. The instance closed T-1377 at 10:35:55Z through an MCP server still running the installed flai 1.39.13, which committed every changed path. S-0362's tests show the current flai task done and MCP task_done leave an open sibling's paths uncommitted, list them in left, and keep them out of the closing task's and the story's touches: TestTaskDoneLeavesASiblingsConventionTemplateAndChangelog (flai/internal/mcpserver/task_test.go) reproduces the instance's layer, and TestRunLeavesOnlyAnOpenSiblingsPaths (flai/internal/taskdone/taskdone_test.go) covers backlog, ready, done, and cancelled siblings. Both fail when commit stages every changed path and pass on the current code, so taskdone.go needed no fix. The guard against an MCP flai older than the project's history is S-0363's, from I-0127.
