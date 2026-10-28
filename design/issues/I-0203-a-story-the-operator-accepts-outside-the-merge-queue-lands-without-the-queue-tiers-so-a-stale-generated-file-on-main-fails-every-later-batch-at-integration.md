---
id: I-0203
title: A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-10-24T00:15:33Z
last_reported: 2026-10-24T00:15:33Z
updated: 2026-10-28T12:54:24Z
---

# I-0203 A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration

## Description
A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration

## Instances

### 2026-10-24T00:15:33Z
Story: S-0490.
S-0491 added `orchestration.permissions.edit_subprojects` to `flai manifest set`'s help without running `make flai-reference`. Its close-out left `integration` to the merge queue (ADR-0174), but the operator accepted it directly (moved to done at 2026-10-23T23:38:50Z), so `integration` never ran on it. Main at `be3af0d` was left with a stale `docs/users/flai-reference.md`. The next batch, 20261024-000635 with S-0490 alone, failed at `TestReferenceIsCurrent` and sent S-0490 back on TH-0632, though S-0490 changes no command. S-0490 regenerated the reference on its branch to get through.

## Remediation

Story S-0529 remediates this issue, created from it at 2026-10-28T12:54:24Z.
