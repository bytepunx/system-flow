---
id: T-2809
type: task
nature: improvement
title: Add a design-routes tier that runs design-routes.test.ts when the dashboard design or a route changes
status: backlog
parent: S-0644
owner: alex
created: 2026-11-09T22:08:30Z
updated: 2026-11-09T22:08:30Z
transitions: []
stream: S-0644
tags: [flaiover, testing]
touches: [system-flow.yaml, flai/internal/verify/select_test.go]
---
# T-2809 Add a design-routes tier that runs design-routes.test.ts when the dashboard design or a route changes

## Work

`flaiover/src/lib/server/design-routes.test.ts` reads `design/system/flaiover-dashboard.md` and `flaiover/src/routes` from disk and imports neither. So the vitest tier's `vitest related` run never selects it for a change to the design or a deleted route (I-0322).

In `system-flow.yaml`'s `tests`, add a tier `design-routes` after `flaiover-flai`, in the same shape:

- `command`: `["../scripts/flaiover-unit.sh", related, "--run", "--reporter=json", "src/lib/server/design-routes.test.ts"]`
- `dir`: `flaiover`
- `paths`: `["design/system/flaiover-dashboard.md", "flaiover/src/routes/**"]`
- `format`: `vitest-json`

Add `design-routes` to the `flaiover` tier's `covers`, so that under `--all` the test runs once.

In `flai/internal/verify/select_test.go`, add a repository-manifest test beside `TestRepositoryManifestChecksTheSettingsIndexWhenItChanges`. It pins that `design/system/flaiover-dashboard.md` alone selects `design-routes` with that argv. It pins that a file under `flaiover/src/routes` selects it too. It pins that under `--all` the `flaiover` tier covers it. Do not gate it under `-short`, so that `make test` catches a manifest that drops the tier.

This task waits for no other: it is the fix.

## Done when

- `flai test design/system/flaiover-dashboard.md` runs `design-routes`, and its output lists the three tests of `design-routes.test.ts` as passed.
- The new test in `select_test.go` fails on the manifest without the tier and passes with it.
- `go test -count=1 -run '^TestRepositoryManifest' ./internal/verify` in `flai/` passes.
- `flai test system-flow.yaml flai/internal/verify/select_test.go` passes.

## Notes
