---
id: I-0062
title: flai guard lets a task sub-agent run flai adr new but refuses flai adr topics
class: defect
status: closed
count: 3
cost: 4m
first_reported: 2026-10-03T07:27:30Z
last_reported: 2026-10-06T06:06:28Z
updated: 2026-10-08T07:22:52Z
---

# I-0062 flai guard lets a task sub-agent run flai adr new but refuses flai adr topics

## Description
flai guard lets a task sub-agent run flai adr new but refuses flai adr topics

## Instances

### 2026-10-03T07:27:30Z
Story: S-0207.
T-0751's task sub-agent ran scripts/flai.sh adr new, which the guard let through and which wrote the ADR and its index row. Then flai adr topics on the same ADR was refused as a sub-agent write, so the story's agent had to set the topics. Either both are writes the guard refuses, or an ADR the task makes in its own touches is the sub-agent's to finish.

### 2026-10-05T00:57:35Z
Story: S-0249.
S-0249's task sub-agent for T-0845 (a general-purpose agent named for the task) was refused flai adr new, even with --print-body, so it wrote ADR-0085 by hand, copying what internal/adr writes: the number, the file name, refines, the index row, and topics.

### 2026-10-06T06:06:28Z
Story: S-0220.
T-0894: the guard refused the task sub-agent's flai adr new, even with --print-body, so the story's agent created ADR-0090 and set its topics; the sub-agent had already cited ADR-0090 by guessing the number.

## Remediation

Story S-0287 remediates this issue, created from it at 2026-10-06T09:56:50Z.
Closed 2026-10-08T07:22:52Z: S-0287 fixed it (ADR-0127): flai guard now lets a sub-agent run flai adr new, flai adr topics, and flai adr accept, and call adr_new, so a task sub-agent records its ADR through flai rather than copying internal/adr by hand or guessing the number. It refuses only --commit, --autocommit, and adr_new's commit, since the commit stays the story's agent's; TestGuardLetsASubAgentWriteAnADRButNotCommitIt in flai/cmd/guard_test.go reproduces the instances. delegation.md and its template copy say so.
