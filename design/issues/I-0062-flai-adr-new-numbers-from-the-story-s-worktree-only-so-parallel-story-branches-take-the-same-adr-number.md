---
id: I-0062
title: flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number
class: efficiency
status: open
count: 2
cost: 8m
first_reported: 2026-10-03T07:58:20Z
last_reported: 2026-10-03T17:52:55Z
updated: 2026-10-03T17:52:55Z
---

# I-0062 flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number

## Description
flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number

## Instances

### 2026-10-03T07:58:20Z
Story: S-0201.
S-0201's ADR: flai adr new in its worktree gives 0075, but story/S-0207 (in review) already has ADR-0075 and story/S-0200 (in progress) ADR-0076. Whichever is accepted later renumbers its ADR at the rebase and fixes every link to it by hand (see I-0034). flai could number past the ADRs on open story branches.

### 2026-10-03T17:52:55Z
Story: S-0200.
S-0200 and S-0207 both made ADR-0075 on their branches; S-0200 renamed its own to ADR-0076 and fixed eight links by hand (TH-0084). The issues have the same defect: S-0200's own record of this went in as I-0062 too, beside two I-0062 already on main, and was folded into this one.

## Remediation
