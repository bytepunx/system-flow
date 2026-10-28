---
id: T-1722
type: task
nature: feature
title: The template and this repository ship design/workflows/roles with planner.md, orchestrator.md, and analyzer.md stating today's rules, and flai upgrade brings them
status: backlog
parent: S-0414
owner: alex
created: 2026-10-16T05:43:05Z
updated: 2026-10-28T13:07:25Z
transitions: []
stream: S-0414
tags: [template]
touches: [template/root/design/workflows/roles/README.md, template/root/design/workflows/roles/planner.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/analyzer.md, design/workflows/roles/README.md, design/workflows/roles/planner.md, design/workflows/roles/orchestrator.md, design/workflows/roles/analyzer.md, template/template.yaml, flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, design/workflows/README.md, template/root/design/workflows/README.md, design/workflows/roles/curator.md, template/root/design/workflows/roles/curator.md]
after: [T-1720]
---
# T-1722 The template and this repository ship design/workflows/roles with planner.md, orchestrator.md, and analyzer.md stating today's rules, and flai upgrade brings them

## Work

Write the three baseline role files in the ADR's format, in `template/root/design/workflows/roles/` and in this repository's `design/workflows/roles/`, each with a `README.md` for the folder.

- Each file states its role's rules as they are today: host action, trigger, item kind, topics, permissions, tools, commands, and files. The sources are `design/system/strategic-agents.md`, ADR-0082, ADR-0087, ADR-0099, ADR-0143, and the rules in `flai/internal/guard/guard.go`.
- Each body is today's prompt text from `planPrompt`, `orchestratePrompt`, and `analyzePrompt` in `flai/internal/harness/harness.go`, followed by the marker.
- Make `flai upgrade` merge a role file at its marker, as it does a convention, and add a test that a missing baseline role arrives and a project's text below the marker stays.
- Bump the template's version.

Waits for T-1720 for the format. Runs beside T-1721: they share no path.

## Done when

The files exist in both places. The upgrade test passes. `flai check --strict` here reports nothing new.

## Notes

- Each file also carries the fields its harness definition is generated from (TH-0471, answer 6), as T-1720's ADR lists them. Those fields state what `.claude/agents/planner.md`, `orchestrator.md`, and `analyzer.md` hold today, so a later generation reproduces them. Put each prompt where the ADR says: the definition's prompt, the start prompt, or both.
- `Compute` and `merge` in `flai/internal/upgrade/upgrade.go` already merge any file that has the marker, whatever its path. The upgrade work may be the test alone, in `upgrade_test.go` beside `TestComputeClassesAndApply`.
- `documentation()` in `flai/internal/check/check.go` walks all of `design/`. Give the role files the `title`, `updated`, and `status` it asks for, unless T-1723 skips `workflows` there.
