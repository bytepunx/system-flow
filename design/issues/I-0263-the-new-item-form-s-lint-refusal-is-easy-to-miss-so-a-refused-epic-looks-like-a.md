---
id: I-0263
title: The new item form's lint refusal is easy to miss, so a refused epic looks like a hang
class: defect
status: open
count: 1
cost: 20m
first_reported: 2026-11-06T21:31:26Z
last_reported: 2026-11-06T21:31:26Z
updated: 2026-11-10T08:55:31Z
---

# I-0263 The new item form's lint refusal is easy to miss, so a refused epic looks like a hang

## Description
The new item form's lint refusal is easy to miss, so a refused epic looks like a hang

## Instances

### 2026-11-06T21:31:26Z
Creating E-0032 from the dashboard twice (21:21 and 21:27 UTC) got a 422 in a few seconds, the markdown lint rejecting two bare URLs (MD034), but the operator saw the form hang indefinitely. NewItemForm renders the refusal alert below the body preview, under the fold on a long body, and nothing scrolls to it or marks the Create button.

## Remediation

Story S-0700 remediates this issue, created from it at 2026-11-10T08:55:31Z.
