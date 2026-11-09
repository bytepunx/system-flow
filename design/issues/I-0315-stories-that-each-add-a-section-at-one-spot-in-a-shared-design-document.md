---
id: I-0315
title: Stories that each add a section at one spot in a shared design document conflict at every sync and acceptance, and the acceptance sends them back
class: efficiency
status: open
count: 2
cost: 35m
first_reported: 2026-11-09T07:27:44Z
last_reported: 2026-11-09T07:53:55Z
updated: 2026-11-09T08:06:42Z
---

# I-0315 Stories that each add a section at one spot in a shared design document conflict at every sync and acceptance, and the acceptance sends them back

## Description
Stories that each add a section at one spot in a shared design document conflict at every sync and acceptance, and the acceptance sends them back

## Instances

### 2026-11-09T07:27:44Z
Story: S-0599.
E-0033's stories each add a "Cause N fixed" section to `design/system/server-performance.md` at the same spot, after the S-0596 Stories table. Each acceptance of another E-0033 story made S-0599's branch conflict there. S-0599 resolved it by hand twice: once in a close-out's sync after S-0605's acceptance, and once after TH-0915, when S-0599's own acceptance stopped on the file and sent the story back to in-progress, costing a close-out of about 20 minutes. The trial merge reported the same conflict with S-0598, S-0601, S-0603, S-0604, and S-0606 too. Every resolution was mechanical: keep both sections, in cause order. The section merge of ADR-0175 settles such a stop only for paths inside `claims.shared`, and this file is not shared.

### 2026-11-09T07:53:55Z
Story: S-0598.
S-0598 was sent back from review twice, once on TH-0914 for check.go and once on TH-0917 for design/system/server-performance.md. In both the cause 8 section met the cause 9 to 19 sections at the same spot after the S-0596 Stories table. Each resolution was mechanical: keep both sections in cause order. Each send-back cost a 10 to 13 minute close-out and a review cycle. The trial merge still reports the same conflict with S-0601, which is in review.

## Remediation

Story S-0637 remediates this issue, created from it at 2026-11-09T08:06:42Z.
