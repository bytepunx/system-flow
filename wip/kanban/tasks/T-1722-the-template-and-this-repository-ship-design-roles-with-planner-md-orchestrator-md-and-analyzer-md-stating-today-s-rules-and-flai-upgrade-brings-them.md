---
id: T-1722
type: task
nature: feature
title: The template and this repository ship design/roles with planner.md, orchestrator.md, and analyzer.md stating today's rules, and flai upgrade brings them
status: backlog
parent: S-0414
owner: alex
created: 2026-10-16T05:43:05Z
updated: 2026-10-16T05:43:05Z
transitions: []
stream: S-0414
tags: [template]
touches: [template/root/design/roles/README.md, template/root/design/roles/planner.md, template/root/design/roles/orchestrator.md, template/root/design/roles/analyzer.md, design/roles/README.md, design/roles/planner.md, design/roles/orchestrator.md, design/roles/analyzer.md, template/template.yaml, flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go]
after: [T-1720]
---
# T-1722 The template and this repository ship design/roles with planner.md, orchestrator.md, and analyzer.md stating today's rules, and flai upgrade brings them

## Work

Write the three baseline role files in the ADR's format, in `template/root/design/roles/` and in this repository's `design/roles/`, each with a `README.md` for the folder.

- Each file states its role's rules as they are today: host action, trigger, item kind, topics, permissions, tools, commands, and files. The sources are `design/system/strategic-agents.md`, ADR-0082, ADR-0087, ADR-0099, ADR-0143, and the rules in `flai/internal/guard/guard.go`.
- Each body is today's prompt text from `planPrompt`, `orchestratePrompt`, and `analyzePrompt` in `flai/internal/harness/harness.go`, followed by the marker.
- Make `flai upgrade` merge a role file at its marker, as it does a convention, and add a test that a missing baseline role arrives and a project's text below the marker stays.
- Bump the template's version.

Waits for T-1720 for the format. Runs beside T-1721: they share no path.

## Done when

The files exist in both places. The upgrade test passes. `flai check --strict` here reports nothing new.

## Notes
