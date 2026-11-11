---
id: T-3012
type: task
nature: improvement
title: Keep the digest rule's per-entry results under .flai-cache, keyed by each entry's size and modification time
status: done
parent: S-0677
owner: alex
created: 2026-11-10T10:31:37Z
updated: 2026-11-10T19:45:47Z
transitions:
  - to: ready
    at: 2026-11-10T19:33:57Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T19:33:58Z
    by: system-flow
  - to: done
    at: 2026-11-10T19:45:47Z
    by: agent-S-0677
stream: S-0677
tags: [flai, performance, check, digest]
touches: [flai/internal/digest/digest.go, flai/internal/digest/digest_test.go, flai/internal/digest/kept.go, flai/internal/digest/kept_test.go, flai/internal/check/check_records.go, flai/internal/check/digest_test.go]
usage:
  source: log
  seconds: 709
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 55
      output: 18450
      cache_read: 7137776
      cache_write: 226930
      cost: 3.1395
---
# T-3012 Keep the digest rule's per-entry results under .flai-cache, keyed by each entry's size and modification time

## Work

The `digest` rule (`flai/internal/check/check_records.go`) reads every entry and runs `digest.Validate` on it, secret scan included. It is the slowest rule run at once, 0.76 to 1.77 s on a loaded host (I-0312).

Keep each entry's results between processes, as S-0600 keeps the lint findings (`flai/internal/mdlint/kept.go`):

- One kept file under the cache folder, keyed per entry by its path, size, and modification time, and as a whole by a format version. Keep only a read made at least two seconds after the file last changed. Write the file again only when what it holds changed, beside the old one and renamed over it, and only into a `.flai-cache` that exists. A missing, damaged, or foreign file is a miss, never an error.
- `Validate` depends on the project through `Env`: `Known` for the IDs an entry cites and `Exists` for its paths. Those results must not be kept. Split `Validate` so the kept part is what the file alone decides: its problems and the IDs and paths it cites. Run `Known` and `Exists` on the cited IDs and paths every check, so that an archived or removed source is still reported.

Waits for no task: it shares no file with the `workItems` task.

## Done when

- A second whole check reads no unchanged digest entry. A test shows it, as `TestASecondFolderLintReadsNoUnchangedFile` does for the lint.
- A test shows that an entry citing an ID or path removed since the kept run is still reported.
- `TestDigestEntriesAreChecked` and the digest package's tests pass, and `flai test` passes on the changed paths.

## Notes

Drafted by the planner for S-0677.
