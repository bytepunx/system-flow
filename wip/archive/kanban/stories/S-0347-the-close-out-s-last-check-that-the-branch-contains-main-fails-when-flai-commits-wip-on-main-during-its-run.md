---
id: S-0347
type: story
nature: improvement
title: The close-out's last check that the branch contains main fails when flai commits wip on main during its run
status: done
owner: alex
created: 2026-10-08T08:08:22Z
updated: 2026-10-08T10:52:25Z
transitions:
  - to: ready
    at: 2026-10-08T10:29:09Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T10:29:15Z
    by: agent-S-0347
  - to: review
    at: 2026-10-08T10:51:15Z
    by: agent-S-0347
  - to: done
    at: 2026-10-08T10:52:25Z
    by: orchestrator
tags: [cli, template]
topics: [testing, git]
touches: [design/adrs/0133-the-close-out-s-sync-check-passes-over-commits-on-main-that-change-only-wip.md, design/adrs/README.md, design/system/flai-cli.md, design/system/devex.md, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/paths.go, flai/internal/verify/paths_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, docs/users/flai.md, docs/users/flai-reference.md, scripts/close-out.sh, template/root/scripts/close-out.sh, scripts/README.md, template/root/scripts/README.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md, design/issues/summary.md, design/adrs/0135-the-sync-step-of-flai-verify-and-the-close-out-s-last-check-pass-over-commits.md, design/issues/I-0126-flai-task-done-commits-the-changed-paths-of-another-open-task-of-the-same-layer.md, docs/operators/settings.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1334
  turns:
    - day: 2026-10-08
      ceremony: 1
      hand_edits: 2
      work: 52
  models:
    - model: claude-opus-5-5
      input: 234
      output: 81878
      cache_read: 13381562
      cache_write: 430862
      cost: 7.1344
  strategic:
    - kind: orchestrator
      seconds: 2275
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 310
          output: 5696
          cache_read: 112364910
          cache_write: 69331
          cost: 27.7017
        - model: claude-sonnet-5-5
          input: 24
          output: 109
          cache_read: 539134
          cache_write: 99286
          cost: 0.5503
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-10-08T08:08:22Z
  value: 75
  by: planner-S-0347
  at: 2026-10-08T08:45:10Z
forecast:
  duration: 1h30m
  delivery: 2026-10-08T12:06:00Z
  basis: "Its own forecast of 1h30m; 3rd in the pull order with an in-progress limit of 5, behind S-0232, S-0338, S-0342 and S-0337."
  by: flai
  at: 2026-10-08T10:29:12Z
finalized:
  by: alex
  at: 2026-10-08T08:39:51Z
---
# S-0347 The close-out's last check that the branch contains main fails when flai commits wip on main during its run

## Goal

This story remediates [I-0119](../../../design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md), "The close-out's last check that the branch contains main fails when flai commits wip on main during its run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0119 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0119 is closed with `flai issue close I-0119 --reason` saying what fixed it

## Tasks
- T-1375 Record in an ADR and the CLI design that the sync check passes over commits on main that change only wip the branch does not change
- T-1377 flai verify's sync step passes when the commits the branch lacks change only wip paths the branch does not change
- T-1379 The definition of done and the close-out rule say the branch contains main but for commits that change only wip the branch does not change
- T-1381 flai verify --sync-only runs the rebase and sync steps alone and stores no record
- T-1382 Both close-out scripts end with flai verify --sync-only instead of git merge-base
- T-1383 Close I-0119 with the reason that names the sync step's wip rule and the close-out's flai verify --sync-only

## Notes

Cost of delay inputs set by flai from I-0119. time_lost_per_cycle 30m: 10m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-08T04:29:28Z, 0.2 days before this story; under one cycle counts as one).

### Planning

