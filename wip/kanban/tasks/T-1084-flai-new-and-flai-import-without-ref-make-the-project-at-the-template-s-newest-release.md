---
id: T-1084
type: task
nature: remediation
title: flai new and flai import without --ref make the project at the template's newest release
status: cancelled
parent: S-0301
owner: alex
created: 2026-10-06T22:52:47Z
updated: 2026-10-06T22:53:27Z
transitions:
  - to: cancelled
    at: 2026-10-06T22:53:27Z
    by: agent-S-0301
stream: S-0301
tags: [cli]
touches: [flai/cmd/new.go, flai/cmd/new_test.go, flai/cmd/import.go, flai/cmd/import_test.go]
after: [T-1078]
---
# T-1084 flai new and flai import without --ref make the project at the template's newest release

## Work

Criterion 2. `flai new` and `flai import` with no `--ref` take the template's newest release when the configured ref follows releases.

- In `flai/cmd/new.go` and `flai/cmd/import.go`, a git template whose ref is unset or follows releases (`FollowsReleases`, T-1078) resolves to the latest version tag from `Releases`. The configured ref is the config's `template.ref`, `main` by default. Another branch or a commit set with `flai template use --ref` is used as given. A local directory is unchanged. With no version tags, the configured ref is used, fetched again when it is a branch.
- The project's manifest and lock record the tag and its version, as they record the ref today.
- Keep the change in these commands. `resolveTemplate` in `new.go` is shared with `upgrade.go`. If T-1081 needs the same resolution there, give it one helper in `new.go` that both call, and say so in your final message.

It waits for T-1078. It shares no path with T-1057 or T-1081.

## Done when

- [ ] Tests in `flai/cmd/new_test.go` and `flai/cmd/import_test.go`, against a tagged git template in a temp directory, show that with config ref `main` and no `--ref` the project is made at the newest tag, and records it in `system-flow.yaml` and the lock.
- [ ] With `--ref` naming an older tag, the project is made at that tag.
- [ ] Existing new and import tests pass.
- [ ] `go test -race ./cmd/ -run 'New|Import'` passes.

## Notes
- 2026-10-06T22:53:27Z: moved to cancelled: duplicate of the planner's draft for S-0301, written before I saw it; the planner's tasks are kept and edited
