---
id: I-0337
title: "flai check finds `adr.index` outside the story at close-out"
class: efficiency
status: open
count: 1
first_reported: 2026-11-10T09:11:58Z
last_reported: 2026-11-10T09:11:58Z
updated: 2026-11-10T09:11:58Z
---

# I-0337 flai check finds `adr.index` outside the story at close-out

## Description
flai check finds `adr.index` outside the story at close-out

## Instances

### 2026-11-10T09:11:58Z
Story: S-0685.
flai check found outside the story:
`design/adrs/README.md`: the ADR index differs from what its ADR files render: rows 0101 and 0122 differ; run flai check --fix or flai adr index, which rewrite it

## Remediation
