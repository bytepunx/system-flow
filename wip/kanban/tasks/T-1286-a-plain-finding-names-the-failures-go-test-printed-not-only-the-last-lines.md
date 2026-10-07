---
id: T-1286
type: task
nature: improvement
title: A plain finding names the failures go test printed, not only the last lines
status: backlog
parent: S-0313
owner: alex
created: 2026-10-07T23:39:26Z
updated: 2026-10-07T23:39:26Z
transitions: []
stream: S-0313
tags: [flai, verify]
touches: [flai/internal/verify/parse.go, flai/internal/verify/parse_test.go]
---
# T-1286 A plain finding names the failures go test printed, not only the last lines

## Work

`plain()` in `flai/internal/verify/parse.go` keeps only the last `plainLines` (20) lines of a failing tier's output. When a plain tier runs `go test` over many packages, the failing package comes early and the tail is a list of `ok` lines and a final `FAIL`, as in S-0215's close-out (I-0113).

- Make `plain()` keep, ahead of the tail, the lines of the output that name a failure and fall outside it: go test's `--- FAIL: <Test>` lines with the indented lines under them, `FAIL<tab><package>` lines, and `panic:` lines. Cap them, as `messageLines` caps a parsed message, so a finding stays short.
- Add the reproducer to `flai/internal/verify/parse_test.go`: a plain tier whose output has a `--- FAIL` block and its `FAIL<tab><package>` line, then more than 20 `ok` lines and a final `FAIL`. Its finding must name the test and the package.
- Keep `TestPlainIsTheLastLinesAndTheExitStatus` true for output with no failure line, or change it and say why in the narrative.

This fixes any project whose plain tier runs go test, the template's among them. It waits for no task.

## Done when

- The new test fails against the old `plain()` and passes now.
- `flai test flai/internal/verify` passes.

## Notes

Drafted by the planner for S-0313.
