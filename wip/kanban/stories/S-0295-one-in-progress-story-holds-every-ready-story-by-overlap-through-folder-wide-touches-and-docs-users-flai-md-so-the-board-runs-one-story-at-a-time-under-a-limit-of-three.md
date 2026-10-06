---
id: S-0295
type: story
nature: improvement
title: One in-progress story holds every ready story by overlap, through folder-wide touches and docs/users/flai.md, so the board runs one story at a time under a limit of three
status: ready
owner: alex
created: 2026-10-06T11:44:50Z
updated: 2026-10-06T18:11:31Z
transitions:
  - to: ready
    at: 2026-10-06T11:47:42Z
    by: alex
tags: [flai]
topics: [cli, mcp, hostapi, dashboard, conventions, planning, template]
touches: [design/adrs, design/system/workflow.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/tech/go-libraries.md, design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, system-flow.yaml, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/hold_i0087_test.go, flai/internal/workitem/boardview.go, flai/internal/serve/hold_test.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/shared.go, flai/internal/mcpserver/shared_test.go, flai/internal/mcpserver/server.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/shared.go, flai/internal/manifest/shared_test.go, flai/go.mod, flai/go.sum, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/cmd/stream_sync_test.go, flai/cmd/shared.go, flai/cmd/shared_test.go, flai/cmd/root.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/itemedit/claim.go, flai/internal/itemedit/claim_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/guard, flai/internal/hostapi/settings.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/contract_test.go, flaiover/src/lib/settings.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/lib/server/agent.ts, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/users/flaiover.md, design/issues/I-0087-one-in-progress-story-holds-every-ready-story-by-overlap-through-folder-wide-touches-and-docs-users-flai-md-so-the-board-runs-one-story-at-a-time-under-a-limit-of-three.md, design/issues/summary.md]
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
    - kind: planner
      seconds: 1861
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 600
          output: 10306
          cache_read: 3291764
          cache_write: 195518
          cost: 0.6717
        - model: claude-opus-5-5
          input: 292
          output: 43426
          cache_read: 11320363
          cache_write: 240530
          cost: 5.622
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5h12m
    by: flai
    at: 2026-10-06T11:44:50Z
  value: 780
  by: planner-S-0295
  at: 2026-10-06T11:49:48Z
forecast:
  duration: 2h
  delivery: 2026-10-07T00:39:00Z
  basis: "Its own forecast of 2h; 7th in the pull order with an in-progress limit of 3, behind S-0226, S-0284, S-0223, S-0224, S-0227, S-0278 and S-0229."
  by: flai
  at: 2026-10-06T18:11:31Z
finalized:
  by: alex
  at: 2026-10-06T11:47:30Z
---
# S-0295 One in-progress story holds every ready story by overlap, through folder-wide touches and docs/users/flai.md, so the board runs one story at a time under a limit of three

## Goal

This story remediates [I-0087](../../../design/issues/I-0087-one-in-progress-story-holds-every-ready-story-by-overlap-through-folder-wide-touches-and-docs-users-flai-md-so-the-board-runs-one-story-at-a-time-under-a-limit-of-three.md), "One in-progress story holds every ready story by overlap, through folder-wide touches and docs/users/flai.md, so the board runs one story at a time under a limit of three". The issue recommends this solution:

Not designed yet. Directions to weigh:

- The planner narrows a folder claim to the files a story will change, where its tasks already name them.
- A document every story edits is split by section or by command, or is generated, so that two stories change different files.
- Whether a story in review still needs to hold others is a question for ADR-0046: its branch is finished and synced, and a later story syncs onto main when it is accepted.

## Acceptance criteria
- [ ] The cause I-0087 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0087 is closed with `flai issue close I-0087 --reason` saying what fixed it

