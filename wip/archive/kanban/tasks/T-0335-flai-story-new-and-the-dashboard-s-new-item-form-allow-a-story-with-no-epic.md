---
id: T-0335
type: task
nature: feature
title: flai story new and the dashboard's new-item form allow a story with no epic
status: done
parent: S-0092
owner: alex
created: 2026-09-23T00:35:59Z
updated: 2026-09-23T00:37:01Z
transitions:
  - to: ready
    at: 2026-09-23T00:36:20Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T00:36:25Z
    by: system-flow
  - to: done
    at: 2026-09-23T00:37:01Z
    by: system-flow
stream: S-0092
tags: []
---

# T-0335 flai story new and the dashboard's new-item form allow a story with no epic

## Work
`--epic` was required for `flai story new` at three layers: the CLI's own `PreRunE`, `workitem.Repo.Create`, and `Item.Validate()` (a task's `--story` stays required at all three). `internal/hostapi/writes.go`'s `item.new` write leaves `--epic` out entirely for a parentless story instead of refusing. `NewItemForm.svelte`'s parent dropdown offers "No epic" as its first, selectable option rather than a disabled placeholder, and the create button no longer waits on a parent being chosen.

## Done when
Both acceptance criteria: the dashboard's form lets "No epic" stay chosen and creates the story; `flai story new` (and `item.new` over the channel) create a parentless story with no `parent:` in its front matter. Go tests for `Create`/`Validate`/the CLI/the write spec's constructed command line; `NewItemForm.svelte.test.ts` for the form; a real end-to-end `writes.test.ts` case through the actual flai binary. `go test ./...`, `pnpm run check`/`test:unit`/`lint`, `flai check --strict` clean. Verified live: `flai story new "..."` with no `--epic` produces front matter with no `parent` key at all.

## Notes
Left `Task`'s `--story` requirement untouched: a task without a story is not what this story asked for, and every existing task-parent code path (narratives, streams, `activity.get`) assumes one.
