---
id: T-1753
type: task
nature: improvement
title: planPrompt, orchestratePrompt, and analyzePrompt are read from the three role files' bodies, with the harness tests' prompt text unchanged
status: backlog
parent: S-0419
owner: alex
created: 2026-10-16T05:46:02Z
updated: 2026-10-16T05:50:22Z
transitions: []
stream: S-0419
tags: [cli, template]
touches: [flai/internal/harness/harness.go, design/roles/planner.md, design/roles/orchestrator.md, design/roles/analyzer.md, template/root/design/roles/planner.md, template/root/design/roles/orchestrator.md, template/root/design/roles/analyzer.md, flai/internal/harness/harness_test.go]
---
# T-1753 planPrompt, orchestratePrompt, and analyzePrompt are read from the three role files' bodies, with the harness tests' prompt text unchanged

## Work

- Remove the prompt text from `planPrompt`, `orchestratePrompt`, and `analyzePrompt` in `flai/internal/harness/harness.go`.
- Read it from the role files' bodies through the defined-role prompt S-0417 built, filling the item, the focus, and the permission text as today.
- Bring the bodies of the three role files, here and in the template, to the exact text the functions produce now. The first story wrote them from the same source; this task makes them the source.

First layer: the serve paths below read the prompt it moves.

## Done when

`flai/internal/harness/harness_test.go` passes unchanged, the prompt assertions of `TestClaudeCodeStartsThePlanner`, `TheOrchestrator`, and `TheAnalyzer` included.

## Notes
