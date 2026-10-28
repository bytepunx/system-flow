---
id: S-0503
type: story
nature: improvement
title: "No test tier selects scripts/**, so a story that changes only shell scripts, and the script tests it adds, pass flai verify and the merge queue without any of them running"
status: done
owner: alex
created: 2026-10-26T04:01:45Z
updated: 2026-10-28T12:42:10Z
transitions:
  - to: ready
    at: 2026-10-26T05:26:43Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-26T05:26:59Z
    by: agent-S-0503
  - to: review
    at: 2026-10-28T12:03:36Z
    by: agent-S-0503
  - to: in-progress
    at: 2026-10-28T12:33:01Z
    by: flai
  - to: review
    at: 2026-10-28T12:40:26Z
    by: agent-S-0503
  - to: done
    at: 2026-10-28T12:42:10Z
    by: alex
tags: []
topics: [cli, code]
touches: [system-flow.yaml, scripts/install-tools.sh, scripts/shellcheck.sh, scripts/flai-test.sh, scripts/README.md, scripts/flaiover-install-test.sh, scripts/install-published-test.sh, scripts/litellm.sh, scripts/release-key-setup.sh, Makefile, design/tech/shellcheck.md, design/tech/README.md, design/conventions/tooling.md, design/system/project-manifest.md, flai/internal/verify/select_test.go, design/issues/I-0199-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md, design/issues/summary.md, scripts/env.sh, design/issues/I-0211-adr-0202-s-status-rejected-is-not-an-adr-status-flaiover-s-search-test-allows-so-the-flaiover-tier-fails-on-main-for-every-story-that-selects-it.md, design/issues/I-0220-testastorywhoseagentendedaskingisblockeduntiltheanswerstartsitagain-reads-the-story-s-blocks-once-the-run-reads-as-ended-before-endrun-writes-the-block-so-it-fails-now-and-then.md, flai/internal/serve/agents_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1915
  priced_by: harness
  turns:
    - day: 2026-10-26
      test_runs: 6
      hand_edits: 3
      work: 50
    - day: 2026-10-28
      ceremony: 5
      work: 38
  models:
    - model: claude-opus-5-5
      input: 340
      output: 91278
      cache_read: 18055263
      cache_write: 548942
      cost: 9.2457
  strategic:
    - kind: planner
      seconds: 29
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 156
          output: 6936
          cache_read: 904141
          cache_write: 110471
          cost: 1.2694
    - kind: orchestrator
      seconds: 14971
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 237
          output: 3748
          cache_read: 48987712
          cache_write: 72707
          cost: 10.7614
        - model: claude-sonnet-5-5
          input: 4
          output: 37
          cache_read: 13176
          cache_write: 26971
          cost: 0.0275
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-10-26T05:25:20Z
  value: 12.5
  by: planner-S-0503
  at: 2026-10-26T05:26:26Z
forecast:
  duration: 1h
  delivery: 2026-10-26T21:20:00Z
  basis: "flai forecast gave 29m (median 89 s per unit over 96 done improvement stories on claude-opus-5-5, size 19: 2 criteria, 17 touches); raised to 1h because the five tasks form five sequential layers with no parallel work, one downloads and checksums a release, one fixes seven lint findings across four scripts, and one adds a Go test that must fail and pass against the manifest; delivery is flai's 20:48 moved by the same 31m, 23rd in the pull order."
  by: planner-S-0503
  at: 2026-10-26T05:25:00Z
finalized:
  by: orchestrator
  at: 2026-10-26T05:26:35Z
---
# S-0503 No test tier selects scripts/**, so a story that changes only shell scripts, and the script tests it adds, pass flai verify and the merge queue without any of them running

## Goal

This story remediates [I-0199](../../../design/issues/I-0199-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md), "No test tier selects scripts/**, so a story that changes only shell scripts, and the script tests it adds, pass flai verify and the merge queue without any of them running". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0199 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0199 is closed with `flai issue close I-0199 --reason` saying what fixed it

## Tasks
- T-2179 scripts/install-tools.sh installs a pinned ShellCheck into bin/, recorded in design/tech
- T-2180 scripts/shellcheck.sh lints the scripts given, or every scripts/*.sh, and the existing scripts pass it
- T-2181 system-flow.yaml gains a shellcheck tier for scripts/**/*.sh, and its flaiover, template, and smoke tiers select the scripts they run
- T-2182 A test in flai/internal/verify asserts this repository's manifest runs shellcheck and the script tests for a change under scripts/
- T-2183 I-0199 is closed with the reason that fixed it, and the story's criteria are ticked
- T-2197 TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain waits for the block endRun writes after the end, so the merge queue's integration tier passes

## Notes

### Planning

Planned by planner-S-0503 on 2026-10-26; plan approved on TH-0661. The solution proposed from I-0199's one instance (S-0481): a `shellcheck` tier that every change under `scripts/` selects, and the `all_only` tiers that run scripts (`flaiover`, `template`, `smoke`) selecting the scripts they run, so `flai verify` and the merge queue's batch run the script tests for a scripts-only story (ADR-0174: the batch runs the tiers its paths select). Tasks in five layers, each waiting for the one before: T-2179, T-2180, T-2181, T-2182, then T-2183 (after T-2181 and T-2182).

Touches, file by file, no folder kept:

- Design and layout, from the issue's recommendation and the code layout: `system-flow.yaml` (the tiers), `scripts/shellcheck.sh` (new), `scripts/install-tools.sh` (the pinned tool), `scripts/flai-test.sh` (keeps in step with the tiers), `Makefile` (the targets), `design/tech/shellcheck.md` (new) and `design/tech/README.md` (the documentation convention: a dependency added is recorded), `design/conventions/tooling.md` (its project addition lists this repository's tiers), `flai/internal/verify/select_test.go` (where the manifest tests for I-0105 and I-0201 live).
- Co-change (`flai touches suggest` over the files above): `scripts/README.md` (20%), `design/system/project-manifest.md` (17%: its `tests` bullet names which paths this repository gives its `all_only` tiers). Left out: `docs/users/flai.md`, `flai/internal/verify/select.go`, `flai/internal/manifest/manifest.go`, `template/root/system-flow.yaml.tmpl` and `template/CHANGELOG.md`, since no command, no code, and no template manifest changes (the template's `test` tier selects `**`).
- Layout, from running ShellCheck 0.10.0 on the tree: `scripts/flaiover-install-test.sh`, `scripts/install-published-test.sh`, `scripts/litellm.sh`, `scripts/release-key-setup.sh` carry the seven findings T-2180 fixes.
- Declared by the criteria: the I-0199 file and `design/issues/summary.md`, which `flai issue close` rewrites.

Forecast: `flai forecast` gave 29m (median 89 s per unit over 96 done improvement stories on claude-opus-5-5, size 19). Set to 1h: the five layers run one after another with nothing in parallel, T-2179 downloads and checksums a release, T-2180 fixes lint in four scripts, and T-2182 must be shown failing and passing against the manifest. Delivery moved by the same 31m from flai's 20:48 to 21:20 UTC, 23rd in the pull order.

Cost of delay: the story had no inputs and no epic to take a share of; the orchestrator set `time_lost_per_cycle: 5m` on TH-0661, the figure recommended there. `flai cod` gives 12.50 USD a week (5m lost per 168h cycle at 150 USD an hour, one cycle a week), written as the value unchanged.

Assumptions: ShellCheck is pinned at v0.10.0, the host's version; `scripts/**/*.sh` selects `scripts/flai.sh` as `flai/**/*.go` selects `flai/main.go` (T-2181 checks); the template's scripts stay out of the default shellcheck set.
- 2026-10-28T12:33:01Z: moved to in-progress: merge queue batch 20261028-122130 fails with S-0503 at 947ecb0aed0a: integration failed: flai/internal/serve/agents_test.go:3123 TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain: the story idles in progress, blocked []
