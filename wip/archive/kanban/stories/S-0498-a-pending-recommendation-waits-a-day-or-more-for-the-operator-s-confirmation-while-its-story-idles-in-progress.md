---
id: S-0498
type: story
nature: improvement
title: A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress
status: done
owner: alex
created: 2026-10-26T04:01:18Z
updated: 2026-10-26T16:01:44Z
transitions:
  - to: ready
    at: 2026-10-26T04:17:53Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-26T04:17:59Z
    by: agent-S-0498
  - to: review
    at: 2026-10-26T05:21:58Z
    by: agent-S-0498
  - to: done
    at: 2026-10-26T16:01:44Z
    by: alex
tags: [flai, template]
topics: [orchestration, cli]
touches: [design/adrs/README.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/check.go, docs/operators/settings.md, design/system/project-manifest.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/users/flai-reference.md, docs/users/flai.md, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/remind.go, flai/internal/serve/remind_test.go, flai/internal/serve/notify.go, flai/internal/serve/notify_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, docs/operators/index.md, design/system/workflow.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/issues/I-0151-a-pending-recommendation-waits-a-day-or-more-for-the-operator-s-confirmation-while-its-story-idles-in-progress.md, design/issues/summary.md, system-flow.yaml, design/adrs/0205-flai-serve-reminds-the-operator-of-a-thread-that-awaits-them-past-threads.md, template/template.yaml, flai/cmd/serve.go, design/issues/I-0210-flai-guard-lets-a-session-flai-serve-started-run-the-flai-serve-sub-commands-that-write-the-host-configuration.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3851
  priced_by: harness
  turns:
    - day: 2026-10-26
      ceremony: 2
      test_runs: 1
      hand_edits: 2
      work: 64
  models:
    - model: claude-opus-5-5
      input: 582
      output: 195901
      cache_read: 33480398
      cache_write: 998881
      cost: 16.3736
  strategic:
    - kind: planner
      seconds: 2
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 2
          output: 7
          cache_read: 14707
          cache_write: 10049
          cost: 0.0337
    - kind: orchestrator
      seconds: 881
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 105
          output: 1733
          cache_read: 20229122
          cache_write: 45683
          cost: 4.4454
        - model: claude-sonnet-5-5
          input: 22
          output: 88
          cache_read: 543920
          cache_write: 70162
          cost: 0.4199
    - kind: analyzer
      seconds: 114
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 98
          output: 6305
          cache_read: 370755
          cache_write: 41688
          cost: 1.2427
        - model: claude-haiku-5-5
          input: 8
          output: 6718
          cache_read: 204809
          cache_write: 23259
          cost: 0.0083
cost_of_delay:
  inputs:
    time_lost_per_cycle: 117h
    by: flai
    at: 2026-10-26T04:01:18Z
  value: 17550
  by: planner-S-0498
  at: 2026-10-26T04:16:50Z
forecast:
  duration: 1h45m
  delivery: 2026-10-26T19:25:00Z
  basis: "flai forecast's 54m (median 89 s per unit over 96 done improvement stories on claude-opus-5-5, size 36) raised, because the plan is seven tasks over five layers worked one after another, each closed with the serve package's tests, the slowest in the tree, and adds two serve files, a CLI command group with host configuration, a manifest block, and an ADR; delivery is flai's 18:31Z moved by the raise, 21st in the pull order behind the E-0023 and E-0024 drafts."
  by: planner-S-0498
  at: 2026-10-26T04:16:50Z
finalized:
  by: orchestrator
  at: 2026-10-26T04:17:45Z
---
# S-0498 A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress

## Goal

This story remediates [I-0151](../../../design/issues/I-0151-a-pending-recommendation-waits-a-day-or-more-for-the-operator-s-confirmation-while-its-story-idles-in-progress.md), "A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress", which the analysis in [2026-10-15-bottlenecks.md](../../../design/analysis/2026-10-15-bottlenecks.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0151 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0151 is closed with `flai issue close I-0151 --reason` saying what fixed it

## Tasks
- T-2127 ADR: flai serve reminds the operator of a thread that awaits them past a threshold through a channel on the host, and blocks a story whose agent ended asking
- T-2128 The manifest's threads.remind_after and threads.remind_every are declared, defaulted, validated by flai check, settable with flai manifest set, and documented
- T-2129 flai serve blocks a story whose agent ended asking, with the thread as the reason, and unblocks it when the answer starts the agent again
- T-2130 flai serve notify keeps a reminder channel, a webhook URL and a host command, in the host configuration, and a sender posts a reminder's JSON to it
- T-2131 flai serve reminds the operator of each thread awaiting them past threads.remind_after, again every threads.remind_every, through the host's channel, with its state in the serve folder and nothing written on the thread
- T-2132 The operator and user docs, the design, and the work-management convention describe the reminders, the notify channel, and the block flai serve writes on an asked story
- T-2133 This project sets threads.remind_after and threads.remind_every, the operator sets the host's channel, and I-0151 is closed naming what fixed it

## Notes

The analyzer found I-0151: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0151. time_lost_per_cycle 117h carried over from I-0151's Impact section. Its Impact time_lost_per_cycle was taken rather than the 19h20m derived from its cost and count.

Strategic usage carried over from I-0151: analyzer 1m54s, 1.2510 USD.

### Planning

Planned by planner-S-0498 on 2026-10-26. The remedy, its alternatives, and the tasks are on the plan's thread. Of I-0151's five instances, three waited on the operator's hands (a cosign key on TH-0324, an `rm` under `.git` on TH-0390, `flai host upgrade` on TH-0461), which no setting may confirm for them; S-0502 (I-0195) takes the design-choice waits with `orchestration.permissions.confirm_after`. This story reaches the operator sooner and marks the idle: `flai serve` reminds them of a thread that awaits them past `threads.remind_after`, through a channel kept on the host (`flai serve notify`), never by writing on the thread, which would end its wait (ADR-0091); and it blocks a story whose agent ended its session asking, with the thread as the reason, until the answer starts the agent again.

Touches. `flai touches suggest S-0498` listed nothing, because the story declared none. Every touch comes from the design and the code layout, read for the two remedies:

- design: `design/adrs/0204-...` (the new ADR; its number and slug are predicted, and `flai task done` widens the touch to the name `flai adr new` gives), `design/adrs/README.md`, `design/system/project-manifest.md` (the `threads` block), `design/system/workflow.md` (§ Blocking), `design/system/flai-cli.md` (flai serve), `design/system/strategic-agents.md` (§ Answering threads), `design/conventions/work-management.md` and its template copy (blocking), `docs/operators/settings.md` (one row per new key, which `TestSettingsIndexIsComplete` requires), `docs/operators/index.md`, `docs/users/flai.md`, `docs/users/flai-reference.md` (regenerated for `flai serve notify`), `template/root/system-flow.yaml.tmpl`, `template/CHANGELOG.md`.
- layout: `flai/internal/manifest/manifest.go` and `settings.go` with their tests (the duration settings, as `digest.compact_after` is declared), `flai/internal/check/check.go` (the manifest's `Errors` hook), `flai/internal/config/config.go` and its test (the host channel, as `providers` and `forges` are kept), `flai/cmd/serve_actions.go` and its test (`flai serve notify`), `flai/internal/serve/remind.go`, `notify.go`, their tests, `serve.go` and `serve_test.go` (the per-minute `look`, as `upgrade.go` runs it, with its state in the serve folder as `upgrade-asks.json` is), `flai/internal/serve/agents.go` and its test (`judgeRun` sets the asked outcome; `resume` starts the agent on the answer), `design/issues/I-0151-...md`, `design/issues/summary.md`, `system-flow.yaml` (this project's `threads` values, set with `flai manifest set`).
- No folder touch is kept: every file a task adds is named.

Forecast: `flai forecast S-0498` gave 3m before the touches (size 2) and 54m after them (median 89 s per unit over 96 done improvement stories on claude-opus-5-5, size 36: 2 criteria, 34 touches), delivery 2026-10-26T18:31Z, 21st in the pull order. Set to 1h45m, delivery 2026-10-26T19:25Z: the median comes mostly from small fixes, while this plan is seven tasks over five layers worked one after another, each closed with the serve package's tests, the slowest in the tree, and adds two serve files, a CLI command group with host configuration, a manifest block, and an ADR.

Cost of delay: `flai cod S-0498` gives 17550 USD a week, 117h of time lost per 168h cycle at 150 USD an hour, one cycle a week; kept as flai computed it. The input is I-0151's W42 confirmed-wait total, which S-0502 also remediates for its part, so the two stories share one loss; the input is the operator's and stands.
