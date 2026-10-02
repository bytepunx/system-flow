---
id: I-0055
title: A story made from the dashboard reached main with a markdown lint error flai task new would have refused
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-02T12:41:22Z
last_reported: 2026-10-02T12:41:22Z
updated: 2026-10-02T12:41:22Z
---

# I-0055 A story made from the dashboard reached main with a markdown lint error flai task new would have refused

## Description
A story made from the dashboard reached main with a markdown lint error flai task new would have refused

## Instances

### 2026-10-02T12:41:22Z
2026-10-02: S-0231, made by the operator from the dashboard at 12:24Z and committed to main, indents its criteria list one space (MD007, four lines). S-0193's close-out stopped on it in its worktree, where the file is the merge base's. Fixed on main in one chore commit (4a2ccf0) under TH-0017's answer A. I-0027 covered items reaching CI unlinted and is closed; this is the dashboard's new-story path letting a body through that flai task new refuses since ADR-0061.

## Remediation
