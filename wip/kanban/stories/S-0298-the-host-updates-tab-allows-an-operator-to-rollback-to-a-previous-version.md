---
id: S-0298
type: story
nature: improvement
title: The Host Updates Tab Allows an Operator to Rollback to A Previous Version
status: ready
owner: alex
created: 2026-10-06T20:53:15Z
updated: 2026-10-07T14:26:01Z
transitions:
  - to: ready
    at: 2026-10-07T14:14:49Z
    by: orchestrator
tags: [dashboard, cli]
topics: [release, security]
touches: [flaiover/src, flai/cmd, flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go, flai/cmd/selfupgrade.go, flai/cmd/selfupgrade_test.go, flai/cmd/dashboard.go, flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_test.go, flai/cmd/host.go, flai/cmd/host_integration_test.go, flai/internal/host/api.go, flai/internal/host/client.go, flai/internal/host/host.go, flai/internal/host/host_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/versions.go, flai/internal/mcpserver/versions_test.go, flai/internal/mcpserver/folder.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/src/routes/api/host/+server.ts, flaiover/src/routes/api/host/host.test.ts, flaiover/src/routes/api/dashboard/+server.ts, flaiover/src/routes/api/dashboard/dashboard.test.ts, flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, flaiover/src/lib/components/HostProcesses.svelte, flaiover/src/lib/components/HostProcesses.svelte.test.ts, design/adrs, design/adrs/README.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flaiover.md, docs/operators/index.md, docs/operators/runbooks/update.md]
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
      seconds: 298
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 113
          cache_read: 1223450
          cache_write: 6286
          cost: 0.3209
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: alex
    at: 2026-10-06T20:53:15Z
  value: 37.5
  by: planner-S-0298
  at: 2026-10-06T21:46:33Z
forecast:
  duration: 1h1m
  delivery: 2026-10-07T15:30:00Z
  basis: "Its own forecast of 1h1m; 1st in the pull order with an in-progress limit of 3, behind S-0293."
  by: flai
  at: 2026-10-07T14:26:01Z
---
# S-0298 The Host Updates Tab Allows an Operator to Rollback to A Previous Version

## Goal

Presently, the host Updates page allows the operator to deploy new versions. There isn't presently any tooling to allow an operator to roll back to a prior or specific version.

Introduce flai commands for the cli, http, mcp protocols so that the operator can discover available older versions and deploy one for both the flai CLI and the dashboard.

## Acceptance criteria
- [ ] New commands are available to allow the operator to view prior versions
- [ ] New commands are available to allow the operator to select a specific version of the CLI and dashboard to deploy
- [ ] Deploying a specific version works like updating to the latest

## Tasks
- T-1038 Record in an ADR and the design that a host action may install a named published release, and name the versions commands
- T-1039 selfupgrade lists the published releases for a tag prefix, newest first
- T-1040 flai self-upgrade --list and flai dashboard versions list the published releases, marking the running and the newest
- T-1041 The MCP read tool versions lists the published flai and dashboard releases
- T-1042 The dashboard's host and dashboard routes pass the versions reads and a chosen version or tag to flai on the host
- T-1043 flai host versions lists the flai releases and flai host upgrade --version installs a chosen one and restarts on it
- T-1044 The Updates page lists the flai and dashboard releases and deploys a chosen one, as Upgrade does for the newest
- T-1045 The host API reads host.versions and dashboard.versions, and its upgrades install a named release only when it is published
- T-1046 The operators' runbook, the user guides, and the generated reference describe listing and deploying an earlier release

## Notes

### Planning

Written by planner-S-0298 on 2026-10-06. The plan is on TH-0202, where the operator answered its two questions: MCP lists versions only and deploys nothing, and a dashboard version chosen on the Updates page applies once, with no pinning. T-1038, T-1041, T-1044, and T-1046 say so.

What exists today: `flai self-upgrade --version` and `flai dashboard upgrade --tag` already install a chosen release from a shell, and the runbook (`docs/operators/runbooks/update.md`) says to go back with them. Nothing lists the published releases, `flai host upgrade` and its `POST /upgrade` install only the newest, the host API's `host.upgrade` and `dashboard.upgrade` take no version or tag from the dashboard by design (S-0081), and no MCP tool touches versions.

Touches, and where each came from:

- Declared by the operator, kept: `flaiover/src` and `flai/cmd`. They are folder touches; the planner keeps every declared touch. Each is covered by the tasks' file touches below, which replace it in the story's claim while the story is in progress (ADR-0096).
- Co-change (`flai touches suggest`, 453 of 1033 commits): `flai/internal/hostapi/writes.go` and `writes_test.go`, `flai/internal/mcpserver/folder.go`, `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`, `docs/users/flaiover.md`, `docs/operators/index.md`, `docs/operators/settings.md`, `design/adrs/README.md`. Its other suggestions (`serve`, `check`, `mcpserver/server.go`, the issues summary, `work-hierarchy.md`, `workflow.md`) are not reached by this story's work, so they are left out.
- Layout, from the upgrade code paths: `flai/internal/selfupgrade/selfupgrade.go` and its test (release lookup), `flai/cmd/selfupgrade.go`, `dashboard.go`, `dashboard_upgrade.go`, `host.go` and their tests, `flai/internal/host/api.go`, `client.go`, `host.go`, `host_test.go` (the host's `POST /upgrade`), `flai/internal/mcpserver/versions.go` and its test (new), `flaiover/src/lib/server/agent.ts` and its test (the channel's method list), `flaiover/src/routes/api/host` and `api/dashboard` `+server.ts` and tests, `flaiover/src/lib/components/HostPanel.svelte` and `HostProcesses.svelte` and tests (the Updates page's two areas).
- Design: `design/adrs` (folder, kept), for the new ADR T-1038 writes: its number is taken when it is written, so no task can name the file yet. `docs/operators/runbooks/update.md`, whose "To go back" lines this story replaces.
- Left out on purpose: `flai/internal/guard/guard.go`. The `versions` tool is not added to `guard.MCPReads`, since a story's sub-agents do not need it, which keeps this story clear of the guard that S-0229 and S-0299 hold.
- Topics added: `release` and `security`, since T-1038 decides what release a dashboard request may install and T-1039 reads the release list.

Forecast: 1h1m, delivery 2026-10-07T09:29Z, as `flai forecast` gives it, and it stands. It is the median over 24 done large-band improvement stories on the same model, and nothing in the work sets it apart from them. The two declared folder touches repeat files already counted, but that moves the size by 2 of 43, too little to adjust for. The plan's five layers are mostly serial, which the median already reflects for stories of this size.

Cost of delay: 37.50 USD a week, from the operator's input of 15m lost per weekly cycle at the project's 150 USD an hour (`flai cod`). It stands: an operator who needs to go back can already do it from a shell (`flai self-upgrade --version`, `flai dashboard upgrade --tag`), so the loss each cycle is the time spent finding and running those by hand.
