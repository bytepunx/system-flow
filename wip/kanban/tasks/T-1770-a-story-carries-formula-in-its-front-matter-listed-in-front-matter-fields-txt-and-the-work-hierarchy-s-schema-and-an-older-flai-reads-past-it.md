---
id: T-1770
type: task
nature: feature
title: A story carries formula in its front matter, listed in front-matter-fields.txt and the work hierarchy's schema, and an older flai reads past it
status: backlog
parent: S-0421
owner: alex
created: 2026-10-16T05:47:18Z
updated: 2026-10-16T05:47:18Z
transitions: []
stream: S-0421
tags: [cli]
touches: [flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/item.go, flai/internal/workitem/fields_test.go, design/system/work-hierarchy.md]
after: [T-1768]
---
# T-1770 A story carries formula in its front matter, listed in front-matter-fields.txt and the work hierarchy's schema, and an older flai reads past it

## Work

- Add `formula` to a story's front matter in `flai/internal/workitem/item.go`, written only by `flai formula run`.
- List it in `front-matter-fields.txt`, so that publishing this flai raises `flai.minimum`.
- Add it to `design/system/work-hierarchy.md` § Front matter schema.
- Test that a story with it round-trips.

Waits for the ADR, which names the field. Runs beside the loader task, since they share no path.

## Done when

`go test ./internal/workitem/...` passes, `fields_test.go` included.

## Notes
