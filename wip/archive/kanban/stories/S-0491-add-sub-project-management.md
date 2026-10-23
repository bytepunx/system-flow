---
id: S-0491
type: story
nature: improvement
title: Add sub-project management
status: done
owner: alex
created: 2026-10-23T20:26:10Z
updated: 2026-10-23T23:38:50Z
transitions:
  - to: ready
    at: 2026-10-23T20:26:18Z
    by: alex
  - to: in-progress
    at: 2026-10-23T20:47:32Z
    by: agent-S-0491
  - to: review
    at: 2026-10-23T23:37:24Z
    by: agent-S-0491
  - to: done
    at: 2026-10-23T23:38:50Z
    by: alex
tags: [cli, dashboard, template]
topics: [manifest, settings, mcp, orchestration]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/projects.go, flai/internal/manifest/projects_test.go, flai/internal/manifest/settings.go, flai/internal/check/check.go, flai/internal/check/projects_test.go, flai/cmd/root.go, flai/cmd/subproject.go, flai/cmd/subproject_test.go, flai/internal/guard/guard.go, flai/internal/guard/subproject_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/subproject.go, flai/internal/mcpserver/subproject_test.go, flai/internal/mcpserver/server_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/settings.ts, flaiover/src/lib/settingspages.ts, flaiover/src/lib/settingspages.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/routes/api/settings/+server.ts, flaiover/src/routes/api/settings/settings.test.ts, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, docs/operators/settings.md, design/system/flai-cli.md, design/system/project-manifest.md, design/system/dashboard-host-channel.md, design/system/flaiover-dashboard.md, design/system/strategic-agents.md, design/conventions/strategic-agents.md, design/adrs, CLAUDE.md, ".claude/agents/orchestrator.md", template/root/CLAUDE.md.tmpl, template/root/system-flow.yaml.tmpl, template/root/design/conventions/strategic-agents.md, template/root/.claude/agents/orchestrator.md, template/CHANGELOG.md, flai/internal/mcpserver/folder_test.go, design/issues/I-0202-task-done-s-answer-carries-every-check-note-outside-the-story-so-it-exceeds-the-40-000-bytes-it-promises-and-overflows-the-tool-result.md, design/issues/summary.md, flai/internal/guard/guard_test.go, flai/internal/guard/orchestrate_test.go, flai/internal/harness/harness.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 7133
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 4
      hand_edits: 2
      work: 70
  models:
    - model: claude-opus-5-5
      input: 752
      output: 308545
      cache_read: 46193046
      cache_write: 1797592
      cost: 25.2833
  strategic:
    - kind: planner
      seconds: 1007
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 546
          output: 34823
          cache_read: 2329949
          cache_write: 190278
          cost: 6.1347
        - model: claude-haiku-5-5
          input: 226
          output: 61931
          cache_read: 25241634
          cache_write: 357653
          cost: 1.5734
    - kind: orchestrator
      seconds: 1301
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 144
          output: 2490
          cache_read: 46525000
          cache_write: 24009
          cost: 10.238
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3h
    by: alex
    at: 2026-10-23T20:26:10Z
  value: 450
  by: planner-S-0491
  at: 2026-10-23T20:42:48Z
forecast:
  duration: 2h
  delivery: 2026-10-23T23:05:00Z
  basis: "flai forecast gave 1h9m from the median of 94 done improvement stories on claude-opus-5-5 in the large band times size 46; raised to 2h because the plan has seven tasks in five layers, each waiting on the one before, across the Go CLI, the host API, the MCP server, the guard, the Svelte dashboard, the docs, and the template, with an ADR and a .claude/ write that waits on the operator's permission prompt; S-0295 and S-0489, which spanned the CLI, the host API, and the dashboard with fewer layers, took 1h15m and 1h29m. Delivery is the start after S-0430's forecast delivery plus 2h at the history's cycle factor of 1.18."
  by: planner-S-0491
  at: 2026-10-23T20:42:48Z
---
# S-0491 Add sub-project management

## Goal

No flai command registers a sub-project, and no document in the project or the template gives the schema of projects. The operator is expected to add the entry to system-flow.yaml by hand. flai's own messages show that an entry has name, path, and kind, but there's little else to go off.

