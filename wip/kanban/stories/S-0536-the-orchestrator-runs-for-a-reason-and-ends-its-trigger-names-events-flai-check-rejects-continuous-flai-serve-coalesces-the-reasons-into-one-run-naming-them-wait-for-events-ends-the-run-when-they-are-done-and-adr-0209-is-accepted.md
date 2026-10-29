---
id: S-0536
type: story
nature: feature
title: "The orchestrator runs for a reason and ends: its trigger names events, flai check rejects continuous, flai serve coalesces the reasons into one run naming them, wait_for_events ends the run when they are done, and ADR-0209 is accepted"
status: ready
parent: E-0030
owner: alex
created: 2026-10-28T21:02:32Z
updated: 2026-10-29T00:42:47Z
transitions:
  - to: ready
    at: 2026-10-28T22:21:22Z
    by: alex
tags: [flai, serve, orchestrator]
topics: [orchestration, cli, template, conventions, mcp]
touches: [flai/internal/serve, flai/internal/mcpserver, flai/internal/harness, flai/internal/manifest, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/adrs, docs, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/check/roles_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, docs/operators/settings.md, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/orchestrate_hold_test.go, flai/internal/serve/orchestrate_timeout.go, flai/internal/serve/orchestrate_timeout_test.go, flai/internal/serve/serve.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, flai/internal/mcpserver/server.go, flai/internal/mcpserver/orchestrator_end_test.go, design/system/agent-narrative.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/users/flai.md, docs/operators/index.md, template/CHANGELOG.md]
after: [S-0532, S-0535]
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
      seconds: 464
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 59
          output: 954
          cache_read: 9373220
          cache_write: 18876
          cost: 2.0666
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h
    by: orchestrator
    at: 2026-10-28T21:42:06Z
  value: 300
  by: planner-S-0536
  at: 2026-10-28T21:47:43Z
forecast:
  duration: 5h
  delivery: 2026-10-29T13:38:00Z
  basis: "Its own forecast of 5h; 5th in the pull order with an in-progress limit of 5, behind S-0519, S-0514, S-0532, S-0537 and S-0534."
  by: flai
  at: 2026-10-29T00:42:47Z
---
# S-0536 The orchestrator runs for a reason and ends: its trigger names events, flai check rejects continuous, flai serve coalesces the reasons into one run naming them, wait_for_events ends the run when they are done, and ADR-0209 is accepted

## Goal

ADR-0209 §1, §2, §4, and §7. The orchestrator is the one role with `trigger: [continuous]`: it primes, decides, holds `wait_for_events`, and decides again without ending, and every wake re-reads a transcript of about 450k tokens, 75 to 80 percent of the time to wait again. This story makes it run for a reason and end, like the planner and the curator. `orchestrator.md`, here and in the template, names the events of the previous story as its trigger; `flai check` rejects `continuous` on any role. flai serve keeps a reason queue per project, starts one run no sooner than `orchestration.coalesce` (a Go duration, `30s` by default) after the first queued reason with every reason seen meanwhile, queues reasons that arrive during a run for the next, and stops a run that holds past `orchestration.run_timeout` (`20m` by default). The start prompt lists each reason with its item or thread; the role's prompt step 15 becomes: act on each reason, log each action, then end. `wait_for_events` answers the orchestrator `end: true`, `why: reasons done` when nothing queued remains for it, as it answers a story's agent that asked (ADR-0154). The guard's hold, the permissions, the activity log and its charge, the orchestrator's own thread, `flai serve orchestrate stop|start`, and the `orchestrate` host action stay as ADR-0087 and ADR-0189 have them. The ADR, drafted as proposed on 2026-10-28, is accepted in this story, and the living design, the convention and its template copy, and the docs change with it.

## Acceptance criteria
- [ ] `orchestrator.md` and its template copy name the trigger events; `flai check` reports `continuous` on any role as a finding, with a test
- [ ] flai serve queues the reasons a project's events give, starts one run per batch after `orchestration.coalesce`, with the reasons in the prompt and the journal entry, and starts the next run for reasons that arrived meanwhile when the run ends; a test drives three events inside one window into one run
- [ ] `wait_for_events` answers the orchestrator `end: true` with `why: reasons done` when nothing queued remains, and a run past `orchestration.run_timeout` is stopped and logged
- [ ] `orchestration.coalesce`, `orchestration.run_timeout`, and `orchestration.schedule` are manifest keys `flai check` validates and `flai manifest set` sets, with the defaults the ADR gives
- [ ] The orchestrator's prompt, the strategic-agents convention and its template copy, `flai-cli.md`, `strategic-agents.md`, `project-manifest.md`, and the user docs say the orchestrator runs for a reason and ends
- [ ] ADR-0209's status is accepted, and the orchestrator in this repository runs under the new trigger for a day with its runs and reasons in the journal, recorded in the narrative with the context tokens per action from `flai stats`

