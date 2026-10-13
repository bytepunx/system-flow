---
id: S-0376
type: story
nature: improvement
title: "flai check finds `markdown.MD012` outside the story at close-out"
status: in-progress
owner: alex
created: 2026-10-12T16:52:58Z
updated: 2026-10-13T17:20:32Z
transitions:
  - to: ready
    at: 2026-10-13T17:16:55Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-13T17:17:01Z
    by: agent-S-0376
tags: [flai, narrative, issues]
topics: [markdown, narrative, planning]
touches: [flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/internal/workitem/archive.go, flai/internal/workitem/archive_test.go, design/system/flai-cli.md, design/issues/I-0139-flai-check-finds-markdown-md012-outside-the-story-at-close-out.md, design/issues/summary.md, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 206
  estimated: true
  turns:
    - day: 2026-10-13
      test_runs: 4
      hand_edits: 1
      work: 27
  models:
    - model: claude-opus-5-5
      input: 64
      output: 585
      cache_read: 4289976
      cache_write: 174971
      cost: 2.0525
  strategic:
    - kind: orchestrator
      seconds: 461
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 69
          output: 1019
          cache_read: 4295721
          cache_write: 28202
          cost: 1.0904
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: orchestrator
    at: 2026-10-13T16:42:54Z
  value: 25
  by: planner-S-0376
  at: 2026-10-13T16:45:11Z
forecast:
  duration: 20m
  delivery: 2026-10-13T17:39:00Z
  basis: "Its own forecast of 20m; 1st in the pull order with an in-progress limit of 5, behind S-0356, S-0377 and S-0379."
  by: flai
  at: 2026-10-13T17:16:57Z
finalized:
  by: orchestrator
  at: 2026-10-13T16:45:39Z
---
# S-0376 flai check finds `markdown.MD012` outside the story at close-out

## Goal

This story remediates [I-0139](../../../design/issues/I-0139-flai-check-finds-markdown-md012-outside-the-story-at-close-out.md), "flai check finds `markdown.MD012` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0139 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0139 is closed with `flai issue close I-0139 --reason` saying what fixed it

## Tasks
- T-1506 flai stream answer leaves one blank line where it removes an open question, and refuses a save the markdown lint rejects
- T-1507 Archiving a narrative collapses each run of blank lines outside fenced code to one, so no archived narrative carries an MD012 finding
- T-1508 flai-cli.md says flai stream answer leaves one blank line and archiving collapses a narrative's blank runs
- T-1509 Close I-0139 saying what fixed it

## Notes

### Planning

Proposed solution, from I-0139's instances. All five are one file, `wip/archive/agents/S-0232.md`, reported between 15:41Z and 16:21Z on 2026-10-12. Its `## Open questions` was left as the heading, two blank lines, and `## Log` in the acceptance commit `7801de00`. Before it, the section held one hand-written bullet between two blank lines, under the mirrored thread block. Thread resolution removed the block cleanly; removing the bullet then left the double blank. `flai stream answer` (`AnswerOpenQuestion`) removes a bullet that way, and saves without the lint guard. A hand edit does too, and nothing clears a finding on an archived narrative but a hand commit on main: `6526d727` broke the front matter, and `8a30de84` repaired it. So the fix is in the two writers: T-1506 fixes `stream answer`, and T-1507 collapses blank runs as `Archive` moves a narrative. The close-out's scope stays as it is: ADR-0123 keeps a markdown finding on an archived narrative a note, because it shows a gap in flai's own lint.

Touches, all files, no folder:

| Touch | From |
|-------|------|
| `flai/internal/workitem/narrative.go`, `narrative_test.go` | layout: `AnswerOpenQuestion` and its tests (explorer); co-change with `threads.go` and `archive.go`, 4 of 19 commits |
| `flai/internal/workitem/archive.go` | layout: `Archive` and `moveFile` (explorer) |
| `flai/internal/workitem/archive_test.go` | layout: a new test file next to it |
| `design/system/flai-cli.md` | design: the `flai stream answer` and `flai archive` rows; co-change 8 of 19 |
| `design/issues/I-0139-…md`, `design/issues/summary.md` | the goal and criterion 2: `flai issue close` writes both |

Left out: `flai/internal/threads/threads.go`. `MirrorNarrative` already keeps one blank line when it removes the block, and `TestMirrorRemovalLeavesOneBlankLine` covers it. `design/system/agent-narrative.md` is left out too, since its Open questions paragraph stays true. `flai/internal/check/scope.go` is left out, per ADR-0123.

Overlap: S-0375, in progress, claims `narrative.go` and `narrative_test.go` (its T-1501). S-0376 is held until S-0375 is accepted. `flai-cli.md` and `design/issues` are in `claims.shared` and hold nothing.

Forecast: flai forecast gave 13m (84 s per unit of size over 12 done medium improvement stories, times size 9). I raised it to 20m. There are four tasks, and T-1507 is a fence-aware rewrite of the archive move with a new test file. S-0346, a three-task story of the same kind, took 13m in progress. Delivery is moved 7m later than flai's 17:57Z to match.

Cost of delay: 25 USD a week, flai cod's value from the 10m of time lost per cycle the orchestrator set on TH-0432, at 150 USD an hour. It stands as flai gave it.
