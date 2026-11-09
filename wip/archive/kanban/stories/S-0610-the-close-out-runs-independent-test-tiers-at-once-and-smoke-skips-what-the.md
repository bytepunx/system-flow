---
id: S-0610
type: story
nature: improvement
title: The close-out runs independent test tiers at once and smoke skips what the close-out already ran, saving about 40 to 150 s a close-out
status: done
parent: E-0033
owner: alex
created: 2026-11-07T22:42:03Z
updated: 2026-11-09T16:55:55Z
transitions:
  - to: ready
    at: 2026-11-08T00:11:23Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T07:38:44Z
    by: agent-S-0610
  - to: review
    at: 2026-11-09T10:43:49Z
    by: agent-S-0610
  - to: done
    at: 2026-11-09T16:55:55Z
    by: alex
tags: [cli]
topics: [server-side, back-end, manifest, template]
touches: [flai/internal/verify/run.go, flai/internal/verify/story.go, scripts/smoke.sh, flai/internal/verify/run_test.go, flai/internal/verify/story_test.go, flai/internal/verify/verify.go, flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, docs/operators/settings.md, design/system/project-manifest.md, template/root/system-flow.yaml.tmpl, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md, system-flow.yaml, design/system/server-performance.md, design/adrs/0266-a-test-tier-marked-parallel-starts-once-the-unmarked-tiers-before-it-have-ended.md, scripts/flai-build.sh, design/issues/I-0319-a-flaiover-props-panel-test-times-out-at-5-000-ms-under-load-and-fails-the.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 11123
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 4
      test_runs: 4
      hand_edits: 3
      work: 88
  models:
    - model: claude-opus-5-5
      input: 474
      output: 126118
      cache_read: 75126298
      cache_write: 6712391
      cost: 53.541
  strategic:
    - kind: orchestrator
      seconds: 99
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 11
          output: 98
          cache_read: 639041
          cache_write: 71055
          cost: 0.1823
cost_of_delay:
  inputs:
    penalty_per_week: 80
    by: orchestrator
    at: 2026-11-07T23:56:05Z
  value: 80
  by: planner-E-0033
  at: 2026-11-08T00:02:46Z
forecast:
  duration: 1h30m
  delivery: 2026-11-09T09:21:00Z
  basis: "Its own forecast of 1h30m; 1st in the pull order with an in-progress limit of 5, behind S-0603, S-0604 and S-0606."
  by: flai
  at: 2026-11-09T07:21:57Z
finalized:
  by: orchestrator
  at: 2026-11-08T00:10:30Z
---
# S-0610 The close-out runs independent test tiers at once and smoke skips what the close-out already ran, saving about 40 to 150 s a close-out

## Goal

Cause 20 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `flai verify` runs its tiers one after another and stops at the first failure (`flai/internal/verify/run.go`). The `flaiover` tier, which is node, waits for `integration`, which is Go, and `smoke` waits for both. `smoke` (`scripts/smoke.sh`) runs the template test, the check, and the markdown lint again after the close-out's own tiers ran them. Passed close-outs in the week to 2026-11-07 took p50 247 s, p90 510 s, max 1,408 s; `integration` mean 238 s, `flaiover` mean 176 s on 12 of 60 jobs, `smoke` mean 44 s.

## Acceptance criteria

- [x] `flai verify` runs at once the selected tiers the manifest marks as independent of each other, cheapest first as now, and cancels the others when one fails, so the verdict and the step it names are the same as a serial run's.
- [x] Under `flai verify`, `smoke` skips the steps a tier of the same run already passed.
- [x] A close-out of a story that selects `integration`, `flaiover`, and `smoke` takes at least 150 s less than its serial run on the same commit, and server-performance.md records the before and after.

## Tasks
- T-2634 Decide and add the manifest key that marks a test tier as able to run beside others
- T-2644 verify.Run starts the tiers the manifest marks independent at once and stops the later ones when one fails
- T-2653 Under flai verify, smoke skips the template test, the check, and the markdown lint a step of the same run already passed
- T-2661 Mark this repository's independent tiers, time a close-out serial and at once on one commit, and record it in server-performance.md

## Notes

- Criterion 3 is left unticked: it was not met on this host. Three rounds of `flai test --all` on one commit, serial (installed flai) against parallel (the tree's), saved 364 s, 8 s, and -14 s. At loads of 39 to 70, `flaiover` (vitest) and `integration` (`go test -race`) each fill the 24 cores, so beside each other they share the CPU rather than use idle CPU. Round 1, at the highest load, met it only because the serial run was slowed most. The figures and the reasons are in `design/system/server-performance.md` § Cause 20. The `parallel` marks are kept: they cost nothing on a host with idle cores and saved nothing here. Taking them off `system-flow.yaml` is a three-line edit.
- How a tier says it may run beside others is a manifest key, which `project-manifest.md` and `docs/operators/settings.md` would gain; that is likely an ADR.
- Splitting `internal/mcpserver`'s tests, 337 of 1,676 package-seconds in the 2026-11-02 batch, would shorten `integration` itself (I-0234).
- Turning the `merge-queue` host action on leaves `integration` and `smoke` to the batch; that is the operator's choice and needs no story.

### Planning

Touches, by where each came from:

| Touch | From |
|-------|------|
| `flai/internal/verify/run.go`, `flai/internal/verify/story.go`, `scripts/smoke.sh` | declared |
| `flai/internal/verify/run_test.go`, `flai/internal/verify/story_test.go`, `flai/internal/verify/verify.go`, `flai/internal/verify/manifest_test.go`, `design/system/flai-cli.md` | co-change (`flai touches suggest`: 19 to 48% of the 27 commits that changed the declared files) |
| `flai/internal/verify/manifest.go` | layout: `FromManifest` copies each manifest key onto `verify.Tier` |
| `flai/internal/manifest/manifest.go`, `flai/internal/manifest/manifest_test.go`, `flai/internal/manifest/settings.go`, `flai/internal/manifest/settings_test.go` | layout: `TestTier`, its validation, and its serialisation, for the key the Notes name |
| `docs/operators/settings.md`, `design/system/project-manifest.md`, `template/root/system-flow.yaml.tmpl` | design: the Notes name the first two; the template documents each `tests` key beside its example |
| `design/adrs/README.md` | design: the ADR the Notes expect; the ADR's own file is named when it is written, and `flai task done` widens T-2634's touches to it, so no folder touch is kept |
| `docs/users/flai.md` | design: `documentation.md` asks for the command's user documentation in the same story |
| `system-flow.yaml` | layout: this repository's `tests` gain the key for criterion 3 |
| `design/system/server-performance.md` | design: criterion 3 records the before and after there |

Left out of the suggestion: `flai/cmd/verify.go` and `flai/internal/verify/select.go`, which keep their shape, since tiers still come in `Select`'s order; `docs/users/flai-reference.md`, since no flag changes; `scripts/README.md` and `Makefile`, since no script is added.

No `after`: the story reads nothing another open story writes. S-0601 (`flai check` at once) shortens the check step, not the tiers, and changes none of these files.

Forecast: `flai forecast` gives 36m from a size of 23 (3 criteria, 20 touches). Raised to 1h30m: the four tasks run one after another, and T-2661 times two to four runs of `integration`, `flaiover`, and `smoke`, each 4 to 8 min on this host, which size does not count. Delivery moved by the same 54 min, from 03:09 to 04:03 on 2026-11-08.

Cost of delay: no value yet. Neither S-0610 nor E-0033 has inputs, so `flai cod` refuses. TH-0888 asks the operator for them; it recommends `penalty_per_week: 80` here, not yet confirmed.
