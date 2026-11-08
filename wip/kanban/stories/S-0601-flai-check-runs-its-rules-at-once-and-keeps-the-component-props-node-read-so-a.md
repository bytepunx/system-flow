---
id: S-0601
type: story
nature: improvement
title: flai check runs its rules at once and keeps the component props node read, so a whole check takes about 1.5 s, not 3.3 to 8 s
status: review
parent: E-0033
owner: alex
created: 2026-11-07T22:39:49Z
updated: 2026-11-08T02:12:07Z
transitions:
  - to: ready
    at: 2026-11-08T00:11:15Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-08T01:10:39Z
    by: agent-S-0601
  - to: review
    at: 2026-11-08T02:12:07Z
    by: agent-S-0601
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/components/props.go, flai/internal/check/concurrent_test.go, flai/internal/components/propscache.go, flai/internal/components/propscache_test.go, flai/internal/components/props.mjs, design/system/server-performance.md, design/system/component-library.md, design/issues/I-0312-a-whole-flai-check-takes-3-3-to-4-9-s-on-a-loaded-host-after-s-0601-not-under-2.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3699
  priced_by: harness
  turns:
    - day: 2026-11-08
      ceremony: 2
      hand_edits: 2
      work: 91
  models:
    - model: claude-opus-5-5
      input: 290
      output: 91759
      cache_read: 40505931
      cache_write: 1369956
      cost: 17.8589
  strategic:
    - kind: orchestrator
      seconds: 141
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 17
          output: 111
          cache_read: 807456
          cache_write: 90001
          cost: 0.2298
        - model: claude-sonnet-5-5
          input: 4
          output: 40
          cache_read: 12683
          cache_write: 29547
          cost: 0.031
cost_of_delay:
  inputs:
    penalty_per_week: 50
    by: orchestrator
    at: 2026-11-07T23:51:41Z
  value: 50
  by: planner-E-0033
  at: 2026-11-08T00:00:08Z
forecast:
  duration: 40m
  delivery: 2026-11-08T01:23:00Z
  basis: "Its own forecast of 40m; 1st in the pull order with an in-progress limit of 5, behind S-0598, S-0599, S-0600, S-0605 and S-0606."
  by: flai
  at: 2026-11-08T00:23:55Z
finalized:
  by: alex
  at: 2026-11-07T23:49:31Z
---
# S-0601 flai check runs its rules at once and keeps the component props node read, so a whole check takes about 1.5 s, not 3.3 to 8 s

## Goal

Cause 11 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `check.Run` runs its rules one after another (`flai/internal/check/check.go`), and `componentExamples` starts `node` on every run (`flai/internal/components/props.go`), 0.6 to 4.2 s. A whole check takes 3.3 s cold and 1.8 s warm without `node`, and 7.4 to 8 s cold and 2.6 s warm with it. `flai check`, `flai verify`, the close-out's check step, and `flai task done` pay it.

## Acceptance criteria

- [x] `check.Run` runs its independent rules at once, gathering their findings under one lock, and gives the same findings in the same order as before, which a test over a fixture project with a finding of every rule shows.
- [x] The props `node` reads are kept under `.flai-cache`, keyed by each component file's path, size, and modification time and by the sub-project's svelte version, and `node` starts only for a component whose key changed.
- [ ] A whole `flai check --strict` on this repository, with nothing changed since the last, takes under 2 s, and server-performance.md records the before and after.

## Tasks
- T-2633 check.Run runs its independent rules at once and gives the same findings in the same order
- T-2643 Keep the props node reads under .flai-cache and start node only for a component whose key changed
- T-2648 Measure a whole flai check --strict before and after and record it in server-performance.md and component-library.md

## Notes

- The rules only read, so running them at once needs no other change. Run the tests under `-race`.
- Criterion 1's "under one lock": each rule writes its findings into its own result, and `Run` merges them in the old order after every worker is done. Nothing is shared while the rules run, so the merge needs no lock.
- Criterion 3 is left unticked. server-performance.md records the before and after, but the warm whole check is not under 2 s on this host. It took 3.3 to 4.9 s wall, against 7.9 to 11.7 s on `main`, at a load average of 84 to 103 on 24 cores. Under 2 s cannot be shown on a host that busy. The critical path that is left, the item list, `workItems` alone, and the digest rule, is recorded as I-0312.

### Planning

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/check/check.go` | declared | `Run` and the rules it calls |
| `flai/internal/components/props.go` | declared | `ReadProps`, which starts `node` |
| `flai/internal/check/concurrent_test.go` | layout | the new test of criterion 1, its fixture built in a temporary folder |
| `flai/internal/components/propscache.go` | layout | the kept props reads of criterion 2, beside `props.go` |
| `flai/internal/components/propscache_test.go` | layout | its test |
| `flai/internal/components/props.mjs` | design | S-0478's skeleton cache stamps every file typescript read; the script names them in its answer so the kept read goes stale with them (see the thread) |
| `design/system/server-performance.md` | design | criterion 3's before and after, under cause 11 |
| `design/system/component-library.md` | design | § Skeletons (S-0449) and § Examples (S-0443) › The check say how props are read; they gain the kept read |

`flai touches suggest` listed `check_test.go` (34%), `flai-cli.md` and `docs/users/flai.md` (22% each), and the manifest files. None is kept: the new test has a file of its own, and `flai check` changes no finding, flag, or setting a user or the manifest sees, only its speed. A task done widens the touches if an existing test needs changing.

No folder touch, and no `after`: S-0598 and S-0600 change `check.go`'s callers and the item list, but this story reads nothing they write; a trial merge judges the shared file.

Forecast: `flai forecast` gave 19m, 103 s per unit of size over 25 medium improvement stories, times size 11. Raised to 40m: the rules share state their order hides (`byID`, `vocab`, and the last-finding writes of `advise`), the kept read spans Go and `props.mjs`, and a close-out that runs `integration` takes about 4 min at p50 (server-performance.md cause 20). Delivery moves by the same 21 minutes, to 2026-11-08T00:40Z.

Cost of delay: no value yet. Neither S-0601 nor E-0033 has inputs, so `flai cod` refuses. The operator has not answered TH-0888, which recommends `penalty_per_week: 50` for this story.
