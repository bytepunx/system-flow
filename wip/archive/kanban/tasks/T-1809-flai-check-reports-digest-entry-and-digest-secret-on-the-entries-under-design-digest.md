---
id: T-1809
type: task
nature: feature
title: flai check reports digest.entry and digest.secret on the entries under design/digest
status: done
parent: S-0428
owner: alex
created: 2026-10-16T06:03:16Z
updated: 2026-10-23T17:00:48Z
transitions:
  - to: ready
    at: 2026-10-23T16:48:36Z
    by: agent-S-0428
  - to: in-progress
    at: 2026-10-23T16:48:36Z
    by: agent-S-0428
  - to: done
    at: 2026-10-23T17:00:48Z
    by: agent-S-0428
stream: S-0428
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/digest_test.go]
after: [T-1808]
usage:
  source: log
  seconds: 732
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 37
      output: 15300
      cache_read: 2060614
      cache_write: 78347
      cost: 1.1796
---
# T-1809 flai check reports digest.entry and digest.secret on the entries under design/digest

## Work

Add `checker.digest` to `flai/internal/check/check.go`, beside `checker.analysis`, and leave the folder out of the design documents' own rules as `documentation()` does for `analysis`. It reports, from `digest.Validate` and `digest.Secrets`:

- `digest.entry`, an error, for each validation problem, on its field's line;
- `digest.secret`, an error, for each secret shape, on its line.

A project with no `design/digest/` has nothing checked there. The folder's `README.md` is left out.

Waits for T-1808, whose `Validate` and `Secrets` it calls. Runs beside the commands' task and the search task, which touch other files.

## Done when

`flai/internal/check/digest_test.go` covers a clean entry, each `digest.entry` cause, a `digest.secret`, and a project with no folder; `flai test flai/internal/check` passes.

## Notes

The check runs on the main checkout too, where the curator and the operator's answers write entries uncommitted, so an entry flai wrote badly is seen there.
