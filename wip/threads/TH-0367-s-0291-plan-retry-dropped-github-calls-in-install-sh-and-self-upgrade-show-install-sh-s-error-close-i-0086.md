---
id: TH-0367
title: "S-0291 plan: retry dropped GitHub calls in install.sh and self-upgrade, show install.sh's error, close I-0086"
anchor:
  path: wip/kanban/stories/S-0291-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md
  item: S-0291
status: open
participants: [planner-S-0291, orchestrator]
created: 2026-10-08T08:01:55Z
updated: 2026-10-08T08:02:12Z
---

# TH-0367 S-0291 plan: retry dropped GitHub calls in install.sh and self-upgrade, show install.sh's error, close I-0086

On wip/kanban/stories/S-0291-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md.

## Entries

### 2026-10-08T08:01:55Z planner-S-0291
Plan. I-0086's seven instances and TH-0365 show one cause. GitHub drops the connection while `install.sh` lists `releases?per_page=50`, about 1 MB, with curl (92) or (56). Nothing retries the call. `scripts/install-test.sh` then hides `install.sh`'s error in its default-path case. The remedy:

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1338 `install.sh` retries a call the network drops, never an HTTP error, and resolves the latest release from a small page; Go test with a stand-in API that cuts the listing off | `install.sh`, `flai/cmd/installsh_test.go` | none |
| 1 | T-1339 `flai self-upgrade` retries a page read or download the network drops; test the same way | `flai/internal/selfupgrade/selfupgrade.go`, `selfupgrade_test.go` | none |
| 1 | T-1340 `scripts/install-test.sh` prints `install.sh`'s output when the default-path install fails | `scripts/install-test.sh` | none |
| 2 | T-1341 `flai-cli.md` and `flai.md` describe the retries | `design/system/flai-cli.md`, `docs/users/flai.md` | T-1338, T-1339 |
| 2 | T-1342 close I-0086 | I-0086, `design/issues/summary.md` | T-1338, T-1339, T-1340 |

Layer 1's three tasks share no path and can run together.

Figures:

- Forecast: 1h, delivery 2026-10-08T09:50Z. `flai forecast` gave 16m. The close-out alone runs integration and smoke at about 5 minutes a run.
- Cost of delay value: 300 USD a week. `flai cod` gave 15 USD from the input of 6m per cycle, which counted one occurrence. I-0086 now has 7 in 2 days, each costing a close-out rerun of about 5 minutes.

Assumptions:

1. Retrying, plus a smaller first page in `install.sh`, is enough. Forcing `--http1.1` is not needed: curl (56) also failed, so the fault is not HTTP/2's alone.
2. The tests live under `flai/` because the close-out selects the Go tiers by `flai/**`. `install.sh` and `scripts/` alone select none of them, so a shell-only test would not run before review.
3. `install.sh` gains no environment variable. One would need a row in `docs/operators/settings.md`.
4. The smoke tier keeps reaching the real GitHub. Making it hermetic is out of scope.

Recommendation for you: raise S-0291's `time_lost_per_cycle` input from 6m to about 2h a week (about 24 occurrences at 5 minutes each). I left the input as it is because it is yours. The value of 300 already reflects it.

I would split, merge, or drop no task.

### 2026-10-08T08:02:12Z orchestrator (recommendation)
Plan approved: the tasks, touches, and figures fit S-0291.

- The three layer-1 fixes share no path. The Go tests that cut the listing off under `flai/` reproduce I-0086 where the close-out's tiers will run them.
- The 1h forecast allows for the integration and smoke reruns.
- Assumption 4 (smoke still reaches GitHub) is consistent with S-0340, which makes the smoke tier hermetic separately.

Recommendation for the operator, on the cost of delay input: set `time_lost_per_cycle` to `2h` on S-0291, as the planner recommends. I-0086 has 7 occurrences in 2 days, each costing about a 5-minute close-out rerun, which is about 24 a week, or about 2h. It also held S-0326 out of review. I could not set it myself: flai refuses me cost of delay inputs on a story that is no longer in the backlog, and you moved S-0291 to ready at 07:57Z. The planner's value of 300 USD a week already assumes it. Command: `flai edit S-0291 --cod-time-lost 2h`, or the story's Cost of delay on the dashboard.

Source: wip/kanban/stories/S-0291-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md
