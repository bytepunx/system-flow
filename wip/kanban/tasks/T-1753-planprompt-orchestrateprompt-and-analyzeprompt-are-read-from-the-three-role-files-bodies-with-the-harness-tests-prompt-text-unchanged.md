---
id: T-1753
type: task
nature: improvement
title: planPrompt, orchestratePrompt, and analyzePrompt are read from the three role files' bodies, with the harness tests' prompt text unchanged
status: backlog
parent: S-0419
owner: alex
created: 2026-10-16T05:46:02Z
updated: 2026-10-21T12:19:48Z
transitions: []
stream: S-0419
tags: [cli, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/workflows/roles/planner.md, design/workflows/roles/orchestrator.md, design/workflows/roles/analyzer.md, template/root/design/workflows/roles/planner.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/analyzer.md, flai/internal/roles/baseline/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/analyzer.md]
after: [T-1754]
---
# T-1753 planPrompt, orchestratePrompt, and analyzePrompt are read from the three role files' bodies, with the harness tests' prompt text unchanged

## Work

- Remove the prompt text from `planPrompt`, `orchestratePrompt`, and `analyzePrompt` in `flai/internal/harness/harness.go`.
- Read it from the role files' bodies through the defined-role prompt S-0417 built, filling the item, the focus, and the permission text as today.
- Bring the bodies of the three role files, here, in the template, and in the embedded copy under `flai/internal/roles/baseline/`, to the exact text the functions produce now. The first story wrote them from the same source; this task makes them the source.

Second layer: waits for the embedded baseline, so that a prompt test in a root with no role files still reads the bodies. The serve paths after it read the prompt it moves.

## Done when

`flai/internal/harness/harness_test.go` passes unchanged, the prompt assertions of `TestClaudeCodeStartsThePlanner`, `TheOrchestrator`, and `TheAnalyzer` included.

## Notes

- Rewritten by planner-S-0419 for the new order, as the orchestrator approved on TH-0579.
