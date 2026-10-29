---
id: T-2338
type: task
nature: improvement
title: The two .claude/agents/orchestrator.md definitions tell the orchestrator to prime with the reasons its start prompt names, as its role file says
status: backlog
parent: S-0537
owner: alex
created: 2026-10-29T15:24:29Z
updated: 2026-10-29T15:24:45Z
transitions: []
stream: S-0537
tags: []
touches: [".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md]
after: [T-2313]
---
# T-2338 The two .claude/agents/orchestrator.md definitions tell the orchestrator to prime with the reasons its start prompt names, as its role file says

## Work

Criterion 3. Until S-0419 generates the definitions from the role files, a change to the orchestrator's prompt is made in `.claude/agents/orchestrator.md` and its template copy as well as in the role file (D-0411). T-2313 changes the role file's step 1 to prime with the reasons the start prompt names. This task makes the same change in both definitions, word for word.

The two files are protected. The story's agent writes them itself, after T-2313's layer is back, through `permission_prompt`, which asks the operator on a thread. No sub-agent writes them. The task is its own so that the rest can close while the write waits (D-0411).

Waits for T-2313, whose role-file words it copies.

## Done when

- Both definitions say what the role file's step 1 says, and `flai check --strict` gives no `roles.definition` finding.

## Notes
