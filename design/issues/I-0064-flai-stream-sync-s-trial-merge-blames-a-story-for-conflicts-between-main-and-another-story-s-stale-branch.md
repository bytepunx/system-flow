---
id: I-0064
title: flai stream sync's trial merge blames a story for conflicts between main and another story's stale branch
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-10-03T18:11:00Z
last_reported: 2026-10-03T18:11:00Z
updated: 2026-10-03T18:11:00Z
---

# I-0064 flai stream sync's trial merge blames a story for conflicts between main and another story's stale branch

## Description
flai stream sync's trial merge blames a story for conflicts between main and another story's stale branch

## Instances

### 2026-10-03T18:11:00Z
Story: S-0203.
S-0203's first sync opened TH-0088, naming five conflicts with story/S-0201 in design/adrs/README.md, design/issues/I-0063, summary.md, flai-cli.md, and work-hierarchy.md. None is in S-0203's diff (only flai/internal/itemnew and workitem/create.go): they are S-0200's accepted changes on main against S-0201, which had not synced, already settled on TH-0087. The trial merge should compare the two stories' own changes, or leave out what main brings.

## Remediation
