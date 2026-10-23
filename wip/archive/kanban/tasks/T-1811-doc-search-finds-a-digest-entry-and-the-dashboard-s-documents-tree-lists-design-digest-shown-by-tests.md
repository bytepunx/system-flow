---
id: T-1811
type: task
nature: feature
title: doc_search finds a digest entry and the dashboard's documents tree lists design/digest, shown by tests
status: done
parent: S-0428
owner: alex
created: 2026-10-16T06:03:27Z
updated: 2026-10-23T17:00:54Z
transitions:
  - to: ready
    at: 2026-10-23T16:48:39Z
    by: agent-S-0428
  - to: in-progress
    at: 2026-10-23T16:48:39Z
    by: agent-S-0428
  - to: done
    at: 2026-10-23T17:00:54Z
    by: agent-S-0428
stream: S-0428
tags: [cli]
touches: [flai/internal/context/sections_test.go, flai/internal/hostapi/docs_test.go]
after: [T-1808]
usage:
  source: log
  seconds: 735
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 41
      output: 16689
      cache_read: 2247595
      cache_write: 85456
      cost: 1.2866
---
# T-1811 doc_search finds a digest entry and the dashboard's documents tree lists design/digest, shown by tests

## Work

`flai/internal/context/sections.go` indexes `design/` and `docs/`, and the host API's documents tree walks `design/`, so an entry should be found and listed without a code change. Prove it:

- In `flai/internal/context/sections_test.go`, a fixture with an entry under `design/digest/` that `SearchSections` returns for a word of its title, and its `README.md` left out as other READMEs are.
- In `flai/internal/hostapi/docs_test.go`, the documents tree of the same fixture lists `design/digest/` with the entry, and `doc.get` on the entry returns its front matter and body.

If either does not hold, fix it in `sections.go` or `docs.go` and add that file to this task's touches.

Waits for T-1808, whose format the fixtures follow. Runs beside the check and the commands' tasks.

## Done when

Both tests pass with `flai test flai/internal/context flai/internal/hostapi`.

## Notes

The dashboard renders what the host API returns, so no flaiover change is expected.
