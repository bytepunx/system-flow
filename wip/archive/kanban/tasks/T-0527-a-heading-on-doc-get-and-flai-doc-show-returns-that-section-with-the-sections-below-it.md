---
id: T-0527
type: task
nature: feature
title: A heading on doc_get and flai doc show returns that section with the sections below it
status: done
parent: S-0147
owner: alex
created: 2026-09-29T05:00:19Z
updated: 2026-09-29T05:06:24Z
transitions:
  - to: ready
    at: 2026-09-29T05:00:26Z
    by: agent-S-0147
  - to: in-progress
    at: 2026-09-29T05:04:09Z
    by: agent-S-0147
  - to: done
    at: 2026-09-29T05:06:24Z
    by: agent-S-0147
stream: S-0147
tags: []
touches: [flai/internal/context/design.go, flai/internal/mcpserver, flai/cmd/doc.go]
---
# T-0527 A heading on doc_get and flai doc show returns that section with the sections below it

## Work

- Find a section by heading: its text, a heading path (`Commands › flai prime`, `>` also taken), or its anchor slug; a path matches at its end; an ambiguous heading is refused naming the matches; an unknown heading is refused naming the headings there are.
- `doc_get` takes `heading` and returns the section and those below it; `flai doc show <path> --heading "<heading>"` prints it, and with `--json` prints path, heading, line, size, and text, never the save hash.
- Without a heading both behave as today.

## Done when

- Behavior tests cover a match by text, path, and slug, ambiguity, an unknown heading, and no heading, and `make test` passes.

## Notes
