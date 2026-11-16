---
id: I-0390
title: "Design sentences the code contradicts: ci.md's race detector on macOS, markdown.md's yaml in flaiover, metrics.md's WIP series and unnamed keys, flai-cli.md's missing dashboard logs row, the analyzer's MultiEdit"
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-11-16T06:16:30Z
last_reported: 2026-11-16T06:16:30Z
updated: 2026-11-16T06:16:30Z
usage:
  strategic:
    - kind: analyzer
      seconds: 99
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 47
          output: 6584
          cache_read: 349360
          cache_write: 34306
          cost: 1.1031
        - model: claude-haiku-5-5
          input: 36
          output: 14843
          cache_read: 2398931
          cache_write: 47318
          cost: 0.1557
---

# I-0390 Design sentences the code contradicts: ci.md's race detector on macOS, markdown.md's yaml in flaiover, metrics.md's WIP series and unnamed keys, flai-cli.md's missing dashboard logs row, the analyzer's MultiEdit

## Description
Design sentences the code contradicts: ci.md's race detector on macOS, markdown.md's yaml in flaiover, metrics.md's WIP series and unnamed keys, flai-cli.md's missing dashboard logs row, the analyzer's MultiEdit

## Instances

### 2026-11-16T06:16:30Z
Report: design/analysis/2026-11-16-all.md.
Analysis 2026-11-16: design/tech/ci.md:37 says the race detector runs on Linux and macOS, .github/workflows/flai.yml:38-41 runs it on ubuntu only; design/tech/markdown.md:13 lists yaml 2.9 in flaiover, which flaiover/package.json does not name; design/system/metrics.md:42 defines WIP as a series where metrics.go:179 emits one count, and § Aggregates names no JSON key for summary, throughput, cfd, burnup, aging, or groups; design/system/flai-cli.md § Commands has no row for flai dashboard logs (flai/cmd/dashboard.go:1038); design/system/strategic-agents.md:940 lists the analyzer's file tools without MultiEdit (flai/internal/guard/claudecode.go:39).

## Remediation

Found by the analysis in [2026-11-16-all.md](../analysis/2026-11-16-all.md).
