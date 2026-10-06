---
id: T-0940
type: task
nature: improvement
title: An issue reads and writes usage with a strategic entry per kind, and front-matter-fields.txt lists the key
status: done
parent: S-0227
owner: alex
created: 2026-10-05T05:45:28Z
updated: 2026-10-06T21:35:10Z
transitions:
  - to: ready
    at: 2026-10-06T21:31:14Z
    by: agent-S-0227
  - to: in-progress
    at: 2026-10-06T21:31:15Z
    by: agent-S-0227
  - to: done
    at: 2026-10-06T21:35:10Z
    by: agent-S-0227
stream: S-0227
tags: [flai]
touches: [flai/internal/issues/issues.go, flai/internal/issues/issues_test.go, flai/internal/issues/usage.go, flai/internal/issues/usage_test.go, flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/usage.go]
usage:
  source: log
  seconds: 235
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 33
      output: 12397
      cache_read: 1747269
      cache_write: 51246
      cost: 0.9037
---
# T-0940 An issue reads and writes usage with a strategic entry per kind, and front-matter-fields.txt lists the key

## Work

Give `issues.Issue` a `usage` of type `usage.Usage` (`flai/internal/usage/usage.go`), carrying only `strategic` entries, read by `Parse` and written by `Marshal` after the issue's other keys, and validated by `Validate`. It waits for T-0935 because the ADR decides the shape.

- Add a function in `flai/internal/issues/usage.go` that adds a charge to an issue's entry for a kind, as `workitem.ChargeStrategic` does for an item, so that the charge and the carry-over share it.
- Add `usage` to the `issue:` line of `flai/internal/workitem/front-matter-fields.txt`, and the issue's usage keys if the file lists them by kind, so that `TestFieldsFileIsTheCode` passes.
- An issue with no usage is written as it is today, with no `usage` key.

## Done when

- An issue with a `strategic` entry survives a read and a write unchanged, and one without usage gains no key: tests pin both.
- Charging an issue twice for one kind sums into one entry, and a second kind gets an entry of its own: tests pin both.
- `scripts/flai-test.sh` passes, `TestFieldsFileIsTheCode` included.

## Notes
