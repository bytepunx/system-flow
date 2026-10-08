---
id: S-0291
type: story
nature: improvement
title: The close-out's install smoke test failed once and passed when run alone, with no cause in its output
status: done
owner: alex
created: 2026-10-06T10:31:55Z
updated: 2026-10-08T08:42:00Z
transitions:
  - to: ready
    at: 2026-10-08T07:57:14Z
    by: alex
  - to: in-progress
    at: 2026-10-08T08:13:36Z
    by: agent-S-0291
  - to: review
    at: 2026-10-08T08:37:02Z
    by: agent-S-0291
  - to: done
    at: 2026-10-08T08:42:00Z
    by: alex
tags: []
topics: [release]
touches: [install.sh, flai/cmd/installsh_test.go, flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md, design/issues/summary.md, design/issues/I-0122-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1420
  turns:
    - day: 2026-10-08
      ceremony: 3
      hand_edits: 1
      work: 36
  models:
    - model: claude-opus-5-5
      input: 146
      output: 53395
      cache_read: 7923083
      cache_write: 319658
      cost: 4.8392
  strategic:
    - kind: orchestrator
      seconds: 1347
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 124
          output: 1995
          cache_read: 24140768
          cache_write: 41752
          cost: 5.9585
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-10-06T10:31:55Z
  value: 300
  by: planner-S-0291
  at: 2026-10-08T08:01:29Z
forecast:
  duration: 1h
  delivery: 2026-10-08T09:51:00Z
  basis: "Its own forecast of 1h; 3rd in the pull order with an in-progress limit of 3, behind S-0232, S-0321, S-0339, S-0342 and S-0322."
  by: flai
  at: 2026-10-08T08:07:59Z
finalized:
  by: alex
  at: 2026-10-07T02:19:00Z
---
# S-0291 The close-out's install smoke test failed once and passed when run alone, with no cause in its output

## Goal

This story remediates [I-0086](../../../design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md), "The close-out's install smoke test failed once and passed when run alone, with no cause in its output". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0086 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0086 is closed with `flai issue close I-0086 --reason` saying what fixed it

## Tasks
- T-1338 install.sh retries a GitHub API call the network drops and resolves the latest release from a small page
- T-1339 flai self-upgrade retries a release listing or download the network drops
- T-1340 scripts/install-test.sh prints install.sh's output when the default-path install fails
- T-1341 flai-cli.md and flai.md say install.sh and self-upgrade retry a call the network drops
- T-1342 Close I-0086 saying the installer and self-upgrade retry a dropped GitHub connection

## Notes

Cost of delay inputs set by flai from I-0086. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T10:22:10Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from I-0086's seven instances and TH-0365. The cause is GitHub dropping the connection while `install.sh` lists `releases?per_page=50`, about 1 MB, with curl (92) `HTTP/2 stream 1 was not closed cleanly` or curl (56) `unexpected eof while reading`. Nothing retries it. `scripts/install-test.sh` then hides `install.sh`'s error in its default-path case. TH-0365 measured 2 of 5 fetches of that listing dropped, against 3 of 3 for `per_page=5`. The fix has three parts:

- `install.sh` retries a call that fails on the network, never an HTTP error, and resolves the latest release from a small page (T-1338).
- `flai self-upgrade`, which the same smoke step runs four times, does the same for its page reads and downloads (T-1339).
- `scripts/install-test.sh` prints `install.sh`'s output when it fails (T-1340).

Go tests with an `httptest` server that cuts the response off mid-body reproduce the drop (T-1338, T-1339). They live under `flai/` because the close-out selects the Go tiers by `flai/**`. `install.sh` and `scripts/` alone select none of them.

Touches: `flai touches suggest S-0291 install.sh scripts/install-test.sh` found 4 of 1434 commits changing them.

| Touch | Source |
|-------|--------|
| `install.sh` | design: I-0086's instances name its release listing |
| `scripts/install-test.sh` | design: I-0086's last instance names its `set -e` exit before `cat "$OUT"` |
| `flai/cmd/installsh_test.go` | layout: new test beside `flai/cmd/settings_doc_test.go`, which already reads `../../install.sh` |
| `flai/internal/selfupgrade/selfupgrade.go`, `selfupgrade_test.go` | layout: `List` and `getBytes`, the Go client the smoke step runs, have no retry |
| `design/system/flai-cli.md` | co-change, 3 of 4 commits; its `flai self-upgrade` row describes both |
| `docs/users/flai.md` | co-change, 2 of 4 commits; its `self-upgrade` paragraph |
| `design/issues/I-0086-…md`, `design/issues/summary.md` | goal: criterion 2 closes I-0086 |

`README.md` co-changed in 2 of 4 commits but only shows the install one-liner, which does not change, so it is left out. No folder touch is kept.

Forecast: `flai forecast` gave 16m at size 11. It is raised to 1h. The story has five tasks, and two of them are Go tests that hijack a stand-in server's connection. Its `flai/` changes make the close-out run integration and smoke, about 5 minutes a run (TH-0365), and smoke reaches GitHub, so a rerun is likely. The delivery is flai's start, about 08:44Z, plus 1h: 09:50Z.

Cost of delay: `flai cod` gave 15 USD a week from the input of 6m per 168h cycle, which counted one occurrence. The value is set to 300 USD a week:

- I-0086 now has 7 occurrences in 2 days, about 24 a week.
- Each stops a close-out after integration has passed, and the rerun costs about 5 minutes (TH-0365).
- That is about 2h a week at 150 USD an hour, which is 306 USD a week.
- It also holds S-0326 out of review (TH-0365).

The input is the operator's and is left as it is. The plan's thread recommends raising it to about 2h.
