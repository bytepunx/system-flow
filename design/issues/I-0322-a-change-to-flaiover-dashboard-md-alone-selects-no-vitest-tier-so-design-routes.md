---
id: I-0322
title: A change to flaiover-dashboard.md alone selects no vitest tier, so design-routes.test.ts does not run when the dashboard design gains a row with no route
class: impression
status: closed
count: 1
cost: 5m
first_reported: 2026-11-09T17:35:44Z
last_reported: 2026-11-09T17:35:44Z
updated: 2026-11-09T23:00:29Z
---

# I-0322 A change to flaiover-dashboard.md alone selects no vitest tier, so design-routes.test.ts does not run when the dashboard design gains a row with no route

## Description
A change to flaiover-dashboard.md alone selects no vitest tier, so design-routes.test.ts does not run when the dashboard design gains a row with no route

## Instances

### 2026-11-09T17:35:44Z
Story: S-0626.
S-0626 added flaiover/src/lib/server/design-routes.test.ts, which fails when design/system/flaiover-dashboard.md's page or API table names a route flaiover lacks. flai test on design/system/flaiover-dashboard.md alone runs only the markdown tier, because the vitest tier in system-flow.yaml's tests selects flaiover paths. So a story that changes only that document closes out without the test. The test runs when a flaiover file changes, under flai test --all, and in the merge queue's batch. A fix would have the vitest tier, or a narrow tier of its own, select design/system/flaiover-dashboard.md too. That changes the manifest's tests, which S-0626 does not touch.

## Remediation

Story S-0644 remediates this issue, created from it at 2026-11-09T21:23:04Z.
Closed 2026-11-09T23:00:29Z: S-0644: a design-routes tier in system-flow.yaml's tests runs flaiover/src/lib/server/design-routes.test.ts when design/system/flaiover-dashboard.md or a file under flaiover/src/routes changes, and the flaiover tier covers it under --all; TestRepositoryManifestChecksTheDashboardDesignsRoutesWhenEitherChanges in flai/internal/verify/select_test.go keeps the tier in the manifest.
