---
id: S-0755
type: story
nature: feature
title: flai reads a project's own workflow file in place of the baseline, the template ships the baseline, flai upgrade keeps the project's file, and the guides say how to change the board
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:02:15Z
updated: 2026-11-12T16:02:15Z
transitions: []
tags: [cli, template]
topics: [cli, template, conventions]
touches: [flai/internal/workflow/load.go, flai/internal/workflow/load_test.go, flai/internal/workflow/baseline_test.go, flai/internal/workitem/repo.go, flai/internal/check/check_workflow.go, flai/internal/check/workflow_test.go, flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, template/root/design/workflows/board.md, template/root/design/workflows/README.md, template/CHANGELOG.md, design/workflows/README.md, design/system/workflow.md, design/system/template.md, docs/operators/index.md, docs/users/flai.md, docs/users/conventions.md]
after: [S-0741]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0755 flai reads a project's own workflow file in place of the baseline, the template ships the baseline, flai upgrade keeps the project's file, and the guides say how to change the board

## Goal

An operator changes the board by editing a file in their project. flai reads the project's workflow file under `design/workflows/` when it has one and the baseline when it has none, `flai check` refuses a project file that fails the workflow's findings before anything acts on it, the template ships the baseline file, `flai upgrade` keeps a project's changes to it, and the operators' and users' guides say how to add a state, a rule, or a timer.

## Acceptance criteria

- [ ] flai reads the project's workflow file at the path the ADR fixes when it exists, and the embedded baseline when it does not, in every command, the MCP server, and flai serve, and a file that does not parse or fails a `workflow.` finding stops each with the finding, falling back to nothing
- [ ] The template ships the baseline workflow file, byte for byte the copy flai embeds, held so by a test as `roles/baseline_test.go` holds the role files
- [ ] `flai upgrade` keeps a project's own changes to the file, as the ADR says, and adds the file to a project that lacks it only where the ADR says so
- [ ] `flai check` advises a project with no workflow file, as `roles.baseline` advises one with no role file, if the ADR asks for it
- [ ] `docs/operators/index.md` gains a runbook for changing the board: adding a state, a rule, and a timer, renaming or removing a state with items in it, and what to check before the story that changes it is accepted
- [ ] `docs/users/flai.md`, `docs/users/conventions.md`, `design/system/workflow.md`, `design/system/template.md`, and `design/workflows/README.md` here and in the template say where the file is and how flai reads it, and `template/CHANGELOG.md` records it

## Tasks

## Notes

- The path `design/workflows/board.md` is assumed, as on S-0737; the story's agent follows the ADR's name.
- The file is code the operator runs: it changes only through a story the operator accepts, and no role's `files` reaches it (ADR-0207 §2 and §7).
