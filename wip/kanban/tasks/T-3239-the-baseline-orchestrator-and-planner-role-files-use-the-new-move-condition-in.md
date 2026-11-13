---
id: T-3239
type: task
nature: feature
title: The baseline orchestrator and planner role files use the new move condition in the template, in flai's embedded copies, and in this project
status: backlog
parent: S-0748
owner: alex
created: 2026-11-12T16:08:51Z
updated: 2026-11-12T16:08:51Z
transitions: []
stream: S-0748
tags: [cli, template]
touches: [template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/planner.md, design/workflows/roles/orchestrator.md, design/workflows/roles/planner.md, template/CHANGELOG.md, template/template.yaml]
after: [T-3231]
---
# T-3239 The baseline orchestrator and planner role files use the new move condition in the template, in flai's embedded copies, and in this project

## Work

Waits for T-3231: the files name the condition that task adds, and `flai check` must accept them. Replace `to-review` and `to-backlog` in the orchestrator's and the planner's role files with the new condition naming `review` and `backlog`, in the template, in `flai/internal/roles/baseline/` byte for byte, and in this project's `design/workflows/roles/`, keeping each project's text below the marker. Record the change in `template/CHANGELOG.md` and bump `template/template.yaml` as the template's versioning asks.

## Done when

`baseline_test.go` passes, `flai check --strict` passes, and the triggers and grants are as they were.

## Notes
