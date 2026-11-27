---
id: I-0419
title: "flai check finds `issues.summary` outside the story at close-out"
class: efficiency
status: open
count: 1
first_reported: 2026-11-27T22:49:52Z
last_reported: 2026-11-27T22:49:52Z
updated: 2026-11-27T22:49:52Z
---

# I-0419 flai check finds `issues.summary` outside the story at close-out

## Description
flai check finds `issues.summary` outside the story at close-out

## Instances

### 2026-11-27T22:49:52Z
Story: S-0850.
flai check found outside the story:
`design/issues/summary.md`: summary.md differs from what its issue files render; run flai check --fix or flai issue summary, which rewrite it

## Remediation
