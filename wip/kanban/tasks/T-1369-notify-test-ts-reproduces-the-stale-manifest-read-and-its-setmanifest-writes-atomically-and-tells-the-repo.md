---
id: T-1369
type: task
nature: remediation
title: notify.test.ts reproduces the stale manifest read, and its setManifest writes atomically and tells the Repo
status: in-progress
parent: S-0288
owner: alex
created: 2026-10-08T08:42:33Z
updated: 2026-10-08T08:59:08Z
transitions:
  - to: ready
    at: 2026-10-08T08:59:08Z
    by: agent-S-0288
  - to: in-progress
    at: 2026-10-08T08:59:08Z
    by: agent-S-0288
stream: S-0288
tags: [flaiover, tests]
touches: [flaiover/src/lib/server/notify.test.ts]
---
# T-1369 notify.test.ts reproduces the stale manifest read, and its setManifest writes atomically and tells the Repo

## Work

First layer: it waits for no task.

- Confirm the cause in the story's `### Planning`: on a fresh `Repo`, `addQuestion`'s `repo.changed(...)` starts a background `manifest()` read that `remember` caches, and `setManifest` rewrites `system-flow.yaml` with a truncating `writeFile` and never tells the `Repo`. If the cause turns out different, say so in `## Decisions` and fix what it is.
- Add a test to `flaiover/src/lib/server/notify.test.ts` that reproduces it without timing luck: report a change to a wip file, await the manifest read it starts (`await repo.manifest()`), then set `dashboard.notify_url` with `setManifest`, and expect `startNotifier` not to return null. Run it against the current `setManifest` first and see it fail.
- Fix `setManifest`: write the new manifest to a temporary file in the same folder and rename it over `system-flow.yaml`, so no reader sees it empty, and then call `repo.changed('system-flow.yaml')`, as flai on the host reports a manifest change.
- Leave `repo.ts` and `notify.ts` as they are unless the reproduction shows the race outside the test's setup.

## Done when

- The new test fails against the old `setManifest` and passes with the fix.
- `flai test flaiover/src/lib/server/notify.test.ts` passes.
- The whole flaiover vitest suite passes five runs in a row, recorded in the narrative's log.

## Notes

The other flaiover tests that rewrite `system-flow.yaml` (`inbox.test.ts`, `writes.test.ts`, `project-events.test.ts`, `projects.test.ts`, `designer.test.ts`) have no recorded instance; they are out of scope unless the reproduction shows they share the cause.
