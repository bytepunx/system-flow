---
id: I-0355
title: A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-13T07:33:55Z
last_reported: 2026-11-13T07:33:55Z
updated: 2026-11-13T07:33:55Z
---

# I-0355 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run

## Description
A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run

## Instances

### 2026-11-13T07:33:55Z
Story: S-0660.
Found by S-0660's planner and agreed out of scope on TH-0987: `design/system/kluster.md` says `flai kluster status` is a read any agent may run, but a sub-agent's call falls through to the generic refusal in `flai/internal/guard/guard.go`, since no read table lists it. Recorded at S-0660's start, as T-2882's notes ask.

## Remediation
