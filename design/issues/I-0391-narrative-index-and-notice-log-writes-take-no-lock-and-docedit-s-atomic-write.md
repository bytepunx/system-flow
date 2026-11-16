---
id: I-0391
title: Narrative, index, and notice-log writes take no lock, and docedit's atomic write does not sync before its rename
class: impression
status: open
count: 1
cost: 5m
first_reported: 2026-11-16T06:16:35Z
last_reported: 2026-11-16T06:16:35Z
updated: 2026-11-16T06:16:35Z
usage:
  strategic:
    - kind: analyzer
      seconds: 99
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 47
          output: 6584
          cache_read: 349359
          cache_write: 34305
          cost: 1.1031
        - model: claude-haiku-5-5
          input: 35
          output: 14842
          cache_read: 2398930
          cache_write: 47318
          cost: 0.1557
---

# I-0391 Narrative, index, and notice-log writes take no lock, and docedit's atomic write does not sync before its rename

## Description
Narrative, index, and notice-log writes take no lock, and docedit's atomic write does not sync before its rename

## Instances

### 2026-11-16T06:16:35Z
Report: design/analysis/2026-11-16-all.md.
Analysis 2026-11-16: no instance observed in the window; a lost Current state write or an empty design document after a crash is the risk.

## Impact
flai/internal/workitem/narrative.go reads and saves a narrative whole at 169/190, 355/400, and 743/729 and rewrites wip/agents/index.md in WriteIndex with no call to lockItems (itemlock.go:41) or any other lock, while the MCP server and flai serve both write it; I-0175 recorded the same race on item files before ADR-0184. flai/internal/itemedit/notice.go:68 truncates .flai-cache/edits.jsonl with _= os.WriteFile against the O_APPEND writer at line 74, errors dropped at 68 and 79. flai/internal/docedit/docedit.go:493 writeAtomic renames without Sync, where flai/internal/atomicfile/atomicfile.go:38 syncs (I-0137).

## Remediation

Found by the analysis in [2026-11-16-all.md](../analysis/2026-11-16-all.md).
