---
id: S-0418
type: story
nature: feature
title: A role a project defines starts on an event, on a schedule, or for as long as its host action is on, as its role file's trigger says, coalesced and logged with its trigger as the planner's runs are
status: ready
parent: E-0024
owner: alex
created: 2026-10-16T05:40:11Z
updated: 2026-10-28T15:18:40Z
transitions:
  - to: ready
    at: 2026-10-28T13:04:42Z
    by: alex
tags: [cli, template]
topics: [cli, planning, orchestration]
touches: [design/system/strategic-agents.md, design/system/dashboard-host-channel.md, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/role.go, flai/cmd/role_test.go, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/serve/activity.go, flai/internal/serve/agents.go, flai/internal/serve/role_continuous.go, flai/internal/serve/role_continuous_test.go, flai/internal/serve/role_events.go, flai/internal/serve/role_events_test.go, flai/internal/serve/role.go, flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flaiover/src/lib/server/agent.ts, template/CHANGELOG.md, design/system/flai-cli.md]
after: [S-0417]
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
      seconds: 230
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 76
          output: 1156
          cache_read: 12961309
          cache_write: 21887
          cost: 2.8667
cost_of_delay:
  value: 49.26
  by: planner-S-0418
  at: 2026-10-21T12:03:24Z
forecast:
  duration: 1h1m
  delivery: 2026-10-28T16:31:00Z
  basis: "Its own forecast of 1h1m; 3rd in the pull order with an in-progress limit of 5, behind S-0415, S-0502, S-0419 and S-0423."
  by: flai
  at: 2026-10-28T15:18:40Z
finalized:
  by: alex
  at: 2026-10-28T13:04:27Z
---
# S-0418 A role a project defines starts on an event, on a schedule, or for as long as its host action is on, as its role file's trigger says, coalesced and logged with its trigger as the planner's runs are

## Goal

A role starts when its file says, not only when the operator asks. A role file's trigger names an event, a schedule, or `continuous`: an event queues the role once per item as the replanner queues the planner; a schedule starts it as `analysis.schedule` starts the analyzer; `continuous` keeps it running while its host action is on, started again when it ends and stopped when the action is turned off, as the orchestrator is. Each run records what started it.

## Acceptance criteria

- [ ] A role whose trigger names events (a story accepted, cancelled, moved to ready or to review, an item edited, a thread opened) is queued by `flai serve` when one happens while its action is on: once per role and item, a second trigger added to the entry, drained one run at a time per role, and dropped with a logged reason when `role.run`'s checks refuse it.
- [ ] A role whose trigger is a schedule, a five-field cron expression in UTC or `daily`, is started each time the schedule comes round while its action is on, through the cadence the replanner and the analysis scheduler share, and a start refused then is logged and not tried again until the schedule next comes round.
- [ ] A role whose trigger is `continuous` is started while its action is on, started again when its run ends and no sooner than a minute after a failed one, and stopped when the action is turned off; `flai role stop <name>` and `flai role start <name>`, and the host API's `role.stop` and `role.start`, hold and release it as the orchestrator's Stop and Start do.
- [ ] Each run records its trigger in `serve/agents.json` and in its activity entry's `- Trigger:` line, as `asked`, `<event> <ID>`, `schedule <spec>`, or `continuous`, joined by semicolons when coalesced.
- [ ] `flai check` reports `roles.file` for a trigger it cannot read: an unknown event, a schedule that does not parse or never comes round, or two triggers that cannot be combined.
- [ ] `design/system/strategic-agents.md`, `docs/users/flai.md` with `flai-reference.md`, `docs/operators/index.md`, and `template/CHANGELOG.md` say how each trigger starts a role.

## Tasks
- T-1741 A role file's trigger is read and checked: events, a schedule, or continuous, with roles.file findings for one flai cannot read
- T-1743 serve queues a role whose trigger names an event when that event happens, once per role and item, and drains the queue one run at a time per role
- T-1744 serve starts a role whose trigger is a schedule each time it comes round, through the cadence the analysis scheduler uses
- T-1746 serve keeps a continuous role running while its action is on, and flai role stop and start and the host API's role.stop and role.start hold and release it
- T-1748 The serve loop runs each role's trigger at every look, and each run records its trigger in agents.json and in its activity entry's Trigger line
- T-1751 The design and the user and operator docs say how each trigger starts a role

## Notes

- Drafted by the planner for E-0024. Waits for the story that starts a role on the operator's word, whose start path every trigger takes.
- The three built-in roles keep their own starters until the next story moves them onto this path.

### Planning

- **TH-0471's answers (2026-10-21).** Only answer 1 changes this story: it gives the cost of delay below. Answer 2's folder, `design/workflows/roles/`, is read through S-0414's loader, so no touch or task here names it. Answers 3 to 6 do not reach the triggers.
- **Touches.** 28, all files; no folder touch. The story declared none of its own; every touch comes from its six tasks.
  - **From layout, kept:** the trigger in `flai/internal/roles` and its check; the three new serve files (`role_events.go`, `role_schedule.go`, `role_continuous.go`) with their tests; `serve.go`, `role.go`, and `activity.go` for the loop and the `- Trigger:` line; `flai/cmd/role.go` and the host API's `writes.go` for the hold.
  - **From layout, added (code read 2026-10-21):**
    - `flai/internal/serve/agents.go` on T-1746 and T-1748. Its `Held` field is set on the orchestrator's run alone, and its `Trigger` field says it is a planner's or an analyzer's. A continuous role's hold and a defined role's trigger both land there. T-1748 waits for T-1746, so the two never share a layer.
    - `flaiover/src/lib/server/agent.ts` on T-1746. `TestTheDashboardAsksForExactlyWhatFlaiOffers` fails when `role.stop` and `role.start` are offered but missing from `REQUIRED_METHODS`.
  - **From layout, dropped:** `flai/internal/serve/analysis_schedule.go` from T-1744. The shared `cadence` lives in `replan.go`, in the same package, so a new file uses it without moving it. T-1744's Work moves it only "if it lives there". Dropping it stops this story holding S-0419 and S-0424, which change that file.
  - **From design:** `strategic-agents.md`, the user and operator docs, the reference, and the changelog. Added `design/system/dashboard-host-channel.md` on T-1751: it lists the host API's methods, `orchestrate.stop` and `start` among them, and `role.stop` and `role.start` join it.
  - **Co-change, kept:** `design/system/flai-cli.md` (shared; holds nothing).
  - **Co-change, not added:** `docs/operators/settings.md`, since a trigger lives in the role file and adds no manifest key. `flai/cmd/serve_actions.go`, since `flai role stop` and `start` live in `flai/cmd/role.go`. `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`, since the role page that shows a trigger is S-0420's.
- **Forecast.** flai's 1h1m is kept: 107 s a unit times size 34, six criteria and 28 touches. Three of the touches added are small (a field comment, two method names, a design line). The hold that `agents.go` now carries for any continuous role is real work that offsets that. The delivery is flai's play-out, 2026-10-21T20:47Z.
- **Cost of delay.** 49.26 USD a week, worked out by hand from answer 1, as S-0414's and S-0417's planners did, because `flai cod` refuses while E-0024 has no inputs.
  - 4h lost per cycle at `planning.hour_rate` 150 over the default 168h cycle gives E-0024 600 USD a week.
  - S-0418 takes the share its 1h1m is of the 12h23m forecast for E-0024's twelve open stories: 61 / 743 × 600.
  - Setting 4h on E-0024 is asked on TH-0567, question 3. Once it is set, `flai cod S-0418` gives the same figure while the forecasts stand.
