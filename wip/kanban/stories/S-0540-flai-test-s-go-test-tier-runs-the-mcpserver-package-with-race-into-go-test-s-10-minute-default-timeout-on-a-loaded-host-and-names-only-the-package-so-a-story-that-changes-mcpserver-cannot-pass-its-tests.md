---
id: S-0540
type: story
nature: remediation
title: flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests
status: backlog
owner: alex
created: 2026-10-29T05:59:43Z
updated: 2026-10-30T00:14:58Z
transitions: []
tags: [flai, testing]
topics: [cli, go]
touches: [flai/internal/verify/gotest.go, flai/internal/verify/parse_test.go, flai/internal/verify/testdata/go-test-timeout.json, design/system/project-manifest.md, system-flow.yaml, flai/internal/verify/select_test.go, design/conventions/tooling.md, flai/internal/mcpserver/prime_test.go, design/issues/I-0227-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md, design/issues/summary.md]
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
    - kind: orchestrator
      seconds: 1052
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 120
          output: 2016
          cache_read: 32511608
          cache_write: 88568
          cost: 7.189
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h12m
    by: flai
    at: 2026-10-29T05:59:43Z
  value: 180
  by: planner-S-0540
  at: 2026-10-29T07:36:26Z
forecast:
  duration: 1h30m
  delivery: 2026-10-30T14:16:00Z
  basis: "Its own forecast of 1h30m; 27th in the pull order with an in-progress limit of 5, behind S-0419, S-0511, S-0527, S-0424, S-0523, S-0524, S-0425, S-0529, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0506, S-0508, S-0509, S-0510, S-0512, S-0516, S-0517, S-0518, S-0520, S-0420, S-0421, S-0422 and S-0423."
  by: flai
  at: 2026-10-30T00:14:58Z
---
# S-0540 flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests

## Goal

This story remediates [I-0227](../../../design/issues/I-0227-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md), "flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0227 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0227 is closed with `flai issue close I-0227 --reason` saying what fixed it

## Tasks
- T-2326 The go-test-json format names each test still running when go test's timeout panics, with its running time
- T-2327 The go-test tier runs go test with -timeout 20m, as the integration tier does, pinned by the repository-manifest test
- T-2328 The mcpserver package's two prime-on-this-repository tests take less time under -race and still check what they check
- T-2329 Close I-0227 with flai issue close, naming the timeout finding, the tier's timeout, and the faster prime tests

## Notes

Cost of delay inputs set by flai from I-0227. time_lost_per_cycle 1h12m: 24m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-29T00:08:41Z, 0.2 days before this story; under one cycle counts as one).

### Planning

The proposed remedy takes all three candidates I-0227's first instance names, one task each:

1. The go-test-json format names the tests still running when go test's timeout panics. Go 1.26's `go test -json`, timed out, writes the panic and its `running tests:` list as output of the running test, sends no `fail` event for that test, and fails the package alone, so `parseGoTest` finds a package failure with no failed test and an empty message. That is why the finding names only the package.
2. The go-test tier passes `-timeout 20m`, as `scripts/integration.sh` does since I-0225, in place of go test's 10-minute default.
3. The two prime-on-this-repository tests in `flai/internal/mcpserver/prime_test.go` take less time without losing what they check. They stay under `-race`, since the integration tier runs only with `-race` too. `TestPrimeOnThisRepository` primes every open story, so its cost grows with the backlog.

Touches, file by file, no folder:

| Touch | From |
|-------|------|
| `flai/internal/verify/gotest.go`, `flai/internal/verify/parse_test.go` | layout: `parseGoTest` and its tests |
| `flai/internal/verify/testdata/go-test-timeout.json` | layout: a new fixture beside `go-test.json` |
| `design/system/project-manifest.md` | design: § `tests`, `format` says what each format reads |
| `system-flow.yaml` | layout: the go-test tier, written with `flai manifest set` as S-0490 did |
| `flai/internal/verify/select_test.go` | co-change: S-0490 pinned the tier's arguments there |
| `design/conventions/tooling.md` | design: its project addition describes the go-test tier |
| `flai/internal/mcpserver/prime_test.go` | the issue: the two slow tests and `thisRepository` live there |
| the I-0227 file, `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

`flai touches suggest` listed nothing, since the story declared no touches. S-0536, in progress, claims `flai/internal/mcpserver` and `design/system/project-manifest.md`, so this story is held while S-0536 is in progress. That is an overlap, not an `after`: nothing here needs S-0536's change on main.

Forecast: `flai forecast`, with the touches written, gave 24m (median 118 s per unit of size over 35 medium remediation stories, size 12), delivery 2026-10-29T23:42Z. Set to 1h30m, delivery 2026-10-30T00:48Z. The extra 66m pays for three `-race` runs of `flai/internal/mcpserver`, about 10 minutes each: before the change, after it, and in the close-out. The median does not include them.

Cost of delay: `flai cod` gives 180 USD a week from the inputs, 1h12m lost per 168h cycle at 150 USD an hour. The value stands as flai computed it. The inputs are the operator's. They likely understate the loss: the three occurrences came within five hours, from three stories. I-0225, the same slowness in the integration tier, is not counted in them.
