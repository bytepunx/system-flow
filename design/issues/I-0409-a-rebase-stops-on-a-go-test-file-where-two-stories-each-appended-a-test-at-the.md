---
id: I-0409
title: A rebase stops on a Go test file where two stories each appended a test at the same place, which flai cannot settle, so an acceptance sends the story back
class: efficiency
status: open
count: 1
cost: 15m
first_reported: 2026-11-21T18:16:18Z
last_reported: 2026-11-21T18:16:18Z
updated: 2026-11-27T22:55:28Z
---

# I-0409 A rebase stops on a Go test file where two stories each appended a test at the same place, which flai cannot settle, so an acceptance sends the story back

## Description
A rebase stops on a Go test file where two stories each appended a test at the same place, which flai cannot settle, so an acceptance sends the story back

## Instances

### 2026-11-21T18:16:18Z
Story: S-0812.
S-0809 and S-0812 each appended new test functions after the same function in flai/internal/verify/select_test.go. The trial merge flagged it, and the two agents agreed in MS-0628 that whichever story rebased second would keep both. But the acceptance's rebase and resolve job (TH-1320) cannot settle a Go file, so it sent S-0812 back to in-progress. The agent then resolved it by hand: it kept both sides and added the closing brace between them. That cost a full round of close-out and review. A merge that keeps both sides when each side only adds whole top-level declarations after the same line would settle this without an agent.

## Remediation

Story S-0857 remediates this issue, created from it at 2026-11-27T22:55:28Z.
