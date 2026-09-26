---
id: I-0048
title: flai's serve tests removed their temp folders while the stub agents they started still wrote to them, so a release run failed
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-09-26T06:32:26Z
last_reported: 2026-09-26T06:32:26Z
updated: 2026-09-26T06:59:35Z
---

# I-0048 flai's serve tests removed their temp folders while the stub agents they started still wrote to them, so a release run failed

## Description
flai's serve tests removed their temp folders while the stub agents they started still wrote to them, so a release run failed

## Instances

### 2026-09-26T06:32:26Z
2026-09-26: the flai/v1.18.3 release run failed in goreleaser's test step on TestAReadyStorysAgentIsStartedOnTheOperatorsWord (TempDir RemoveAll cleanup: directory not empty). Reproduced locally within 5 runs. Stubs started on the operator's word are handed over and waited for by nobody, and the launcher records an ended agent after the test returns. Rerunning the job passed; T-0450 (S-0122) makes newAgentLab wait for them.

## Remediation
Closed 2026-09-26T06:59:35Z: S-0122 T-0450 (cbd211a): newAgentLab waits for its stubs and the launcher before its temp folders go; 20 race runs pass
