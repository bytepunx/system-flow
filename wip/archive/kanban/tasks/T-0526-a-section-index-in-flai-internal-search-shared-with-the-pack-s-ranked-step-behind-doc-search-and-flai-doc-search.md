---
id: T-0526
type: task
nature: feature
title: A section index in flai/internal/search, shared with the pack's ranked step, behind doc_search and flai doc search
status: done
parent: S-0147
owner: alex
created: 2026-09-29T05:00:18Z
updated: 2026-09-29T05:04:08Z
transitions:
  - to: ready
    at: 2026-09-29T05:00:26Z
    by: agent-S-0147
  - to: in-progress
    at: 2026-09-29T05:00:26Z
    by: agent-S-0147
  - to: done
    at: 2026-09-29T05:04:08Z
    by: agent-S-0147
stream: S-0147
tags: []
touches: [flai/internal/search, flai/internal/context/rank.go, flai/internal/mcpserver, flai/cmd/doc.go]
---
# T-0526 A section index in flai/internal/search, shared with the pack's ranked step, behind doc_search and flai doc search

## Work

- Move the section cut (down to the next heading of any level, the heading path without the title) and the stopword terms into `flai/internal/search` as a section index, and rank the pack's candidates with it so the pack and search cut the same way.
- Load the sections of every markdown file under the manifest's `design` and `docs` folders (conventions included; READMEs and the ADR template left out).
- The MCP tool `doc_search` and `flai doc search <query>`: at most twenty hits, each with path, heading path, first lines, size, and score; `--limit` and `limit` below twenty.
- The MCP tool list test and the server's instructions name `doc_search`.

## Done when

- Behavior tests cover the index, the tool, and the command, and `make test` passes.

## Notes
