---
id: I-0287
title: The merge queue builds one-story batches and runs three verifications per story landed
class: efficiency
status: open
count: 1
first_reported: 2026-11-07T08:17:35Z
last_reported: 2026-11-07T08:17:35Z
updated: 2026-11-10T08:34:53Z
usage:
  strategic:
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
---

# I-0287 The merge queue builds one-story batches and runs three verifications per story landed

## Description
The merge queue builds one-story batches and runs three verifications per story landed

## Instances

### 2026-11-07T08:17:35Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: workflow.md says the queue verifies the stories in review as one batch; every batch held one story, and main moving dropped most batches, so each landing ran the integration and smoke tiers about three times.

## Impact
queue.weeks[] W44: landed 19, sent_back 6, batches 59, batch_size_mean 1.0, batch_size_max 1, bisections 7, verifications 56, verifications_per_landed 2.95; W43 batch_size_mean 1.08; review.mean_seconds W44 16,227

- time_lost_per_cycle: 7h

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0675 remediates this issue, created from it at 2026-11-10T08:34:53Z.
