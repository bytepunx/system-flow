---
id: T-1317
type: task
nature: remediation
title: flai's MD034 reports a bare www. literal as markdownlint does
status: done
parent: S-0324
owner: alex
created: 2026-10-08T00:29:29Z
updated: 2026-10-08T00:33:36Z
transitions:
  - to: ready
    at: 2026-10-08T00:31:34Z
    by: agent-S-0324
  - to: in-progress
    at: 2026-10-08T00:31:34Z
    by: agent-S-0324
  - to: done
    at: 2026-10-08T00:33:36Z
    by: agent-S-0324
stream: S-0324
tags: [mdlint]
touches: [flai/internal/mdlint/inline.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases/www.md, flai/internal/mdlint/testdata/cases/expected.txt]
usage:
  source: log
  seconds: 122
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 20
      output: 6310
      cache_read: 1340525
      cache_write: 31236
      cost: 0.6443
---
# T-1317 flai's MD034 reports a bare `www.` literal as markdownlint does

## Work

First layer: it waits for no task.

- Write `flai/internal/mdlint/testdata/cases/www.md`, a fixture of bare `www.` literals and their edges, in the style of `email.md`:
  - at the start of a line, after a space, and after each character the `c == 'w'` case accepts before it (`(`, `*`, `_`, `~`, `[`, `]`)
  - after a letter, a slash, or a colon, where it is not a literal
  - with trailing punctuation that `trimURL` cuts
  - in a code span, in an autolink, as a link destination, in link text, and after a `[` nothing closes
  - `www.` with no domain after it.
- Regenerate `flai/internal/mdlint/testdata/cases/expected.txt` with `scripts/mdlint-fixtures.sh`, so markdownlint-cli2 0.20.0 says which lines are MD034.
- In `parseRange` in `flai/internal/mdlint/inline.go`, record the literal's line and column in `out.urls` in the `c == 'w'` case, as the `c == 'h'` case does, only when `!link && !unclosed`. Fix any edge where flai and the fixture's expected lines differ.
- Add a test to `flai/internal/mdlint/mdlint_test.go`, after `TestBareEmailOfI0056`, that reproduces I-0110: `Guard` refuses a thread entry with a bare `www.example.com` with `MD034/no-bare-urls`, and passes one in a code span or an autolink.

## Done when

- `www.md` and the regenerated `expected.txt` are committed, and flai's findings match markdownlint's for every fixture.
- The I-0110 test passes, and failed before the change to `inline.go`.
- `flai test` on `flai/internal/mdlint` passes.

## Notes

Drafted by the planner for S-0324. The fix is the one I-0110's instance proposes.
