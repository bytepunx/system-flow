---
id: T-1720
type: task
nature: feature
title: "An ADR records the role file: where it lives, its front matter, its body as the prompt, the marker, and the named conditions it may reference"
status: backlog
parent: S-0414
owner: alex
created: 2026-10-16T05:42:52Z
updated: 2026-10-16T05:42:52Z
transitions: []
stream: S-0414
tags: [cli]
touches: [design/adrs, design/adrs/README.md]
---
# T-1720 An ADR records the role file: where it lives, its front matter, its body as the prompt, the marker, and the named conditions it may reference

## Work

Write the ADR with `flai adr new`, on TH-0471's answers to questions 2, 4, and 6. It records:

- Where a role file lives (`design/roles/<name>.md` by default) and that the manifest may name the folder.
- Each front matter field and its meaning: `name` (the activity document's and the page's), `role` (the verb `flai prime --role` and `FLAI_ROLE` take), `definition`, `action`, `trigger`, `item`, `topics`, `permissions`, `tools`, `commands`, `files`, `works_story`, and `charge`, which S-0425 uses.
- The body as the harness prompt, and the baseline marker.
- The named conditions flai implements for rules that are not plain allow-lists. Draw the list from `flai/internal/guard/guard.go`'s planner, orchestrator, and analyzer rules.

First layer: every other task of the story reads the format this decides.

## Done when

The ADR is accepted and listed in `design/adrs/README.md`. `flai check --strict` reports nothing on it.

## Notes

The ADR's file name is not known until `flai adr new` makes it, so the task touches the shared `design/adrs` folder.