## Tasks
- T-2300 flai check reports continuous on any role as a finding, and orchestrator.md and its template copy name ADR-0209's trigger events
- T-2301 orchestration.coalesce, orchestration.run_timeout, and orchestration.schedule are manifest keys flai check validates and flai manifest set sets, with ADR-0209's defaults
- T-2302 flai serve queues the orchestrator's reasons per project and starts one run per batch after orchestration.coalesce, naming the reasons in its journal entry, and the next run for reasons that came meanwhile
- T-2303 The orchestrator's start prompt lists each reason with its item or thread, and its role prompt says act on each reason, log each action, then end
- T-2304 wait_for_events answers the orchestrator end true with why reasons done when nothing queued remains for it
- T-2305 flai serve stops an orchestrator run that holds past orchestration.run_timeout and logs and journals it
- T-2306 The strategic-agents convention and its template copy, flai-cli.md, strategic-agents.md, project-manifest.md, the roles README, and the user and operator docs say the orchestrator runs for a reason and ends
- T-2307 The orchestrator in this repository runs under the new trigger for a day, and the narrative records its runs, their reasons, and the context tokens per action

## Notes

Evidence: ADR-0209's context section. The budget hold is the second story of E-0030; reasons queued while it is set wait. The reason-specific pack is the next story: until it lands, a run primes the general `orchestrate` pack as today, which the per-action figure will show.

### Planning

Planned by planner-S-0536 on 2026-10-28. The tasks, their layers, and the assumptions are on the plan's thread.

ADR-0209 is already `accepted` on main (commit `docs: ADR-0209 accepted`), so criterion 6 needs only the day-long run for that half. The orchestrator's loop lives in `flai/internal/serve/orchestrate.go` (`look`, `orchestrateRetry`); its start prompt is `orchestratePrompt` in `flai/internal/harness/harness.go`; its role prompt is the body of `design/workflows/roles/orchestrator.md`, which `.claude/agents/orchestrator.md` copies; `continuous` is checked in `checker.triggers` in `flai/internal/roles/roles.go`; the `end` answer is `endWhy` in `flai/internal/mcpserver/server.go`; the orchestration keys are in `flai/internal/manifest/manifest.go` and `settings.go`.

Touches, and where each came from:

- Folders kept, all declared: `flai/internal/serve`, `flai/internal/mcpserver`, `flai/internal/harness`, `flai/internal/manifest`, and `docs`. The tasks name every file below them the plan expects, so the story's claim narrows each folder to those files (ADR-0096); the folders stay for files the reason queue may add that no task can name yet.
- `design/adrs`: declared, and kept as declared. No task names a file under it, since ADR-0209 is accepted already; the plan's thread proposes dropping it, so that it does not hold S-0537 and other stories that add an ADR.
- `design/workflows/roles/orchestrator.md` and its template copy, the strategic-agents convention and its template copy, `flai-cli.md`, `strategic-agents.md`, `project-manifest.md`: declared; criteria 1 and 5.
- `flai/internal/roles/roles.go`, `roles_test.go`, `flai/internal/check/roles_test.go`: layout; criterion 1.
- `flai/internal/manifest/manifest.go`, `manifest_test.go`, `settings.go`, `settings_test.go`, `docs/operators/settings.md`: layout and design; criterion 4, and the settings index the documentation convention requires.
- `flai/internal/serve/orchestrate.go`, `orchestrate_test.go`, `orchestrate_hold.go`, `orchestrate_hold_test.go`, `serve.go`, and the new `orchestrate_reasons.go`, `orchestrate_reasons_test.go`, `orchestrate_timeout.go`, `orchestrate_timeout_test.go`: layout; criteria 2 and 3.
- `flai/internal/harness/harness.go`, `harness_test.go`: layout; the reasons in the start prompt.
- `flai/internal/mcpserver/server.go` and the new `orchestrator_end_test.go`: layout; criterion 3.
- `.claude/agents/orchestrator.md`, `template/root/.claude/agents/orchestrator.md`: layout; the definitions Claude Code runs copy the role file's body.
- `design/system/agent-narrative.md`, `design/workflows/roles/README.md` and its template copy, `docs/users/flai.md`, `docs/operators/index.md`, `template/CHANGELOG.md`: co-change and design; criterion 5.
- Left out of the suggestion: `flaiover-dashboard.md`, `workflow.md`, `work-management.md`, the host API, `serve_actions.go`, `root.go`, `check.go`, `guard.go`, and `flaiover`. They co-change with past serve stories, but nothing in them changes here: no command, host write, guard rule, or page is added.

Forecast: `flai forecast` gave 35m (median 107 s per unit of size over 95 done feature stories on claude-opus-5-5 in the large band, size 19: 6 criteria, 13 touches). Raised to 5h: eight tasks in five layers, with a reason queue that replaces the orchestrator's loop and its tests, the end and the timeout, three manifest keys, two prompts, and eleven documents, beside S-0535's 4h for nine tasks. Delivery is flai's 2026-10-29T17:02Z moved by the added 4h25m and by the day criterion 6 runs the orchestrator for: 2026-10-30T21:27Z.

Cost of delay: `flai cod` gave 300.00 USD a week from the input the orchestrator set on TH-0716 (2h of time lost per 168h cycle at 150 USD an hour, 1.00 cycles a week). The value stands as computed; the input is not changed.
