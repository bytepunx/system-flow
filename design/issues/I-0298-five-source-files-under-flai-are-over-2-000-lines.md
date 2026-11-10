---
id: I-0298
title: Five source files under flai are over 2,000 lines
class: impression
status: closed
count: 1
first_reported: 2026-11-07T08:17:37Z
last_reported: 2026-11-07T08:17:37Z
updated: 2026-11-10T06:49:15Z
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

Story S-0631 remediates this issue, created from it at 2026-11-07T23:40:53Z.
Closed 2026-11-10T06:49:15Z: S-0631 split the six non-test Go files over 2,000 lines by responsibility within their packages, moving declarations unchanged: flai/internal/hostapi/writes.go (3,626 lines, now 32, itemSpecs cut into itemWriteSpecs, streamSpecs, and hostSpecs), flai/internal/serve/agents.go (3,318, now 419), flai/cmd/serve_actions.go (2,754, now 601), flai/internal/manifest/manifest.go (2,449, now 591), flai/internal/check/check.go (2,235, now 284), and flai/internal/guard/guard.go (2,151, now 878); no new file is over 1,000 lines. TestNoSourceFileIsOver2000Lines in flai/internal/check/sourcesize_test.go fails on any non-test Go file under flai over 2,000 lines, naming it and its count; it failed on main before the splits and runs under -short. The nolint:nilerr in flai/internal/watch/watch.go now carries its reason on the directive.
