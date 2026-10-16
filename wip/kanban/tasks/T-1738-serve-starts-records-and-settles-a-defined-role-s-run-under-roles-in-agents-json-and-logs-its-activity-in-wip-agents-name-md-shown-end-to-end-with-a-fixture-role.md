---
id: T-1738
type: task
nature: feature
title: serve starts, records, and settles a defined role's run under roles in agents.json and logs its activity in wip/agents/name.md, shown end to end with a fixture role
status: backlog
parent: S-0417
owner: alex
created: 2026-10-16T05:45:01Z
updated: 2026-10-16T05:45:01Z
transitions: []
stream: S-0417
tags: [cli]
touches: [flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/agents.go, flai/internal/serve/activity.go]
after: [T-1735, T-1736, T-1737]
---
# T-1738 serve starts, records, and settles a defined role's run under roles in agents.json and logs its activity in wip/agents/name.md, shown end to end with a fixture role

## Work

Add `flai/internal/serve/role.go`, modelled on `plan.go`:

- `RunRole(name, item, inputs)` checks the host action the file names, the item kind and state, one run per role and item, and a harness or command.
- It starts the run through the harness with `RoleAgent`, and records it in `AgentState.Roles`, by name and item, with its log `<key>-<name>-<start>.log`.
- It settles the run `asked`, `failed`, or `worked`, as `planEnded` does, and logs the activity with `LogRunEnd` in `wip/agents/<name>.md`, charging nothing.
- It journals each start, failure, and end under the role's action and notifies the dashboard with the role.
- Test the fixture role `reviewer` end to end with a fake harness: its environment, prompt, record, activity entry, and each refusal.

Waits for T-1735, T-1736, and T-1737, whose agent, request, and activity kinds it uses.

## Done when

`go test ./internal/serve/...` passes. The existing planner, orchestrator, and analyzer tests are unchanged.

## Notes
