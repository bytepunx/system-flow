---
id: T-2710
type: task
nature: remediation
title: mdlint.Fenced quotes tool output in a code block the markdown lint passes
status: done
parent: S-0614
owner: alex
created: 2026-11-09T09:35:42Z
updated: 2026-11-09T09:46:42Z
transitions:
  - to: ready
    at: 2026-11-09T09:40:18Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T09:40:19Z
    by: system-flow
  - to: done
    at: 2026-11-09T09:46:42Z
    by: agent-S-0614
stream: S-0614
tags: [flai, markdown-lint]
touches: [flai/internal/mdlint/quote.go, flai/internal/mdlint/quote_test.go]
usage:
  source: log
  seconds: 383
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 20
      output: 148
      cache_read: 2252180
      cache_write: 19899
      cost: 0.8261
---
# T-2710 mdlint.Fenced quotes tool output in a code block the markdown lint passes

## Work

Add `Fenced(text string) string` to `flai/internal/mdlint/quote.go`, beside `QuoteBareURLs`: `text` in a fenced code block of language `text` that flai's own lint passes under the project's configuration. Line ends `\r\n` and `\r` become `\n`, tabs become four spaces (MD010), ANSI escape sequences and other control characters but `\n` are dropped, trailing spaces are trimmed, runs of blank lines fold to one and leading and trailing blank lines go, and the fence is a backtick longer than the longest backtick run in the text and at least three. An empty text gives an empty block. `fenced` in `flai/internal/digest/seed.go` does most of this already: follow it, but leave it in place.

This task waits for nothing: it is the first layer, and the two send-back tasks quote with it.

## Done when

- `mdlint.Fenced` exists, with a doc comment saying what it changes and why (I-0242).
- `flai/internal/mdlint/quote_test.go` lints a document holding `Fenced` of go test's summary line (`FAIL<tab>github.com/…/serve<tab>131.960s`), a CRLF log, coloured output with ANSI escapes, trailing spaces and tabs, three blank lines, and a run of four backticks, with the repository's `.markdownlint.yaml`, and finds nothing; and asserts the summary line is quoted with its tabs as spaces.
- `flai test flai/internal/mdlint` passes.

## Notes
