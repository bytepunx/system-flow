---
id: I-0063
title: flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number
class: efficiency
status: open
count: 2
cost: 13m
first_reported: 2026-10-03T07:58:20Z
last_reported: 2026-10-03T17:55:12Z
updated: 2026-10-03T17:55:12Z
---

# I-0063 flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number

## Description
flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number

## Instances

### 2026-10-03T07:58:20Z
Story: S-0201.
S-0201's ADR: flai adr new in its worktree gives 0075, but story/S-0207 (in review) already has ADR-0075 and story/S-0200 (in progress) ADR-0076. Whichever is accepted later renumbers its ADR at the rebase and fixes every link to it by hand (see I-0034). flai could number past the ADRs on open story branches.

### 2026-10-03T17:52:55Z
Story: S-0200.
S-0200 and S-0207 both made ADR-0075 on their branches; S-0200 renamed its own to ADR-0076 and fixed eight links by hand (TH-0084). The issues have the same defect: S-0200's own record of this went in as I-0062 too, beside two I-0062 already on main, and was folded into this one.

### 2026-10-03T17:56:54Z
Story: S-0200.
Again in S-0200: its second ADR, for TH-0082's answer, took 0077, which story/S-0201 already held. Renamed to ADR-0078 and fixed ten files' links by hand (TH-0085).

### 2026-10-03T17:55:12Z
Story: S-0201.
Issue IDs collide the same way. S-0201's flai issue new, run in the main checkout, took I-0062 while story/S-0207 had already recorded its own I-0062. S-0207's acceptance committed both, so flai check errored on main with issues.duplicate-id until S-0201 renumbered its issue to I-0063 by hand. Also S-0201's ADR-0075 was renumbered to ADR-0077 at its rebase.

## Remediation
