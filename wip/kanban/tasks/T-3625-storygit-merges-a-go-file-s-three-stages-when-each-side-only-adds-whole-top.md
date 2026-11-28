---
id: T-3625
type: task
nature: improvement
title: storygit merges a Go file's three stages when each side only adds whole top-level declarations at one place
status: backlog
parent: S-0857
owner: alex
created: 2026-11-27T23:58:11Z
updated: 2026-11-27T23:58:11Z
transitions: []
stream: S-0857
tags: [go, git]
touches: [flai/internal/storygit/gomerge.go, flai/internal/storygit/gomerge_test.go]
after: [T-3623]
---
# T-3625 storygit merges a Go file's three stages when each side only adds whole top-level declarations at one place

## Work

Add `MergeGo(base, ours, theirs []byte) ([]byte, error)` in `flai/internal/storygit/gomerge.go`, the rule the ADR of T-3623 decides; it waits for that task because it builds that rule. Reuse `lineHunks` and `merge3` from `mdmerge.go` for the line merge, keeping ours, then theirs, at one insertion point. Around it, check with `go/parser` and `go/ast` that each side's hunks are pure insertions of whole top-level declarations with their doc comments, that neither touches the import block, that the two sides declare no name in common, and that the merged file parses and is gofmt-clean (`go/format`). Any other case returns an error and writes nothing.

## Done when

- `gomerge_test.go` covers, in table tests: two sides each appending test functions after the same function, the I-0409 case from `flai/internal/verify/select_test.go`, which merges to both functions whole, ours first, gofmt-clean; and the refusals: a side that edits an existing function, a side that adds an import, both sides declaring one name, and a result that does not parse.
- `flai test flai/internal/storygit` passes.

## Notes
