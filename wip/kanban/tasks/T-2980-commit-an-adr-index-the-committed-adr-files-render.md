---
id: T-2980
type: task
nature: remediation
title: Commit an ADR index the committed ADR files render
status: backlog
parent: S-0710
owner: alex
created: 2026-11-10T10:29:22Z
updated: 2026-11-10T10:29:22Z
transitions: []
stream: S-0710
tags: [flai, adr]
touches: [flai/internal/adr/adr.go, flai/internal/adr/adr_test.go]
---
# T-2980 Commit an ADR index the committed ADR files render

## Work

Fix the cause of I-0336 in `flai/internal/adr`. `Accept`, `New`, and `SetTopics` write `design/adrs/README.md` from the working tree's ADR files, and `Result.commit` commits that index with only the ADR files the call wrote. When another ADR file has uncommitted changes, as ADR-0101 and ADR-0122 had in the main checkout on 2026-11-10 (commits `859f16842` and `4a247f7a9`), the commit records rows the committed ADR files do not render. Every worktree made from that main then finds `adr.index` outside its story.

Make the commit record the index that HEAD's ADR files, with the files the commit records, render under ADR-0228 § 1. For example, render it from that tree, write it to git's index as the README's blob, and commit, leaving the working tree's README as its own files render it. Cover `--autocommit` and the story worktree's `--commit` alike, whichever code path each takes. With no other ADR file changed, the commit is what it is today.

Add a test in `flai/internal/adr/adr_test.go` on a real git repository that reproduces the cause: leave one ADR accepted in the working tree and uncommitted, accept another with autocommit, and check that HEAD's README is what HEAD's ADR files render (`adrindex.Render` over HEAD's tree, or `flai check` in a clean checkout of HEAD). The test fails before the fix. If the render needs a helper over git's tree, put it in `flai/internal/adrindex/adrindex.go`; `flai task done` widens the touches.

This task waits for nothing: it is the first layer.

## Done when

- The new test fails on the code before the fix and passes after it.
- An autocommitted `flai adr accept`, `new`, or `topics` with another ADR file left uncommitted records an index that matches HEAD's ADR files, and the working tree's index still matches the working tree's files.
- `flai test` on the changed paths passes.

## Notes
