---
id: I-0289
title: flai guard refuses the orchestrator git tag --list and git worktree list, which only read
class: efficiency
status: open
count: 1
first_reported: 2026-11-07T08:17:36Z
last_reported: 2026-11-07T08:17:36Z
updated: 2026-11-07T08:17:36Z
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

# I-0289 flai guard refuses the orchestrator git tag --list and git worktree list, which only read

## Description
flai guard refuses the orchestrator git tag --list and git worktree list, which only read

## Instances

### 2026-11-07T08:17:36Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: strategic-agents.md allows the orchestrator git's reads with no permission; the analyzer's own git tag --sort read was refused the same way.

## Impact
wip/agents/orchestrator.md Refusals: git tag --sort=-creatordate refused 18 times, git worktree list 4 times, 2026-10-25 to 2026-11-07; flai/internal/guard/guard.go:330 gitReads has no tag, worktree, or branch, and line 403 lists tag under gitPublishes

- time_lost_per_cycle: 30m

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).
