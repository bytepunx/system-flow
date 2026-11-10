---
id: I-0328
title: A checked edit of a work item writes, checks, and restores it with no lock, so a reader in between sees an edit the check refuses
class: impression
status: open
count: 1
first_reported: 2026-11-09T23:20:38Z
last_reported: 2026-11-09T23:20:38Z
updated: 2026-11-10T05:45:23Z
---

# I-0328 A checked edit of a work item writes, checks, and restores it with no lock, so a reader in between sees an edit the check refuses

## Description
A checked edit of a work item writes, checks, and restores it with no lock, so a reader in between sees an edit the check refuses

## Instances

### 2026-11-09T23:20:38Z
Story: S-0633.
Found while fixing I-0301 in S-0633. A checked edit (item_edit, flai touches, flai edit) writes the item, runs the check, and puts it back with Repo.Restore when the check refuses it. ADR-0272 closed that window for a checked creation only: the creation holds the checked-write lock, and flai serve's look reads under it. An edit that widens a ready story's touches, or a task's, and is then refused could still be read by a look in the window, which would judge the claim by touches that never stood. No instance shows it yet. Taking LockCheckedWrites from the edit's write until it keeps or restores the item would close it.

## Remediation

Story S-0654 remediates this issue, created from it at 2026-11-10T05:45:23Z.
