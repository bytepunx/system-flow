---
id: S-0074
type: story
nature: feature
title: The dashboard reads documents, ADRs, narratives, the inbox, and search through the channel, and parses no project file
status: ready
parent: E-0003
owner: alex
created: 2026-09-20T07:26:50Z
updated: 2026-09-20T07:27:45Z
transitions:
  - to: ready
    at: 2026-09-20T07:27:45Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, design/system]
---
# S-0074 The dashboard reads documents, ADRs, narratives, the inbox, and search through the channel, and parses no project file

## Goal
The rest of the reads leave the mount: the document tree and one document, the ADR list, narratives and the activity view, the designer's inbox, and search. After this story flaiover's server reads no file of the project.

## Acceptance criteria
- [ ] Methods for the document tree with each file's front matter, one document (as `flai doc show` gives it), the ADR list with status and relations, narratives with their last log entry and open questions, and the activity view; paths are repository-relative Markdown under the manifest's folders and anything else is refused by flai
- [ ] The designer's inbox (threads awaiting them, open questions, stories in review, blocked items, overlapping touches) is composed in flai and answered as one method; it is not the agent's MCP `inbox`, and the two are named so that nobody confuses them
- [ ] Search runs in flai: an index over the same Markdown the dashboard indexes today, kept current by the watcher, answering a query with the fields and snippets the search page shows; ranking is compared with today's on this repository and differences that matter are recorded
- [ ] `repo.ts`, `search.ts`, `inbox.ts`, and `activity.ts` no longer read the filesystem, chokidar and the front matter parser leave flaiover's dependencies, and a test fails if server code imports `node:fs` for project files
- [ ] Answers of a few megabytes (the document tree of this repository, a long document) arrive within the message cap or are paged
- [ ] The design documents, `design/tech`, and the users' documentation are updated

## Tasks

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Depends on the work-item reads story. Search in Go is the largest unknown here: say in the narrative which library or approach was chosen and why.
