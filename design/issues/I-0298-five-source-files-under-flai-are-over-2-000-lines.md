---
id: I-0298
title: Five source files under flai are over 2,000 lines
class: impression
status: open
count: 1
first_reported: 2026-11-07T08:17:37Z
last_reported: 2026-11-07T08:17:37Z
updated: 2026-11-07T08:17:37Z
usage:
  strategic:
    - kind: analyzer
      seconds: 41
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759049
          cache_write: 13854
          cost: 0.278
---

# I-0298 Five source files under flai are over 2,000 lines

## Description
Five source files under flai are over 2,000 lines

## Instances

### 2026-11-07T08:17:37Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: code-quality.md asks for SRP at every level; no TODO, FIXME, or HACK anywhere in Go, flaiover/src, or scripts.

## Impact
flai/internal/hostapi/writes.go 3,582 lines, flai/internal/serve/agents.go 2,935, flai/cmd/serve_actions.go 2,421, flai/internal/manifest/manifest.go 2,343, flai/internal/check/check.go 2,125; with tests 12 files over 2,000 (agents_test.go 4,333, writes_test.go 3,824, harness_test.go 3,439); 26 nolint comments, one without a reason at flai/internal/watch/watch.go:54

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).
