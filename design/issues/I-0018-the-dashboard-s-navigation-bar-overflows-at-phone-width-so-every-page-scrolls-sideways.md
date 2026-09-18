---
id: I-0018
title: The dashboard's navigation bar overflows at phone width, so every page scrolls sideways
class: defect
status: open
count: 1
cost: 3m
first_reported: 2026-09-18T20:48:07Z
last_reported: 2026-09-18T20:48:07Z
updated: 2026-09-18T20:48:07Z
---

# I-0018 The dashboard's navigation bar overflows at phone width, so every page scrolls sideways

## Description
The dashboard's navigation bar overflows at phone width, so every page scrolls sideways

## Instances

### 2026-09-18T20:48:07Z
S-0048: looking at /board at 390px, the document was 544px wide. The offenders are the top navigation links ADRs and Search and the theme button in the root layout; nothing in the board or its cards overflowed. Found while checking the parent indicator; not caused by it and not fixed in S-0048.

## Remediation
