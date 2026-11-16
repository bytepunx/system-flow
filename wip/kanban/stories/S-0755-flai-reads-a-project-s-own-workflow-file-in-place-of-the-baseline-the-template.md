---
id: S-0755
type: story
nature: feature
title: flai reads a project's own workflow file in place of the baseline, the template ships the baseline, flai upgrade keeps the project's file, and the guides say how to change the board
status: review
parent: E-0034
owner: alex
created: 2026-11-12T16:02:15Z
updated: 2026-11-15T23:16:28Z
transitions:
  - to: ready
    at: 2026-11-15T21:35:42Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-15T21:46:30Z
    by: agent-S-0755
  - to: review
    at: 2026-11-15T23:16:28Z
    by: agent-S-0755
tags: [cli, template]
topics: [cli, template, conventions]
touches: [flai/internal/workflow/load.go, flai/internal/workflow/load_test.go, flai/internal/workitem/repo.go, flai/internal/check/check_workflow.go, flai/internal/check/workflow_test.go, template/root/design/workflows/board.md, flai/internal/workflow/baseline_test.go, flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, template/CHANGELOG.md, docs/operators/runbooks/change-the-board.md, docs/operators/runbooks/README.md, docs/operators/index.md, docs/users/flai.md, docs/users/conventions.md, design/system/workflow.md, design/system/template.md, design/workflows/README.md, template/root/design/workflows/README.md, design/adrs/0307-flai-reads-the-workflow-file-of-the-checkout-it-opens-refuses-to-open-a.md, design/adrs/README.md, flai/cmd/check.go, flai/cmd/project.go, flai/internal/check/concurrent_test.go, flai/internal/workflow/validate.go, flai/internal/workitem/repo_flow_test.go, flai/internal/check/check_items.go, flai/cmd/accept.go, flai/cmd/accept_overlap.go, flai/cmd/accept_overlap_test.go, flai/cmd/cod.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/cmd/task_done.go, flai/internal/inbox/inbox.go, flai/internal/itemedit/forecasts.go, flai/internal/itemedit/itemedit.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/workflow_test.go, flai/internal/planning/cod.go, flai/internal/planning/forecast.go, flai/internal/preview/accept.go, flai/internal/serve/flow.go, flai/internal/serve/flow_test.go, flai/internal/serve/kluster.go, flai/internal/serve/replan.go, flai/internal/serve/role_events.go, design/system/flai-cli.md, design/issues/I-0380-adr-0301-rules-for-a-project-s-own-workflow-file-that-no-code-enforces-yet.md, design/issues/summary.md, flai/cmd/guard.go, flai/internal/check/check_roles.go, flai/internal/check/roles_test.go, flai/internal/guard/roles.go]
after: [S-0741]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5424
  priced_by: harness
  turns:
    - day: 2026-11-15
      ceremony: 1
      hand_edits: 1
      work: 69
  models:
    - model: claude-haiku-5-5
      input: 34
      output: 24870
      cache_read: 991725
      cache_write: 98044
      cost: 0.0346
    - model: claude-opus-5-5
      input: 630
      output: 215367
      cache_read: 44413000
      cache_write: 1303149
      cost: 20.8616
  strategic:
    - kind: orchestrator
      seconds: 99
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 21
          output: 111
          cache_read: 760184
          cache_write: 93160
          cost: 0.2556
        - model: claude-sonnet-5-5
          input: 4
          output: 27
          cache_read: 13864
          cache_write: 36378
          cost: 0.0398
cost_of_delay:
  value: 83.73
  by: system-flow
  at: 2026-11-14T17:44:15Z
forecast:
  duration: 47m
  delivery: 2026-11-15T22:42:00Z
  basis: "Its own forecast of 47m; 1st in the pull order with an in-progress limit of 3, behind S-0750 and S-0790."
  by: flai
  at: 2026-11-15T21:46:10Z
finalized:
  by: orchestrator
  at: 2026-11-14T17:45:58Z
---
# S-0755 flai reads a project's own workflow file in place of the baseline, the template ships the baseline, flai upgrade keeps the project's file, and the guides say how to change the board

## Goal

An operator changes the board by editing a file in their project. flai reads the project's workflow file under `design/workflows/` when it has one and the baseline when it has none, `flai check` refuses a project file that fails the workflow's findings before anything acts on it, the template ships the baseline file, `flai upgrade` keeps a project's changes to it, and the operators' and users' guides say how to add a state, a rule, or a timer.

## Acceptance criteria

