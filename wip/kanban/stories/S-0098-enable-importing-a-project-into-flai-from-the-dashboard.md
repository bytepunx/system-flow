---
id: S-0098
type: story
nature: feature
title: Enable importing a project into flai from the dashboard
status: backlog
owner: alex
created: 2026-09-23T03:58:10Z
updated: 2026-09-23T03:58:10Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# S-0098 Enable importing a project into flai from the dashboard

## Goal

When picking a project in the board that hasn't been initialized into system flow using flai, there should be a prompt asking whether the user wants to migrate that project into system-flow.

## Acceptance criteria
- [ ] a project that was not created with flai can be migrated through the board
- [ ] on completion of migration, run available tests in the project, then commit the introduction of the new file structure and files
- [ ] when a newly migrated project has been selected from the dashboard, the operator should not receive any more prompts about importing
- [ ] if the user declines the import, notify them that the will need to select a different project
- [ ] if a project the operator declined to import earlier is selected again, the operator should be prompted again

## Tasks

## Notes
