---
id: S-0024
type: story
nature: feature
title: Prime agent sessions from CLAUDE.md
status: done
parent: E-0005
owner: alex
created: 2026-09-15T18:26:32Z
updated: 2026-09-16T23:03:57Z
transitions:
  - to: ready
    at: 2026-09-16T22:44:06Z
    by: agent
  - to: in-progress
    at: 2026-09-16T22:44:06Z
    by: agent
  - to: review
    at: 2026-09-16T22:46:09Z
    by: agent
  - to: done
    at: 2026-09-16T23:03:57Z
    by: alex
tags: [conventions, template]
---

# S-0024 Prime agent sessions from CLAUDE.md

## Goal
The template's CLAUDE.md instructs every agent session to load the conventions first, so norms are in context before any work starts, and the root CLAUDE.md of this repository does the same.

## Acceptance criteria
- [x] template/root/CLAUDE.md.tmpl opens with a "Prime your session" section: read design/conventions/README.md and every file it lists in order, then wip/agents/index.md, then the board, before any change
- [x] The section says what to do when a convention conflicts with a user instruction (the user wins, the conflict is logged in the narrative and raised as an open question) and when a convention seems wrong (propose an edit, do not silently deviate)
- [x] Baseline CLAUDE.md shrinks: layout table, priming section, pointers; the norms themselves are in design/conventions/
- [x] This repository's CLAUDE.md is re-rendered from the template with its project section preserved
- [x] Rendering the template with flai new produces a project whose CLAUDE.md and conventions pass flai check

## Tasks
- T-0048 Rewrite template CLAUDE.md.tmpl around a priming section
- T-0049 Re-render this repository's CLAUDE.md with its project section preserved
- T-0050 Verify render and check; update design notes

## Notes
- Rendered projects are now markdownlint-clean out of the box; three templates gained conditional description lines and the README URL is wrapped.
