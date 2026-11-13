---
id: I-0358
title: flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-13T08:42:15Z
last_reported: 2026-11-13T08:42:15Z
updated: 2026-11-13T08:42:15Z
---

# I-0358 flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs

## Description
flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs

## Instances

### 2026-11-13T08:42:15Z
Story: S-0660.
Found by S-0660's T-2882: guard's subcommands() (flai/internal/guard/guard_refusals.go) skips a value only for flags in flaiValues, but cobra reads a subcommand's local flag wherever it stands and treats an unknown flag before the command as taking the next word. So `flai --wait 5m kluster deploy S-0001` read 5m as the command and passed the kluster rule while flai would run a deploy. S-0660 added --wait to flaiValues; other value flags, such as --by or --reason before the command, have the same gap for every rule.

## Remediation
