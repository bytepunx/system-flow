---
id: T-1811
type: task
nature: feature
title: doc_search finds a digest entry and the dashboard's documents tree lists design/digest, shown by tests
status: backlog
parent: S-0428
owner: alex
created: 2026-10-16T06:03:27Z
updated: 2026-10-16T06:03:27Z
transitions: []
stream: S-0428
tags: [cli]
touches: [flai/internal/context/sections_test.go, flai/internal/hostapi/docs_test.go]
after: [T-1808]
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
