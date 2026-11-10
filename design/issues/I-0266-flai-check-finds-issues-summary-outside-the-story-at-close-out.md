---
id: I-0266
title: "flai check finds `issues.summary` outside the story at close-out"
class: efficiency
status: open
count: 1
first_reported: 2026-11-07T00:11:02Z
last_reported: 2026-11-07T00:11:02Z
updated: 2026-11-10T08:55:37Z
---

# I-0266 flai check finds `issues.summary` outside the story at close-out

## Description
flai check finds `issues.summary` outside the story at close-out

## Instances

### 2026-11-07T00:11:02Z
Story: S-0576.
flai check found outside the story:
`design/issues/summary.md`: summary.md differs from what its issue files render; run flai check --fix or flai issue summary, which rewrite it

## Remediation

Story S-0702 remediates this issue, created from it at 2026-11-10T08:55:37Z.
