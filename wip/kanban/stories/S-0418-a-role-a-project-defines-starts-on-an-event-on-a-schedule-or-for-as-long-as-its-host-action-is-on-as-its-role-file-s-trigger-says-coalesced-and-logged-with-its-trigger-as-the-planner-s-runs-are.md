---
id: S-0418
type: story
nature: feature
title: A role a project defines starts on an event, on a schedule, or for as long as its host action is on, as its role file's trigger says, coalesced and logged with its trigger as the planner's runs are
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:40:11Z
updated: 2026-10-20T18:53:34Z
transitions: []
tags: [cli, template]
topics: [cli, planning, orchestration]
touches: [design/system/strategic-agents.md, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/role.go, flai/cmd/role_test.go, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/serve/activity.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/role_continuous.go, flai/internal/serve/role_continuous_test.go, flai/internal/serve/role_events.go, flai/internal/serve/role_events_test.go, flai/internal/serve/role.go, flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, template/CHANGELOG.md, design/system/flai-cli.md]
after: [S-0417]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 51m
  delivery: 2026-10-21T02:22:00Z
  basis: "Its own forecast of 51m; 13th in the pull order with an in-progress limit of 5, behind S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416 and S-0417."
  by: flai
  at: 2026-10-20T18:53:34Z
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

- **Touches.** The story declared none. Every touch comes from its six tasks.
  - **From layout:** the trigger in `flai/internal/roles`, the check, three new serve files (`role_events.go`, `role_schedule.go`, `role_continuous.go`) with their tests, `analysis_schedule.go` for the shared cadence, `serve.go`, `role.go`, and `activity.go` for the loop and the `- Trigger:` line, `flai/cmd/role.go`, and the host API's `writes.go`.
  - **From design:** `strategic-agents.md`, the user and operator docs, and the changelog.
  - **Co-change, added:** `design/system/flai-cli.md` (46%); shared.
  - **Co-change, not added:** `flai/cmd/serve_actions.go`. Its Settings section shows the planner's triggers, and no criterion asks for a role's there; the role page in S-0420 shows them.
  - No folder touch.
- **Forecast.** Kept at flai's 51m. Each trigger is modelled on a starter that exists. The forecast's basis gives the figures, and its delivery is flai's play-out moved by any raise times the cycle factor.
- **Cost of delay.** Not set. Neither the story nor E-0024 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0471, question 1, and the value follows its answer.
