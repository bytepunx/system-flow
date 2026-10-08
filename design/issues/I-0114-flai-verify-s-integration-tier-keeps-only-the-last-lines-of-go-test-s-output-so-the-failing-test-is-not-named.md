---
id: I-0114
title: flai verify's integration tier keeps only the last lines of go test's output, so the failing test is not named
class: efficiency
status: open
count: 2
cost: 13m
first_reported: 2026-10-07T09:20:56Z
last_reported: 2026-10-08T00:43:55Z
updated: 2026-10-08T00:43:55Z
---

# I-0114 flai verify's integration tier keeps only the last lines of go test's output, so the failing test is not named

## Description
flai verify's integration tier keeps only the last lines of go test's output, so the failing test is not named

## Instances

### 2026-10-07T09:20:56Z
Story: S-0246.
S-0246's close-out stopped twice at the integration tier with findings that held only the last lines of go test's output, every package shown passing, so the failing test, TestReferenceIsCurrent in flai/cmd (docs/users/flai-reference.md stale after a hand edit), was found only by running scripts/integration.sh by hand.

### 2026-10-08T00:43:55Z
Story: S-0324.
S-0324's close-out failed its integration tier after 3m49s. Neither the output, which ended in a list of passing packages and FAIL, nor `flai verify S-0324 --last` named the failure. A hand run of the tier's go test, filtered to failures, named TestRepositoryLintsClean in flai/internal/mdlint.

## Remediation

Story S-0327 remediates this issue, created from it at 2026-10-07T18:59:59Z.
