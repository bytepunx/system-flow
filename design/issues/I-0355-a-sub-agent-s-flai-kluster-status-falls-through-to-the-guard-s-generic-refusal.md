---
id: I-0355
title: A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-11-13T07:33:55Z
last_reported: 2026-11-13T07:33:55Z
updated: 2026-11-17T13:12:10Z
---

# I-0355 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run

## Description
A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run

## Instances

### 2026-11-13T07:33:55Z
Story: S-0660.
Found by S-0660's planner and agreed out of scope on TH-0987: `design/system/kluster.md` says `flai kluster status` is a read any agent may run, but a sub-agent's call falls through to the generic refusal in `flai/internal/guard/guard.go`, since no read table lists it. Recorded at S-0660's start, as T-2882's notes ask.

## Remediation

Story S-0829 remediates this issue, created from it at 2026-11-17T08:17:21Z.
Closed 2026-11-17T13:12:10Z: S-0829 fixed it: flai guard's cliReads, the table of flai commands a sub-agent may run, now lists kluster status and deployments, so every sub-agent and strategic role reads the cluster as kluster.md says, where a sub-agent's call fell through to the generic refusal; TestEverySubAgentAndRoleReadsTheKluster in flai/internal/guard/kluster_test.go reproduces it on main
