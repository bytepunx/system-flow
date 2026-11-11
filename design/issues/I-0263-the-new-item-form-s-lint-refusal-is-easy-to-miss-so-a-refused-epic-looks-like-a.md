---
id: I-0263
title: The new item form's lint refusal is easy to miss, so a refused epic looks like a hang
class: defect
status: closed
count: 1
cost: 20m
first_reported: 2026-11-06T21:31:26Z
last_reported: 2026-11-06T21:31:26Z
updated: 2026-11-11T08:24:13Z
---

# I-0263 The new item form's lint refusal is easy to miss, so a refused epic looks like a hang

## Description
The new item form's lint refusal is easy to miss, so a refused epic looks like a hang

## Instances

### 2026-11-06T21:31:26Z
Creating E-0032 from the dashboard twice (21:21 and 21:27 UTC) got a 422 in a few seconds, the markdown lint rejecting two bare URLs (MD034), but the operator saw the form hang indefinitely. NewItemForm renders the refusal alert below the body preview, under the fold on a long body, and nothing scrolls to it or marks the Create button.

## Remediation

Story S-0700 remediates this issue, created from it at 2026-11-10T08:55:31Z.
Closed 2026-11-11T08:24:13Z: S-0700: NewItemForm.svelte now scrolls every refusal into view and focuses it (tabindex="-1"), so a refusal drawn below the body and preview is no longer missed on a submit with Enter from the title, and the Create button reads "Refused: fix it and create <type> again", described by the alert, until the title or body changes or the next submit (T-2966). The component test "scrolls to and focuses a refusal, and marks the Create button refused until the text changes" reproduces the missed refusal: it fails on main's component and passes with the fix.
