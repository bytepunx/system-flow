---
id: I-0187
title: The orchestrator cannot set a cost of delay input the operator decided on an epic or a story past the backlog
class: efficiency
status: open
count: 1
first_reported: 2026-10-21T16:01:28Z
last_reported: 2026-10-21T16:01:28Z
updated: 2026-10-28T12:53:43Z
---

# I-0187 The orchestrator cannot set a cost of delay input the operator decided on an epic or a story past the backlog

## Description
The orchestrator cannot set a cost of delay input the operator decided on an epic or a story past the backlog

## Instances

### 2026-10-21T16:01:28Z
On TH-0579 the operator asked that the orchestrator be allowed to update the cost of delay. flai guard refuses the orchestrator item_edit on an epic (E-0024 at 2026-10-21T16:0xZ: it gives a story cost of delay inputs only), and flai refused it inputs on in-progress S-0464, so S-0473's 190 on TH-0574 went back to the operator. Confirmed inputs still unset: 4h on E-0024 (TH-0567, TH-0473) and 4h on E-0025 (TH-0476). Until they are set, E-0024's stories carry values the planners worked out by hand, and E-0025's ten stories have none, so flai promote --drafts lists them incomplete.

## Remediation

Story S-0523 remediates this issue, created from it at 2026-10-28T12:53:43Z.
