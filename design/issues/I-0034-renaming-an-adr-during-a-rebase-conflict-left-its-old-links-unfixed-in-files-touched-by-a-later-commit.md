---
id: I-0034
title: Renaming an ADR during a rebase conflict left its old links unfixed in files touched by a later commit
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-09-22T21:34:58Z
last_reported: 2026-09-22T21:34:58Z
updated: 2026-09-22T21:34:58Z
---

# I-0034 Renaming an ADR during a rebase conflict left its old links unfixed in files touched by a later commit

## Description
Renaming an ADR during a rebase conflict left its old links unfixed in files touched by a later commit

## Instances

### 2026-09-22T21:34:58Z
S-0080's ADR was renamed 0032 to 0033 while resolving its rebase conflict onto main (S-0087 had already claimed 0032). The README and the two files touched by the commit under conflict were fixed then, but a later commit in the same rebase (the docs update) introduced three more references to the old number and path, found 2026-09-22 while updating docs for S-0081 and fixed there.

## Remediation
