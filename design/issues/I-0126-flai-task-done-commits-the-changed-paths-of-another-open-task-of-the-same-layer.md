---
id: I-0126
title: flai task done commits the changed paths of another open task of the same layer
class: defect
status: open
count: 1
cost: 3m
first_reported: 2026-10-08T10:36:24Z
last_reported: 2026-10-08T10:36:24Z
updated: 2026-10-08T10:36:24Z
---

# I-0126 flai task done commits the changed paths of another open task of the same layer

## Description
flai task done commits the changed paths of another open task of the same layer

## Instances

### 2026-10-08T10:36:24Z
Story: S-0347.
S-0347's T-1377 and T-1379 ran together in one layer with disjoint touches. Closing T-1377 with the MCP tool task_done (installed flai 1.39.13) committed T-1379's three changed paths in T-1377's commit and widened T-1377's touches to them, though ADR-0128 says task done leaves a path only another open task covers. T-1379 then closed with nothing to commit, and T-1377's touches were narrowed back by hand with flai touches --remove.

## Remediation
