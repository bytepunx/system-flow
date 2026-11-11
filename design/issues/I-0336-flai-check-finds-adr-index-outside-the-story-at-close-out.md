---
id: I-0336
title: "flai check finds `adr.index` outside the story at close-out"
class: efficiency
status: closed
count: 2
first_reported: 2026-11-10T09:08:17Z
last_reported: 2026-11-10T09:11:58Z
updated: 2026-11-11T14:17:32Z
---

# I-0336 flai check finds `adr.index` outside the story at close-out

## Description
flai check finds `adr.index` outside the story at close-out

## Instances

### 2026-11-10T09:08:17Z
Story: S-0680.
flai check found outside the story:
`design/adrs/README.md`: the ADR index differs from what its ADR files render: rows 0101 and 0122 differ; run flai check --fix or flai adr index, which rewrite it

### 2026-11-10T09:11:58Z
Story: S-0685.
flai check found outside the story:
`design/adrs/README.md`: the ADR index differs from what its ADR files render: rows 0101 and 0122 differ; run flai check --fix or flai adr index, which rewrite it

## Remediation

Story S-0710 remediates this issue, created from it at 2026-11-10T09:56:02Z.
Closed 2026-11-11T14:17:32Z: S-0710: flai adr new, accept, and topics with --autocommit now commit the index HEAD's ADR files render with the ones the commit records in their place (adr.stageCommittedIndex), so an ADR file left uncommitted in the main checkout no longer puts a row in the committed index that its committed files do not render; TestAcceptCommitsTheIndexHEADsFilesRender reproduces the cause.
