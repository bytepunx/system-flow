---
id: I-0288
title: The board stood empty for three days while the backlog held only drafts waiting on operator decisions nobody asked for again
class: impression
status: open
count: 1
first_reported: 2026-11-07T08:17:35Z
last_reported: 2026-11-07T08:17:35Z
updated: 2026-11-10T08:49:02Z
usage:
  strategic:
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3103
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
---

# I-0288 The board stood empty for three days while the backlog held only drafts waiting on operator decisions nobody asked for again

## Description
The board stood empty for three days while the backlog held only drafts waiting on operator decisions nobody asked for again

## Instances

### 2026-11-07T08:17:35Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: E-0023's seven stories wait behind S-0406, a draft held for an OpenCode install and an OpenRouter key; the orchestrator can promote no draft and nothing reminded the operator.

## Impact
cfd[] 2026-11-03 to 2026-11-05: backlog 26, ready 0, in-progress 0, review 0; turns.days 106, 0, 0; concurrency.share[0] 0.357; TH-0532 and TH-0533 have one entry each since 2026-10-20; 16 ready-empty events in .flai-cache/events.jsonl

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0691 remediates this issue, created from it at 2026-11-10T08:49:02Z.
