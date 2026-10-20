---
id: I-0169
title: message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new
class: efficiency
status: closed
count: 1
cost: 10m
first_reported: 2026-10-20T00:34:27Z
last_reported: 2026-10-20T00:34:27Z
updated: 2026-10-20T19:06:26Z
---

# I-0169 message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new

## Description
message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new

## Instances

### 2026-10-20T00:34:27Z
Story: S-0403.
S-0403: S-0402's agent asked S-0403 to record a split of flai/cmd/accept.go with message_share on MS-0148, since flai task new refused S-0402's task on the overlap; message_share refused because neither story was ready, so the agreed split could not be recorded and S-0402's task stayed refused.

## Remediation

Story S-0451 remediates this issue, created from it at 2026-10-20T18:59:50Z.
Closed 2026-10-20T19:06:26Z: fixed by S-0451: the refusal came from itemnew.Create, behind flai task new and item_new, which refused the advisory wip.overlap a new task brings, while flai edit and flai touches keep it as a warning (ADR-0019). It now keeps it as a warning too, so two stories in progress that agree a split need no share, which lifts only a ready story's hold (ADR-0134)
