---
id: S-0755
type: story
nature: feature
title: flai reads a project's own workflow file in place of the baseline, the template ships the baseline, flai upgrade keeps the project's file, and the guides say how to change the board
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:02:15Z
updated: 2026-11-13T15:47:58Z
transitions: []
tags: [cli, template]
topics: [cli, template, conventions]
touches: [flai/internal/workflow/load.go, flai/internal/workflow/load_test.go, flai/internal/workitem/repo.go, flai/internal/check/check_workflow.go, flai/internal/check/workflow_test.go, template/root/design/workflows/board.md, flai/internal/workflow/baseline_test.go, flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, template/CHANGELOG.md, template/template.yaml, docs/operators/runbooks/change-the-board.md, docs/operators/runbooks/README.md, docs/operators/index.md, docs/users/flai.md, docs/users/conventions.md, design/system/workflow.md, design/system/template.md, design/workflows/README.md, template/root/design/workflows/README.md]
after: [S-0741]
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
      seconds: 10
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 14
          cache_read: 121381
          cache_write: 20507
          cost: 0.0405
draft: true
forecast:
  duration: 47m
  delivery: 2026-11-13T22:30:00Z
  basis: "Its own forecast of 47m; 22nd in the pull order with an in-progress limit of 3, behind S-0751, S-0760, S-0759, S-0757, S-0241, S-0730, S-0736, S-0737, S-0738, S-0739, S-0741, S-0742, S-0743, S-0744, S-0746, S-0747, S-0748, S-0749, S-0750, S-0752, S-0753 and S-0754."
  by: flai
  at: 2026-11-13T15:47:58Z
---
# S-0755 flai reads a project's own workflow file in place of the baseline, the template ships the baseline, flai upgrade keeps the project's file, and the guides say how to change the board

## Goal

An operator changes the board by editing a file in their project. flai reads the project's workflow file under `design/workflows/` when it has one and the baseline when it has none, `flai check` refuses a project file that fails the workflow's findings before anything acts on it, the template ships the baseline file, `flai upgrade` keeps a project's changes to it, and the operators' and users' guides say how to add a state, a rule, or a timer.

## Acceptance criteria

- [ ] flai reads the project's workflow file at the path the ADR fixes when it exists, and the embedded baseline when it does not, in every command, the MCP server, and flai serve, and a file that does not parse or fails a `workflow.` finding stops each with the finding, falling back to nothing
- [ ] The template ships the baseline workflow file, byte for byte the copy flai embeds, held so by a test as `roles/baseline_test.go` holds the role files
- [ ] `flai upgrade` keeps a project's own changes to the file, as the ADR says, and adds the file to a project that lacks it only where the ADR says so
- [ ] `flai check` advises a project with no workflow file, as `roles.baseline` advises one with no role file, if the ADR asks for it
- [ ] A runbook under `docs/operators/runbooks/`, listed in `docs/operators/index.md`, says how to change the board: adding a state, a rule, and a timer, renaming or removing a state with items in it, and what to check before the story that changes it is accepted
- [ ] `docs/users/flai.md`, `docs/users/conventions.md`, `design/system/workflow.md`, `design/system/template.md`, and `design/workflows/README.md` here and in the template say where the file is and how flai reads it, and `template/CHANGELOG.md` records it

## Tasks

- T-3253 flai reads the project's own workflow file when it has one and the baseline when it has none, everywhere a repository is opened
- T-3259 flai check runs the workflow findings on the project's own file and advises a project with none, if the ADR asks
- T-3260 The template ships the baseline workflow file byte for byte as flai embeds it, and flai upgrade keeps a project's changes to it
- T-3262 A runbook says how to change the board, and the guides and design say where the workflow file is and how flai reads it

## Notes

- The path `design/workflows/board.md` is assumed, as on S-0737; the story's agent follows the ADR's name.
- The file is code the operator runs: it changes only through a story the operator accepts, and no role's `files` reaches it (ADR-0207 §2 and §7).

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
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
