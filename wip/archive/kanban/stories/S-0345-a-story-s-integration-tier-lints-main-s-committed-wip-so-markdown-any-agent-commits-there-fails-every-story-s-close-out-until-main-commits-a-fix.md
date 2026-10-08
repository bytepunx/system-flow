---
id: S-0345
type: story
nature: remediation
title: A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix
status: done
owner: alex
created: 2026-10-08T08:08:19Z
updated: 2026-10-08T09:18:44Z
transitions:
  - to: ready
    at: 2026-10-08T08:40:14Z
    by: alex
  - to: in-progress
    at: 2026-10-08T09:08:32Z
    by: agent-S-0345
  - to: review
    at: 2026-10-08T09:17:47Z
    by: agent-S-0345
  - to: done
    at: 2026-10-08T09:18:44Z
    by: orchestrator
tags: [flai]
topics: [testing]
touches: [flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/repo_test.go, scripts/lint-md.sh, scripts/README.md, design/system/flai-cli.md, design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/summary.md, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 566
  turns:
    - day: 2026-10-08
      test_runs: 3
      hand_edits: 1
      work: 29
  models:
    - model: claude-opus-5-5
      input: 106
      output: 33675
      cache_read: 5740005
      cache_write: 303885
      cost: 3.9535
  strategic:
    - kind: orchestrator
      seconds: 63
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 330
          cache_read: 5524560
          cache_write: 9716
          cost: 1.3636
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h12m
    by: flai
    at: 2026-10-08T08:08:19Z
  value: 180
  by: planner-S-0345
  at: 2026-10-08T08:41:04Z
forecast:
  duration: 30m
  delivery: 2026-10-08T09:43:00Z
  basis: "Its own forecast of 30m; 3rd in the pull order with an in-progress limit of 5, behind S-0232, S-0288, S-0321, S-0322, S-0341, S-0290 and S-0344."
  by: flai
  at: 2026-10-08T08:59:49Z
finalized:
  by: alex
  at: 2026-10-08T08:40:10Z
---
# S-0345 A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix

## Goal

This story remediates [I-0117](../../../design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md), "A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0117 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0117 is closed with `flai issue close I-0117 --reason` saying what fixed it

## Tasks
- T-1363 TestRepositoryLintsClean scopes itself to the story in a close-out, so main's committed wip is a note, not a failure
- T-1364 lint-md.sh's whole-repository run leaves main's wip out in a close-out and lints only the wip files the story changes
- T-1365 Document that the repository lint test and lint-md.sh scope themselves to the story in a close-out
- T-1366 Close I-0117 with the reason that names the scoped lint test and lint-md.sh

## Notes

Cost of delay inputs set by flai from I-0117. time_lost_per_cycle 1h12m: 18m per occurrence × 4 occurrences ÷ 1 cycle of 168h (first reported 2026-10-08T04:12:50Z, 0.2 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0117's instances: the two whole-repository markdown checks a close-out runs scope themselves to the story when `CLOSE_OUT_STORY` is set, as `TestMonorepoIsClean` and `scripts/check.sh` already do (S-0249, [ADR-0085](../../../design/adrs/0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md)). A finding in a `wip/` file the story does not change is a note, not a failure. Without the variable, as in `make lint-md` and CI, every file still counts, so main's `wip/` stays linted there. The two checks:

- `TestRepositoryLintsClean` in `flai/internal/mdlint`, which the integration tier runs: the failure in all four instances.
- `scripts/lint-md.sh` with no files, which `scripts/smoke.sh` runs right after integration: it lints main's `wip/` the same way, and failed close-outs on it under I-0027.

`flai check` in the close-out is already scoped, and `flai test`'s markdown tier lints only the files the branch changes, so neither needs a change.

Touches, file by file; no folder touch was kept:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/mdlint/mdlint_test.go` | design (I-0117 names the test) | holds `TestRepositoryLintsClean` today |
| `flai/internal/mdlint/repo_test.go` | layout | new external test package: `workitem` imports `mdlint`, so the internal test cannot import `storygit` |
| `scripts/lint-md.sh` | layout | smoke's whole-repository lint |
| `scripts/README.md` | layout | documents `lint-md.sh`; S-0340 also touches it, so this story is held while S-0340 is in progress |
| `design/system/flai-cli.md` | co-change, design | the `flai verify` row names the checks that read `CLOSE_OUT_STORY`; a shared path |
| `design/issues/I-0117-…md` | criterion 2 | closed by `flai issue close` |
| `design/issues/summary.md` | criterion 2, co-change | regenerated when the issue closes |

`flai touches suggest` found no paths before these were declared. Seeded with them, it lists mostly `docs/users/flai.md` (33%), operator and dashboard docs, and other issues. `docs/users/flai.md` already says the tiers' checks count only what is the story's, and the rest are not this change, so none was added. `scripts/smoke.sh` was left out on purpose: S-0340 changes it, and the variable reaches `lint-md.sh` through the environment.

Forecast: `flai forecast` gave 17m, from 110 s per unit of size times size 9. Adjusted to 30m. Four close-out scoping stories like this one, S-0279, S-0280, S-0318, and S-0323, took 12 to 49 minutes of agent time, median about 24m. This change also runs the integration and smoke tiers at close-out because it changes `flai/`. The delivery is flai's, 2026-10-08T13:51Z, moved by the 13m added.

Cost of delay: 180 USD a week, as `flai cod` worked it out from the operator's input of 1h12m lost per 168h cycle at 150 USD an hour. It stands. Every story that changes `flai/` can meet this failure while any agent writes to `wip/` on main, so the input may undercount, but it is the operator's to change.

### Accepted by the orchestrator

- Verified: 9859fa3ec6b6acaad1db2cf8f3f7df6e78c63ec2
- At: 2026-10-08T09:18:44Z

Verdict: accept. flai verify passed every step at the branch head 9859fa3e, and the verifier matched both criteria to the diff. Without CLOSE_OUT_STORY the full lint still counts every file, so CI and make lint-md are not weakened.
- 1: flai/internal/mdlint/repo_test.go, flai/internal/mdlint/mdlint_test.go, scripts/lint-md.sh, scripts/README.md, design/system/flai-cli.md
- 2: design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/summary.md
