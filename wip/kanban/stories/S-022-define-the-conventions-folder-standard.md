---
id: S-022
type: story
nature: feature
title: Define the conventions folder standard
status: backlog
parent: E-005
owner: alex
created: 2026-09-15T18:26:31Z
updated: 2026-09-15T18:26:31Z
transitions: []
tags: [conventions]
---

# S-022 Define the conventions folder standard

## Goal
Decide where conventions live, how they are structured, and how they relate to CLAUDE.md, design/system, and the template, and record it so the other stories build on a settled base.

## Acceptance criteria
- [ ] An ADR decides the folder (proposed: a fourth top-level documentation folder `conventions`, renameable through `system-flow.yaml` layout, superseding ADR-0002 in part) and the one-file-per-topic structure with a README index in read order
- [ ] design/system has a `conventions.md` living document covering purpose, topic list, file format (front matter with title, updated, audience: agent, priority), the project-extension marker, and how CLAUDE.md, conventions, and design/system divide responsibility (CLAUDE.md is the short map, conventions are the norms, design/system is the design)
- [ ] repository-layout.md, project-manifest.md, template.md, and documentation-standard.md updated
- [ ] The topic list is agreed with the operator before S-023 starts

## Tasks

## Notes
- Candidate topics: session-start (priming and recovery), communication (how to report, when to ask, when to proceed), work-management (pull, WIP, sizing, narratives, definitions of ready and done), decisions (when an ADR, when a living-doc edit), documentation (where things go, front matter, style), code-quality (tests, lint, dependencies), git (branches, commit messages with story IDs, when to commit, never force), safety (secrets, destructive actions, external services), tooling (use flai, never hand-edit front matter).
- Relationship to existing docs: CLAUDE.md currently carries some norms inline (narrative obligations, definition of done). Those move to conventions and CLAUDE.md points at them.
