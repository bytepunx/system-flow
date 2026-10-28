---
id: D-0412
title: "Plan for S-0535: nine tasks in five layers, an ADR for the trigger vocabulary first: accept recommendation and resolve"
kind: recommendation
sources: [TH-0715]
topics: [orchestration, cli, mcp, template, dashboard, flai, go, conventions, code, all]
paths: [flai/internal/serve, flai/internal/mcpserver, design/system/flai-cli.md, design/workflows/roles/README.md, design/adrs, design/adrs/README.md, flai/internal/inbox/events.go, flai/internal/inbox/events_test.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/inbox/threadevents.go, flai/internal/inbox/threadevents_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/check/roles_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/readyempty.go, flai/internal/serve/readyempty_test.go, flai/internal/serve/agents.go, flai/cmd/accept.go, flai/cmd/accept_release_due_test.go, flai/internal/release/evaluate.go, flai/internal/release/evaluate_test.go, template/root/design/workflows/roles/README.md, template/CHANGELOG.md, design/system/agent-narrative.md, docs/users/flai.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0412 Plan for S-0535: nine tasks in five layers, an ADR for the trigger vocabulary first: accept recommendation and resolve

The operator, alex, answered TH-0715 at 2026-10-28T22:41:35Z.

What was asked:

```text
S-0535 is planned: nine tasks in five layers, its touches, a 4h forecast, and a cost of delay of 37.50 USD a week (from the 15m the orchestrator set on TH-0714).
```

The recommendation it overruled:

```text
**Recommended:** yes to both.

1. **Fold `to-backlog` as a condition on `moved` into T-2293.** It is one line beside `to-review`. ADR-0209 §1's own trigger, `{event: moved, when: to-backlog}`, is refused as `roles.file` without it, so S-0536, the next story of E-0030, could not pass `flai check`. Leaving it out only moves the same line into S-0536.
2. **Drop `flai/internal/workitem` from S-0535's touches.** No task changes a file under it. While S-0535 is in progress, that folder touch holds every ready story touching a workitem file, which costs pull time for nothing (ADR-0096). The story's agent widens its touches by itself, through `flai task done`, if a file there turns out to change.

**Or:** say no to either here. The planner then leaves `to-backlog` to S-0536, or keeps the touch.

Source: design/adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md
```

The operator's answer:

```text
accept recommendation and resolve
```

Source: TH-0715.
