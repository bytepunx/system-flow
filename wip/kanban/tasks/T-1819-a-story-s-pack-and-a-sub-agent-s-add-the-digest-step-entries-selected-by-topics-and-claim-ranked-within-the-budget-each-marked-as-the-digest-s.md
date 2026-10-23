---
id: T-1819
type: task
nature: feature
title: "A story's pack and a sub-agent's add the digest step: entries selected by topics and claim, ranked within the budget, each marked as the digest's"
status: backlog
parent: S-0430
owner: alex
created: 2026-10-16T06:04:41Z
updated: 2026-10-23T16:50:28Z
transitions: []
stream: S-0430
tags: [cli]
touches: [flai/internal/context/digest.go, flai/internal/context/digest_test.go, flai/internal/context/items.go, flai/internal/context/design.go, flai/internal/context/context.go, flai/internal/context/rank.go, flai/internal/context/story.go]
after: [T-1818]
---
# T-1819 A story's pack and a sub-agent's add the digest step: entries selected by topics and claim, ranked within the budget, each marked as the digest's

## Work

- `flai/internal/context/digest.go`: select the entries the ADR names for a story (through `digest.List`), rank them with `Rank` against the story's query, and add them within the budget as items with `Step: digest`, a reason (`topics: ...`, `path: ...`, `rank n`), and a size; a compacted entry loads as its title.
- `design.go` gains `StepDigest`; `items.go` orders the step where the ADR puts it; `context.go` prints the step under its own heading in text and keeps it in JSON; `story.go` adds the step to `ForStory`, and to `ForRole` as briefs only.

Waits for T-1818, which decides the selection and the budget.

## Done when

`digest_test.go` packs a fixture story whose topics select one entry, whose touches select another, and whose budget leaves out a third, plus an absorbed entry that never loads, a compacted one that loads as its title, and a sub-agent's pack with briefs only; the existing context tests pass unchanged; `flai test flai/internal/context` passes.

## Notes

`strategic.go` is the orchestrator task's, which follows this one.
