---
id: T-1419
type: task
nature: improvement
title: mdlint quotes the bare URLs md034 finds in a line, so flai can make text it takes lint clean
status: done
parent: S-0346
owner: alex
created: 2026-10-08T08:53:52Z
updated: 2026-10-08T09:32:49Z
transitions:
  - to: ready
    at: 2026-10-08T09:30:51Z
    by: agent-S-0346
  - to: in-progress
    at: 2026-10-08T09:30:51Z
    by: agent-S-0346
  - to: done
    at: 2026-10-08T09:32:49Z
    by: agent-S-0346
stream: S-0346
tags: [flai, mdlint]
touches: [flai/internal/mdlint/inline.go, flai/internal/mdlint/quote.go, flai/internal/mdlint/quote_test.go, flai/internal/mdlint/rules.go]
usage:
  source: log
  seconds: 118
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 11
      output: 2551
      cache_read: 815363
      cache_write: 27337
      cost: 0.4328
---
# T-1419 mdlint quotes the bare URLs md034 finds in a line, so flai can make text it takes lint clean

## Work

- `parseRange` in `flai/internal/mdlint/inline.go` records each bare URL, `www.` literal, and email address in `out.urls` by its start only. Record where each ends as well, without changing what `md034` reports.
- Add `QuoteBareURLs(line string) string` in a new `flai/internal/mdlint/quote.go`. It parses the line as `md034` does and wraps each bare URL it finds in a code span. The fence is one backtick longer than the longest backtick run inside the URL. Text `md034` does not flag is left as it is.
- Test it in `flai/internal/mdlint/quote_test.go`: a bare `https://` URL, a bare `www.` literal, a bare email address, one already in a code span or link, and a line with none. Each result must lint clean under `md034`.
- Waits for nothing: it is the first layer.

## Done when

- `QuoteBareURLs` returns a line on which `md034` reports nothing, for every case in its test.
- `md034`'s findings on the existing fixtures are unchanged.
- `flai test flai/internal/mdlint` passes.

## Notes
