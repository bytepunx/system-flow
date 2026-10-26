---
id: S-0502
type: story
nature: improvement
title: A thread with a pending orchestrator recommendation waits on the operator however long it takes, with no setting to confirm it after a threshold
status: backlog
owner: alex
created: 2026-10-26T04:01:40Z
updated: 2026-10-26T05:14:27Z
transitions: []
tags: [flai, template]
topics: [orchestration, cli, analysis]
touches: [design/adrs/0206-flai-serve-confirms-a-recommendation-the-operator-left-unanswered-past-confirm-after.md, design/adrs/README.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/orchestration_test.go, docs/operators/settings.md, design/system/project-manifest.md, flai/internal/threads/threads.go, flai/internal/threads/threads_test.go, flai/cmd/thread.go, flai/cmd/thread_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, flai/internal/serve/confirm.go, flai/internal/serve/confirm_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/internal/metrics/waiting.go, flai/internal/metrics/waiting_test.go, design/system/metrics.md, flaiover/src/lib/viz/charts.ts, design/system/strategic-agents.md, design/system/flai-cli.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, template/CHANGELOG.md, system-flow.yaml, design/issues/I-0195-a-thread-with-a-pending-orchestrator-recommendation-waits-on-the-operator-however-long-it-takes-with-no-setting-to-confirm-it-after-a-threshold.md, design/issues/summary.md]
after: [S-0498]
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
      seconds: 999
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 131
          output: 2168
          cache_read: 33785734
          cache_write: 50997
          cost: 7.4183
cost_of_delay:
  inputs:
    time_lost_per_cycle: 60h
    by: orchestrator
    at: 2026-10-26T04:40:45Z
  value: 9000
  by: planner-S-0502
  at: 2026-10-26T04:42:28Z
forecast:
  duration: 2h
  delivery: 2026-10-26T21:05:00Z
  basis: "Its own forecast of 2h; 22nd in the pull order with an in-progress limit of 5, behind S-0498, S-0504, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425 and S-0499."
  by: flai
  at: 2026-10-26T04:57:43Z
finalized:
  by: orchestrator
  at: 2026-10-26T05:14:27Z
---
# S-0502 A thread with a pending orchestrator recommendation waits on the operator however long it takes, with no setting to confirm it after a threshold

## Goal

This story remediates [I-0195](../../../design/issues/I-0195-a-thread-with-a-pending-orchestrator-recommendation-waits-on-the-operator-however-long-it-takes-with-no-setting-to-confirm-it-after-a-threshold.md), "A thread with a pending orchestrator recommendation waits on the operator however long it takes, with no setting to confirm it after a threshold". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0195 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0195 is closed with `flai issue close I-0195 --reason` saying what fixed it

## Tasks
- T-2141 ADR: flai serve confirms an orchestrator recommendation the operator left unanswered past orchestration.permissions.confirm_after, never one marked for the operator alone
- T-2142 The manifest's orchestration.permissions.confirm_after is declared as a duration, off by default, validated by flai check, settable with flai manifest set, shown in the settings catalog, and indexed in the operator docs
- T-2143 The threads package marks a recommendation for the operator alone, tells whether the owner wrote since, and writes flai serve's confirmation by timeout, with --operator-only on flai thread reply and operator_only on thread_reply
- T-2144 flai stats tells a wait that flai serve's confirmation by timeout ended from one the operator confirmed, as waiting.weeks[].threads.timeout and awaited flai serve, by the ADR's change to metrics.md, and the dashboard's types follow
- T-2145 The orchestrator's prompt and its agent definition tell it to mark a recommendation on money or a host change for the operator alone, and that flai serve confirms the rest after confirm_after
- T-2146 flai serve confirms, on its per-minute look while orchestrate is on, each orchestrator recommendation that waited confirm_after with no owner entry since, on a thread an agent opened and not marked for the operator alone, and logs it
- T-2147 The design, the strategic-agents convention in both copies, the user docs, the command reference, and the changelog describe confirm_after, the operator-alone mark, and flai serve's confirmation by timeout
- T-2148 This project sets orchestration.permissions.confirm_after, the acceptance criteria are verified and ticked, and I-0195 is closed naming what fixed it

## Notes

### Planning

