---
id: I-0392
title: "flai file show follows a symlink under design or docs out of the repository: docedit/file.go confines by string prefix and os.Stat, unlike hostapi/docs.go's EvalSymlinks"
class: impression
status: open
count: 1
cost: 5m
first_reported: 2026-11-16T06:16:39Z
last_reported: 2026-11-16T06:16:39Z
updated: 2026-11-16T06:16:39Z
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

# I-0392 flai file show follows a symlink under design or docs out of the repository: docedit/file.go confines by string prefix and os.Stat, unlike hostapi/docs.go's EvalSymlinks

## Description
flai file show follows a symlink under design or docs out of the repository: docedit/file.go confines by string prefix and os.Stat, unlike hostapi/docs.go's EvalSymlinks

## Instances

### 2026-11-16T06:16:39Z
Report: design/analysis/2026-11-16-all.md.
Analysis 2026-11-16: no instance observed; found by reading the code.

## Impact
flai/internal/docedit/file.go:52-80 FileReach checks the path by string prefix, :344-346 os.Stat and IsRegular follow links, :425 os.ReadFile(abs). flai/internal/hostapi/docs.go:145-155 resolves with EvalSymlinks and checks the real path is under the root; flai/internal/hostapi/logs.go:329-355 re-checks with SameFile (ADR-0300). A committed link under design/ or docs/ from an agent's branch or an imported repository (ADR-0035, ADR-0045) would be read and saved over through file.get and file.save.

## Remediation

Found by the analysis in [2026-11-16-all.md](../analysis/2026-11-16-all.md).
