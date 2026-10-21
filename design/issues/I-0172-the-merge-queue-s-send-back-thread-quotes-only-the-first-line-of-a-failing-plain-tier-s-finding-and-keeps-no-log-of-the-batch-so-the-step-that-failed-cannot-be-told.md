---
id: I-0172
title: The merge queue's send-back thread quotes only the first line of a failing plain tier's finding and keeps no log of the batch, so the step that failed cannot be told
class: efficiency
status: open
count: 1
cost: 30m
first_reported: 2026-10-21T01:54:16Z
last_reported: 2026-10-21T01:54:16Z
updated: 2026-10-21T02:13:46Z
---

# I-0172 The merge queue's send-back thread quotes only the first line of a failing plain tier's finding and keeps no log of the batch, so the step that failed cannot be told

## Description
The merge queue's send-back thread quotes only the first line of a failing plain tier's finding and keeps no log of the batch, so the step that failed cannot be told

## Instances

### 2026-10-21T01:54:16Z
Story: S-0454.
TH-0545 said S-0454's batch 20261021-014447 failed smoke at 'smoke: template render and check', which is only the first line smoke.sh prints. The finding's tail and the batch worktree were gone, and serve.log records only stopped_at smoke. Reproducing the template step in a fresh worktree and with flai serve's config passed both times; the cause was the gateway check's retired model (I-0171), found only by reading flai serve's environment and scripts/env.sh. Recommend the thread quote the finding's failure lines and tail, as flai verify prints them, and the queue keep each batch's tier output under .flai-cache/merge-queue.

## Remediation

Story S-0457 remediates this issue, created from it at 2026-10-21T02:13:46Z.
