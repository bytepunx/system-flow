---
id: T-2810
type: task
nature: improvement
title: Describe the design-routes tier in devex.md and tooling.md and close I-0322
status: backlog
parent: S-0644
owner: alex
created: 2026-11-09T22:08:41Z
updated: 2026-11-09T22:08:41Z
transitions: []
stream: S-0644
tags: [docs, testing]
touches: [design/system/devex.md, design/conventions/tooling.md, design/issues/I-0322-a-change-to-flaiover-dashboard-md-alone-selects-no-vitest-tier-so-design-routes.md, design/issues/summary.md]
after: [T-2809]
---
# T-2810 Describe the design-routes tier in devex.md and tooling.md and close I-0322

## Work

In `design/system/devex.md`, in the Test tiers row of the table under "Used by this monorepo", name `design-routes` beside `settings-index` and `flaiover-flai`. Say that it runs `design-routes.test.ts` alone, that `design/system/flaiover-dashboard.md` and `flaiover/src/routes/**` select it, and why: the vitest tier's `related` run follows the import graph, and the test imports neither. Cite I-0322.

In `design/conventions/tooling.md`, below the marker, add `design-routes` to the list of tiers `flai test` runs, and to the tiers the `flaiover` tier covers under `--all`. Edit only the project additions.

Close I-0322 with `flai issue close I-0322 --reason`, saying that the `design-routes` tier now runs the test when the dashboard design or a route changes, and naming the repository-manifest test that pins it. `flai issue close` writes `design/issues/summary.md`.

This task waits for T-2809: the documents describe the tier as T-2809 builds it, and the issue closes only once the fix is in.

## Done when

- `devex.md` and `tooling.md` name the `design-routes` tier, its paths, and the reason, matching `system-flow.yaml`.
- I-0322's status is `closed`, with the reason, and `design/issues/summary.md` no longer lists it.
- `flai test design/system/devex.md design/conventions/tooling.md` passes.
- `flai check --strict` is clean.

## Notes
