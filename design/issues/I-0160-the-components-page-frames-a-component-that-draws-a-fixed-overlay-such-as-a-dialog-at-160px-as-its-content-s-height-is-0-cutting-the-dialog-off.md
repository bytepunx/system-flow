---
id: I-0160
title: The Components page frames a component that draws a fixed overlay, such as a dialog, at 160px, as its content's height is 0, cutting the dialog off
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-10-17T05:32:30Z
last_reported: 2026-10-17T05:32:30Z
updated: 2026-10-17T05:32:30Z
---

# I-0160 The Components page frames a component that draws a fixed overlay, such as a dialog, at 160px, as its content's height is 0, cutting the dialog off

## Description
The Components page frames a component that draws a fixed overlay, such as a dialog, at 160px, as its content's height is 0, cutting the dialog off

## Instances

### 2026-10-17T05:32:30Z
Story: S-0444.
S-0444's pass over flaiover's components found that each dialog, such as AcceptConfirm, AgentStopConfirm, LaneMoveDialog, and WipLimitDialog, draws a fixed inset-0 overlay, so the harness's size message reports a height of 0. ComponentFrame.svelte takes a height of 0 as no height yet and frames the harness at UNSIZED, 160px, so a dialog centred in the overlay is cut off at its top and bottom when it is taller than that. The harness could measure the overlay's content, or the page could take a zero height from a rendered component as a viewport's.

## Remediation
