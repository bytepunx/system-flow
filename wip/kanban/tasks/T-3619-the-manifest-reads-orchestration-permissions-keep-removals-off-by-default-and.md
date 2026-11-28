---
id: T-3619
type: task
nature: improvement
title: The manifest reads orchestration.permissions.keep_removals, off by default and needing accept_reviews, and the settings index lists it
status: backlog
parent: S-0852
owner: alex
created: 2026-11-27T23:57:31Z
updated: 2026-11-27T23:57:31Z
transitions: []
stream: S-0852
tags: [flai]
touches: [flai/internal/manifest/manifest_roles.go, flai/internal/manifest/settings.go, flai/internal/manifest/manifest_test.go, docs/operators/settings.md, design/system/project-manifest.md]
after: [T-3618]
---
# T-3619 The manifest reads orchestration.permissions.keep_removals, off by default and needing accept_reviews, and the settings index lists it

## Work

Add the permission the ADR of T-3618 names, `keep_removals` by the planner's recommendation, beside `PermitAcceptReviews` in `flai/internal/manifest/manifest_roles.go`. It is off by default, and `OrchestratorPermits` reads it as it reads the others. A manifest that turns it on with `accept_reviews` off is a finding of the manifest check, naming both keys. Add its row to the settings in `flai/internal/manifest/settings.go`, to `docs/operators/settings.md`, which `TestSettingsIndex` checks, and to the permissions table in `design/system/project-manifest.md`.

It waits for T-3618, whose ADR fixes the permission's name and its rule.

## Done when

- `orchestration.permissions.keep_removals` parses, is off when absent, and a test shows it reported when on without `accept_reviews`.
- `docs/operators/settings.md` and `design/system/project-manifest.md` list it, and `flai test` passes on the changed paths.

## Notes
