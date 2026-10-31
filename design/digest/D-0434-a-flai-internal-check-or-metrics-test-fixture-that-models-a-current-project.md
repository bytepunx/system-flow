---
id: D-0434
title: A flai/internal/check or metrics test fixture that models a current project carries the template's role files, through writeTemplateRoles, so the roles.baseline advisory stays out of its findings
kind: decision
sources: [S-0419]
topics: [go, code, cli, flai]
paths: [flai/internal/check/roles_test.go, flai/internal/check/activity_test.go, flai/internal/check/epic_lags_test.go, flai/internal/check/check_test.go, flai/internal/metrics/testdata/good/design/workflows/roles/planner.md, flai/internal/metrics/testdata/good/design/workflows/roles/orchestrator.md, flai/internal/metrics/testdata/good/design/workflows/roles/analyzer.md, flai/internal/metrics/testdata/good/design/workflows/roles/curator.md, flai/cmd/check_stats_test.go]
date: 2026-10-31
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0434 A flai/internal/check or metrics test fixture that models a current project carries the template's role files, through writeTemplateRoles, so the roles.baseline advisory stays out of its findings

S-0419 added `roles.baseline`, an advisory warning `flai check` gives for each of `planner.md`, `orchestrator.md`, `analyzer.md`, and `curator.md` a project lacks. Every test project without those files then gained up to four advisories, which broke tests that count or list a check's findings.

S-0419 decided: "Test fixtures that model a current project now carry the template's role files: `flai/internal/metrics/testdata/good`, plus `activityRepo` and two advisory-count tests in `flai/internal/check` through `writeTemplateRoles`. The alternative was to filter `roles.baseline` out of every assertion." The summary test in `flai/cmd/check_stats_test.go` gives its project the role files too.

So a new test that builds a project and asserts on `flai check`'s findings or advisory counts calls `writeTemplateRoles(t, root)` (defined in `flai/internal/check/roles_test.go`, it copies the template's role files into the project), or starts from `testdata/good`, rather than filtering `roles.baseline` out.

Source: S-0419 (wip/archive/agents/S-0419.md, Decisions).
