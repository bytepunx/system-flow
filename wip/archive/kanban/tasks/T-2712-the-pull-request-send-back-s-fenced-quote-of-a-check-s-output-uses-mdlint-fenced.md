---
id: T-2712
type: task
nature: remediation
title: The pull request send-back's fenced quote of a check's output uses mdlint.Fenced
status: done
parent: S-0614
owner: alex
created: 2026-11-09T09:35:59Z
updated: 2026-11-09T10:02:15Z
transitions:
  - to: ready
    at: 2026-11-09T09:52:29Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T09:52:29Z
    by: system-flow
  - to: done
    at: 2026-11-09T10:02:15Z
    by: agent-S-0614
stream: S-0614
tags: [flai, forge]
touches: [flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go]
after: [T-2710]
usage:
  source: log
  seconds: 586
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 17
      output: 4809
      cache_read: 2046180
      cache_write: 37764
      cost: 0.789
---
# T-2712 The pull request send-back's fenced quote of a check's output uses mdlint.Fenced

## Work

In `flai/internal/serve/forgestatus.go`, `quote(s, fenced)` with `fenced` true puts the forge's text in a code block as it is, tabs and all, so a check log that carries go test's summary line fails the send-back thread's lint as I-0242's did. Make the fenced branch return `mdlint.Fenced` of the text trimmed to `quoteLimit`, and leave the blockquote branch as it is.

Waits for T-2710, whose `mdlint.Fenced` it quotes with. It shares no file with T-2711, so the two run together.

## Done when

- A test in `flai/internal/serve/forgestatus_test.go` gives `sendBackText` a failing check whose text carries tabs, CRLF line ends, and a run of three backticks, with `fenced` true, lints the result with the project's `.markdownlint.yaml`, and finds nothing.
- The existing tests of `quote` and `sendBackText` pass.
- `flai test flai/internal/serve/forgestatus.go flai/internal/serve/forgestatus_test.go` passes.

## Notes
