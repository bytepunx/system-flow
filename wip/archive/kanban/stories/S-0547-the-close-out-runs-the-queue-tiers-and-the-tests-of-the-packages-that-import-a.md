---
id: S-0547
type: story
nature: improvement
title: The close-out runs the queue tiers and the tests of the packages that import a changed package, so what the merge queue would send back is found before review
status: done
parent: E-0031
owner: alex
created: 2026-10-30T01:52:02Z
updated: 2026-11-01T23:56:43Z
transitions:
  - to: ready
    at: 2026-10-31T07:24:13Z
    by: alex
  - to: in-progress
    at: 2026-10-31T08:40:01Z
    by: agent-S-0547
  - to: review
    at: 2026-11-01T23:41:47Z
    by: agent-S-0547
  - to: done
    at: 2026-11-01T23:56:43Z
    by: alex
blocked:
  - from: 2026-10-31T10:12:10Z
    until: 2026-10-31T10:12:16Z
    reason: awaiting the operator on TH-0789
  - from: 2026-10-31T10:12:35Z
    until: 2026-10-31T10:18:06Z
    reason: awaiting the operator on TH-0789
  - from: 2026-10-31T10:18:27Z
    until: 2026-10-31T10:20:48Z
    reason: awaiting the operator on TH-0789
  - from: 2026-10-31T10:21:07Z
    until: 2026-10-31T10:35:11Z
    reason: awaiting the operator on TH-0789
  - from: 2026-10-31T10:35:51Z
    until: 2026-10-31T10:40:59Z
    reason: awaiting the operator on TH-0789
  - from: 2026-10-31T10:41:16Z
    until: 2026-11-01T05:22:59Z
    reason: awaiting the operator on TH-0789
  - from: 2026-11-01T05:23:25Z
    until: 2026-11-01T06:04:34Z
    reason: awaiting the operator on TH-0789
  - from: 2026-11-01T06:05:00Z
    until: 2026-11-01T23:29:54Z
    reason: awaiting the operator on TH-0789
