---
id: S-0724
type: story
nature: remediation
title: A go build in a story's worktree stamps vcs.revision of the main checkout's HEAD, so a worktree's bin/flai names the wrong commit
status: done
owner: alex
created: 2026-11-11T08:46:42Z
updated: 2026-11-11T18:15:36Z
transitions:
  - to: ready
    at: 2026-11-11T08:57:35Z
    by: alex
  - to: in-progress
    at: 2026-11-11T14:36:52Z
    by: agent-S-0724
  - to: review
    at: 2026-11-11T14:43:42Z
    by: agent-S-0724
  - to: done
    at: 2026-11-11T18:15:36Z
    by: alex
tags: [go, build, worktree]
touches: [scripts/flai-build.sh, scripts/flai-build-stamp-test.sh, scripts/smoke.sh, scripts/README.md, design/system/release-signing.md, docs/users/flai.md, design/issues/I-0348-a-go-build-in-a-story-s-worktree-stamps-vcs-revision-of-the-main-checkout-s.md, design/issues/summary.md, design/issues/I-0347-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 437
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 4
      hand_edits: 1
      work: 36
  models:
    - model: claude-opus-5-5
      input: 114
      output: 26684
      cache_read: 8161237
      cache_write: 257526
      cost: 4.0889
  strategic:
    - kind: planner
      seconds: 8
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 1518
          cache_read: 475818
          cache_write: 10774
          cost: 0.2045
    - kind: orchestrator
      seconds: 34
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 11
          output: 86
          cache_read: 369491
          cache_write: 91458
          cost: 0.1292
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-11T08:46:42Z
  value: 25
  by: planner-S-0724
  at: 2026-11-11T08:55:22Z
forecast:
  duration: 1h
  delivery: 2026-11-11T15:45:00Z
  basis: "Its own forecast of 1h; 1st in the pull order with an in-progress limit of 3, behind S-0694, S-0710 and S-0719."
  by: flai
  at: 2026-11-11T14:23:13Z
finalized:
  by: alex
  at: 2026-11-11T08:50:05Z
---
# S-0724 A go build in a story's worktree stamps vcs.revision of the main checkout's HEAD, so a worktree's bin/flai names the wrong commit

## Goal

This story remediates [I-0348](../../../design/issues/I-0348-a-go-build-in-a-story-s-worktree-stamps-vcs-revision-of-the-main-checkout-s.md), "A go build in a story's worktree stamps vcs.revision of the main checkout's HEAD, so a worktree's bin/flai names the wrong commit". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0348 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0348 is closed with `flai issue close I-0348 --reason` saying what fixed it

## Tasks
- T-3093 The smoke tier's flai-build-stamp-test.sh reproduces I-0348 in a scratch main checkout with a nested story worktree
- T-3100 flai-build.sh links the built tree's own commit and date and builds with -buildvcs=false
- T-3101 release-signing.md and flai.md say a tree build links its own commit, and I-0348 is closed

## Notes

Cost of delay inputs set by flai from I-0348. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-11T07:50:38Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0724 on 2026-11-11. The orchestrator asked for a plan again, naming a thread `asked`. No thread of that name exists, and none is on S-0724 or I-0348. So this is the story's first plan, made from I-0348.

The proposed fix: `scripts/flai-build.sh` builds with `-buildvcs=false` and links `buildinfo.Commit` and `buildinfo.Date` by `-ldflags -X`, taken from git run in the tree being built. git resolves a worktree's `.git` file, so the values are the story branch's. Go's own lookup cannot be fixed here: it takes only a `.git` folder for a repository's root, and a story's worktree sits inside the main checkout. Only `flai version` reads the VCS stamp, through `buildinfo.Get()`, and only when `Commit` is empty. Every check that a binary is current compares `Version`, never `Commit`. Release builds set `Commit` through GoReleaser already.

The story declared no touches. `flai touches suggest` was given `scripts/flai-build.sh` and `design/system/release-signing.md` as start paths. It found 19 commits that changed them.

| Touch | Source | Why |
|-------|--------|-----|
| `scripts/flai-build.sh` | layout | Builds `bin/flai` with no `ldflags`, so a worktree's build falls back to the VCS stamp of the main checkout (T-3100) |
| `scripts/flai-build-stamp-test.sh` | layout | New: the smoke tier's test that reproduces I-0348 in a scratch main checkout with a nested story worktree (T-3093) |
| `scripts/smoke.sh` | layout | Runs the new test, beside `agent-vcs-test.sh` (T-3093) |
| `scripts/README.md` | layout | Its table lists every script (T-3093, T-3100) |
| `design/system/release-signing.md` | design | Its "How flai is released" section says a build without `ldflags` falls back to `debug.ReadBuildInfo` (T-3101) |
| `docs/users/flai.md` | design, co-change (47%) | Its `Version` section says what `flai version` prints (T-3101) |
| `design/issues/I-0348-…checkout-s.md` | design | The issue the goal names, closed by T-3101 |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it (T-3101) |

Co-change also listed `design/system/flai-cli.md` (47%), `docs/operators/index.md` (42%), and other release documents (21% or less). They are left out, since they change with the release process, not with the local build. `scripts/install-test.sh` builds a test flai with no `Commit` either. It is left out, because that binary's commit is never read. There are no folder touches.

Figures:

- Forecast: `flai forecast` gave 6m with no touches. Re-run with the eight touches, it gave 23m: 135 s a unit over 59 done remediation stories, times size 10 (2 criteria, 8 touches). Raised to 1h, since T-3093's test builds a module in a scratch repository and the smoke tier, which runs `install-test.sh`, is slow. Delivery is moved from flai's 22:32Z by the extra 37 minutes.
- Cost of delay: `flai cod` gave 25.00 USD a week from flai's input: 10m lost per 168h cycle at 150 USD an hour. Kept as given. I-0348 has one occurrence, and nothing but `flai version` reads the wrong commit.
- No `after`. No other story writes output this one reads.
