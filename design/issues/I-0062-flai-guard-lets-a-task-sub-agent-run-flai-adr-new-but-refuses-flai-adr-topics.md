---
id: I-0062
title: flai guard lets a task sub-agent run flai adr new but refuses flai adr topics
class: defect
status: open
count: 1
cost: 2m
first_reported: 2026-10-03T07:27:30Z
last_reported: 2026-10-03T07:27:30Z
updated: 2026-10-03T07:27:30Z
---

# I-0062 flai guard lets a task sub-agent run flai adr new but refuses flai adr topics

## Description
flai guard lets a task sub-agent run flai adr new but refuses flai adr topics

## Instances

### 2026-10-03T07:27:30Z
Story: S-0207.
T-0751's task sub-agent ran scripts/flai.sh adr new, which the guard let through and which wrote the ADR and its index row. Then flai adr topics on the same ADR was refused as a sub-agent write, so the story's agent had to set the topics. Either both are writes the guard refuses, or an ADR the task makes in its own touches is the sub-agent's to finish.

## Remediation
