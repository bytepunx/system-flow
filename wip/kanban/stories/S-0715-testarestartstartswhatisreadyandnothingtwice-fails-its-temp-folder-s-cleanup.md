---
id: S-0715
type: story
nature: remediation
title: "TestARestartStartsWhatIsReadyAndNothingTwice fails its temp folder's cleanup with \"directory not empty\", as I-0048's and I-0233's serve tests did"
status: backlog
owner: alex
created: 2026-11-10T10:55:22Z
updated: 2026-11-11T14:52:10Z
transitions: []
tags: []
topics: [testing]
touches: [flai/internal/serve/agents_test.go, design/issues/I-0343-testarestartstartswhatisreadyandnothingtwice-fails-its-temp-folder-s-cleanup.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 163
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 48
          output: 189
          cache_read: 1856599
          cache_write: 126488
          cost: 0.1355
        - model: claude-opus-5-5
          input: 164
          output: 11719
          cache_read: 6499039
          cache_write: 475131
          cost: 2.5108
    - kind: orchestrator
      seconds: 66
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 18
          output: 144
          cache_read: 733681
          cache_write: 87893
          cost: 0.2297
cost_of_delay:
  inputs:
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-11-10T10:55:22Z
  value: 30
  by: planner-S-0715
  at: 2026-11-11T09:04:21Z
forecast:
  duration: 40m
  delivery: 2026-11-12T03:34:00Z
  basis: "Its own forecast of 40m; 19th in the pull order with an in-progress limit of 3, behind S-0689, S-0716, S-0719, S-0699, S-0705, S-0681, S-0718, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0701, S-0703, S-0707, S-0709, S-0711, S-0712 and S-0714."
  by: flai
  at: 2026-11-11T14:52:10Z
finalized:
  by: alex
  at: 2026-11-11T09:02:10Z
---
# S-0715 TestARestartStartsWhatIsReadyAndNothingTwice fails its temp folder's cleanup with "directory not empty", as I-0048's and I-0233's serve tests did

## Goal

This story remediates [I-0343](../../../design/issues/I-0343-testarestartstartswhatisreadyandnothingtwice-fails-its-temp-folder-s-cleanup.md), "TestARestartStartsWhatIsReadyAndNothingTwice fails its temp folder's cleanup with "directory not empty", as I-0048's and I-0233's serve tests did". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0343 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0343 is closed with `flai issue close I-0343 --reason` saying what fixed it

## Tasks
- T-3112 The agent lab waits for every launcher it makes before its temp folders go, so the restart test's fresh launcher no longer writes into a folder being removed
- T-3113 Close I-0343 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0343. time_lost_per_cycle 12m: 12m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T10:41:09Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed cause, from the instance and the code. The failing path, `.../001`, is the lab's first temp folder, the project root `newAgentLab` makes. The test makes a second launcher, `fresh := newLauncher(lab.o, lab.l.entry)` in `flai/internal/serve/agents_test.go`, and `fresh` starts the missed story's agent and waits for it. The lab's cleanup, `settle`, waits only for `lab.l.waiting` to empty and for each stub's `ended-<pid>` file. The stub writes that file before it exits. So `settle` can return while `fresh`'s `await` goroutine is still in `end`, where `endRun` and `measure` write the story's run and usage under the root. That is the I-0048 cause, a launcher nobody waits for, on a launcher S-0122's fix did not reach.

Proposed fix: the lab waits for every launcher it makes, not only `lab.l`. The test gets `fresh` from the lab, and `settle` waits until no launcher of the lab has a PID in `waiting`. `settle` also reports, with `t.Logf`, the PIDs it still waited for if its deadline passes, so that the next failure of this kind names its launcher. The fix stays in the test file: no production code changes. If the story's agent finds the cause is elsewhere, it widens the touches when it closes the task.

Touches:

- `flai/internal/serve/agents_test.go`: design. It holds the test, `newAgentLab`, and `settle`. `flai touches suggest S-0715 flai/internal/serve/agents_test.go` lists `flai/internal/serve/agents.go` at 84% co-change; it is left out because the fix needs no change to the launcher.
- `design/issues/I-0343-...md` and `design/issues/summary.md`: layout. `flai issue close` writes both.

Figures:

- Forecast: `flai forecast` gave 6m, from a median of 152 s per unit over 22 done remediation stories, at size 2 with no touches. It is raised to 40m. Reproducing and confirming the fix takes repeated `-race` runs of the serve package, and the close-out's integration tier took 8m21s on this package in the instance. Delivery stays at flai's 2026-11-11T22:03:00Z from the pull order.
- Cost of delay: 30 USD a week, as `flai cod` gives it from the operator's input of 12m lost per 168h cycle at 150 USD an hour. It stands as given.

No `after`: the story reads nothing another open story writes.