- [x] flai reads the project's workflow file at the path the ADR fixes when it exists, and the embedded baseline when it does not, in every command, the MCP server, and flai serve, and a file that does not parse or fails a `workflow.` finding stops each with the finding, falling back to nothing
- [x] The template ships the baseline workflow file, byte for byte the copy flai embeds, held so by a test as `roles/baseline_test.go` holds the role files
- [x] `flai upgrade` keeps a project's own changes to the file, as the ADR says, and adds the file to a project that lacks it only where the ADR says so
- [x] `flai check` advises a project with no workflow file, as `roles.baseline` advises one with no role file, if the ADR asks for it
- [x] A runbook under `docs/operators/runbooks/`, listed in `docs/operators/index.md`, says how to change the board: adding a state, a rule, and a timer, renaming or removing a state with items in it, and what to check before the story that changes it is accepted
- [x] `docs/users/flai.md`, `docs/users/conventions.md`, `design/system/workflow.md`, `design/system/template.md`, and `design/workflows/README.md` here and in the template say where the file is and how flai reads it, and `template/CHANGELOG.md` records it

## Tasks

- T-3253 flai reads the project's own workflow file when it has one and the baseline when it has none, everywhere a repository is opened
- T-3259 flai check reports a project's workflow file it refuses, judges the items against the file it runs, and gives no advisory to a project without one
- T-3260 The template ships the baseline workflow file byte for byte as flai embeds it, and flai upgrade keeps a project's changes to it
- T-3262 A runbook says how to change the board, and the guides and design say where the workflow file is and how flai reads it
- T-3390 The MCP server, flai serve, and the commands that still read the baseline run the project's workflow
- T-3391 flai guard and the role files' check judge a to condition against the project's workflow file

## Notes

- The path `design/workflows/board.md` is assumed, as on S-0737; the story's agent follows the ADR's name.
- The file is code the operator runs: it changes only through a story the operator accepts, and no role's `files` reaches it (ADR-0207 §2 and §7).
- ADR-0307 settles what ADR-0301 left open. The file read is the opened checkout's. A file flai cannot run refuses every open but `flai check`'s. The MCP server starts on a broken file and refuses each call until it is fixed. Upgrade adds the file only when it is absent. No advisory is given for a project without one, which is how criterion 4 is met.
- Criterion 1 does not reach `flai stats`. Its figures read the workflow through `metrics.Options.Workflow`, which S-0750 adds, so wiring `repo.Flow()` into it is left to S-0750, by agreement on MS-0583, or to S-0752. The schedulers of `flai serve` stop on a refused open without logging it. The passes that go through `openRepo`, the replanner, and the role events do log it.
- I-0380 lists the ADR-0301 §6 and §7 rules for a project's file that no code applies yet. They are operator-only acceptance of a change to the file, the dashboard's refusal to save it on main, `board.wip-state`, and limits on a state the project adds.
- S-0750 changes `flai/internal/workflow/baseline/board.md`, and the template's copy must stay equal to it byte for byte. Whichever of S-0750 and S-0755 is accepted second updates `template/root/design/workflows/board.md`, since `TestBaselineIsTheTemplatesBoard` holds the two equal.

### Planning

- Touches, all files:
  - `flai/internal/workflow/load.go`, `load_test.go`, `flai/internal/workitem/repo.go`: layout. The load and where a repository is opened.
  - `flai/internal/check/check_workflow.go`, `workflow_test.go`: layout, the check S-0741 adds.
  - `template/root/design/workflows/board.md`, `flai/internal/workflow/baseline_test.go`, `flai/internal/upgrade/upgrade.go`, `upgrade_test.go`: layout. `Compute` is in `upgrade.go`; the byte test mirrors `roles/baseline_test.go`.
  - Co-change: `template/template.yaml` (8% in `flai touches suggest`) and `template/CHANGELOG.md`, which a template change bumps and records.
  - Design (criteria 5 and 6): the runbook, a new file `docs/operators/runbooks/change-the-board.md` beside the other runbooks, with the runbooks' README, and the guides and design documents criterion 6 names.
  - Left out from `flai touches suggest`: `flai-cli.md` (63%, shared) and `docs/operators/settings.md` (16%), since the file is no setting.
- `after`: S-0741. A project's own file is checked by the `workflow.` findings S-0741 adds before anything acts on it, and the first criterion stops on them.
- Tasks: T-3253 (the load) first; then T-3259 (the check on the project's file) and T-3260 (template and upgrade) together; then T-3262 (the runbook and documents), which describes the file as T-3260 ships it.
- Forecast 47m, `flai forecast`'s rate for size 26, up from its 42m for the first draft's 23 once the runbook and the template version were added.
- Cost of delay: 83.73 USD a week, as `flai cod` gives it: this story's share of E-0034's 1500.00 USD a week, the operator's `time_lost_per_cycle: 10h` (TH-1115) at 150 USD an hour, by its forecast of 47m over the 14h2m of the epic's 19 open stories. Kept: the epic's outcome needs every story, so a share by size is the split that fits.
