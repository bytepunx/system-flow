---
id: S-0341
type: story
nature: improvement
title: flai verify resumes at the tier that failed when the branch head and its base are unchanged, so a retry re-runs the failure and not the tiers that passed
status: review
owner: alex
created: 2026-10-08T07:59:12Z
updated: 2026-10-08T09:41:28Z
transitions:
  - to: ready
    at: 2026-10-08T07:59:29Z
    by: system-flow
  - to: in-progress
    at: 2026-10-08T08:37:17Z
    by: agent-S-0341
  - to: review
    at: 2026-10-08T09:41:28Z
    by: agent-S-0341
tags: [cli]
topics: [testing]
touches: [flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/run.go, flai/internal/verify/run_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, scripts/close-out.sh, docs/users/flai.md, design/system/flai-cli.md, flai/internal/verify/verify.go, docs/users/flai-reference.md, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, docs/users/flaiover.md, docs/operators/settings.md, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md, design/issues/summary.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4392
  turns:
    - day: 2026-10-08
      ceremony: 14
      test_runs: 7
      hand_edits: 3
      work: 105
  models:
    - model: claude-opus-5-5
      input: 388
      output: 111770
      cache_read: 23988291
      cache_write: 605737
      cost: 11.2066
  strategic:
    - kind: orchestrator
      seconds: 758
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 71
          output: 1186
          cache_read: 22442834
          cache_write: 64030
          cost: 5.5484
        - model: claude-sonnet-5-5
          input: 24
          output: 128
          cache_read: 470711
          cache_write: 59648
          cost: 0.4572
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: alex
    at: 2026-10-08T07:59:12Z
  value: 12.5
  by: planner-S-0341
  at: 2026-10-08T08:04:45Z
forecast:
  duration: 1h30m
  delivery: 2026-10-08T10:13:00Z
  basis: "Its own forecast of 1h30m; 2nd in the pull order with an in-progress limit of 5, behind S-0232, S-0321, S-0322, S-0340 and S-0342."
  by: flai
  at: 2026-10-08T08:37:06Z
---
# S-0341 flai verify resumes at the tier that failed when the branch head and its base are unchanged, so a retry re-runs the failure and not the tiers that passed

## Goal

`flai verify` runs every selected tier from the start on every call. When a late tier fails for a reason outside the branch, the agent fixes nothing and runs the close-out again, and pays the whole suite again before reaching the tier that failed. On 2026-10-08 S-0326's agent ran the close-out five times at the same head: each run passed rebase, sync, narrative, check, gofmt, vet, golangci-lint, go-test, markdown, and integration, about four minutes, and then smoke failed on a network fault. Twenty minutes of Go tests were spent to learn five times what the first run had shown.

A verify run at the same branch head, against the same main commit, with the same selected tiers as the last recorded run, can trust that run's passed tiers: nothing they read has changed. Such a run starts at the first tier the record shows as failed or not reached, after the cheap checks that guard the premise, and says which tiers it reused.

## Acceptance criteria

- [x] `flai verify S-nnnn` reads the story's last record and, when the branch head, the main commit it was verified against, and the selected tiers are the same, re-runs the rebase, sync, narrative, and check steps and then only the tiers from the first one the record shows as failed or not reached; a tier it did not run is reported as `reused` with the time of the run it comes from, in the text and in `--json`.
- [x] Any change to the head, the base, the selected tiers, the manifest's `tests`, or a tier's command makes the next run a full run, and `--fresh` forces one; the record written by a resumed run holds every tier's state, reused ones included, so `--last` and the review page show a whole result.
- [x] `scripts/close-out.sh` runs `flai verify` as it does, so a second close-out at an unchanged head after a smoke failure reaches smoke within the cheap checks' time; the close-out's last line says when tiers were reused.
- [x] Tests cover a resumed run, each condition that forces a full run, and `--fresh`.
- [x] `docs/users/flai.md` and `design/system/flai-cli.md` describe when a run resumes and how to force a full one.

## Tasks
- T-1343 flai verify resumes a story's run from its last record when the head, the base, and the tiers are unchanged
- T-1344 flai verify takes --fresh and prints reused tiers and a last line that says it resumed
- T-1346 The review page shows a reused tier with the time of the run it comes from
- T-1347 The close-out's last line says when flai verify reused tiers
- T-1351 The user guide, the dashboard guide, and the CLI design say when flai verify resumes and how to force a full run

## Notes

The record is `.flai-cache/verify/S-nnnn.json`, which `flai verify --last` prints; it already carries the head and the main commit a run was verified against. A commit in the close-out between verify and the next run changes the head and so forces a full run, which is right: the tiers ran against other code.

### How criterion 3 was verified

- `scripts/close-out.sh` still runs `flai verify --record-issues` with no `--fresh`. A close-out that stops at flai verify commits nothing, so the next close-out runs at the same head and resumes, as `TestVerifyResumesAtTheTierTheLastRunFailed` shows for verify.
- The last line's `(reused N tiers from <time>)` was checked with a hand-run harness (T-1347): a stub git and a fake `scripts/flai.sh`, under `sh`, `bash`, and `dash`, for a pass, a stop, one tier, a malformed record, nothing reused, exit 2, and TERM.
- The story's own close-out passed in full at c7c37c90 and ended with the unchanged form. It was not run a second time to show a real resumed close-out, since a passing run is not repeated.
- When a close-out loses the race with main (I-0119), the sync it needs moves the head and the base, so the run after it is full. Resuming does not help that case.

### Planning

Touches, by where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/verify/story.go`, `story_test.go`, `run.go`, `run_test.go` | declared | The record, the steps before the tiers, and the tier run |
| `flai/cmd/verify.go`, `verify_test.go` | declared | `--fresh` and the text lines, which `verifyText` prints |
| `scripts/close-out.sh` | declared | Criterion 3 |
| `docs/users/flai.md`, `design/system/flai-cli.md` | declared | Criterion 5 |
| `flai/internal/verify/verify.go` | layout | The `State` constants, where `Reused` joins `Passed`, `Failed`, and `NotReached` |
| `docs/users/flai-reference.md` | co-change (26%) and layout | Generated from the flags by `make flai-reference`; `--fresh` changes it |
| `flaiover/src/lib/review.ts`, `review.test.ts`, `components/Review.svelte` | design | Criterion 2's review page: `VerifyState` is a closed union and the row colours only `passed` and `failed` |
| `docs/users/flaiover.md` | design | "Reviewing a story" lists the states a step shows |

No folder touch is kept. `flai touches suggest` found no co-change above 30% beyond the general documents (`docs/operators/index.md`, `design/system/flaiover-dashboard.md`), which nothing in the criteria reaches.

Figures:

- Forecast 1h30m, raised from flai's 19m. flai's median is over all improvement stories of the size; the two latest `flai verify` stories, S-0326 (81 min, 2 criteria) and S-0270 (100 min, 5 criteria, with the review page), took four to five times it. This one has five criteria and five tasks in three layers across Go, shell, Svelte, and docs. Delivery moves out by the same 71 minutes, to 2026-10-08T10:30:00Z.
- Cost of delay 12.50 USD a week, as `flai cod` gives it from the operator's input of 5m lost per cycle, at 150 USD an hour and one 168h cycle a week. It stands: the input is the operator's, and the five-run morning of S-0326 is one occurrence, not yet a weekly rate.
