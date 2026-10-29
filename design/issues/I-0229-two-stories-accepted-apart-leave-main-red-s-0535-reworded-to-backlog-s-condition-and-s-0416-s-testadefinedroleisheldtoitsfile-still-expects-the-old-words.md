---
id: I-0229
title: "Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words"
class: defect
status: open
count: 1
cost: 30m
first_reported: 2026-10-29T03:27:16Z
last_reported: 2026-10-29T03:27:16Z
updated: 2026-10-29T03:27:16Z
---

# I-0229 Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words

## Description
Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words

## Instances

### 2026-10-29T03:27:16Z
Story: S-0519.
Merge queue batch 20261029-031456 failed S-0519 at integration (TH-0730) on flai/internal/guard/roles_test.go:139 TestADefinedRoleIsHeldToItsFile, a package S-0519 does not touch. The test fails on main alone (flai test flai/internal/guard in the main checkout at c24575c2e). S-0535 (ba74a6f94) changed the to-backlog condition's Allows in flai/internal/roles/conditions.go to "on a call, a move to backlog and nowhere else; on a trigger, only a move to backlog". S-0416, accepted after it (afbc75b15), added a test that expects the old "(a move to backlog and nowhere else)". Neither story's close-out saw the other's change, and the queue sent back the next story it batched. S-0519 T-2323 updates the test's expected text. The queue then failed a story that did not cause the break, since every batch carries main's failure.

## Remediation
