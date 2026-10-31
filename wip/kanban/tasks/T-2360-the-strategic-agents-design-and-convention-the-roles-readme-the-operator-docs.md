---
id: T-2360
type: task
nature: improvement
title: The strategic-agents design and convention, the roles README, the operator docs, flai-cli.md, and the changelog say a built-in role's start prompt and triggers are read from its file
status: backlog
parent: S-0554
owner: alex
created: 2026-10-30T20:49:25Z
updated: 2026-10-30T20:49:25Z
transitions: []
stream: S-0554
tags: [template]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/operators/index.md, template/CHANGELOG.md]
after: [T-2357, T-2359]
---
# T-2360 The strategic-agents design and convention, the roles README, the operator docs, flai-cli.md, and the changelog say a built-in role's start prompt and triggers are read from its file

## Work

Rewrite the parts of `design/system/strategic-agents.md` that name `planPrompt`, `orchestratePrompt`, `analyzePrompt`, and `curatePrompt`, the replanner's triggers, and the analysis and digest schedules, to say each is read from the role file. Do the same in `design/system/flai-cli.md` where it says when `flai serve` starts each role. Say in the strategic-agents convention and the roles README, here and in the template, how a project changes a built-in role's start prompt and triggers. Update `docs/operators/index.md` and add the change to `template/CHANGELOG.md`.

Fifth layer: it waits for the start prompt and the trigger tasks, whose behaviour it describes.

## Done when

No document here or in the template says a built-in role's start prompt or triggers are fixed in Go. `flai check --strict` and the markdown lint report nothing.

## Notes

Drafted by the planner of E-0024 on 2026-10-30.
