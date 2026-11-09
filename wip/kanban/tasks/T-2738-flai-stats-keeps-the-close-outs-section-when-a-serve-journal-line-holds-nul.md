---
id: T-2738
type: task
nature: remediation
title: flai stats keeps the close-outs section when a serve journal line holds NUL bytes or is not an entry, skipping and warning of it
status: done
parent: S-0624
owner: alex
created: 2026-11-09T17:04:44Z
updated: 2026-11-09T17:17:26Z
transitions:
  - to: ready
    at: 2026-11-09T17:12:37Z
    by: agent-S-0624
  - to: in-progress
    at: 2026-11-09T17:12:38Z
    by: agent-S-0624
  - to: done
    at: 2026-11-09T17:17:26Z
    by: agent-S-0624
stream: S-0624
tags: [flai, stats]
touches: [flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go]
after: [T-2737]
usage:
  source: log
  seconds: 288
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 1253
      cache_read: 1081733
      cache_write: 17476
      cost: 0.3812
---
# T-2738 flai stats keeps the close-outs section when a serve journal line holds NUL bytes or is not an entry, skipping and warning of it

## Work

Build what T-2737's ADR decides. It waits for T-2737 because the ADR settles what is skipped and what the warning says.

- In `flai/internal/statsread/statsread.go`, make `journalEntries` remove the NUL bytes from each line before it trims and decodes the line. When a line still does not decode as an entry, skip it, count it, and keep the first skipped line's number and error. Return the entries with the count. A scanner error, such as a line over the buffer, still fails the read.
- Make `CloseOuts` return the entries it read when lines were skipped, and log one warning naming the journal's path, how many lines it skipped, and the first one's number and error. A journal that cannot be opened or read is still nil, warned of as now; a missing one is still nil with no warning.
- In `flai/internal/statsread/statsread_test.go`, beside `TestCloseOutsReadsTheServeJournal`, add a test that reproduces I-0290: a journal whose close-out end entry follows a line of NUL bytes and a valid entry, and a line that is not JSON. Assert that every close-out entry is read, the NUL-prefixed entry included, that one warning names the skipped line, and that a journal of valid lines logs no warning.

## Done when

- The new test fails on the current code and passes with the change.
- `TestCloseOutsReadsTheServeJournal` and the golden stats test still pass.
- `flai test` on the changed files passes.

## Notes
