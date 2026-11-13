---
id: I-0360
title: flai stream sync reports a rebased branch undoes main's line when both changed one long markdown line and the branch keeps every word of main's
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-13T09:21:49Z
last_reported: 2026-11-13T09:21:49Z
updated: 2026-11-13T11:42:36Z
---

# I-0360 flai stream sync reports a rebased branch undoes main's line when both changed one long markdown line and the branch keeps every word of main's

## Description
flai stream sync reports a rebased branch undoes main's line when both changed one long markdown line and the branch keeps every word of main's

## Instances

### 2026-11-13T09:21:49Z
Story: S-0761.
S-0761's acceptance was sent back (TH-1147) because its rebase stopped on `design/system/flai-cli.md`. Both S-0735 (accepted first) and S-0761 had changed line 54, the `flai check` row of § Commands, which is one table row on a single line, and the section merge could not settle that. I resolved the stop by taking main's line and adding S-0761's sentence and list entry to it, and `git diff --word-diff main` shows only additions. The sync then warned that `story/S-0761` "removes 1 line main added to design/system/flai-cli.md since d54d693", and recorded it in `.flai-cache/undone/S-0761.json`. `storygit.FindUndone` compares a multiset of whole line texts. So any branch that changes a line main also changed counts as undoing main's change, even when it keeps every word of main's version, and ADR-0155 then holds the acceptance until the operator passes `--keep-removal` or restores the path. Restoring would drop the story's own edit. A word- or sentence-level check, like ADR-0287's sentence merge, would tell an extension apart from an undo.

## Remediation

Story S-0767 remediates this issue, created from it at 2026-11-13T11:42:36Z.
