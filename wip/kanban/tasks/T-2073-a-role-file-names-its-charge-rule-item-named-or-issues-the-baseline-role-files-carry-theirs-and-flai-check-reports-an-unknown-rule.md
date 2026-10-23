---
id: T-2073
type: task
nature: feature
title: A role file names its charge rule, item, named, or issues, the baseline role files carry theirs, and flai check reports an unknown rule
status: backlog
parent: S-0425
owner: alex
created: 2026-10-23T15:34:05Z
updated: 2026-10-23T15:34:05Z
transitions: []
stream: S-0425
tags: [cli, template]
touches: [flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/check/check.go, flai/internal/check/roles_test.go, design/workflows/roles/planner.md, design/workflows/roles/orchestrator.md, design/workflows/roles/analyzer.md, template/root/design/workflows/roles/planner.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/analyzer.md, template/CHANGELOG.md]
after: [T-1790]
---
# T-2073 A role file names its charge rule, item, named, or issues, the baseline role files carry theirs, and flai check reports an unknown rule

## Work

The role file S-0414 defines has no field that says how a role's activity is charged. Add one, as the ADR of T-1790 names it:

- `flai/internal/roles` parses an optional `charge` front matter field with one of three values: `item`, the item the run was started for and every item above it; `named`, split evenly over the work items the activity names; `issues`, split evenly over the issues it names. A role without `charge` charges nothing: its activities stay in its project strategic total.
- `flai check` reports `roles.file` for a `charge` that is none of the three, and for `charge: item` on a role whose file takes no item.
- The baseline role files here and in the template carry `charge: item` on the planner, `charge: named` on the orchestrator, and `charge: issues` on the analyzer, so that the three keep the rules ADR-0083, ADR-0095, and ADR-0100 give them. `template/CHANGELOG.md` records the field.

Waits for T-1790, which fixes the field's name and values. Runs beside T-1792 in its layer: they share no path.

## Done when

`go test ./internal/roles/... ./internal/check/...` passes, with a test per rule and one for an unknown rule, and `flai check --strict` reports nothing on the six role files.

## Notes

The template's three role files and this repository's are written by S-0414 (T-1722); this task adds one line to each.