tags: [flai]
topics: [testing, cli]
touches: [flai/internal/check, flai/internal/verify, flai/cmd, system-flow.yaml, design/system/flai-cli.md, design/conventions/testing.md, docs/users/flai.md, design/adrs, design/adrs/README.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, docs/operators/settings.md, docs/users/flai-reference.md, design/system/project-manifest.md, design/conventions/code-quality.md, design/conventions/tooling.md, design/issues/I-0221-the-close-out-s-go-test-tier-runs-only-the-changed-packages-so-a-change-that-breaks-a-caller-package-s-tests-passes-before-review.md, design/issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md, design/issues/summary.md, design/issues/I-0226-two-flai-serve-tests-check-before-what-they-wait-for-has-happened-and-fail-a-loaded-merge-queue-batch-on-timing.md, design/issues/I-0237-the-smoke-tier-s-openrouter-gateway-check-fails-on-a-402-because-the-operator-s.md, design/issues/I-0235-flai-check-finds-narrative-missing-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 6325
  priced_by: harness
  turns:
    - day: 2026-10-31
      ceremony: 4
      hand_edits: 3
      work: 69
    - day: 2026-11-01
      ceremony: 8
      test_runs: 4
      hand_edits: 2
      work: 28
  models:
    - model: claude-opus-5-5
      input: 600
      output: 160718
      cache_read: 53647838
      cache_write: 2097301
      cost: 28.4691
  strategic:
    - kind: orchestrator
      seconds: 143
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 59
          output: 405
          cache_read: 1648903
          cache_write: 305280
          cost: 0.4453
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3h
    by: alex
    at: 2026-10-31T09:08:43Z
  value: 450
  by: planner-S-0547
  at: 2026-10-31T09:09:21Z
forecast:
  duration: 4h
  delivery: 2026-10-31T13:44:00Z
  basis: "Its own forecast of 4h; 2nd in the pull order with an in-progress limit of 5, behind S-0548, S-0550, S-0552 and S-0546."
  by: flai
  at: 2026-10-31T08:20:33Z
---
# S-0547 The close-out runs the queue tiers and the tests of the packages that import a changed package, so what the merge queue would send back is found before review

## Goal

The close-out's go-test tier runs only the changed packages, and ADR-0174 leaves the `queue: true` tiers, integration and smoke, to the merge-queue batch. In the fortnight 24 of 104 queue verifications failed and 14 stories were sent back, 22 percent of landings, each costing the agent a restart, a fix, another close-out, a re-enqueue, and another 10-minute verification. I-0221 records a change that broke a caller package's tests and passed its close-out; I-0203 a stale generated file that failed every later batch at integration. This story makes the close-out find what the queue finds: the go-test tier runs the changed packages and every package that imports them (`go list -deps` in reverse), and the queue tiers run in the close-out too, in the worktree, when the story's touches reach what they test, with the batch keeping them for the landing check. With the close-out as a serve job, their time no longer hits a tool timeout.

## Acceptance criteria
- [x] The go-test tier's `{packages}` expands to the changed packages and the packages that import them, and a test shows a change that breaks a caller's tests fails the close-out
- [x] A tier with `queue: true` runs in the close-out when a `when:` the manifest gives it matches the story's touches, and in the batch as today; the manifest schema, `flai check`, and the reference describe `when`
- [x] system-flow.yaml gives integration and smoke a `when` that reaches them from the flai and template paths, and I-0221 and I-0203 are closed by the story that proves it
- [x] The testing convention, `flai-cli.md`, and the user docs say which tiers the close-out runs and which the batch repeats

## Tasks

The planner drafts them.
- T-2375 An ADR refining ADR-0174 and ADR-0221 says a queue tier whose when matches the story runs in the close-out and again in the batch
- T-2376 A test tier may carry when, path patterns flai check validates and allows only beside queue: true, read into verify's Tier and described in the settings index
- T-2377 The go-test tier runs the changed packages and every package that imports them, and a test shows a change that breaks a caller's tests fails verify
- T-2378 While the queue can batch, flai verify and the MCP verify tool run a queue tier whose when matches the story's touches or changed paths, and leave the rest to the batch
- T-2379 system-flow.yaml gives integration and smoke a when on the flai and template paths, a close-out shows both run, and I-0221 is closed
- T-2380 The testing convention, flai-cli.md, the manifest design, and the user docs say which tiers the close-out runs and which the batch repeats

## Notes

Evidence: `flai stats` merge-queue figures for 2026-10-15 to 2026-10-29 and the serve journal's send-backs, in the operator's conversation on 2026-10-29.

### Proof run

The close-out at 9794966, 2026-10-31, with the merge queue able to batch:

| Tier | Result |
|------|--------|
| go-test | skipped, covered by integration |
| integration | passed, 9m41.5s, "ran here, its when flai/** matches flai/internal/check/tests_test.go" |
| smoke | failed, 2m55.2s, "ran here, its when flai/** matches flai/internal/check/tests_test.go" |

Smoke failed only on `TestGatewayOpenRouter`. The OpenRouter key is out of credits and answered 402 (I-0237, TH-0789).

The close-out at 3b8598b, 2026-11-01, after the rebase onto S-0548, S-0546, S-0550, and S-0551, passed every step. It ran without the OpenRouter key, as the operator chose on TH-0789:

| Tier | Result |
|------|--------|
| go-test | skipped, covered by integration |
| integration | passed, 3m9.7s, "ran here, its when flai/** matches flai/internal/check/tests_test.go" |
| smoke | passed, 33.7s, "ran here, its when flai/** matches flai/internal/check/tests_test.go" |

### Planning

Planned by planner-S-0547 on 2026-10-31. The tasks form three layers: T-2375, T-2376, and T-2377 wait for none; T-2378 waits for T-2376; T-2379 waits for T-2377 and T-2378, and T-2380 for T-2375 and T-2378.

Touches:

- Declared, all kept: `flai/internal/check`, `flai/internal/verify`, `flai/cmd`, `system-flow.yaml`, `design/system/flai-cli.md`, `design/conventions/testing.md`, and `docs/users/flai.md`. The three folders `flai/internal/check`, `flai/internal/verify`, and `flai/cmd` are the operator's; the tasks name the files in them (`check/tests_test.go`; `verify/verify.go`, `manifest.go`, `manifest_test.go`, `select.go`, `select_test.go`, `story.go`, `story_test.go`, `importers_test.go`; `cmd/verify.go`, `verify_test.go`), and flai replaces each folder in the claim with those files (ADR-0096). `design/conventions/testing.md` does not exist; the testing tiers are in the project additions of `code-quality.md` and `tooling.md`, which T-2380 changes, as S-0548's plan found.
- Folder kept: `design/adrs`, because T-2375 adds an ADR whose number and slug `flai adr new` gives when it is written, so no task can name the file yet; T-2375 names `design/adrs/README.md`, which narrows the folder in the claim.
- Design: `design/adrs/README.md` (ADR-0174 and ADR-0221 are refined); `design/system/project-manifest.md`, which describes the tier keys and this repository's tiers; `docs/operators/settings.md`, whose test fails on a manifest key without a row; `docs/users/flai-reference.md`, regenerated when `flai verify`'s help changes (I-0203).
- Layout: `flai/internal/manifest/manifest.go`, `manifest_test.go`, `settings.go`, `settings_test.go`, where the tier keys and `Manifest.TestErrors` live; `flai/internal/mcpserver/verify.go` and `verify_test.go`, the MCP `verify` tool's second `LeaveToQueue` call site; the I-0221 and I-0203 issue files and `design/issues/summary.md`, which criterion 3 closes.
- Co-change: `flai touches suggest` ranked `docs/users/flai-reference.md` (24%) and `docs/operators/settings.md` (15%) first, both taken; `docs/operators/index.md`, `flaiover-dashboard.md`, `flaiover.md`, `hostapi`, `mcpserver/server.go`, and the rest no criterion reaches, so they are left out.

Forecast: `flai forecast` gave 18m, 96 s per unit of size over 22 improvement stories at size 11 (4 criteria, 7 touches). Raised to 4h: the size counts neither the six tasks, nor the selection change through `flai verify` and the MCP tool with their reports, nor T-2377's real two-package `go test` fixture, nor T-2379's proving run of integration and smoke, about twelve minutes a round. The delivery I first wrote, 22:00 UTC, assumed the story would wait for S-0548 and S-0546 to reach review. It started at 08:40 on a share instead, and flai's replan now gives the delivery.

Cost of delay: 450 USD a week, from the operator's input `time_lost_per_cycle: 3h` (set 2026-10-31T09:08:43Z, as TH-0774 recommended), × `planning.hour_rate` 150 × one 168h cycle a week, as `flai cod S-0547` works it out. It stands unadjusted: the input is the 7 send-backs a week of about 25 minutes each that the goal's figures give.

After: none. S-0547 needs no other story's accepted change on main. S-0548 changes the same `go-test` line in `system-flow.yaml` (`-timeout`, `-short`) and `flai/internal/check/tests_test.go`; T-2377 and T-2376 say to keep its change.
