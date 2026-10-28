---
id: S-0504
type: story
nature: remediation
title: The curator has never run on system-flow, so no thread is in any entry's asked_again and the digest's repeat measure reads 0 whatever was asked
status: done
owner: alex
created: 2026-10-26T04:01:51Z
updated: 2026-10-28T12:40:12Z
transitions:
  - to: ready
    at: 2026-10-26T04:11:21Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-26T04:11:31Z
    by: agent-S-0504
  - to: review
    at: 2026-10-28T12:37:42Z
    by: agent-S-0504
  - to: done
    at: 2026-10-28T12:40:12Z
    by: alex
blocked:
  - from: 2026-10-28T12:03:53Z
    until: 2026-10-28T12:13:40Z
    reason: awaiting the operator on TH-0649
tags: [flai, flaiover, digest, metrics]
topics: [digest, metrics]
touches: [design/adrs, design/adrs/README.md, design/system/metrics.md, flai/internal/metrics/digest.go, flai/internal/metrics/digest_test.go, flai/internal/metrics/metrics.go, flai/cmd/stats.go, flai/cmd/stats_test.go, flai/cmd/check_stats_test.go, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/explain.ts, flaiover/src/lib/viz/explain.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, design/issues/I-0207-the-curator-has-never-run-on-system-flow-so-no-thread-is-in-any-entry-s-asked-again-and-the-digest-s-repeat-measure-reads-0-whatever-was-asked.md, design/issues/summary.md, template/template.yaml, design/issues/I-0212-a-curator-run-that-ends-on-the-harness-s-usage-limit-is-recorded-failed-logs-no-activity-and-is-not-started-again-until-the-next-acceptance-or-schedule.md, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3801
  priced_by: harness
  turns:
    - day: 2026-10-26
      ceremony: 7
      hand_edits: 1
      work: 44
    - day: 2026-10-28
      ceremony: 4
      hand_edits: 2
      work: 42
  models:
    - model: claude-opus-5-5
      input: 352
      output: 102060
      cache_read: 26488336
      cache_write: 831196
      cost: 13.075
  strategic:
    - kind: planner
      seconds: 512
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 738
          output: 27869
          cache_read: 3010033
          cache_write: 245476
          cost: 6.1224
        - model: claude-haiku-5-5
          input: 134
          output: 20796
          cache_read: 6290295
          cache_write: 280859
          cost: 0.4182
    - kind: orchestrator
      seconds: 657
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 112
          output: 1813
          cache_read: 17768750
          cache_write: 37270
          cost: 3.9039
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-26T04:01:51Z
  value: 50
  by: planner-S-0504
  at: 2026-10-26T04:08:23Z
forecast:
  duration: 2h30m
  delivery: 2026-10-26T21:00:00Z
  basis: "flai forecast gave 38m (median 91 s per unit over 37 done remediation stories on claude-opus-5-5 in the large band, size 25: 2 criteria, 23 touches; delivery 2026-10-26T18:26Z, 27th in the pull order). Raised to 2h30m of agent time because the size band counts touches, not what they hold: this story writes an ADR, changes a flai stats section with its Go tests and printed text, changes the dashboard's Digest chart with its vitest and page tests, updates four documents and a convention with its template copy, and closes the issue. The first curator run depends on the operator enabling curate and on an acceptance or flai curate, which is wall-clock, not agent time, so delivery is set 2h40m after flai's: 2026-10-26T21:00Z, and later if the operator's answer or the first run waits."
  by: planner-S-0504
  at: 2026-10-26T04:08:23Z
finalized:
  by: alex
  at: 2026-10-26T04:02:38Z
---
# S-0504 The curator has never run on system-flow, so no thread is in any entry's asked_again and the digest's repeat measure reads 0 whatever was asked

## Goal

This story remediates [I-0207](../../../design/issues/I-0207-the-curator-has-never-run-on-system-flow-so-no-thread-is-in-any-entry-s-asked-again-and-the-digest-s-repeat-measure-reads-0-whatever-was-asked.md), "The curator has never run on system-flow, so no thread is in any entry's asked_again and the digest's repeat measure reads 0 whatever was asked". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0207 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0207 is closed with `flai issue close I-0207 --reason` saying what fixed it

## Tasks
- T-2120 Confirm I-0207's cause on the host and propose the remedy to the operator on a thread
- T-2121 An ADR and metrics.md define the curator's runs in the digest section and leave asked_again out of a week no run has measured
- T-2122 The curate host action is on for system-flow and the first curator run on it is recorded
- T-2123 flai stats reports the curator's runs in the digest section and leaves asked_again out of a week no run has measured, with a test that reproduces I-0207
- T-2124 The user docs, the analyzer's convention, and its template copy read the curator's runs beside the digest's repeats
- T-2125 The dashboard's Digest chart shows the curator's runs and marks a week whose repeats were not measured
- T-2126 I-0207 is closed with the reason that fixed it, once stats measures the repeats and the curator has run here

## Notes

Cost of delay inputs set by flai from I-0207. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-24T23:48:27Z, 1.2 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0504 on 2026-10-26. The cause has two parts, and the plan remedies both:

1. **The operational cause.** The `curate` host action was never turned on for system-flow (`design/experiments/S-0435-digest-priming-and-threads-to-the-operator.md` § What else changed between the windows), so `flai serve` never started the curator here and `wip/agents/curator.md` does not exist. Only the operator turns a host action on (`flai serve enable curate`). The story asks for it on a thread (T1) and records the first run (T6); no code can make it.
2. **The measure's defect.** `flai stats` reports `digest.weeks[].asked_again` as 0 whether the curator counted no repeat or never ran, because `metrics.digestSection` reads the threads and the entries alone and nothing in `digest` says whether a curator ran. The plan has it read the curator's activity document, which `metrics.Options.Activities` already loads from `wip/agents/curator.md`, and report the curator's runs per week beside the repeats, leaving `asked_again` out of a week no curator run has ended since, as `per_story` is left out of a week with no completed story. `metrics.md` changes only with an ADR, so T2 writes one, and the Digest chart and its docs follow (T4, T5). The test that reproduces I-0207 is a stats fixture with threads to the operator and entries listing none of them in `asked_again`, and no curator activity: today it reads `asked_again: 0`; after T3 it reads no `asked_again` and `curator_runs: 0`.

Touches, and where each came from:

- **Declared:** none. The story was made from the issue with no touches, so `flai touches suggest S-0504` refused without a path; it was run from `flai/internal/metrics/digest.go`, `flai/cmd/stats.go`, `flaiover/src/lib/viz/charts.ts`, and `design/system/metrics.md` (81 of 2422 commits).
- **Criteria:** `design/issues/I-0207-....md` and `design/issues/summary.md`, which `flai issue close` writes.
- **Layout:** `flai/internal/metrics/digest.go` and `digest_test.go` (`digestSection`, `DigestWeek`), `flai/internal/metrics/metrics.go` (the call that builds the section, which must pass `opt.Activities`), `flai/cmd/stats.go` and `stats_test.go` (`printDigest`), the Digest chart in `flaiover/src/lib/viz/charts.ts`, `charts.test.ts`, `explain.ts`, `explain.test.ts`, and `flaiover/src/routes/charts/[kind]/+page.svelte` (its table).
- **Co-change:** `flai/cmd/check_stats_test.go` (19%, asserts the digest's printed text), `flaiover/src/routes/charts/[kind]/charts.svelte.test.ts` (17%, the page's digest notes), `design/adrs/README.md` (23%, the ADR index), `docs/users/flai.md` (22%, the stats sample and "Running the curator"), `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` (14% and 12%, the Digest chart's rows). Left out: `flai/internal/statsread/statsread.go` (7%; it already loads the activity documents), `docs/users/flai-reference.md` and `design/system/flai-cli.md` (no flag or command changes), `design/tech/charts.md` (no chart library change).
- **Design:** `design/system/metrics.md` § Digest (the measure's contract), `design/conventions/strategic-agents.md` § As the analyzer and its template copy `template/root/design/conventions/strategic-agents.md` (the analyzer reads the curator's runs beside `asked_again`), with `template/CHANGELOG.md` for the template change.
- **Folder kept:** `design/adrs`, because the ADR's file name is not known until `flai adr new` allocates its number; it is a shared path, so it holds nothing (ADR-0096). Every task names files only.

Figures:

- **Forecast:** `flai forecast` gave 5m from size 2 (2 criteria, 0 touches) before the touches were written, and is run again with them below. See the forecast's basis for what stands.
- **Cost of delay:** `flai cod` gives 50.00 USD a week from the operator's input, 20m of time lost per 168h cycle at 150 USD an hour, one cycle a week. It stands: the inputs are flai's from the issue's one instance and are the operator's to change. The figure does not count that S-0435's re-run from 2026-11-09 cannot measure repeats until the curator has run; that is an argument for the pull order, not a change to the inputs.
- **After:** none. The story needs no other story's change on main.
