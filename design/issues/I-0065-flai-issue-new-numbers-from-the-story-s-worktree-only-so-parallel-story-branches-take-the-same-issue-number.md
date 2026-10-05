---
id: I-0065
title: flai issue new numbers from the story's worktree only, so parallel story branches take the same issue number
class: defect
status: closed
count: 3
cost: 4m
first_reported: 2026-10-03T18:11:17Z
last_reported: 2026-10-04T04:54:27Z
updated: 2026-10-05T00:04:31Z
---

# I-0065 flai issue new numbers from the story's worktree only, so parallel story branches take the same issue number

## Description
flai issue new numbers from the story's worktree only, so parallel story branches take the same issue number

## Instances

### 2026-10-03T18:11:17Z
Story: S-0203.
S-0203's flai issue new took I-0063, which story/S-0201 already holds for its own issue; renumbered to I-0064 by hand after checking every story branch. I-0063 (ADR numbers) is the same defect for ADRs.

### 2026-10-03T19:00:10Z
Story: S-0206.
S-0206 and S-0204 each recorded a different issue as I-0066 from their own worktrees within minutes; TH-0093's reply named the collision, and S-0206 renumbered its own to I-0067 by hand.

### 2026-10-04T04:54:27Z
Story: S-0209.
S-0209 numbered an issue I-0070 in its worktree while another story's I-0070 reached main; the rebase onto S-0225's acceptance stopped on design/issues/summary.md, and the issue was renumbered I-0071 by hand.

## Remediation
Closed 2026-10-05T00:04:31Z: S-0252: flai issue new numbers one past the highest issue on main, in every story worktree, and on every story branch (storygit.FolderNames), with a test reproducing the collision
