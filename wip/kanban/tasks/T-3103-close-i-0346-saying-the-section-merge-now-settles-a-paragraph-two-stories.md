---
id: T-3103
type: task
nature: improvement
title: Close I-0346, saying the section merge now settles a paragraph two stories change in different sentences
status: backlog
parent: S-0719
owner: alex
created: 2026-11-11T08:55:04Z
updated: 2026-11-11T08:55:04Z
transitions: []
stream: S-0719
tags: [issues]
touches: [design/issues/I-0346-a-trial-merge-that-finds-two-open-stories-editing-different-sentences-of-one.md, design/issues/summary.md]
after: [T-3102, T-3099]
---
# T-3103 Close I-0346, saying the section merge now settles a paragraph two stories change in different sentences

## Work

Run `flai issue close I-0346 --reason "<why>"`, the reason naming the new ADR, T-3097's sentence merge in `flai/internal/storygit/mdmerge.go`, and T-3102's two reproductions, and saying that a paragraph two stories change in different sentences no longer stops sync or acceptance, while the trial merge still reports the pair. `flai issue close` writes `design/issues/summary.md` again.

Waits for T-3102 and T-3099: the reason cites the reproductions that show the fix, and the issue closes only once the documents say what changed.

## Done when

- I-0346 is closed with that reason, and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` passes.

## Notes
