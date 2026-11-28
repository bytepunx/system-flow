---
id: I-0420
title: The dashboard's review page shows operator_only only beside protected paths, so it does not show that only the operator accepts a story whose branch changes the workflow file
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-28T00:24:02Z
last_reported: 2026-11-28T00:24:02Z
updated: 2026-11-28T00:24:02Z
---

# I-0420 The dashboard's review page shows operator_only only beside protected paths, so it does not show that only the operator accepts a story whose branch changes the workflow file

## Description
The dashboard's review page shows operator_only only beside protected paths, so it does not show that only the operator accepts a story whose branch changes the workflow file

## Instances

### 2026-11-28T00:24:02Z
Story: S-0844.
S-0844 made the acceptance preview set operator_only and a new workflow_file field when a story's branch changes design/workflows/board.md (ADR-0301 §7). flaiover/src/lib/components/Review.svelte (around lines 81 and 1774) shows operator_only only when protected is non-empty, and its type comment says the two always appear together. So the review page shows nothing for a story that changes only the workflow file. The refusal still holds: the orchestrator and the queue cannot accept it. The fix is in flaiover: show operator_only whenever it is set, list workflow_file beside the protected paths, and add the field to the preview's type.

## Remediation
