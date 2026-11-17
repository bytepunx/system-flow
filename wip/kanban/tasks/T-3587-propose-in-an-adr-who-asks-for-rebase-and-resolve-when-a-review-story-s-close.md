---
id: T-3587
type: task
nature: improvement
title: Propose in an ADR who asks for rebase and resolve when a review story's close-out stops after main moved, and write it into the design
status: backlog
parent: S-0845
owner: alex
created: 2026-11-17T12:33:56Z
updated: 2026-11-17T12:33:56Z
transitions: []
stream: S-0845
tags: [flai, orchestrator]
touches: [design/adrs/README.md, design/system/strategic-agents.md, design/system/workflow.md]
---
# T-3587 Propose in an ADR who asks for rebase and resolve when a review story's close-out stops after main moved, and write it into the design

## Work

Read I-0386's instances in `design/analysis/2026-11-16-all.md` § Finding 1 (S-0739 `undone_main`, S-0742 stopped at vet, S-0747 stopped at sync) against ADR-0313, which made rebase and resolve a flai serve job that sends a story back only on conflicts or a failed step, but left the request the operator's alone (`flai/cmd/accept_resolve.go`, `resolveFlagsRefused`). Propose the fix in a new ADR with `flai adr new`, refining ADR-0313. The plan assumes the orchestrator, under `orchestration.permissions.accept_reviews`, asks flai serve for rebase and resolve through a new MCP tool when a story's last close-out did not pass and main has moved past its base, in place of a thread that awaits the operator; weigh against it a close-out job of a review story that sends back on a stop. Say what stays the operator's: a story held by `undone_main` (ADR-0155) is not sent back by this, since only the operator restores or keeps a removal. Ask the operator on a thread on the story to accept the ADR before the code tasks start. Write the decision into `design/system/strategic-agents.md` § Accepting a story (S-0221) and `design/system/workflow.md` § Transitions and who makes them and § What waits on an acceptance (S-0296).

The first layer: every other task builds what this ADR decides.

## Done when

- The ADR is accepted by the operator and names who asks for rebase and resolve, when, and what stays the operator's.
- `design/system/strategic-agents.md` and `design/system/workflow.md` say it, linking the ADR.
- `flai check --strict` passes.

## Notes
