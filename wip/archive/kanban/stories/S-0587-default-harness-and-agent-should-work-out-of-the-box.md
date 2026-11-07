---
id: S-0587
type: story
nature: remediation
title: Default harness and agent should work out of the box
status: done
owner: alex
created: 2026-11-07T08:40:55Z
updated: 2026-11-07T15:15:56Z
transitions:
  - to: ready
    at: 2026-11-07T08:41:12Z
    by: alex
  - to: in-progress
    at: 2026-11-07T10:22:00Z
    by: agent-S-0587
  - to: review
    at: 2026-11-07T12:05:01Z
    by: system-flow
  - to: done
    at: 2026-11-07T15:15:56Z
    by: alex
tags: [cli, dashboard, template]
topics: [agents, settings]
touches: [design/adrs/README.md, flai/internal/manifest/agent.go, flai/internal/manifest/agent_test.go, flai/internal/workitem/agent.go, flai/internal/workitem/agent_test.go, flai/internal/serve/default_agent_test.go, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, flai/cmd/new_test.go, flai/cmd/import_test.go, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go, docs/operators/settings.md, docs/users/flai-reference.md, flaiover/src/lib/settings.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flai/internal/serve/start.go, flai/internal/serve/restart.go, flai/internal/serve/commit.go, flai/internal/serve/agents.go, flai/internal/serve/plan.go, flai/internal/serve/analyze.go, flai/internal/serve/curate.go, flai/internal/serve/orchestrate.go, flai/internal/serve/role.go, flai/internal/serve/start_test.go, flai/internal/serve/agents_test.go, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, flaiover/src/lib/components/HostAgentNotice.svelte, flaiover/src/lib/components/HostAgentNotice.svelte.test.ts, flaiover/src/lib/activity.ts, flaiover/src/lib/activity.test.ts, docs/operators/index.md, docs/users/flai.md, docs/users/flaiover.md, docs/operators/runbooks/install.md, design/system/agent-adapters.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/template.md, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, design/system/flaiover-dashboard.md, design/system/strategic-agents.md, design/adrs/0261-a-story-or-strategic-agent-that-names-no-harness-falls-back-after-the-project-s.md, flai/cmd/agent.go, flai/cmd/analyze_test.go, flai/cmd/curate_test.go, flai/cmd/plan_test.go, flai/cmd/role.go, flai/cmd/role_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/mcpserver/role.go, flai/internal/serve/analyze_test.go, flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go, flai/internal/serve/curate_test.go, flai/internal/serve/hold_test.go, flai/internal/serve/orchestrate_hold_test.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/plan_test.go, flai/internal/serve/review_wait_test.go, flai/internal/serve/role_events_test.go, flai/internal/serve/role_test.go, flai/internal/serve/serve_test.go, flai/internal/serve/stop_test.go, flaiover/src/lib/server/writes.test.ts, design/issues/I-0304-flai-test-on-a-change-to-flai-s-host-api-runs-none-of-flaiover-s-server-tests.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 6202
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 6
      hand_edits: 2
      work: 102
  models:
    - model: claude-haiku-5-5
      input: 42
      output: 15957
      cache_read: 1343593
      cache_write: 99819
      cost: 0.0339
    - model: claude-opus-5-5
      input: 766
      output: 226393
      cache_read: 147712245
      cache_write: 1277443
      cost: 41.8591
    - model: claude-sonnet-5-5
      input: 18
      output: 6547
      cache_read: 401569
      cache_write: 78320
      cost: 0.3015
  strategic:
    - kind: planner
      seconds: 1121
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 117
          output: 1870
          cache_read: 6920686
          cache_write: 195606
          cost: 0.4687
        - model: claude-opus-5-5
          input: 316
          output: 39653
          cache_read: 24729896
          cache_write: 392598
          cost: 7.9702
    - kind: orchestrator
      seconds: 230
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 38
          output: 289
          cache_read: 1392508
          cache_write: 149569
          cost: 0.3871
        - model: claude-sonnet-5-5
          input: 6
          output: 35
          cache_read: 43557
          cache_write: 40779
          cost: 0.061
forecast:
  duration: 3h
  delivery: 2026-11-07T13:53:00Z
  basis: "flai forecast's 1h31m (105 s per unit over 57 large remediation stories, size 52), raised for eight tasks in five layers across Go, the template, and flaiover, with a host-config migration; S-0573, the narrower fix of the same failure, took 31m for three tasks."
  by: planner-S-0587
  at: 2026-11-07T10:13:32Z
---
# S-0587 Default harness and agent should work out of the box

## Goal

In new projects, once an agent has been enabled, new stories that an operator tries to start the agent for report the following error message:

```text
S-#### names no harness, and no command is set on the host (flai serve agent set -- <program> [args...])
```

This indicates a failure of flai to provide the default settings (the system already has a default harness and agent selected).

Finally, if it were missing, the error message to direct the user to the configuration pages or sub-pages necessary to complete configuration.

## Acceptance criteria
- [x] A default agent should always be provided to stories
- [x] A default harness should always be provided to stories
- [x] Error messages must provide a link to the page(s) where the user can enter the missing information if no default exists
- [x] Harness and agent settings should be independent of other projects (setting these for one should never change another project's harness or agent)

## Tasks
- T-2577 An ADR decides the built-in default agent, the template's default, per-project host agent settings, and the refusal's links to the settings pages
- T-2578 flai's built-in default agent, claude-code with no model, is what an agent falls back to when neither it, the project's default, nor the host's command names a harness
- T-2579 The template's system-flow.yaml ships a default agent, so flai new and flai import write one into every new project
- T-2580 The host's agent command, name, and harness programs are kept per project, with a host-wide fallback, so setting them for one project changes no other
- T-2581 The Harnesses settings page edits this project's command and harness programs, and changes every project only when the operator chooses everywhere
- T-2582 Every start falls back to the built-in default agent, and each refusal and waiting reason that remains names the missing setting, the dashboard page that sets it, and the command
- T-2583 The story page and the host agent notice render the refusal's settings pages as links
- T-2584 The design and the user and operator docs say every story has a default agent and harness, and that host agent settings are per project

## Notes

### Planning

Cause, read from the code (planner's explorer, 2026-11-07):

- `template/root/system-flow.yaml.tmpl` ships `agent:` commented out, and neither `flai new` nor `flai import` writes one. So a new project has no default agent.
- `workitem.Repo.AgentOf` (ADR-0251) falls back to that default, which is empty. flai has no built-in default for a story's agent; `builtinAgents` in `flai/internal/manifest/manifest.go` covers only the strategic roles.
- The `claude-code` the dashboard shows is a placeholder in `AgentFields.svelte`, not a stored default.
- The host command and the harness programs (`flai serve agent set`, `flai serve agent harness`) live in `~/.flai/config.json` for every project on the host. Only the host actions are keyed by project root. That breaks criterion 4.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/serve/start.go`, `restart.go`, `commit.go`, `agents.go` | design: ADR-0251 and the refusal's code |
| `flai/internal/serve/plan.go`, `analyze.go`, `curate.go`, `orchestrate.go`, `role.go` | layout: the same refusal for each role |
| `flai/internal/workitem/agent.go`, `create.go`, `flai/internal/manifest/agent.go` and tests | design: ADR-0037, ADR-0251 |
| `template/root/system-flow.yaml.tmpl`, `template/CHANGELOG.md`, `flai/cmd/new_test.go`, `flai/cmd/import_test.go` | layout: the template manifest; co-change |
| `flai/internal/config/config.go`, `flai/cmd/serve_actions.go`, `flai/internal/hostapi/settings.go` and tests | co-change (40%, 27%, 35%); layout |
| `docs/operators/settings.md`, `docs/users/flai-reference.md` | co-change (37%, 31%); the settings-index rule |
| `flaiover/src/lib/settings.ts`, `SettingsPanel.svelte`, `StoryAgent.svelte`, `HostAgentNotice.svelte`, `activity.ts` and tests | layout: the settings pages and where the refusal shows |
| `docs/operators/index.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, `docs/operators/runbooks/install.md`, `design/system/agent-adapters.md`, `project-manifest.md`, `flai-cli.md`, `template.md` | co-change; the docs that state no default today |
| `design/adrs`, `design/adrs/README.md` | design: the new ADR |

The folder `design/adrs` stays. The ADR's file name comes from `flai adr new`, which allocates its number when the story's agent runs it, so no task can name the file yet.

No `after`: S-0587 reads no other open story's output.

Plan, in five layers:

1. T-2577
2. T-2578, T-2579, T-2580
3. T-2581, T-2582
4. T-2583
5. T-2584

No two tasks in the same layer share a touch.

Forecast: 3h, raised from flai's 1h31m. flai's figure is 105 s per unit of size over 57 large remediation stories, times size 52, and it counts the touches but not the work. The work is eight tasks in five layers across Go, the template, and flaiover, which run three test tiers, plus a host-config migration that must keep old configs working. For comparison, S-0573 fixed the same failure more narrowly: three tasks and 17 touches in 31 minutes.

Cost of delay: no value yet. S-0587 has no inputs and no epic, so `flai cod S-0587` refuses. I asked for an input on TH-0872 and recommended `penalty_per_week: 50`, as on S-0573.
