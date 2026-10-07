---
id: I-0112
title: flai verify --record-issues opens an issue on the story branch that another branch opened under the same title meanwhile
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-07T09:09:31Z
last_reported: 2026-10-07T09:09:31Z
updated: 2026-10-07T09:09:31Z
---

# I-0112 flai verify --record-issues opens an issue on the story branch that another branch opened under the same title meanwhile

## Description
flai verify --record-issues opens an issue on the story branch that another branch opened under the same title meanwhile

## Instances

### 2026-10-07T09:09:31Z
Story: S-0275.
S-0275's close-out opened I-0112 for its narrative.state notes; S-0265 had opened I-0111 with the same title on main; after a sync the close-out bumped I-0111, leaving I-0112 a duplicate, removed by hand.

## Remediation
