---
id: T-0635
type: task
nature: remediation
title: What flai writes in wip is refused or cleaned before it breaks the lint
status: ready
parent: S-0179
owner: arobson
created: 2026-10-01T08:23:14Z
updated: 2026-10-01T08:23:32Z
transitions:
  - to: ready
    at: 2026-10-01T08:23:32Z
    by: claude-opus-5-5
stream: S-0179
tags: []
touches: [flai/internal/workitem, flai/internal/threads, flai/internal/itemedit, flai/internal/mcpserver]
---
# T-0635 What flai writes in wip is refused or cleaned before it breaks the lint

## Work

Titles lose trailing heading punctuation (MD026) when an item or thread is made or retitled. `workitem` Create lints the document it is about to write, and threads and narrative log entries lint the document they are about to save, refusing with the rule and line anything they would introduce, so `item_new`, `flai story new`, `flai task new`, `thread_open`, `thread_reply`, `thread_resolve`, and `flai stream log` cannot write what the lint rejects.

## Done when

- [ ] A title ending in a period is written without it, in the front matter and the heading
- [ ] A body, thread entry, or log entry that breaks a rule is refused with the rule and line and nothing is written
- [ ] `make test` passes

## Notes

Part of S-0179.
