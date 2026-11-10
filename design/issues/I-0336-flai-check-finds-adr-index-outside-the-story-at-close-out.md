---
id: I-0336
title: "flai check finds `adr.index` outside the story at close-out"
class: efficiency
status: open
count: 1
first_reported: 2026-11-10T09:08:17Z
last_reported: 2026-11-10T09:08:17Z
updated: 2026-11-10T09:08:17Z
---

# I-0336 flai check finds `adr.index` outside the story at close-out

## Description
flai check finds `adr.index` outside the story at close-out

## Instances

### 2026-11-10T09:08:17Z
Story: S-0680.
flai check found outside the story:
`design/adrs/README.md`: the ADR index differs from what its ADR files render: rows 0101 and 0122 differ; run flai check --fix or flai adr index, which rewrite it

## Remediation
