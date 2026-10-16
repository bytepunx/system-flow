---
id: T-1754
type: task
nature: improvement
title: A project with no design/roles runs the three roles from baseline files embedded in flai, and flai check warns roles.baseline naming flai upgrade
status: backlog
parent: S-0419
owner: alex
created: 2026-10-16T05:46:09Z
updated: 2026-10-16T05:46:09Z
transitions: []
stream: S-0419
tags: [cli]
touches: [flai/internal/roles/baseline.go, flai/internal/roles/baseline_test.go, flai/internal/check/check.go, flai/internal/check/roles_test.go]
after: [T-1753]
---
# T-1754 A project with no design/roles runs the three roles from baseline files embedded in flai, and flai check warns roles.baseline naming flai upgrade

## Work

- Embed the template's three role files in flai with `go:embed`, from a copy the build keeps in step with `template/root/design/roles/`. A test fails when they differ.
- `roles.Load` falls back to them for a role a project's folder lacks.
- `flai check` warns `roles.baseline` once per missing built-in role, naming `flai upgrade`.

Waits for the prompt task, so that the embedded bodies are the prompt text.

## Done when

`go test ./internal/roles/... ./internal/check/...` passes, including a project with no roles folder whose planner, orchestrator, and analyzer still load.

## Notes
