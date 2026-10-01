---
id: S-0185
type: story
nature: remediation
title: flai new refuses an empty required variable before writing, and flai upgrade renders a fork's own template variables
status: ready
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:32:13Z
transitions:
  - to: ready
    at: 2026-10-01T08:32:13Z
    by: alex
tags: [flai, template]
touches: [flai/cmd/new.go, flai/cmd/upgrade.go, flai/internal/template, flai/internal/manifest, docs/contributors/template.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0185 flai new refuses an empty required variable before writing, and flai upgrade renders a fork's own template variables

## Goal

- **I-0041:** in `flai/cmd/new.go` (~190-192), a value given with `--var` is stored and `continue` skips the `v.Required && val == ""` check (~211), so an explicitly empty required variable such as `--var project_name=` passes, `flai new` writes the files, and only then fails on validating the manifest, leaving a half-made project.
- **I-0040:** `flai upgrade` renders with the five standard variables only (`manifestVars`, `flai/cmd/upgrade.go` ~196-201), so a template variable a fork adds fails every upgrade. `docs/contributors/template.md` documents this as a limitation.

## Acceptance criteria
- [ ] `flai new` applies the required check to given values too, and refuses before writing anything; a test covers an empty required `--var`
- [ ] The values `flai new` rendered with are recorded in the project (in `system-flow.yaml`, or a file `flai` owns), and `flai upgrade` renders with them, with the template's defaults for variables added since, and `--var` to supply or change one
- [ ] `flai upgrade` names a required variable it has no value for and refuses before writing, rather than failing mid-render
- [ ] Tests cover an upgrade with a fork's added variable, with its default, and with `--var`
- [ ] The design (`design/system/template.md`, `design/system/project-manifest.md`) and `docs/contributors/template.md` describe it, and the limitation is removed
- [ ] I-0040 and I-0041 are closed with what fixed them

## Tasks

## Notes
