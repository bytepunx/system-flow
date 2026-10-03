---
id: I-0062
title: Two stories in progress at once take the same ADR number, and one renumbers by hand
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-10-03T07:37:37Z
last_reported: 2026-10-03T07:37:37Z
updated: 2026-10-03T07:37:37Z
---

# I-0062 Two stories in progress at once take the same ADR number, and one renumbers by hand

## Description
Two stories in progress at once take the same ADR number, and one renumbers by hand

## Instances

### 2026-10-03T07:37:37Z
Story: S-0200.
flai adr new numbers from the files on the story's branch, so S-0200 and S-0207 both made ADR-0075. flai stream sync's trial merge found it only as a conflict in design/adrs/README.md (TH-0084); S-0200 renamed its ADR to 0076 and fixed eight links by hand.

## Remediation
