---
id: D-0191
title: "S-0298 plan: nine tasks in five layers, and two questions on scope (MCP deploy, dashboard pin): list only, 2. apply once, no pinning"
kind: answer
sources: [TH-0202]
topics: [release, security, flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src, flai/cmd, flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go, flai/cmd/selfupgrade.go, flai/cmd/selfupgrade_test.go, flai/cmd/dashboard.go, flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_test.go, flai/cmd/host.go, flai/cmd/host_integration_test.go, flai/internal/host/api.go, flai/internal/host/client.go, flai/internal/host/host.go, flai/internal/host/host_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/versions.go, flai/internal/mcpserver/versions_test.go, flai/internal/mcpserver/folder.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/src/routes/api/host/+server.ts, flaiover/src/routes/api/host/host.test.ts, flaiover/src/routes/api/dashboard/+server.ts, flaiover/src/routes/api/dashboard/dashboard.test.ts, flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, flaiover/src/lib/components/HostProcesses.svelte, flaiover/src/lib/components/HostProcesses.svelte.test.ts, design/adrs, design/adrs/README.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flaiover.md, docs/operators/index.md, docs/operators/runbooks/update.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/summary.md, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/server_test.go, flai/cmd/dashboard_versions.go, flai/cmd/dashboard_versions_test.go, flai/internal/buildinfo/buildinfo.go, flai/internal/buildinfo/buildinfo_test.go, flai/cmd/host_versions_test.go, design/issues/I-0115-the-flaiover-guide-says-check-for-updates-always-works-but-the-dashboard-area-hides-it-while-the-dashboard-host-action-is-off.md]
date: 2026-10-06
status: absorbed
absorbed_by: "design/adrs/0117-the-host-api-installs-a-release-the-dashboard-names-only-when-flai-lists-it-as.md#decision"
asked_again: []
cited: []
---

# D-0191 S-0298 plan: nine tasks in five layers, and two questions on scope (MCP deploy, dashboard pin): list only, 2. apply once, no pinning

The operator, alex, answered TH-0202 at 2026-10-06T21:51:45Z.

What was asked:

```text
Plan for S-0298, by planner-S-0298. The touches, forecast (1h1m, delivery 2026-10-07T09:29Z) and cost of delay (37.50 USD a week, from your 15m per cycle) are on the story. Where each touch came from is under `### Planning` in its Notes.
```

The answer:

```text
1. list only, 2. apply once, no pinning
```

Source: TH-0202.