Planned by planner-S-0502 on 2026-10-26. The remedy is the one the operator confirmed on TH-0598, which I-0195 records: `orchestration.permissions.confirm_after`, a duration, off by default; `flai serve` confirms an orchestrator recommendation that has waited that long with no entry from the operator, marked as a confirmation by timeout in the thread, in `wip/agents/orchestrator.md`, and in `flai stats` beside the operator's confirmations (ADR-0091); never for money, host changes, permission prompts, `.claude/` writes, or a thread the operator wrote on since. What flai can tell (the opener, the permission prompt's title, the owner's entries) it excludes itself; what only the orchestrator can tell (money, a host change) the orchestrator marks on the recommendation as the operator's alone. The tasks and the assumptions are on the plan's thread.

Touches. `flai touches suggest S-0502` listed nothing, because the story declared none. Every touch comes from the design and the code layout, read for the remedy:

- design: `design/adrs/0206-...` (the new ADR; 0204 is the next free number on main, S-0504 may take it, and S-0498 has taken 0205 on its branch, so the name is predicted and `flai task done` widens the touch to the name `flai adr new` gives), `design/adrs/README.md`, `design/system/project-manifest.md` (the permission's row), `design/system/strategic-agents.md` (§ Answering threads, § Its permissions and the guard), `design/system/metrics.md` (§ Waiting, by the ADR), `design/system/flai-cli.md` (the reply flag), `design/conventions/strategic-agents.md` and its template copy, `.claude/agents/orchestrator.md` and its template copy, `docs/operators/settings.md` (the row `TestSettingsIndexIsComplete` requires), `docs/users/flai.md`, `docs/users/flaiover.md` (the permissions table), `docs/users/flai-reference.md` (regenerated), `template/CHANGELOG.md`.
- layout: `flai/internal/manifest/manifest.go` and `settings.go` with their tests (the permission, as `answer_threads` is declared and `digest.compact_after` is a duration), `flai/internal/check/orchestration_test.go`, `flai/internal/threads/threads.go` and its test (the mark, `WroteSince`, the confirmation by timeout beside `Confirm`), `flai/cmd/thread.go` and its test (`--operator-only`), `flai/internal/mcpserver/server.go` and its test (`thread_reply`'s `operator_only`), `flai/internal/harness/harness.go` and its test (`orchestratePrompt`), `flai/internal/serve/confirm.go`, `confirm_test.go`, `serve.go`, and `serve_test.go` (the per-minute look, as `orchestrate.go`, `analysis_schedule.go`, and `curate_schedule.go` run theirs), `flai/internal/metrics/waiting.go` and its test (the `timeout` ender beside `byOrchestrator` and `byConfirmation`), `flaiover/src/lib/viz/charts.ts` (the `WaitWeek` type), `design/issues/I-0195-...md`, `design/issues/summary.md`, `system-flow.yaml` (this project's value, set with `flai manifest set`).
- No folder touch is kept: every file a task adds is named.

After: S-0498. Its T-2128 declares the manifest's `threads` block in the same manifest, settings, and check files and the same docs rows, its T-2131 adds `remind.go`'s look to the same per-minute block of `serve.go`, and its T-2132 rewrites the same § Answering threads and the convention; the confirmer scans the threads awaiting the operator past a threshold as the reminder does, so it builds on S-0498's change rather than meeting it on a branch.

Forecast: `flai forecast S-0502` gave 3m before the touches (size 2) and 1h after them (median 89 s per unit over 96 done improvement stories on claude-opus-5-5, size 40: 2 criteria, 38 touches), delivery 2026-10-26T19:37Z, 22nd in the pull order. Set to 2h, delivery 2026-10-26T20:37Z: the median comes mostly from small fixes, while this plan is eight tasks over five layers worked one after another, each closed with the serve, threads, manifest, or metrics package tests, and adds a manifest permission of a new kind, a thread entry mark across three surfaces, a serve look, a metrics field by ADR, and the ADR itself; and the story is held in ready until S-0498 is accepted.

Cost of delay: the story had no inputs and no epic, so I asked on TH-0656 with 60h recommended; the orchestrator set `time_lost_per_cycle` 60h under `plan_backlog_stories`. `flai cod S-0502` gives 9000 USD a week, 60h of time lost per 168h cycle at 150 USD an hour, one cycle a week; kept as flai computed it. The figure is two design-choice recommendations a week left pending, as TH-0583 and TH-0588 are now, at I-0151's mean confirmed wait of about 29 hours each; S-0498 carries I-0151's own 117h, so the two stories do not count the same loss twice.
