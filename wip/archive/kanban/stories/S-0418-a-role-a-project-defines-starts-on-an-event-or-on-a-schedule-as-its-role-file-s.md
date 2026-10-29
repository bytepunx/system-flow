---
id: S-0418
type: story
nature: feature
title: A role a project defines starts on an event or on a schedule, as its role file's trigger says, coalesced and logged with its trigger as the planner's runs are
status: done
parent: E-0024
owner: alex
created: 2026-10-16T05:40:11Z
updated: 2026-10-29T18:11:13Z
transitions:
  - to: ready
    at: 2026-10-28T13:04:42Z
    by: alex
  - to: backlog
    at: 2026-10-28T22:20:39Z
    by: alex
  - to: ready
    at: 2026-10-29T14:42:10Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-29T14:42:16Z
    by: agent-S-0418
  - to: ready
    at: 2026-10-29T16:37:42Z
    by: agent-S-0418
  - to: in-progress
    at: 2026-10-29T16:46:12Z
    by: agent-S-0418
  - to: review
    at: 2026-10-29T17:57:15Z
    by: agent-S-0418
  - to: done
    at: 2026-10-29T18:11:13Z
    by: alex
blocked:
  - from: 2026-10-29T14:43:49Z
    until: 2026-10-29T15:17:18Z
    reason: awaiting the operator on TH-0749
  - from: 2026-10-29T15:17:46Z
    until: 2026-10-29T15:17:50Z
    reason: awaiting the operator on TH-0749
  - from: 2026-10-29T15:18:01Z
    until: 2026-10-29T15:24:32Z
    reason: awaiting the operator on TH-0749
  - from: 2026-10-29T15:24:45Z
    until: 2026-10-29T15:31:50Z
    reason: awaiting the operator on TH-0749
  - from: 2026-10-29T15:32:01Z
    until: 2026-10-29T15:38:13Z
    reason: awaiting the operator on TH-0749
  - from: 2026-10-29T15:38:26Z
    until: 2026-10-29T15:39:31Z
    reason: awaiting the operator on TH-0749
  - from: 2026-10-29T15:39:44Z
    until: 2026-10-29T15:42:15Z
    reason: awaiting the operator on TH-0749
  - from: 2026-10-29T15:42:43Z
    until: 2026-10-29T16:26:52Z
    reason: awaiting the operator on TH-0749
  - from: 2026-10-29T16:27:11Z
    until: 2026-10-29T16:35:26Z
    reason: awaiting the operator on TH-0749
tags: [cli, template]
topics: [cli, planning, orchestration]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/serve/activity.go, flai/internal/serve/agents.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/role_events.go, flai/internal/serve/role_events_test.go, flai/internal/serve/role.go, flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, template/CHANGELOG.md, design/issues/I-0232-flai-serve-wakes-a-story-s-agent-that-is-paused-on-a-thread-for-every-task-done-message-from-another-story-each-needing-only-a-settled-reply.md, design/issues/summary.md, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, design/issues/I-0233-testachangedfilereachesthedashboardasanotification-fails-its-temp-folder-s-cleanup-with-directory-not-empty-as-i-0048-s-serve-tests-did.md, design/system/project-manifest.md, docs/operators/runbooks/backup.md, flai/cmd/role.go]
after: [S-0417, S-0536]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4618
  priced_by: harness
  turns:
    - day: 2026-10-29
      ceremony: 8
      test_runs: 1
      hand_edits: 4
      work: 141
  models:
    - model: claude-opus-5-5
      input: 620
      output: 230620
      cache_read: 48151966
      cache_write: 1201422
      cost: 21.992
  strategic:
    - kind: orchestrator
      seconds: 455
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 173
          output: 2515
          cache_read: 40061178
          cache_write: 142871
          cost: 8.8701
        - model: claude-sonnet-5-5
          input: 12
          output: 60
          cache_read: 197880
          cache_write: 63993
          cost: 0.1804
cost_of_delay:
  value: 49.26
  by: planner-S-0418
  at: 2026-10-21T12:03:24Z
forecast:
  duration: 1h1m
  delivery: 2026-10-29T15:58:00Z
  basis: "Its own forecast of 1h1m; 1st in the pull order with an in-progress limit of 5, with nothing ahead of it."
  by: flai
  at: 2026-10-29T14:42:13Z
finalized:
  by: alex
  at: 2026-10-28T13:04:27Z
---
# S-0418 A role a project defines starts on an event or on a schedule, as its role file's trigger says, coalesced and logged with its trigger as the planner's runs are

## Goal

