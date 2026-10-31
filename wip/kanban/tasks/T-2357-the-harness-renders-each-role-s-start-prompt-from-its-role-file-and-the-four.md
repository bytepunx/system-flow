---
id: T-2357
type: task
nature: feature
title: The harness renders each role's start prompt from its role file, and the four baseline role files carry today's start prompts
status: backlog
parent: S-0554
owner: alex
created: 2026-10-30T20:49:01Z
updated: 2026-10-30T20:49:01Z
transitions: []
stream: S-0554
tags: [cli, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/roles/baseline/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/analyzer.md, flai/internal/roles/baseline/curator.md, design/workflows/roles/planner.md, design/workflows/roles/orchestrator.md, design/workflows/roles/analyzer.md, design/workflows/roles/curator.md, template/root/design/workflows/roles/planner.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/analyzer.md, template/root/design/workflows/roles/curator.md]
after: [T-2355]
---
# T-2357 The harness renders each role's start prompt from its role file, and the four baseline role files carry today's start prompts

## Work

Move the text of `planPrompt`, `orchestratePrompt`, `analyzePrompt`, and `curatePrompt` in `flai/internal/harness/harness.go` into the start prompt of the four baseline role files, embedded in flai, here, and in the template, keeping the three copies equal. Their branches become the template's conditionals over the run's fields. `Prompt` renders the role file's start prompt for every role. `definedPrompt` becomes the start prompt a defined role gets when its file has none. Delete the Go text that the files now carry.

Third layer: it waits for the start prompt task, whose renderer it calls. It shares no path with the planner's triggers task.

## Done when

`TestClaudeCodeStartsThePlanner`, `TestClaudeCodeStartsTheOrchestrator`, `TestClaudeCodeStartsTheAnalyzer`, `TestClaudeCodeStartsTheCurator`, and `TestClaudeCodeStartsADefinedRole` pass with the same start prompt text, for every branch they cover. A new case shows a project's role file whose additions change the planner's start prompt. `flai check --strict` reports nothing, and the baseline test finds the three copies equal.

## Notes

Drafted by the planner of E-0024 on 2026-10-30. In S-0419's worktree that day, `Prompt` is at `harness.go:444`, and the prompts at `:591`, `:712`, `:806`, `:852`, and `:886`.
