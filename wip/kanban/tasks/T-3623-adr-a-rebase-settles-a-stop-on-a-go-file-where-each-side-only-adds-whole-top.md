---
id: T-3623
type: task
nature: improvement
title: "ADR: a rebase settles a stop on a Go file where each side only adds whole top-level declarations at one place"
status: backlog
parent: S-0857
owner: alex
created: 2026-11-27T23:58:03Z
updated: 2026-11-27T23:58:03Z
transitions: []
stream: S-0857
tags: [adr, git]
touches: [design/adrs/README.md]
---
# T-3623 ADR: a rebase settles a stop on a Go file where each side only adds whole top-level declarations at one place

## Work

Write a new ADR with `flai adr new`, refining ADR-0175, whose text says a shared path that is not markdown is never merged. It decides that every rebase flai runs settles a stop on a `.go` file when, against the merge base, each side only inserts whole top-level declarations (func, type, var, const) at the same place, and nothing else in the file differs. The result keeps ours, then theirs, as the section merge does. The merge works from the three stages, not from git's markers, since git's hunk split the closing brace between the two sides in I-0409. It refuses, so the stop waits as today, when a side changes anything else, the import block included, when both add a declaration of the same name, or when the result does not parse or is not gofmt-clean. Name I-0409, S-0809, S-0812, MS-0628, and D-0644 in its context. The ADR file takes the next free number, which is why the story keeps the `design/adrs` folder touch. It waits for nothing.

## Done when

- The ADR is written with status accepted, names ADR-0175 as the decision it refines, and `design/adrs/README.md` lists it.
- `flai check --strict` and the markdown lint pass on both files.

## Notes
