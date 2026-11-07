---
id: S-0570
type: story
nature: remediation
title: A story branch rebased across a release brings back the changelog's old Unreleased heading and the lines the release moved, since union merge keeps both sides
status: done
owner: alex
created: 2026-11-02T21:41:52Z
updated: 2026-11-07T05:07:04Z
transitions:
  - to: ready
    at: 2026-11-07T03:56:09Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T04:14:43Z
    by: agent-S-0570
  - to: review
    at: 2026-11-07T05:05:28Z
    by: agent-S-0570
  - to: done
    at: 2026-11-07T05:07:04Z
    by: orchestrator
tags: [cli, git, release]
topics: [release, git]
touches: [design/adrs, design/adrs/README.md, flai/internal/release/changelog.go, flai/internal/release/changelog_test.go, flai/internal/release/release.go, flai/internal/storygit/sync.go, flai/internal/storygit/undone.go, flai/internal/storygit/sync_test.go, flai/internal/issues/generated.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, design/system/workflow.md, design/conventions/git.md, template/root/design/conventions/git.md, docs/users/flai.md, docs/contributors/template.md, template/CHANGELOG.md, design/issues/I-0251-a-story-branch-rebased-across-a-release-brings-back-the-changelog-s-old.md, design/issues/summary.md, flai/cmd/stream.go, design/adrs/0257-a-rebase-of-a-story-branch-puts-the-branch-s-own-changelog-lines-back-under-the.md, design/issues/I-0248-testtheoperatorsstopisonestoppedagentevent-s-waiting-for-an-answer-case-fails.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3053
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 2
      test_runs: 9
      hand_edits: 1
      work: 52
  models:
    - model: claude-opus-5-5
      input: 196
      output: 60812
      cache_read: 26819052
      cache_write: 405132
      cost: 9.6454
  strategic:
    - kind: orchestrator
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 12
          cache_read: 167036
          cache_write: 1708
          cost: 0.041
cost_of_delay:
  inputs:
    time_lost_per_cycle: 26m
    by: flai
    at: 2026-11-02T21:41:52Z
  value: 65
  by: planner-S-0570
  at: 2026-11-07T03:49:48Z
forecast:
  duration: 2h
  delivery: 2026-11-07T11:40:00Z
  basis: "flai's 39m from the remediation median times size 23, raised to 2h for an ADR, a changelog merge, the same step in three rebase paths, and three real-git tests; delivery is flai's play-out moved by the raise"
  by: planner-S-0570
  at: 2026-11-07T03:49:48Z
finalized:
  by: orchestrator
  at: 2026-11-07T03:54:24Z
---
# S-0570 A story branch rebased across a release brings back the changelog's old Unreleased heading and the lines the release moved, since union merge keeps both sides

## Goal

This story remediates [I-0251](../../../design/issues/I-0251-a-story-branch-rebased-across-a-release-brings-back-the-changelog-s-old.md), "A story branch rebased across a release brings back the changelog's old Unreleased heading and the lines the release moved, since union merge keeps both sides". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0251 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0251 is closed with `flai issue close I-0251 --reason` saying what fixed it

## Tasks
- T-2524 An ADR proposes from I-0251's two instances how a rebase keeps a story's changelog lines under main's top Unreleased, refining ADR-0228 §4
- T-2531 The release package writes a changelog as the base branch has it with a branch's own added lines under its top Unreleased
- T-2535 flai stream sync and flai accept put a rebased branch's changelog lines back under main's top Unreleased before checking what the rebase undoes
- T-2536 The merge queue's batch puts each rebased story's changelog lines back under the top Unreleased, so a batch over a release lints clean
- T-2537 The design, git.md and its template copy, and the guides say a rebase puts a story's changelog lines back under the top Unreleased
- T-2538 I-0251 is closed with flai issue close, its reason naming the ADR and the step that fixed it

## Notes

Cost of delay inputs set by flai from I-0251. time_lost_per_cycle 26m: 13m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-02T03:57:28Z, 0.7 days before this story; under one cycle counts as one).

