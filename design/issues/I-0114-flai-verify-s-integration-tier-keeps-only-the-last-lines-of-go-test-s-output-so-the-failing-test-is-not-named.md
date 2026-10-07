---
id: I-0114
title: flai verify's integration tier keeps only the last lines of go test's output, so the failing test is not named
class: efficiency
status: open
count: 1
cost: 15m
first_reported: 2026-10-07T09:20:56Z
last_reported: 2026-10-07T09:20:56Z
updated: 2026-10-07T18:59:59Z
---

# I-0114 flai verify's integration tier keeps only the last lines of go test's output, so the failing test is not named

## Description
flai verify's integration tier keeps only the last lines of go test's output, so the failing test is not named

## Instances

### 2026-10-07T09:20:56Z
Story: S-0246.
S-0246's close-out stopped twice at the integration tier with findings that held only the last lines of go test's output, every package shown passing, so the failing test, TestReferenceIsCurrent in flai/cmd (docs/users/flai-reference.md stale after a hand edit), was found only by running scripts/integration.sh by hand.

## Remediation

Story S-0327 remediates this issue, created from it at 2026-10-07T18:59:59Z.
