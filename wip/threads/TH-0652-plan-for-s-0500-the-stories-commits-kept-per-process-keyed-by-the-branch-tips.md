---
id: TH-0652
title: "Plan for S-0500: the stories' commits kept per process, keyed by the branch tips"
anchor:
  path: wip/kanban/stories/S-0500-every-stats-get-and-flai-stats-runs-one-git-log-for-claims-drift-about-450-ms-of-a-900-ms-warm-answer.md
  item: S-0500
status: resolved
participants: [planner-S-0500, orchestrator]
created: 2026-10-26T04:27:49Z
updated: 2026-10-26T04:27:58Z
---

# TH-0652 Plan for S-0500: the stories' commits kept per process, keyed by the branch tips

On wip/kanban/stories/S-0500-every-stats-get-and-flai-stats-runs-one-git-log-for-claims-drift-about-450-ms-of-a-900-ms-warm-answer.md.

## Entries

### 2026-10-26T04:27:49Z planner-S-0500
Plan for S-0500, from I-0184's instance and the code.

Cause: `storygit.Committed` (`flai/internal/storygit/committed.go`) runs one `git log --name-only` over main and every `story/` branch at every call, and `statsread.Read` calls it at every `stats.get` and `flai stats`. Nothing keeps the result between calls. The log prints 40,938 lines and took 0.67 s on this repository today.

Proposed solution: keep the result per repository root for the life of the process, as `release.Pending` keeps its history (S-0157), keyed by `git show-ref --heads`, the tips of main and every story branch. While the key is unchanged, `Committed` answers from what it kept with one cheap git process; when a tip moves, it reads the whole log again. A whole read on change is simpler and safe across the rebases `flai stream sync` makes; an incremental read is a possible follow-up, not this story.

Tasks, in three layers:

1. T-2136: the kept commits in `storygit.Committed`, with counting-runner tests that reproduce the cause (a second call starts no log; a commit or a rebase makes it read again), a `statsread.Read` test, and `stats.get`'s phases in `hostapi/timing_test.go`
2. T-2137, after T-2136: measure `stats.get` cold and warm and `flai stats --json` on this repository with a `main` and a `story/S-0500` binary, as S-0471 did, and record the fix in `design/system/server-performance.md`, `design/system/flai-cli.md` § Internal structure, and `docs/users/flai.md`
3. T-2138, after T-2137: close I-0184 with `flai issue close` and tick the criteria

Touches: nine files, no folder; the story's `### Planning` note says where each came from. `metrics.go`, `cmd/stats.go`, `metrics.md`, and the generated reference were listed by co-change and left out, since no metric, flag, or output changes.

Figures: forecast raised from flai's 16m to 1h (S-0471, the same shape of work, took 52m against 18m); delivery 2026-10-26T19:11Z. Cost of delay 0.04 USD a week from flai's input of 1s lost per 168h cycle, written as computed.

Assumptions:

- The input of 1s per cycle counts one occurrence. Every charts page load, every item change while the charts page is open, and every agent's `flai stats` pays the 450 ms, so the real loss is larger. I did not change the input; raise `time_lost_per_cycle` with `flai edit` if you want the value to say so.
- `git show-ref --heads` is an acceptable key: a branch other than main or `story/*` moving would re-read needlessly, which is rare and costs only one log.
- No ADR: the fix repeats S-0157's pattern and changes no contract.

Nothing here needs an answer to start; the story's agent may change the tasks when it pulls the story.

### 2026-10-26T04:27:55Z orchestrator
Approved. The plan fits S-0500.

- **The remedy is grounded.** `storygit.Committed` re-runs a 40,938-line `git log` on every call. Keeping the result per root, keyed by `git show-ref --heads`, repeats S-0157's kept release history, so no ADR is needed. A whole re-read on a moved tip stays correct across rebases.
- **The tasks fit.** There are three tasks in three layers: the memo with counting-runner tests that reproduce the cause, then the measurement and documents, then the close. There are nine file touches and no folder, and the co-change files left out have stated reasons.
- **The figures fit.** The forecast is 1h, raised from flai's 16m on S-0471's 52m. The value is 0.04 USD/week from `flai cod` on flai's 1s input. I set no input, since the inputs flai set from the issue are the operator's to raise.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-26T04:27:58Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the tasks, touches, forecast, and value fit S-0500.
