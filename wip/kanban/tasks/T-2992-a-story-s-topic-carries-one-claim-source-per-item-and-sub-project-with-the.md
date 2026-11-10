---
id: T-2992
type: task
nature: improvement
title: A story's topic carries one claim source per item and sub-project, with the count of the paths that reach it
status: backlog
parent: S-0699
owner: alex
created: 2026-11-10T10:30:27Z
updated: 2026-11-10T10:30:27Z
transitions: []
stream: S-0699
tags: [flai, cli, go]
touches: [flai/internal/topics/story.go, flai/internal/topics/story_test.go, flai/internal/context/context_test.go, flai/internal/mcpserver/prime_test.go]
---
# T-2992 A story's topic carries one claim source per item and sub-project, with the count of the paths that reach it

## Work

`topics.OfStory` (`flai/internal/topics/story.go`) adds one `claim` source for each touches entry that reaches a sub-project, under every topic that sub-project gives. A story with 50 touches across `flai` and `flaiover` gets about 130 sources. Those sources took 12 to 25 KB of part 1 of S-0425's and S-0586's packs (I-0261). Later parts carry an empty `topics` list already (`flai/internal/context/parts.go`, `open`), so part 1 is the only part this cost.

Fold the claim sources of one topic into one source per item and sub-project. Keep its `kind`, `item`, and `project`, and replace `via` with a count of the touches entries that reached the sub-project, as a new field such as `paths`. Keep the first entry as `via` only if a reader needs it: find the readers with `grep` before you choose. Every other kind of source stays as it is. The text summary at the end of `story.go`, which `flai show` and the text pack print, already counts the touches per sub-project. Keep its words the same, now computed from the count.

`flai show --json` gives the same `topics`, so its shape changes the same way. Nothing outside flai reads the sources: flaiover reads none.

This task waits for no other task. No other task of S-0699 touches these files, so it runs beside the briefs tasks.

## Done when

- A test in `flai/internal/topics/story_test.go` gives a story 40 touches under one sub-project that reaches three topics. Each topic then has one claim source for the story, with a count of 40.
- A pack for that story encodes its `topics` in under 2 KB of JSON, and a test in `flai/internal/context/context_test.go` or `flai/internal/mcpserver/prime_test.go` asserts that bound.
- The text that `flai show` and `flai prime` print for a story's topics is unchanged.
- `flai test` passes on the changed files.

## Notes
