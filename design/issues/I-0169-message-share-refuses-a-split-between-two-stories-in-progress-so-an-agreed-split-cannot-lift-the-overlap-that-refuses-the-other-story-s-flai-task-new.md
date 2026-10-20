---
id: I-0169
title: message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new
class: efficiency
status: open
count: 1
cost: 10m
first_reported: 2026-10-20T00:34:27Z
last_reported: 2026-10-20T00:34:27Z
updated: 2026-10-20T00:34:27Z
---

# I-0169 message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new

## Description
message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new

## Instances

### 2026-10-20T00:34:27Z
Story: S-0403.
S-0403: S-0402's agent asked S-0403 to record a split of flai/cmd/accept.go with message_share on MS-0148, since flai task new refused S-0402's task on the overlap; message_share refused because neither story was ready, so the agreed split could not be recorded and S-0402's task stayed refused.

## Remediation
