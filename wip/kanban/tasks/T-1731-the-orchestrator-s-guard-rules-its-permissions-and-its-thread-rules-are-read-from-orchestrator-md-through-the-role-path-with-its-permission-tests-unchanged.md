---
id: T-1731
type: task
nature: feature
title: The orchestrator's guard rules, its permissions and its thread rules, are read from orchestrator.md through the role path, with its permission tests unchanged
status: backlog
parent: S-0416
owner: alex
created: 2026-10-16T05:44:15Z
updated: 2026-10-16T05:44:15Z
transitions: []
stream: S-0416
tags: [cli, template]
touches: [flai/internal/guard/guard.go, design/roles/orchestrator.md, template/root/design/roles/orchestrator.md]
after: [T-1730]
---
# T-1731 The orchestrator's guard rules, its permissions and its thread rules, are read from orchestrator.md through the role path, with its permission tests unchanged

## Work

Route `orchestrate` through the role path. The orchestrator's file lists each permission with the calls it allows, as `design/system/strategic-agents.md` § Its permissions and the guard tabulates them. Keep in Go, as named conditions, the rules a list cannot state:

- `finalizesDraft`, and cost of delay inputs alone;
- the thread rules by opener and `answer_threads`;
- `--by orchestrator` on `accept`;
- the refusal of a placement by hand.

Keep `AppendRefusal`'s entry and the refusal texts word for word.

Waits for T-1730, which changes `guard.go` first.

## Done when

`flai/internal/guard/orchestrate_test.go`, `thread_test.go`, `flai/cmd/guard_permissions_change_test.go`, and `flai/cmd/orchestrate_permissions_test.go` pass unchanged.

## Notes
