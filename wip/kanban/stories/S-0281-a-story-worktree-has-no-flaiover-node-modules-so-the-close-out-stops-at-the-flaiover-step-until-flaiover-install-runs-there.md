---
id: S-0281
type: story
nature: improvement
title: A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there
status: backlog
owner: alex
created: 2026-10-05T07:09:06Z
updated: 2026-10-06T23:17:28Z
transitions: []
tags: []
touches: [scripts/flaiover-install.sh, scripts/flaiover-test.sh, scripts/flaiover-unit.sh, scripts/README.md, design/system/devex.md, docs/operators/index.md, design/issues/I-0080-a-story-worktree-has-no-flaiover-node-modules-so-the-close-out-stops-at-the-flaiover-step-until-flaiover-install-runs-there.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-05T07:09:06Z
  value: 12.5
  by: planner-S-0281
  at: 2026-10-06T22:53:49Z
forecast:
  duration: 30m
  delivery: 2026-10-07T10:26:00Z
  basis: "Its own forecast of 30m; 32nd in the pull order with an in-progress limit of 3, behind S-0261, S-0300, S-0301, S-0228, S-0269, S-0270, S-0302, S-0271, S-0212, S-0213, S-0214, S-0215, S-0216, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254, S-0265, S-0272, S-0273, S-0274, S-0275, S-0277, S-0279 and S-0280."
  by: flai
  at: 2026-10-06T23:17:28Z
finalized:
  by: alex
  at: 2026-10-06T22:49:26Z
---
# S-0281 A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there

## Goal

This story remediates [I-0080](../../../design/issues/I-0080-a-story-worktree-has-no-flaiover-node-modules-so-the-close-out-stops-at-the-flaiover-step-until-flaiover-install-runs-there.md), "A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0080 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0080 is closed with `flai issue close I-0080 --reason` saying what fixed it

## Tasks
- T-1088 The flaiover scripts install flaiover's dependencies in a story worktree that lacks them, so close-out's flaiover step and vitest run there
- T-1095 The scripts README, devex design, and operators' guide say a story worktree installs flaiover's dependencies, and I-0080 is closed

## Notes

Cost of delay inputs set by flai from I-0080. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-05T07:07:27Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0281 on 2026-10-06.

**Proposed fix.** The issue gives no solution yet, so the plan proposes one. The cause is that `flai stream open` makes the worktree with `git worktree add` (`flai/cmd/branch.go`), which leaves out the git-ignored `flaiover/node_modules`. The fix goes in the project's scripts, not in flai, because only this project has a `flaiover/`:

- `scripts/flaiover-install.sh` gets a mode that installs only when the dependencies are missing or older than the lock file.
- `scripts/flaiover-test.sh` always calls it.
- `scripts/flaiover-unit.sh` calls it in a story worktree. That script now skips vitest there without a word, a second effect of the same cause.

The plan thread asks the operator to confirm this.

**Touches.** The story declared none, so `flai touches suggest` had nothing to start from. Each touch is a file, and none is a folder.

- `scripts/flaiover-install.sh`, `scripts/flaiover-test.sh`, `scripts/flaiover-unit.sh` (layout). These are the scripts close-out runs for flaiover's step, and the ones the issue's instances name.
- `scripts/README.md` (layout). It describes those three scripts.
- `design/system/devex.md` (design). Its Test tiers row says flaiover's vitest runs only when `flaiover/node_modules` is present.
- `docs/operators/index.md` (design). It says the same of `scripts/flai-test.sh` at close-out.
- The I-0080 file and `design/issues/summary.md` (criteria). `flai issue close` writes them.
- Co-change, from `flai touches suggest` run from these eight: it lists `design/system/flai-cli.md`, `docs/users/flai.md`, `design/system/flaiover-dashboard.md`, and others, at 31% or less. Each one comes from summary.md and devex.md changing alongside unrelated stories, so none was added. `scripts/close-out.sh` and `scripts/env.sh` are left out. T-1088 widens its touches to them only if the fix needs it.
- The template has no flaiover scripts (`template/root/scripts/`), so nothing under `template/` is touched.

**Forecast.** It is 30m, against `flai forecast`'s 15m: a median 86 s per unit of size over 7 done medium-band improvement stories on claude-opus-5-5, times size 10. I doubled it because T-1088's reproduction runs flaiover's prettier, eslint, svelte-check, and vitest at least twice in a scratch worktree. Each run takes minutes, which the size count does not see. The delivery, 2026-10-07T09:52Z, is flai's 09:37Z moved by the same 15m. The story is 33rd in the pull order.

**Cost of delay.** The value is 12.50 USD a week, as `flai cod` gives it: 5m lost per 168h cycle at 150 USD an hour. The issue now has two occurrences, S-0217 and S-0295, both in under two days, so the inputs, which counted one, may understate it. The inputs are the operator's, so they are left as they are. The plan thread mentions it.
