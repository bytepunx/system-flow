---
id: T-0500
type: task
nature: feature
title: flai prime --story prints the pack with a header, --json, and refuses unknown or archived stories
status: done
parent: S-0136
owner: alex
created: 2026-09-28T22:57:55Z
updated: 2026-09-28T23:02:00Z
transitions:
  - to: ready
    at: 2026-09-28T22:57:59Z
    by: agent-S-0136
  - to: in-progress
    at: 2026-09-28T22:59:59Z
    by: agent-S-0136
  - to: done
    at: 2026-09-28T23:02:00Z
    by: agent-S-0136
stream: S-0136
tags: []
touches: [flai/cmd/prime.go]
---
# T-0500 flai prime --story prints the pack with a header, --json, and refuses unknown or archived stories

## Work

- `flai prime --story S-nnnn`: resolve the story (unknown, archived, or not a story fails naming the ID), work out its topics with `topics.ForStory`, build the conventions part with `flai/internal/context`, print a header with the story, its topics and their sources, and the size of the pack, then the pack.
- `--json`: the story, its topics with sources, the size, and per convention the sections kept and left out with their topics.
- Plain `flai prime` unchanged.
- Command tests: a fixture project with narrowed conventions shows sections dropped and listed; with every convention `[all]` the text after the header equals `flai prime --cat`; unknown and archived IDs fail naming them.

## Done when

- `go test -race -short ./cmd/` passes and the tests above exist.

## Notes
