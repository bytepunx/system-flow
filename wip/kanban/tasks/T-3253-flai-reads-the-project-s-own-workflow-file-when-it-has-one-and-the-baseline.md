---
id: T-3253
type: task
nature: feature
title: flai reads the project's own workflow file when it has one and the baseline when it has none, everywhere a repository is opened
status: backlog
parent: S-0755
owner: alex
created: 2026-11-12T16:09:46Z
updated: 2026-11-12T16:09:46Z
transitions: []
stream: S-0755
tags: [cli]
touches: [flai/internal/workflow/load.go, flai/internal/workflow/load_test.go, flai/internal/workitem/repo.go]
---
# T-3253 flai reads the project's own workflow file when it has one and the baseline when it has none, everywhere a repository is opened

## Work

Add `load.go`: read the workflow file at the path the ADR fixes under the manifest's design folder, parse it, and fall back to the baseline only when the file is absent. A file that does not parse is an error naming the finding, never a silent fallback. Make `workitem`'s repository open use it, so every command, the MCP server, and flai serve run the project's workflow. Test a project with a file, without one, and with a broken one.

## Done when

`load_test.go` passes for the three cases, and every existing test passes on projects without the file.

## Notes
