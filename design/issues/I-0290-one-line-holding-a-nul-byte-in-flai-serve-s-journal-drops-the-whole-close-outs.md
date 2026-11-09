---
id: I-0290
title: One line holding a NUL byte in flai serve's journal drops the whole close-outs section of flai stats
class: defect
status: closed
count: 1
first_reported: 2026-11-07T08:17:36Z
last_reported: 2026-11-07T08:17:36Z
updated: 2026-11-09T17:17:43Z
usage:
  strategic:
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
---

# I-0290 One line holding a NUL byte in flai serve's journal drops the whole close-outs section of flai stats

## Description
One line holding a NUL byte in flai serve's journal drops the whole close-outs section of flai stats

## Instances

### 2026-11-07T08:17:36Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: line 874 is an orchestrate entry of 2026-10-08 preceded by NUL bytes from a write cut short; a month of close-out verdicts and durations is unreported, and the tree's flai reads .flai-cache/serve, which has no journal, so it never shows the section.

## Impact
flai 1.84.0 stats --json warns: cannot read flai serve's journal /home/alex/.flai/serve/journal.jsonl, line 874: invalid character '\x00' looking for beginning of value; the close-outs section is left out. The journal holds 166 lines naming verify.job; metrics.md § Close-outs says an unreadable journal leaves the section absent

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0624 remediates this issue, created from it at 2026-11-07T23:40:09Z.
Closed 2026-11-09T17:17:43Z: S-0624 fixed it (ADR-0268, in `flai/internal/statsread/statsread.go`): flai stats reads each serve journal line with its NUL bytes removed, skips a line that is still not an entry with one warning naming the count and the first line, and keeps the close-outs section; TestCloseOutsSkipsADamagedJournalLineAndKeepsTheSection reproduces line 874's NUL-prefixed entry.
