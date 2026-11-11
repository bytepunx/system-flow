---
id: S-0722
type: story
nature: improvement
title: "A close-out's Go tests run under -race on every core, so three stories closing out at once took the host to a load of 84: run them on one core without -race"
status: in-progress
owner: alex
created: 2026-11-11T05:50:20Z
updated: 2026-11-11T06:14:00Z
transitions:
  - to: ready
    at: 2026-11-11T05:50:21Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T06:02:34Z
    by: agent-S-0722
tags: [tests, close-out]
touches: [system-flow.yaml, scripts/integration.sh, scripts/test.sh, scripts/flai-test.sh, design/system/workflow.md, scripts/race.sh, Makefile, flai/internal/verify/select_test.go, design/system/project-manifest.md, design/system/devex.md, design/conventions/code-quality.md, design/conventions/tooling.md, design/tech/go.md, scripts/README.md, docs/operators/index.md, docs/contributors/index.md, flai/README.md]
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
      seconds: 53
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 91
          cache_read: 281857
          cache_write: 63942
          cost: 0.096
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: alex
    at: 2026-11-11T06:14:00Z
forecast:
  duration: 2h
  delivery: 2026-11-11T08:03:00Z
  basis: "flai forecast's 20m (107 s a unit of size times 11) raised to 2h, because the story runs the full Go suite at least six times: before and after measurements, three runs in a row at one core, and the race tier under -race, about 10 minutes each."
  by: planner-S-0722
  at: 2026-11-11T06:07:37Z
---
# S-0722 A close-out's Go tests run under -race on every core, so three stories closing out at once took the host to a load of 84: run them on one core without -race

## Goal

The Go tests a story's close-out runs use one core and no race detector, so several stories can close out at once without saturating the host. The race detector still runs, by hand and where the whole suite is run on purpose, not in every close-out.

## Acceptance criteria
- [ ] The `go-test` tier in `system-flow.yaml` and the `integration` tier's `scripts/integration.sh`, as a close-out runs them, run `go test` without `-race`, with `-p 1` (one package at a time, which also bounds the compiler), and with the test binaries held to one core (`GOMAXPROCS=1`); `scripts/test.sh` matches.
- [ ] The race detector is kept as a tier of its own, `race`, marked `manual: true` like `gateway`: `flai test race` runs the full suite under `-race` as `integration.sh` does today, and `scripts/flai-test.sh` (`make flai-test`) still runs it. No close-out and no merge queue batch runs it unasked.
- [ ] The narrative records, measured on this host for a change that touches `flai/internal/mcpserver`: the close-out's Go tests' wall time, user CPU, and peak cores before and after. After, one close-out's Go tests hold no more than about one core.
- [ ] Every test passes at one core, three runs in a row. A test that only passes with more is fixed if the fix is small, and otherwise the story stops and says which on a thread; it is not skipped and the cap is not raised to hide it.
- [ ] The timeouts sized for `-race` (20 minutes a package, and the comment in `integration.sh` citing S-0507 and I-0225) are set again from the measured figures, or kept with the reason restated.
- [ ] `design/system/workflow.md` (or where the tiers are described) says what a close-out runs, that it is one core by design, and when to run `flai test race`; the template's copy of anything changed here is changed with it.

## Tasks
- T-3086 Run the go-test tier, integration.sh, and test.sh without -race, one package at a time, on one core
- T-3087 Keep the race detector as a manual race tier that flai test race and make flai-test run
- T-3088 Measure the close-out's Go tests at one core, pass the suite three runs in a row, and set the timeouts again
- T-3089 Say where the tiers are described that a close-out's Go tests run on one core without -race, and when to run flai test race
- T-3090 The close-out's Go tests run without -race, one package at a time, on one core, and a manual race tier keeps the detector
- T-3091 Measure the close-out's Go tests before and after, run the suite three times at one core, and set the timeouts again from the figures
- T-3092 The design and docs say what a close-out's Go tests run, that it is one core by design, and when to run flai test race

