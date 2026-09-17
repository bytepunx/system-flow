---
id: I-004
title: Hand-edited work items drifted from the schema
class: defect
status: open
count: 3
cost: 10m
first_reported: 2026-09-15T18:00:57Z
last_reported: 2026-09-15T18:11:49Z
updated: 2026-09-15T22:40:33Z
---

# I-004 Hand-edited work items drifted from the schema

## Description
Before flai owned work items, front matter was written by hand and drifted: an unquoted colon in a title broke parsing, an archived story kept an unchecked criterion, and the board order listed backlog stories against the rule of the time.

## Instances

### 2026-09-15T18:00:57Z
S-0007: round-trip test failed on `title: Work item package: parse, write, list, IDs`; quoted it and made item templates quote titles.

### 2026-09-15T18:11:49Z
S-0008: first `flai check` run found the unchecked criterion in archived S-0002.
S-0008: same run flagged backlog stories in the board order; the rule was relaxed instead.

## Remediation
flai now creates and moves items and `flai check --strict` runs before every hand-over (tooling.md). Close when no new instance appears for an epic.
