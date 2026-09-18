---
id: S-0038
type: story
nature: feature
title: Comment threads and answers as files under wip
status: done
parent: E-0006
owner: alex
created: 2026-09-17T19:46:28Z
updated: 2026-09-18T05:50:24Z
transitions:
  - to: ready
    at: 2026-09-18T04:56:34Z
    by: alex
  - to: in-progress
    at: 2026-09-18T04:56:34Z
    by: alex
  - to: review
    at: 2026-09-18T05:50:06Z
    by: alex
  - to: done
    at: 2026-09-18T05:50:24Z
    by: alex
tags: [dashboard, cli]
touches: [flai/internal/threads, flai/internal/check, flai/cmd/thread.go, flaiover/src/lib/components, flaiover/src/routes/api/threads, scripts/env.sh]
---

# S-0038 Comment threads and answers as files under wip

## Goal
Comments, questions, and answers between the designer and agents are files under `wip/threads/`, one per thread, anchored to a document, heading, or work item, validated by `flai check`, rendered by the dashboard, and readable by an agent with nothing but the repository.

## Acceptance criteria
- [x] `wip/threads/<T-id or slug>.md` with front matter: `id`, `anchor` (path plus optional heading or item ID), `status` (open, answered, resolved), `participants`, `created`, `updated`; body is a log of dated entries with an author (`alex`, or the agent and session)
- [x] `flai thread new|reply|resolve|list` create and update threads; the anchor is validated against an existing file, heading, or item; `flai check` validates threads and warns on open threads anchored to archived items
- [x] The dashboard shows threads inline next to their anchor on document and item pages and lets the designer post; open threads on a story appear in the agent narrative's `## Open questions` mirror
- [x] ADR-0020 accepted (files plus MCP as the communication model; threads are the durable channel); layout gains `wip/threads/` in the manifest and template; docs/users describes threads

## Tasks
- T-0124 threads package: TH-nnnn files under wip/threads with anchor, status, participants, entries; flai thread new|reply|resolve|list|show; narrative Open questions mirror
- T-0125 flai check: thread schema, file name, anchor existence, heading presence, open threads on archived items; layout subfolder; template and design docs
- T-0126 flaiover: threads reader and API (list, new, reply, resolve via flai), Threads component on document and item pages
- T-0127 Tests on both sides; docs/users threads section; worktree-aware scripts/env.sh

## Notes
Threads are the human-to-agent channel that MCP (S-0039) exposes with low latency. Anchoring by heading survives edits better than by line.
