---
id: I-0056
title: flai's wip markdown lint does not flag a bare email address, so a thread entry with one reached main and failed a story's close-out
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-02T16:11:21Z
last_reported: 2026-10-02T16:11:21Z
updated: 2026-10-05T00:03:14Z
---

# I-0056 flai's wip markdown lint does not flag a bare email address, so a thread entry with one reached main and failed a story's close-out

## Description
flai's wip markdown lint does not flag a bare email address, so a thread entry with one reached main and failed a story's close-out

## Instances

### 2026-10-02T16:11:21Z
2026-10-02: TH-0067's entries of 15:35Z (alex) and 15:40Z (the operator's session) carry the bare address `alex@robsonandmilligan.com` four times; S-0231's acceptance (f1b6e13) committed them to main, and S-0191's close-out stopped at the markdown lint on them (MD034, lines 33 and 36), outside the story. flai's MD034 (flai/internal/mdlint/inline.go, bareURL) knows only http and https literals, while markdownlint also flags a bare email address, so flai took the entries without a finding and flai check reports none. Fixed on main in one chore commit under TH-0017's answer A, as TH-0056 asks for findings outside a story. Remedy to consider: bareURL also matches GFM's extended email autolink.

## Remediation

Story S-0265 remediates this issue, created from it at 2026-10-05T00:03:14Z.
