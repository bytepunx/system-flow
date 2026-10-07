---
id: T-1169
type: task
nature: feature
title: The manifest gains orchestration.permissions.plan_backlog_stories, off by default, in the settings catalog and the operators' index
status: cancelled
parent: S-0328
owner: alex
created: 2026-10-07T19:39:27Z
updated: 2026-10-07T19:42:15Z
transitions:
  - to: cancelled
    at: 2026-10-07T19:42:15Z
    by: agent-S-0328
stream: S-0328
tags: [flai, manifest, orchestrator]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, docs/operators/settings.md]
after: [T-1168]
---
# T-1169 The manifest gains orchestration.permissions.plan_backlog_stories, off by default, in the settings catalog and the operators' index

## Work

Add `PlanBacklogStories` to the orchestrator's `Permissions` in `flai/internal/manifest/manifest.go`, with the constant `PermitPlanBacklogStories = "plan_backlog_stories"`, in `PermissionNames` after `plan_backlog_epics`, and in `Allows`. Give it its meaning and risk in `permissionText` in `settings.go`, so the settings catalog, and with it the dashboard's Orchestrator settings panel, offers it. Add its row to `docs/operators/settings.md`, which the settings test requires of every manifest key.

It waits for T-1168, whose ADR names the permission and what it allows.

## Done when

- `flai check` accepts `plan_backlog_stories: true` and reports a bad value as `manifest.orchestration`.
- `flai manifest set orchestration.permissions.plan_backlog_stories true` writes it, and the settings catalog lists it with its meaning and risk.
- The manifest and settings tests cover it, and `flai test` passes on the paths changed.

## Notes

The project's own `system-flow.yaml` is not changed: turning the permission on is the operator's.
- 2026-10-07T19:42:15Z: moved to cancelled: duplicates the story agent's plan (T-1174, T-1176 to T-1179), written at start while the planner ran; its refinements are folded into T-1176 to T-1178
