---
id: T-2759
type: task
nature: improvement
title: I-0293 is closed, naming the corrected documents and the route-table test
status: backlog
parent: S-0626
owner: alex
created: 2026-11-09T17:13:28Z
updated: 2026-11-09T17:13:28Z
transitions: []
stream: S-0626
tags: [issues]
touches: [design/issues/I-0293-design-documents-describe-a-conventions-page-a-streams-page-a-mounted.md, design/issues/summary.md]
after: [T-2755, T-2756, T-2757, T-2758]
---
# T-2759 I-0293 is closed, naming the corrected documents and the route-table test

## Work

Close I-0293 with `flai issue close I-0293 --reason "<reason>"`. The reason names the documents T-2755 to T-2758 corrected, and `flaiover/src/lib/server/design-routes.test.ts`, which now fails when the dashboard design names a route flaiover lacks. flai regenerates `design/issues/summary.md`.

It waits for T-2755, T-2756, T-2757, and T-2758, because the issue is closed only once every instance it lists is corrected.

## Done when

- I-0293's status is `closed`, with the reason recorded.
- `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes
