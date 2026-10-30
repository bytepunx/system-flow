---
id: I-0203
title: A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration
class: defect
status: closed
count: 2
cost: 18m
first_reported: 2026-10-24T00:15:33Z
last_reported: 2026-10-29T03:38:03Z
updated: 2026-10-30T00:57:55Z
---

# I-0203 A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration

## Description
A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration

## Instances

### 2026-10-24T00:15:33Z
Story: S-0490.
S-0491 added `orchestration.permissions.edit_subprojects` to `flai manifest set`'s help without running `make flai-reference`. Its close-out left `integration` to the merge queue (ADR-0174), but the operator accepted it directly (moved to done at 2026-10-23T23:38:50Z), so `integration` never ran on it. Main at `be3af0d` was left with a stale `docs/users/flai-reference.md`. The next batch, 20261024-000635 with S-0490 alone, failed at `TestReferenceIsCurrent` and sent S-0490 back on TH-0632, though S-0490 changes no command. S-0490 regenerated the reference on its branch to get through.

### 2026-10-29T03:38:03Z
Story: S-0538.
S-0416 was accepted outside the merge queue (afbc75b) after S-0535 had changed the to-backlog condition's wording in flai/internal/roles/conditions.go. TestADefinedRoleIsHeldToItsFile in flai/internal/guard/roles_test.go still expected the old wording, so main failed integration. Merge queue batch 20261029-032429 failed S-0538 at integration, though S-0538 changes only flai/internal/adr; S-0519 was sent back the same way. Both stories fixed the expectation with the same one-line change (S-0538's T-2324, S-0519's T-2323).

## Remediation

Story S-0529 remediates this issue, created from it at 2026-10-28T12:54:24Z.
Closed 2026-10-30T00:57:55Z: S-0529: while the merge queue can batch, a hand acceptance (flai accept, flai move <S> done, the board, the dashboard) is refused before anything is merged until the story's last flai verify passed at its branch's head with every tier marked queue, naming the tiers and flai verify <S> --all, which now reuses the tiers the close-out passed and runs only the queue tiers; --without-queue-tiers lands it anyway and records the tiers in its Notes (ADR-0221). Reproduced by TestAHandAcceptanceWaitsForTheQueueTiersWhileTheQueueBatches in flai/cmd/accept_queue_tiers_test.go.
