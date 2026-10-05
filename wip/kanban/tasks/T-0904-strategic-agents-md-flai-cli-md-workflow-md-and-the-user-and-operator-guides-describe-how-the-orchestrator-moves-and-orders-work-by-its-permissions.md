---
id: T-0904
type: task
nature: feature
title: strategic-agents.md, flai-cli.md, workflow.md, and the user and operator guides describe how the orchestrator moves and orders work by its permissions
status: backlog
parent: S-0219
owner: alex
created: 2026-10-05T04:47:42Z
updated: 2026-10-05T04:47:42Z
transitions: []
stream: S-0219
tags: [flai, docs]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators]
after: [T-0899, T-0901]
---
# T-0904 strategic-agents.md, flai-cli.md, workflow.md, and the user and operator guides describe how the orchestrator moves and orders work by its permissions

## Work

Describe what this story built, where a reader looks for it:

- `design/system/strategic-agents.md`, under the orchestrator S-0218 describes: what it does with each of `plan_backlog_epics`, `finalize_drafts`, `promote_to_ready`, and `order_ready`; the commands it reads; the refusals `flai move`, `flai edit`, and `plan` give it; the planner runs it starts, triggered `orchestrator`; and how each action is logged with its policy figure.
- `design/system/flai-cli.md`: `flai plan --candidates`, `flai promote --drafts`, `flai order --by --apply`'s window for hand placements, the placement record, and the guard's rows for the orchestrator's four permissions.
- `design/system/workflow.md`, under "The pull order": that `flai order` records who placed a story and when, linking T-0884's ADR, and that the orchestrator keeps a hand placement of the last day.
- `docs/users/flai.md` and `docs/users/flai-reference.md` (regenerated with `make flai-reference`): the new commands and flags.
- The operator's guide S-0218 writes in `docs/operators`: what turning each of the four permissions on lets the orchestrator do, and what it never does.

This task waits for T-0899 and T-0901, so that it describes what the prompt says and the tests pin.

## Done when

- each document above says what the code does, and links the ADR and the stories
- `flai check --strict` and the markdown lint pass on the documents

## Notes
