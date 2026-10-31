---
id: T-2382
type: task
nature: improvement
title: The unfixed-markdown findings move from flai/cmd into flai/internal/verify, so the merge queue reports them as flai accept does
status: backlog
parent: S-0549
owner: alex
created: 2026-10-31T07:40:40Z
updated: 2026-10-31T07:40:40Z
transitions: []
stream: S-0549
tags: [flai]
touches: [flai/cmd/accept_markdown.go, flai/internal/verify/findings.go, flai/internal/verify/findings_test.go]
---
# T-2382 The unfixed-markdown findings move from flai/cmd into flai/internal/verify, so the merge queue reports them as flai accept does

## Work

`flai accept` lints the main checkout's uncommitted files in `refuseUnfixedUncommitted` (`flai/cmd/accept_markdown.go`), through `verify.CheckoutFix` and `verify.CheckoutUnfixed` (`flai/internal/verify/fix.go`). The code that turns the check's output into findings, `unfixedFindings`, `findingLine`, and the markdownlint line regex, is in package `cmd`, which `flai/internal/mergequeue` cannot import.

Move that parsing into a new `flai/internal/verify/findings.go`, exported, with its tests in `findings_test.go`, and make `accept_markdown.go` call it, so that `flai accept`'s refusal and its `--json` data stay as they are. Leave `fix.go` alone: S-0547 changes `flai/internal/verify`.

First layer: it waits for nothing. The preflight task waits for it.

## Done when

`flai/internal/verify` exports the findings of a fixing tier's check (file, line, rule, message), `accept_markdown.go` uses them, and `go test ./cmd ./internal/verify` passes with the accept tests (`accept_uncommitted_test.go`, `accept_fix_test.go`) unchanged.

## Notes

Drafted by the planner of S-0549 on 2026-10-31.