A role starts when its file says, not only when the operator asks. A role file's trigger names events or a schedule (ADR-0207 §3, ADR-0211). An event queues the role through the reason queue that S-0536 gives the orchestrator, coalesced into one run naming its reasons (ADR-0209 §2). A schedule starts it as `analysis.schedule` starts the analyzer. Each run records what started it. `continuous` is no trigger (ADR-0209 §1): S-0536's `flai check` reports it, and nothing here keeps a role running.

## Acceptance criteria

- [x] A role whose trigger names an event in ADR-0207's and ADR-0211's vocabulary, under any condition flai implements, is queued by `flai serve` through S-0536's reason queue when the event happens while its action is on: reasons coalesced into one run per role, reasons that arrive during a run starting the next, and a reason dropped with a logged reason when `role.run`'s checks refuse it.
- [x] A role whose trigger is a schedule setting, such as `roles.<name>.schedule`, a five-field cron expression in UTC or `daily` in the manifest, is started each time the schedule comes round while its action is on, through the `cadence` the replanner and the analysis scheduler share, and a start refused then is logged and not tried again until the schedule next comes round.
- [x] Each run records its trigger in `serve/agents.json` and in its activity entry's `- Trigger:` line, as `asked`, `<event> <ID>`, or `schedule <spec>`, joined by semicolons when coalesced.
- [x] `flai check` reports `roles.file` for a trigger it cannot read, an unknown event or condition, and reports a `roles.<name>.schedule` value that does not parse or never comes round.
- [x] `design/system/strategic-agents.md`, `docs/users/flai.md` with `flai-reference.md`, `docs/operators/index.md` with `settings.md`, and `template/CHANGELOG.md` say how each trigger starts a role.

## Tasks
- T-1741 A role file's trigger is read and checked: events, a schedule, or continuous, with roles.file findings for one flai cannot read
- T-1743 serve queues a role whose trigger names an event through the orchestrator's reason queue, coalesced into one run per role and item
- T-1744 A role's schedule is roles.<name>.schedule in the manifest, checked, and serve starts the role each time it comes round through cadence
- T-1746 serve keeps a continuous role running while its action is on, and flai role stop and start and the host API's role.stop and role.start hold and release it
- T-1748 RunRole and a role's start take the trigger that started it, and each run records it in agents.json and in its activity entry's Trigger line
- T-1751 The design and the user and operator docs say how each trigger starts a role
- T-2340 The serve loop looks at each defined role's event queue and schedule on every look of the project
- T-2343 flai serve waits for a served project's looks before it stops, so that no look writes after it, which failed the new loop test's cleanup

## Notes

- Drafted by the planner for E-0024. Waits for the story that starts a role on the operator's word, whose start path every trigger takes.
- The three built-in roles keep their own starters until the next story moves them onto this path.

### Re-scope (TH-0749, 2026-10-29)

- The operator confirmed option A on TH-0749 at 2026-10-29T16:35:26Z: the story predated ADR-0209 and ADR-0211.
- Dropped the criterion that kept a `continuous` role running with `flai role stop|start` and `role.stop|start`: ADR-0209 §1 rejects `continuous`, and S-0536 reports it.
- Cancelled T-1741 (S-0414 and S-0536 already read and check triggers) and T-1746 (`continuous`).
- The event trigger reuses S-0536's reason queue rather than a queue of its own. The schedule comes from a new `roles.<name>.schedule` manifest key.
- S-0418 waits for S-0536 (`after`), whose files it builds on.

### Planning

- **TH-0471's answers (2026-10-21).** Only answer 1 changes this story: it gives the cost of delay below. Answer 2's folder, `design/workflows/roles/`, is read through S-0414's loader, so no touch or task here names it. Answers 3 to 6 do not reach the triggers.
- **Forecast.** flai's 1h1m was set before the re-scope, for six criteria and 28 touches; the re-scope removes one criterion and two tasks.
- **Cost of delay.** 49.26 USD a week, worked out by hand from answer 1, as S-0414's and S-0417's planners did, because `flai cod` refuses while E-0024 has no inputs.
  - 4h lost per cycle at `planning.hour_rate` 150 over the default 168h cycle gives E-0024 600 USD a week.
  - S-0418 takes the share its 1h1m is of the 12h23m forecast for E-0024's twelve open stories: 61 / 743 × 600.
  - Setting 4h on E-0024 is asked on TH-0567, question 3. Once it is set, `flai cod S-0418` gives the same figure while the forecasts stand.
- 2026-10-29T16:37:42Z: moved to ready: TH-0749 option A: re-scoped to ADR-0209 and ADR-0211; waits for S-0536 (after)