The proposed fix: the `sync` step of `flai verify`, and the close-out's last check through a new `flai verify --sync-only`, pass when every commit of main the branch lacks changes only paths under the manifest's `wip` folder that the branch does not change, and name those commits in a note. `wip/` stays on main (ADR-0019), the tiers that read the repository scope themselves to the story (ADR-0085), and acceptance merges the branch, so such commits change nothing the run verified. A commit outside `wip/`, such as the acceptance in I-0119's third instance, still needs a sync and a new run. T-1375 records it in an ADR before it is built.

Touches, by where each came from. The story declared none; `flai touches suggest` needs a path to start from, and found no co-change above 12% from these.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/verify/story.go`, `story_test.go` | layout | `storyRun.sync`, the step that fails, and the I-0119 reproduction |
| `flai/internal/verify/paths.go`, `paths_test.go` | layout | `ChangedPaths` lives there; the helper that sorts the missing commits joins it |
| `flai/cmd/verify.go`, `verify_test.go` | layout | `--sync-only` |
| `docs/users/flai.md`, `docs/users/flai-reference.md` | design | "Verify a story before review" and its table of steps; the reference is generated from the flags |
| `scripts/close-out.sh`, `template/root/scripts/close-out.sh` | layout | The last check, `git merge-base --is-ancestor`, which I-0119 names |
| `scripts/README.md`, `template/root/scripts/README.md` | design | Their `close-out.sh` rows describe the last check |
| `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` | design | The definition of done and the close-out rule say "contains the main branch" |
| `design/system/flai-cli.md`, `design/system/devex.md` | design | The `flai verify` row and the close-out row describe the `sync` step and the last check |
| `design/adrs/0133-the-close-out-s-sync-check-passes-over-commits-on-main-that-change-only-wip.md`, `design/adrs/README.md` | design | The ADR that changes the definition of done; its number and slug are a prediction, and `flai task done` widens the touches to what flai gives |
| `design/issues/I-0119-...md`, `design/issues/summary.md` | declared by the goal | Criterion 2 |

No folder touch is kept. `design/system/workflow.md` (8% co-change) is left out: it does not say "contains the main branch".

Overlaps: S-0341, in progress, touches `flai/internal/verify/story.go`, `story_test.go`, `flai/cmd/verify.go`, `verify_test.go`, `scripts/close-out.sh`, `docs/users/flai.md`, `docs/users/flai-reference.md`, and `design/system/flai-cli.md`, so this story is held until S-0341 leaves progress; its tasks build on S-0341's record. S-0322 touches `design/adrs/README.md`.

Figures:

- Forecast 1h30m, raised from flai's 30m (78 s per unit of size × 23). S-0341, the nearest story, is forecast 1h30m in the same package and script; S-0326 and S-0270 on `flai verify` took 81 and 100 min. This one has six tasks in five layers across an ADR, Go, two shell scripts, the template, and docs. Delivery moves out by the same hour, to 2026-10-08T15:28:00Z.
- Cost of delay 75.00 USD a week, as `flai cod` gives it from flai's input of 30m lost per cycle, at 150 USD an hour and one 168h cycle a week. It stands: the input is the issue's three occurrences at 10m each, and the instances describe four to ten minutes each.

### Accepted by the orchestrator

- Verified: 9b15341e2318357ac21cb072508495e8ae88621d
- At: 2026-10-08T10:52:25Z

Verdict: accept. flai verify passed every step at the branch head 9b15341e, and the verifier matched both criteria to the diff. Only wip-only main commits are passed over, the manifest's wip folder is used, and --sync-only stores no record. Minor gap noted: a wip path the branch deleted is not counted as shared; the acceptance rebase still catches such a conflict.
- 1: flai/internal/verify/story.go, flai/internal/verify/paths.go, flai/internal/verify/story_test.go, flai/internal/verify/paths_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, scripts/close-out.sh, template/root/scripts/close-out.sh, design/adrs/README.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/devex.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md
- 2: design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md, design/issues/summary.md
