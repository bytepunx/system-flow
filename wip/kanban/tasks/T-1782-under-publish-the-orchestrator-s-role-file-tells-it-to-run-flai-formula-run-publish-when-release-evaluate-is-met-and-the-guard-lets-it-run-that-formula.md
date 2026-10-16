---
id: T-1782
type: task
nature: improvement
title: Under publish, the orchestrator's role file tells it to run flai formula run publish when release_evaluate is met, and the guard lets it run that formula
status: backlog
parent: S-0423
owner: alex
created: 2026-10-16T05:48:37Z
updated: 2026-10-16T05:48:37Z
transitions: []
stream: S-0423
tags: [cli, template]
touches: [design/roles/orchestrator.md, template/root/design/roles/orchestrator.md, flai/cmd/orchestrate_permissions_test.go]
after: [T-1780]
---
# T-1782 Under publish, the orchestrator's role file tells it to run flai formula run publish when release_evaluate is met, and the guard lets it run that formula

## Work

In `orchestrator.md`, template and here:

- list `formula run publish` under the `publish` permission;
- tell the orchestrator to run it when `release_evaluate` is met, and to log the instance's ID with `activity_log`;
- keep `release_publish` for a project with no `publish` formula.

Add a case to `orchestrate_permissions_test.go`: the command passes with `publish` on and is refused with it off, and `formula run` of any other formula is refused.

Waits for the formula file. Runs beside the publish step task, since they share no path.

## Done when

`go test ./cmd/...` passes. The existing orchestrator permission cases are unchanged.

## Notes
