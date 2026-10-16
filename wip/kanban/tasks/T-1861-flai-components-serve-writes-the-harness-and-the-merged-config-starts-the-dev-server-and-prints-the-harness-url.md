---
id: T-1861
type: task
nature: feature
title: flai components serve writes the harness and the merged config, starts the dev server, and prints the harness URL
status: backlog
parent: S-0438
owner: alex
created: 2026-10-16T21:33:52Z
updated: 2026-10-16T21:33:52Z
transitions: []
stream: S-0438
tags: [cli]
touches: [flai/cmd/components.go, flai/cmd/components_test.go, docs/users/flai-reference.md]
after: [T-1859, T-1860]
---
# T-1861 flai components serve writes the harness and the merged config, starts the dev server, and prints the harness URL

## Work

Add `serve [--project <name>]` to the `components` group S-0437 made in `flai/cmd/components.go`. It picks the sub-project (the only supported one, or `--project`), writes the harness with T-1859, starts the dev server with T-1860, prints the harness URL, and runs in the foreground until interrupted, stopping the server. Regenerate `docs/users/flai-reference.md`.

Waits for T-1859 for the files and T-1860 for the process.

## Done when

Command tests with a fake `vite` cover the URL printed, `--project`, an unsupported or ambiguous sub-project, and the stop on interrupt. `flai test flai/cmd` passes.

## Notes

None.
