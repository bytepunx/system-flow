---
id: S-0316
type: story
nature: remediation
title: The dashboard upgrade stops the old container and cannot start the new one, because the name flaiover is still in use
status: done
owner: alex
created: 2026-10-07T18:59:46Z
updated: 2026-10-08T04:49:57Z
transitions:
  - to: ready
    at: 2026-10-07T23:55:54Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T23:59:51Z
    by: agent-S-0316
  - to: review
    at: 2026-10-08T04:42:08Z
    by: agent-S-0316
  - to: done
    at: 2026-10-08T04:49:57Z
    by: orchestrator
tags: []
touches: [flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_test.go, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/runbooks/update.md, design/system/flai-cli.md, design/issues/I-0094-the-dashboard-upgrade-stops-the-old-container-and-cannot-start-the-new-one-because-the-name-flaiover-is-still-in-use.md, design/issues/summary.md, design/issues/I-0116-flai-dashboard-restart-starts-the-container-s-tag-again-so-a-newer-image-a-check-pulled-under-a-floating-tag-is-started-without-an-upgrade.md, design/issues/I-0109-flai-check-finds-narrative-state-outside-the-story-at-close-out.md, design/issues/I-0110-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md, design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4656
  turns:
    - day: 2026-10-07
      work: 4
    - day: 2026-10-08
      ceremony: 9
      test_runs: 6
      hand_edits: 2
      work: 73
  models:
    - model: claude-opus-5-5
      input: 262
      output: 68329
      cache_read: 19614682
      cache_write: 590241
      cost: 9.5925
  strategic:
    - kind: orchestrator
      seconds: 12792
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 88
          output: 1308
          cache_read: 27106589
          cache_write: 55587
          cost: 6.6922
cost_of_delay:
  inputs:
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-10-07T18:59:46Z
  value: 30
  by: planner-S-0316
  at: 2026-10-07T23:54:58Z
forecast:
  duration: 30m
  delivery: 2026-10-08T06:54:00Z
  basis: "flai forecast's 16m (94 s per unit over 17 done medium remediation stories, size 10) raised to 30m: the fake docker runner needs a name held after stop, and two commands plus a rollback change; delivery is flai's, moved by the 14m added"
  by: planner-S-0316
  at: 2026-10-07T23:54:58Z
finalized:
  by: orchestrator
  at: 2026-10-07T23:55:50Z
---
# S-0316 The dashboard upgrade stops the old container and cannot start the new one, because the name flaiover is still in use

## Goal

This story remediates [I-0094](../../../design/issues/I-0094-the-dashboard-upgrade-stops-the-old-container-and-cannot-start-the-new-one-because-the-name-flaiover-is-still-in-use.md), "The dashboard upgrade stops the old container and cannot start the new one, because the name flaiover is still in use". The issue recommends this solution:

Directions to weigh: after stopping the old container, wait until Docker no longer lists the name before running the new one, with a short limit; or run the new container under a temporary name and rename it; and when the start fails, start the previous image again so that the operator is not left without a dashboard. A test with a stand-in for docker that keeps the name for a moment after `stop` would reproduce it.

## Acceptance criteria
- [x] The cause I-0094 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0094 is closed with `flai issue close I-0094 --reason` saying what fixed it

## Tasks
- T-1295 Wait for the container name to be released before the swap, and start the previous image again when the new one fails
- T-1296 Say in the operator docs and the design how the upgrade waits for the name and falls back, and close I-0094

## Notes

Cost of delay inputs set by flai from I-0094. time_lost_per_cycle 12m: 6m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-06T21:01:06Z, 0.9 days before this story; under one cycle counts as one).

### Planning

Touches, file by file; no folder touch is kept.

| Touch | From | Why |
|-------|------|-----|
| `flai/cmd/dashboard_upgrade.go` | layout | `runDashboardUpgrade` runs `docker stop` and then `startDashboard` under the same name at once; `runDashboardRestart` does the same |
| `flai/cmd/dashboard_test.go` | layout | `fakeRunner` and the upgrade tests live here; the stand-in that keeps the name after `stop` goes here |
| `docs/users/flai-reference.md` | design | generated from the upgrade command's help by `make flai-reference`; changes when its `Long` text does |
| `docs/operators/index.md` | design | its paragraph on the upgrade says what happens when the new image fails |
| `docs/operators/runbooks/update.md` | design | step 2 of the dashboard update says what the operator sees on failure |
| `design/system/flai-cli.md` | design | the builders' description of how `flai dashboard upgrade` swaps the container |
| `design/issues/I-0094-…md` | declared by the goal | `flai issue close` writes it |
| `design/issues/summary.md` | declared by the goal | `flai issue close` regenerates it |

`flai touches suggest` lists `docs/users/flai.md` (33%), `design/system/flaiover-dashboard.md` (15%), and `flai/cmd/dashboard.go` (3%) by co-change. None is predicted: `docs/users/flai.md` only names the command in a comment line that stays true, and `containerRunning` in `flai/cmd/dashboard.go` can be called as it is. The story's agent widens the touches if the fix needs either.

Forecast: `flai forecast` gave 16m, 94 s per unit of size over 17 done remediation stories in the medium band, times size 10 (2 criteria, 8 touches). Raised to 30m. `fakeRunner` has no notion of a name held after `stop`, so the reproducing test needs new fixture behaviour; the fix changes both `upgrade` and `restart` and adds a rollback path. Delivery is flai's, 2026-10-08T06:40Z, moved by the 14m added.

Cost of delay: `flai cod` gives 30.00 USD a week, from the inputs flai set (12m lost per 168h cycle at 150 USD an hour). It stands. The issue's own text counts four occurrences in the journal against the two in its `count`, so the inputs understate the loss. They are the operator's to raise.

### Accepted by the orchestrator

- Verified: cf8169d9f498e5bbaa3f29f3a4cf147408fd566e
- At: 2026-10-08T04:49:57Z

Verdict: accept; both criteria met (verifier at cf8169d9f498e5bbaa3f29f3a4cf147408fd566e; flai verify passed every step at that commit). Minor: I-0094's close reason cites 66f2113e, a pre-rebase hash of the fix (003bd05a on the branch). The issue bumps on the branch are the close-outs' own records, all within the touches.

- 1: flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_test.go, design/system/flai-cli.md, docs/operators/index.md, docs/operators/runbooks/update.md, docs/users/flai-reference.md
- 2: design/issues/I-0094-the-dashboard-upgrade-stops-the-old-container-and-cannot-start-the-new-one-because-the-name-flaiover-is-still-in-use.md, design/issues/summary.md
