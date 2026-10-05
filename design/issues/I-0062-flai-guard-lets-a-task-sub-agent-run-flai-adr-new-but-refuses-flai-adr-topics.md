---
id: I-0062
title: flai guard lets a task sub-agent run flai adr new but refuses flai adr topics
class: defect
status: open
count: 2
cost: 4m
first_reported: 2026-10-03T07:27:30Z
last_reported: 2026-10-05T00:57:35Z
updated: 2026-10-05T00:57:35Z
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

## Remediation
