---
id: I-0062
title: flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number
class: efficiency
status: open
count: 1
cost: 10m
first_reported: 2026-10-03T07:58:20Z
last_reported: 2026-10-03T07:58:20Z
updated: 2026-10-03T07:58:20Z
---

# I-0062 flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number

## Description
flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number

## Instances

### 2026-10-03T07:58:20Z
Story: S-0201.
S-0201's ADR: flai adr new in its worktree gives 0075, but story/S-0207 (in review) already has ADR-0075 and story/S-0200 (in progress) ADR-0076. Whichever is accepted later renumbers its ADR at the rebase and fixes every link to it by hand (see I-0034). flai could number past the ADRs on open story branches.

## Remediation
