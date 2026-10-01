---
id: S-0185
type: story
nature: remediation
title: flai new refuses an empty required variable before writing, and flai upgrade renders a fork's own template variables
status: in-progress
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:53:51Z
transitions:
  - to: ready
    at: 2026-10-01T08:32:13Z
    by: alex
  - to: in-progress
    at: 2026-10-01T08:48:23Z
    by: agent-S-0185
tags: [flai, template]
touches: [flai/cmd/new.go, flai/cmd/upgrade.go, flai/cmd/new_test.go, flai/cmd/upgrade_test.go, flai/internal/template, flai/internal/manifest, flai/internal/lock, flai/internal/upgrade, docs/contributors/template.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/system/template.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/repository-layout.md, design/issues]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 560
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 166
      output: 1100
      cache_read: 9247556
      cache_write: 200817
      cost: 3.7375
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
- T-0645 flai new refuses an empty required --var before writing anything
- T-0646 flai records the variables a project was rendered with, and flai upgrade renders with them
- T-0647 The design and guides describe recorded variables, and I-0040 and I-0041 are closed

## Notes
