---
id: I-0014
title: Charts page keeps showing the previous chart when switching chart type until a manual refresh
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-09-18T16:15:39Z
last_reported: 2026-09-18T16:15:39Z
updated: 2026-09-18T16:25:56Z
---

# I-0014 Charts page keeps showing the previous chart when switching chart type until a manual refresh

## Description
Charts page keeps showing the previous chart when switching chart type until a manual refresh

## Instances

### 2026-09-18T16:15:39Z
Reported by the operator on 2026-09-18 against flaiover 0.10.x. Tracked by S-0045.

## Remediation
Closed 2026-09-18T16:25:56Z: fixed in S-0045: the chart redraw effect now tracks its option; empty report lists no longer throw