## Notes

Asked for by the operator on 2026-11-10: change the test runs from `-race` to single threaded, because more than one story usually runs at a time and the close-outs together saturate the CPU.

`-race` and threads are two settings, and this story changes both:

- `-race` is Go's data race detector. It does not add threads; it makes each one cost several times more. `./internal/mcpserver` alone, `-short`, on this host: 608 s of user CPU with it, 65 s without (56 s and 10 s wall).
- How many cores a run takes is `-p` (packages built and tested at once, default one per core, 24 here), `GOMAXPROCS` (threads running Go code in each test binary, default one per core), and `-parallel` (tests marked `t.Parallel` at once, default `GOMAXPROCS`).

Where they are set today: `system-flow.yaml` tier `go-test` (`go test -json -race -short ... {importers}`), `scripts/integration.sh` (`go test -race -count=1 -timeout 20m ... ./...`, tier `integration`, `queue: true`), and `scripts/test.sh`.

What is given up: a data race in a story's change is no longer caught by its close-out. This project is not on the merge queue, so the close-out was the only place the detector ran unasked. The `race` tier keeps it one command away; whether something should run it on a schedule or before a release is not this story's.

What to expect: at one core the plain mcpserver suite is about its user CPU in wall time, a minute or so, and less once S-0721's schema cache lands. Tests that lean on timing (held waits with 10 to 20 ms heartbeats and a 3 s limit) are the likeliest to object to one core, which is what the fourth criterion is for.

`system-flow.yaml`'s tiers are flai's to write: change them the way the manifest's owner allows, not by hand if a command does it.

Related: S-0721 (share one SchemaCache) cuts the work itself; this story bounds what is left.

### Planning

Two task sets were written at once: the planner's T-3086 to T-3089 and the story's agent's T-3090 to T-3092, for the same work. Only one set should stay; the plan's thread on S-0722 proposes which, and the story's agent decides.

Touches, file by file, and where each came from:

| Touch | From |
|-------|------|
| `system-flow.yaml`, `scripts/integration.sh`, `scripts/test.sh`, `scripts/flai-test.sh`, `design/system/workflow.md` | declared |
| `scripts/race.sh`, `Makefile` | layout: the `race` tier needs a script of its own beside `gateway-smoke.sh`, and a Make target calls it |
| `flai/internal/verify/select_test.go` | layout and co-change (11 of 152 commits): `TestRepositoryManifestRunsChangedGoPackagesShortAndIntegrationWhole` pins the `go-test` flags with `-race`, and the race tier's test sits beside the gateway tier's |
| `design/system/project-manifest.md`, `design/system/devex.md`, `docs/operators/index.md`, `scripts/README.md` | design and co-change: each quotes the tiers or the `-race` commands |
| `design/conventions/code-quality.md`, `design/conventions/tooling.md`, `design/tech/go.md` | design: each says this repository's Go tests run `-race`, below the marker in the two conventions |
| `docs/contributors/index.md`, `flai/README.md` | the story's agent's T-3092, which names them |

No folder touch: the test files a fix at one core may change cannot be named yet, and `flai task done` widens the story's touches with what the measuring task commits. The template's `scripts/test.sh` and `scripts/integration.sh` are placeholders and its conventions do not mention `-race`, so no template path is predicted. `flai touches suggest` also listed `design/system/flai-cli.md`, `docs/users/flai.md`, and the ADR index at 20 to 54% co-change; none describes these tiers' flags, so they are left out.

No `after`: the story reads nothing another open story writes. S-0721, which it names as related, is done.

Forecast: `flai forecast` gave 20m (107 s a unit of size over 35 done improvement stories, times size 11). Raised to 2h, because the story runs the full Go suite at least six times, before and after, three runs in a row at one core, and the race tier, about 10 minutes each under `-race` and several minutes at one core.

Cost of delay: the story has no inputs and no epic, so `flai cod` refuses; the inputs are asked for on the plan's thread.
