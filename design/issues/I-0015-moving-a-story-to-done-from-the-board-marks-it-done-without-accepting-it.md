---
id: I-0015
title: Moving a story to done from the board marks it done without accepting it
class: defect
status: closed
count: 1
cost: 15m
first_reported: 2026-09-18T16:33:03Z
last_reported: 2026-09-18T16:33:03Z
updated: 2026-09-18T16:49:26Z
---

# I-0015 Moving a story to done from the board marks it done without accepting it

## Description
Moving a story to done from the board marks it done without accepting it

## Instances

### 2026-09-18T16:33:03Z
The operator dragged S-0045 to done on 2026-09-18 expecting acceptance. flai move only changed the state: the story branch stayed unmerged, nothing was archived or released, and flai accept then refused the item as already done. Landed by hand with accept's own steps. Tracked by S-0046.

## Remediation
Closed 2026-09-18T16:49:26Z: S-0046: moving a story to done runs acceptance; accept completes stranded stories; check flags them
