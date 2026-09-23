---
id: S-0098
type: story
nature: feature
title: Enable importing a project into flai from the dashboard
status: done
owner: alex
created: 2026-09-23T03:58:10Z
updated: 2026-09-23T05:04:00Z
transitions:
  - to: ready
    at: 2026-09-23T03:58:13Z
    by: alex
  - to: in-progress
    at: 2026-09-23T03:59:07Z
    by: system-flow
  - to: review
    at: 2026-09-23T04:27:25Z
    by: system-flow
  - to: done
    at: 2026-09-23T05:04:00Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# S-0098 Enable importing a project into flai from the dashboard

## Goal

When picking a project in the board that hasn't been initialized into system flow using flai, there should be a prompt asking whether the user wants to migrate that project into system-flow.

## Acceptance criteria
- [x] a project that was not created with flai can be migrated through the board
- [x] on completion of migration, run available tests in the project, then commit the introduction of the new file structure and files
- [x] when a newly migrated project has been selected from the dashboard, the operator should not receive any more prompts about importing
- [x] if the user declines the import, notify them that the will need to select a different project
- [x] if a project the operator declined to import earlier is selected again, the operator should be prompted again

## Tasks
- T-0350 flai import --commit refuses a dirty tree, runs the tests it detects, and commits the migration only when they pass
- T-0351 flai serve offers the board every git repository without a manifest under the folders the operator names
- T-0352 Picking a repository that is not imported prompts to import it, and a decline says to pick another
- T-0353 Tried with scratch repositories and a browser; ADR and docs say how importing from the board works

## Notes