## Tasks
- T-1023 An ADR refines ADR-0046: a story in review holds nothing, shared paths never hold, and tasks narrow a folder claim, and the design says so
- T-1024 A story in review no longer holds a ready story; only one in progress does
- T-1025 The manifest carries the shared paths as glob patterns, validated, with a matcher and an edit function the CLI, HTTP, and MCP share
- T-1026 The planner and the story's agent are told to declare file-level touches and to narrow a folder touch once tasks name its files
- T-1027 A story's folder touch is narrowed in its claim to the files its tasks name inside it, done tasks included
- T-1028 flai shared lists, adds, removes, and checks the shared paths from the CLI
- T-1029 MCP tools list and check the shared paths for any agent, and change them only for the operator's session
- T-1030 An overlap wholly inside a shared path neither holds a ready story nor counts as a wip.overlap or a grown claim
- T-1031 The host API reports and checks the shared paths, and settings.shared adds or removes them through flai shared
- T-1032 The dashboard's project settings list the shared paths, add and remove patterns, and check a path against them
- T-1033 The users' and operators' guides and the dashboard design describe the new hold rules, the shared paths, and how to edit them
- T-1034 A test rebuilds I-0087's board and shows three ready stories running beside the open one, and I-0087 is closed

## Notes

Cost of delay inputs set by flai from I-0087. time_lost_per_cycle 5h12m: 1h44m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-06T10:32:27Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Scope comes from TH-0180. The operator took all three of I-0087's directions as one ADR refining ADR-0046:

- a story in review holds nothing;
- shared paths, as glob patterns in the manifest, never hold;
- tasks narrow a folder claim.

The operator added that the shared list is editable through the CLI, HTTP, MCP, and the dashboard's project settings, with globs. Splitting `docs/users/flai.md` by command was not chosen; the shared list covers it.

**Touches.** The story declared none, so `flai touches suggest` had nothing to start from. Every touch comes from design or layout, and each is the union of its tasks' touches:

- *Design*: the ADR and the sections it changes. These are `design/adrs` (kept as a folder because the ADR's number is not known until `flai adr new`), `workflow.md` § Branches and collisions, `project-manifest.md`, `flai-cli.md`, `flaiover-dashboard.md`, and the two conventions with their template copies. The docs are `docs/users/flai.md`, `flai-reference.md`, `flaiover.md`, and `docs/operators/settings.md`. The issue is I-0087 with `summary.md`.
- *Layout*: read from the code by the explorer.
  - The hold and claim: `flai/internal/workitem/hold.go` and `boardview.go`, and their readers `flai/internal/check/check.go`, `flai/internal/itemedit/claim.go`, and `flai/cmd/stream_sync.go`'s test.
  - The manifest: `flai/internal/manifest`, `system-flow.yaml`, and `template/root/system-flow.yaml.tmpl`.
  - The settings surfaces: `flai/internal/hostapi/settings.go`, `flai/cmd/serve_actions.go`, `flaiover/src/lib/settings.ts`, `SettingsPanel.svelte`, and `flaiover/src/lib/server/agent.ts`.
  - The MCP server's `server.go` and `flai/internal/guard`, kept as a folder because the guard's rule file is not yet chosen.
  - The planner's prompt in `flai/internal/harness/harness.go`.
- *Co-change*: once the touches were set, `flai touches suggest` proposed `docs/operators/index.md` (11%) and `flai/internal/hostapi/writes.go` (7%), among others. Neither was taken. `index.md` only links `settings.md`, and the settings methods live in `settings.go`, not `writes.go`.
- Files are named, not folders, wherever a task can name them, as the story itself asks.

**Forecast.** `flai forecast` gave 1h24m: 89 s per unit over 21 large improvement stories, times size 56. I raised it to 2h. The twelve tasks run in seven layers, mostly one after another, across flai, flaiover, and the template, and the per-touch rate does not see that chain. Delivery moved by the same 36m, to 18:34Z from flai's 17:58Z, as 8th in the pull order.

**Cost of delay.** The value of 780 USD a week stands as `flai cod` gave it: 5h12m lost per 168h cycle at 150 USD an hour. Only three occurrences on one day back the input, so it may understate. Raising it is the operator's call.
