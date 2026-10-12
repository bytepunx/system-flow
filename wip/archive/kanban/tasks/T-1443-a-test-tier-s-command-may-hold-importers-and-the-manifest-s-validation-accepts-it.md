---
id: T-1443
type: task
nature: improvement
title: "A test tier's command may hold {importers}, and the manifest's validation accepts it"
status: done
parent: S-0364
owner: alex
created: 2026-10-12T14:40:26Z
updated: 2026-10-12T15:31:17Z
transitions:
  - to: ready
    at: 2026-10-12T15:26:57Z
    by: agent-S-0364
  - to: in-progress
    at: 2026-10-12T15:26:58Z
    by: agent-S-0364
  - to: done
    at: 2026-10-12T15:31:17Z
    by: agent-S-0364
stream: S-0364
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go]
usage:
  source: log
  seconds: 259
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 17
      output: 5156
      cache_read: 767939
      cache_write: 30704
      cost: 0.4563
---
# T-1443 A test tier's command may hold {importers}, and the manifest's validation accepts it

## Work

Add the placeholder `{importers}` beside `{packages}` and `{files}` in `flai/internal/manifest/manifest.go`, as `PlaceholderImporters`. It stands for the packages `{packages}` names and every package of the same module that imports one of them, directly or through other packages, in its own files or its test files. Accept it in `commandProblems` wherever `{packages}` is accepted, refuse it in `all_command` as `{packages}` is refused, and name all three placeholders in each refusal that lists them. Update the comments on `TestTier.Command` and on the placeholder constants. Cover the acceptance and the refusals in `manifest_test.go`.

This task waits for none: it is the constant the verify task fills.

## Done when

- `flai check` accepts a tier whose `command` holds `"{importers}"` as an argument of its own, and refuses it inside an argument and in `all_command`.
- Every refusal that lists the placeholders names `{packages}`, `{files}`, and `{importers}`.
- `flai test flai/internal/manifest` passes.

## Notes

Drafted by planner-S-0364.
