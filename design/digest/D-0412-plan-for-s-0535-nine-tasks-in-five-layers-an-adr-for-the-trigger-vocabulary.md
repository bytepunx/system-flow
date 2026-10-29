---
id: D-0412
title: "Plan for S-0535: nine tasks in five layers, an ADR for the trigger vocabulary first: accept recommendation and resolve"
kind: recommendation
sources: [TH-0715, S-0535]
topics: [orchestration, cli, mcp, template, dashboard, flai, go, conventions, code, all]
paths: [flai/internal/serve, flai/internal/mcpserver, design/system/flai-cli.md, design/workflows/roles/README.md, design/adrs, design/adrs/README.md, flai/internal/inbox/events.go, flai/internal/inbox/events_test.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/inbox/threadevents.go, flai/internal/inbox/threadevents_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/check/roles_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/readyempty.go, flai/internal/serve/readyempty_test.go, flai/internal/serve/agents.go, flai/cmd/accept.go, flai/cmd/accept_release_due_test.go, flai/internal/release/evaluate.go, flai/internal/release/evaluate_test.go, template/root/design/workflows/roles/README.md, template/CHANGELOG.md, design/system/agent-narrative.md, docs/users/flai.md]
date: 2026-10-28
status: absorbed
absorbed_by: "design/adrs/0211-the-orchestrator-s-trigger-vocabulary-gains-the-events-moved-thread-planned.md#decision"
asked_again: []
cited: []
---

# D-0412 Plan for S-0535: nine tasks in five layers, an ADR for the trigger vocabulary first: accept recommendation and resolve

The operator, alex, answered TH-0715 at 2026-10-28T22:41:35Z, the plan thread of S-0535 (nine tasks in five layers, an ADR for the trigger vocabulary first), with "accept recommendation and resolve". They took the orchestrator's recommendation whole, on both proposals:

1. **`to-backlog` qualifies `{event: moved}` as a condition, added in S-0535 beside `to-review` (T-2293).** Without it, ADR-0209 §1's own trigger, `{event: moved, when: to-backlog}`, is refused by `flai check` as `roles.file`, so S-0536 could not pass. It had been only a call condition.
2. **A declared folder touch no task changes a file under is dropped from the story.** `flai/internal/workitem` left S-0535's touches: while the story is in progress such a touch holds every ready story touching a file there, costing pull time for nothing (ADR-0096). The story's agent widens its touches through `flai task done` if a file there turns out to change.

The plan's other assumptions stood with it: the vocabulary gets an ADR of its own since ADR-0209 is accepted, keeping ADR-0207's `review` and `thread-opened` beside `moved` and `thread`; `release-due` is an event of its own, not `{event: accepted, when: release-due}`; `planned` fires only when the run wrote something beyond its item.

One assumption did not stand. The plan had `ready-empty` fire when the ready column is empty. S-0535's agent followed criterion 2's wording instead: it fires when the ready count, held stories counted, falls below the in-progress limit's room, and T-2295 was rewritten to match. S-0535's narrative records that it took this plan with the orchestrator's recommendation. [ADR-0211](../adrs/0211-the-orchestrator-s-trigger-vocabulary-gains-the-events-moved-thread-planned.md) now states the vocabulary as built.

Source: TH-0715 and S-0535's narrative (`wip/archive/agents/S-0535.md`), citing [ADR-0096](../adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md).
