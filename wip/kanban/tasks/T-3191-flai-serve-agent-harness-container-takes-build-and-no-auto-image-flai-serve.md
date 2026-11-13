---
id: T-3191
type: task
nature: feature
title: flai serve agent harness container takes --build and --no-auto-image, flai serve settings shows them, and flai guard refuses them to served agents
status: backlog
parent: S-0759
owner: alex
created: 2026-11-12T16:05:38Z
updated: 2026-11-12T16:05:38Z
transitions: []
stream: S-0759
tags: [cli]
touches: [flai/cmd/serve_harness.go, flai/cmd/serve_settings.go, flai/cmd/serve_settings_test.go, flai/internal/guard/hostconfig.go, flai/internal/guard/hostconfig_test.go]
after: [T-3190]
---
# T-3191 flai serve agent harness container takes --build and --no-auto-image, flai serve settings shows them, and flai guard refuses them to served agents

## Work

In `flai/cmd/serve_harness.go`, add `--build <command>` (repeatable words, or one word list, as the flag style there has it) and `--clear-build`. Add `--no-auto-image` and `--auto-image`, refused on a harness other than `container` (line 255). Show both in `flai serve settings` (`serve_settings.go`) and `flai serve agent show`. Test them in `serve_settings_test.go`.

`flai guard` already refuses served agents the commands that write the host's configuration (ADR-0220, `flai/internal/guard/hostconfig.go`). Add the new flags to `hostconfig_test.go`'s cases, and widen `hostconfig.go` only if a case passes.

It waits for T-3190, whose fields it sets.

## Done when

- `go test ./cmd/ -run 'ServeSettings|Harness' ./internal/guard/ -run HostConfig` passes.
- `flai serve agent harness container --help` names both settings.

## Notes
