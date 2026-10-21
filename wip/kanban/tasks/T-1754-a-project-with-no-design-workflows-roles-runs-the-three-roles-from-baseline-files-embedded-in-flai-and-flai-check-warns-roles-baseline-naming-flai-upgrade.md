---
id: T-1754
type: task
nature: improvement
title: A project with no design/workflows/roles runs the three roles from baseline files embedded in flai, and flai check warns roles.baseline naming flai upgrade
status: backlog
parent: S-0419
owner: alex
created: 2026-10-16T05:46:09Z
updated: 2026-10-21T12:19:38Z
transitions: []
stream: S-0419
tags: [cli]
touches: [flai/internal/roles/baseline.go, flai/internal/roles/baseline_test.go, flai/internal/roles/baseline/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/analyzer.md, flai/internal/check/check.go, flai/internal/check/roles_test.go]
---
# T-1754 A project with no design/workflows/roles runs the three roles from baseline files embedded in flai, and flai check warns roles.baseline naming flai upgrade

## Work

- Embed the template's three role files in flai with `go:embed`, from a copy under `flai/internal/roles/baseline/` that the build keeps in step with `template/root/design/workflows/roles/`. A test fails when they differ.
- `roles.Load` falls back to them for a role a project's folder lacks.
- `flai check` warns `roles.baseline` once per missing built-in role, naming `flai upgrade`.

First layer. The prompt task waits for it: the harness's prompt tests start roles in a temporary root with no role files, and they pass unchanged only when the prompt falls back to these embedded bodies.

## Done when

`go test ./internal/roles/... ./internal/check/...` passes, including a project with no roles folder whose planner, orchestrator, and analyzer still load.

## Notes

- Rewritten by planner-S-0419 for the folder in TH-0471's answer 2 and the new order, as the orchestrator approved on TH-0579.
