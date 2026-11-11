---
id: I-0286
title: Declared touches match what a story changes in 16 percent of completed stories, so claims mislead the pull hold, the trial merge, and the conversations
class: efficiency
status: closed
count: 1
first_reported: 2026-11-07T08:17:35Z
last_reported: 2026-11-07T08:17:35Z
updated: 2026-11-11T02:05:38Z
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
          cache_write: 13855
          cost: 0.278
---

# I-0286 Declared touches match what a story changes in 16 percent of completed stories, so claims mislead the pull hold, the trial merge, and the conversations

## Description
Declared touches match what a story changes in 16 percent of completed stories, so claims mislead the pull hold, the trial merge, and the conversations

## Instances

### 2026-11-07T08:17:35Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: the planner's predicted touches and task done's widening leave both files outside the claim and touches nothing changed; each conflict met at acceptance is a send-back and a second close-out.

## Impact
claims.weeks[].exact_share W43 0.651, W44 0.360, W45 0.163 (7 of 43); coordination.weeks[] W44 conflicts.found 36 on 94 paths, met_accept 4; W45 found 28, met_accept 3

- time_lost_per_cycle: 4h

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0674 remediates this issue, created from it at 2026-11-10T08:34:52Z.
Closed 2026-11-11T02:05:38Z: Fixed by S-0674 (ADR-0284). T-2939: storygit.Committed and the touches drift leave out flai's folders, wip and the design folder's digest, analysis, and issues, on both sides, so flai's bookkeeping commits (accept and archive, create story, record I-nnnn) no longer count as files outside a story's touches; TestClaimsDriftLeavesOutFlaiFoldersOnBothSides reproduces the instance. T-2940: flai task done removes from the closing task's touches each file the branch does not change, and from the story's when no other open task names it. Measured 2026-11-11 with flai stats --since 30d, stories from S-0560 on: exact 27 of 114 before, 78 of 103 after T-2939 (11 stories changed only flai's folders and have no drift entry now); 1 story still has a file outside its touches (S-0664, two ADR files), and 24 of the 25 with unchanged touches had every one named by a task, which T-2940's narrowing drops at the task's close. Weekly exact_share W46 0.303 to 0.771, W45 0.123 to 0.593.
