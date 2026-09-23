---
id: S-0102
type: story
nature: remediation
title: No project selector is visible
status: review
owner: alex
created: 2026-09-23T05:56:16Z
updated: 2026-09-23T06:04:28Z
transitions:
  - to: ready
    at: 2026-09-23T05:56:19Z
    by: alex
  - to: in-progress
    at: 2026-09-23T05:57:35Z
    by: system-flow
  - to: review
    at: 2026-09-23T06:04:28Z
    by: system-flow
tags: [dashboard]
touches: [flaiover/src]
---
# S-0102 No project selector is visible

## Goal

A control in the menu bar of the board should not only display the name of the project, but allow for switching to or importing a subdirectory of the current `flai serve` working directory.

## Acceptance criteria
- [x] The board should always show a control for switching the active project
- [x] The control should allow for importing subdirectories that are not yet imported/migrated into system-flow via flai

## Tasks
- T-0359 flai serve started in a folder serves the projects below it and offers its other repositories for import
- T-0360 The switcher is always a control, with every project and every repository offered for import
- T-0361 Tried with a flai serve started in a folder of repositories and a browser; docs

## Notes