### Planning

- **Touches.** None were declared, and `flai touches suggest` needs a starting path when a story has none, so every touch is the planner's.
  - **Design:** the ADR that proposes the fix and refines [ADR-0228](../../../design/adrs/0228-flai-generates-the-adr-index-from-the-adr-files-and-checks-the-issues-summary.md) §4, with the index `design/adrs/README.md` that `flai adr new` writes again. The fix is then described in `design/system/workflow.md` § Hub files: what S-0550 changed, `design/conventions/git.md` and its template copy `template/root/design/conventions/git.md`, `docs/users/flai.md` § Story branches, and `docs/contributors/template.md` § Versioning. The template copy's change gets a line in `template/CHANGELOG.md`.
  - **Layout:** `flai/internal/release/release.go` already parses the `## Unreleased` section (`cutUnreleased`, `bumpVersionAndChangelog`), so the merge of a changelog goes beside it in `changelog.go`, new, with `changelog_test.go`. `flai/internal/storygit/sync.go` `Rebase` is the one rebase `flai stream sync` and `flai accept` share. `undone.go` `checkUndone` runs there before any fold, which is why instance 2 reported `undoes_main`. `flai/internal/issues/generated.go` `FilesIn` builds the files a rebase settles. `flai/internal/mergequeue/batch.go` `rebase` runs its own `git rebase` and needs the same step. `flai/cmd/stream_sync.go` prints what a sync did. The real-git tests go in `flai/internal/storygit/sync_test.go`, `flai/cmd/stream_sync_test.go` beside `unionStories`, and `flai/internal/mergequeue/batch_test.go`.
  - **Issue:** `design/issues/I-0251-…md` and `design/issues/summary.md`, which `flai issue close` writes.
  - **Co-change:** with the touches above as its start, `flai touches suggest` lists hub documents (`design/system/flai-cli.md`, `docs/operators/index.md`, `docs/users/flai-reference.md`, `flaiover-dashboard.md`). None was added. The fix adds no command, flag, or setting, and `flai-cli.md` has no section on the sync's rebase.
  - **Folder kept:** `design/adrs`, because the ADR's file name is not known until `flai adr new` writes it. It is a shared path (`claims.shared`), so it holds no ready story.
  - **Left out:** `.gitattributes` and `template/root/.gitattributes`, on the assumption that `merge=union` stays. If T-2524's ADR picks a merge driver or changelog fragments, the story's agent widens the touches.
- **After.** None. The story reads no output of another open story.
- **Forecast.** 2h against flai's 39m. flai's figure is the remediation median, 101 s per unit, times size 23: two criteria and 21 touches. The story adds an ADR, a changelog merge, the same step in three rebase paths (sync, accept, and the merge queue's batch), and three real-git tests that reproduce both instances, which is about three times the median remediation. Delivery 2026-11-07T11:40Z is flai's play-out, 10:15Z, moved by the raise.
- **Cost of delay.** 65.00 USD a week, as `flai cod` gives it: 26m lost per 168h cycle, from I-0251's two occurrences at 13m each, at the project's 150 USD an hour. Kept. The inputs are flai's, taken from the issue, and nothing in the instances says the loss is larger. Each release while a story with a changelog line is open costs another fix by hand.

### Accepted by the orchestrator

- Verified: 012fcd7d96902357751997002e968494c170f7dd
- At: 2026-11-07T05:07:04Z

Verdict: pass
- 1: flai/internal/release/changelog.go, flai/internal/release/changelog_test.go, flai/internal/release/release.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, design/adrs/0257-a-rebase-of-a-story-branch-puts-the-branch-s-own-changelog-lines-back-under-the.md
- 2: design/issues/I-0251-a-story-branch-rebased-across-a-release-brings-back-the-changelog-s-old.md, design/issues/summary.md
