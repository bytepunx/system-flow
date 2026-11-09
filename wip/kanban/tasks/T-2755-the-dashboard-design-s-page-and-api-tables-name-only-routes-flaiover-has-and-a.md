---
id: T-2755
type: task
nature: improvement
title: The dashboard design's page and API tables name only routes flaiover has, and a test fails on a row with no route
status: in-progress
parent: S-0626
owner: alex
created: 2026-11-09T17:13:02Z
updated: 2026-11-09T17:27:13Z
transitions:
  - to: ready
    at: 2026-11-09T17:27:12Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T17:27:13Z
    by: system-flow
stream: S-0626
tags: [flaiover, design, test]
touches: [design/system/flaiover-dashboard.md, flaiover/src/lib/server/design-routes.test.ts]
---
# T-2755 The dashboard design's page and API tables name only routes flaiover has, and a test fails on a row with no route

## Work

Write the reproducing test first: `flaiover/src/lib/server/design-routes.test.ts`, a vitest test in the manner of `no-project-files.test.ts` and `flai/cmd/reference_test.go`'s `TestReferenceIsCurrent`. It reads `../design/system/flaiover-dashboard.md`, takes the route of each row of the page table under `## Views` and of the API table under `## API (S-0011)`, and fails, naming the row, when `flaiover/src/routes` has no `+page.svelte` or `+server.ts` for it. Map SvelteKit parameters as the doc writes them, such as `[id]` and `[...path]`. Check that it fails on the rows as main has them: `/conventions` (line 126) and `/streams` (line 129).

Then correct the tables:

- Replace the `/conventions` row with the page that serves it now, the docs explorer at `/docs/design/conventions` (`flaiover/src/lib/sitemenu.ts:293-297`), or drop the row where line 46 already says so.
- Drop the `/streams` row and its `/api/streams` API: the active streams are `/activity` (lines 43 and 91).
- Add API rows for `api/items/[id]/checks` (GET, `checks.status` and `checks_enabled`) and `api/items/[id]/checks/tail` (ADR-0182).

Bump the document's `updated`. S-0611, in progress, also touches this document; a trial merge judges the overlap, and `flai stream sync` says when you must agree with it.

It waits for no task: no other task changes these files.

## Done when

- `design-routes.test.ts` fails against the tables as main has them and passes after the correction.
- Every page and API row of `flaiover-dashboard.md` names a route under `flaiover/src/routes`, and the two `checks` routes have rows.
- `flai test` on the two paths passes, `vitest` and `flaiover-lint` included.

## Notes
