---
id: D-0251
title: "S-0291 plan: retry dropped GitHub calls in install.sh and self-upgrade, show install.sh's error, close I-0086: Plan approved: the tasks, touches, and figures fit S-0291"
kind: answer
sources: [TH-0367]
topics: [release, flai, cli, go, code, all]
paths: [install.sh, flai/cmd/installsh_test.go, flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md, design/issues/summary.md, design/issues/I-0122-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0251 S-0291 plan: retry dropped GitHub calls in install.sh and self-upgrade, show install.sh's error, close I-0086: Plan approved: the tasks, touches, and figures fit S-0291

The operator, alex, answered TH-0367 at 2026-10-08T08:09:46Z.

What was asked:

```text
Plan. I-0086's seven instances and TH-0365 show one cause. GitHub drops the connection while `install.sh` lists `releases?per_page=50`, about 1 MB, with curl (92) or (56). Nothing retries the call. `scripts/install-test.sh` then hides `install.sh`'s error in its default-path case. The remedy:
```

The answer:

```text
Plan approved: the tasks, touches, and figures fit S-0291.

- The three layer-1 fixes share no path. The Go tests that cut the listing off under `flai/` reproduce I-0086 where the close-out's tiers will run them.
- The 1h forecast allows for the integration and smoke reruns.
- Assumption 4 (smoke still reaches GitHub) is consistent with S-0340, which makes the smoke tier hermetic separately.

Recommendation for the operator, on the cost of delay input: set `time_lost_per_cycle` to `2h` on S-0291, as the planner recommends. I-0086 has 7 occurrences in 2 days, each costing about a 5-minute close-out rerun, which is about 24 a week, or about 2h. It also held S-0326 out of review. I could not set it myself: flai refuses me cost of delay inputs on a story that is no longer in the backlog, and you moved S-0291 to ready at 07:57Z. The planner's value of 300 USD a week already assumes it. Command: `flai edit S-0291 --cod-time-lost 2h`, or the story's Cost of delay on the dashboard.

Source: wip/kanban/stories/S-0291-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md
```

Source: TH-0367.
