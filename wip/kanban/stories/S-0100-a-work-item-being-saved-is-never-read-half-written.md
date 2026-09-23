---
id: S-0100
type: story
nature: remediation
title: A work item being saved is never read half-written
status: review
owner: alex
created: 2026-09-23T05:10:18Z
updated: 2026-09-23T05:15:16Z
transitions:
  - to: ready
    at: 2026-09-23T05:10:28Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T05:10:28Z
    by: system-flow
  - to: review
    at: 2026-09-23T05:15:16Z
    by: system-flow
tags: []
touches: [flai/internal]
---

# S-0100 A work item being saved is never read half-written

## Goal

A work item's file is replaced in one step when it is saved, so that whatever reads it at the same moment (the MCP server's inbox and wait_for_events, flai serve's watcher, the dashboard through flai) sees the old content or the new, never an empty or half-written file.

## Acceptance criteria
- [x] Saving a work item never exposes an empty or partial file to a reader
- [x] The same holds for the other files flai rewrites that agents and the dashboard read while they change: narratives and threads
- [x] A test reads a file continuously while it is saved over and over, and fails on the old way of writing

## Tasks
- T-0355 Items, narratives, and threads are replaced in one step when saved

## Notes
Found in CI (flai workflow, publish of flai 1.12.1): TestWaitForEventsReportsAnEventWhileHeld failed with "S-0001-story.md: no front matter". `Repo.Save` wrote with `os.WriteFile`, which truncates the file before writing it, and the held wait read it in between. Reported by the operator on 2026-09-23.