## Acceptance criteria
- [x] The projects page supports sub-project management in addition to the projects themselves
- [x] If no sub-folder exists, creating a sub-project adds that sub-folder to the project
- [x] Add supporting API and MCP calls for this as well so that agents (orchestrator and story agents) can handle this directly without requiring the user to take direct steps

## Tasks
- T-2092 Validate and edit the manifest's projects list in flai's manifest package, and report manifest.projects in flai check
- T-2093 Add flai subproject list, add, set, and remove, creating the folder an added sub-project lacks, and document the command
- T-2094 Add the host API write settings.subproject, gated by the settings action, that runs flai subproject for the dashboard
- T-2095 Add the MCP tools subprojects and subproject_edit, open to a story's agent, refused to sub-agents, the planner, and the analyzer by flai guard
- T-2096 Give the orchestrator subproject_edit behind a new permission edit_subprojects, off by default, recorded in an ADR, the conventions, its agent definition, and the template
- T-2097 Add a Sub-projects section to the dashboard's Projects configuration page that lists, adds, edits, and removes a project's sub-projects through settings.subproject
- T-2098 Document the projects schema and who writes it in project-manifest.md, the template's manifest comments, CLAUDE.md, and the template changelog

## Notes

### Planning

Planned by the planner on 2026-10-23; TH-0626 on this story names the tasks, their layers, and the assumptions.

Touches, none declared, all predicted from the goal, the criteria, and the code layout (`flai touches suggest` had no path to start from):

- Layout, flai: `flai/internal/manifest/manifest.go`, `settings.go`, and the new `projects.go` with its test (the `projects[]` model, its validation, and the line editor beside `EditComponents`); `flai/internal/check/check.go` and `projects_test.go` (the `manifest.projects` finding beside `manifest.components`); `flai/cmd/root.go`, `subproject.go`, and `subproject_test.go` (the command group); `flai/internal/hostapi/settings.go` and `writes_test.go` (the write `settings.subproject`, whose test demands a case per method); `flai/internal/mcpserver/folder.go`, `subproject.go`, `subproject_test.go`, and `server_test.go` (the tools and the tool list); `flai/internal/guard/guard.go`, `subproject_test.go`, and `orchestrate_test.go` (who may call them).
- Layout, dashboard: `flaiover/src/lib/server/agent.ts` (the methods contract test), `settings.ts`, `settingspages.ts` and its test, `components/SettingsPanel.svelte` and its test, `routes/api/settings/+server.ts` and its test (the Sub-projects section of the Projects configuration page).
- Design and docs: `design/system/project-manifest.md` (the schema the goal says is missing), `flai-cli.md`, `dashboard-host-channel.md`, `flaiover-dashboard.md`, `strategic-agents.md`, and `design/conventions/strategic-agents.md` (the new permission); `docs/users/flai.md`, `flai-reference.md` (generated), `flaiover.md`, and `docs/operators/settings.md` (the flag index and the permission row, each kept by a test).
- Template: `template/root/CLAUDE.md.tmpl`, `system-flow.yaml.tmpl`, `design/conventions/strategic-agents.md`, `.claude/agents/orchestrator.md`, and `template/CHANGELOG.md`, with the project's `CLAUDE.md` and `.claude/agents/orchestrator.md`, which the template copies mirror; the `.claude/` files are protected paths, so the operator accepts the story (ADR-0106).
- Folder kept: `design/adrs/`, because T-2096 records the orchestrator's permission with `flai adr new`, whose file name is not known until it is made; `design/adrs/README.md` is named by the task.

Forecast: `flai forecast` gave 1h9m (median 89 s per unit of size over 94 done improvement stories on claude-opus-5-5 in the large band, size 46); set to 2h, delivery 2026-10-23T23:05:00Z, because the seven tasks form five layers that cannot run together, span three codebases, add an orchestrator permission with an ADR, and write under `.claude/`, which waits on the operator's permission prompt; the nearest comparable stories, S-0295 and S-0489, took 1h15m and 1h29m with fewer layers.

Cost of delay: 450 USD a week, as `flai cod` computes it from the operator's input of 3h lost per 168h cycle at 150 USD an hour, one cycle a week; it stands.
