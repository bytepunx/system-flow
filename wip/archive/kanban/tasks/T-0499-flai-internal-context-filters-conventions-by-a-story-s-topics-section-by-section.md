---
id: T-0499
type: task
nature: feature
title: flai/internal/context filters conventions by a story's topics, section by section
status: done
parent: S-0136
owner: alex
created: 2026-09-28T22:57:54Z
updated: 2026-09-28T22:59:59Z
transitions:
  - to: ready
    at: 2026-09-28T22:57:58Z
    by: agent-S-0136
  - to: in-progress
    at: 2026-09-28T22:57:59Z
    by: agent-S-0136
  - to: done
    at: 2026-09-28T22:59:59Z
    by: agent-S-0136
stream: S-0136
tags: []
touches: [flai/internal/context]
---
# T-0499 flai/internal/context filters conventions by a story's topics, section by section

## Work

- Add `flai/internal/context`: for each convention (README first, then read order), parse it with `topics.Parse` (default `[all]`), keep the sections whose effective topics include `all` or any of the story's topics, keep front matter, the baseline marker, and the `## Project additions` heading whatever their section, and keep the heading line of a dropped ancestor above a kept section.
- Record every section kept and left out with its heading path, line, and topics; a left-out line reads `code-quality.md § Go (cli)`, one per outermost section left out.
- Render the pack body: the same headers as `flai prime --cat`, the open-issues table, then the left-out lines when there are any.
- Behavior tests with in-memory fixture conventions: every section `[all]` reproduces the file byte for byte; narrowed headings drop and list sections; marker and additions heading survive.

## Done when

- `go test -race -short ./internal/context/` passes and lint is clean for the package.

## Notes
