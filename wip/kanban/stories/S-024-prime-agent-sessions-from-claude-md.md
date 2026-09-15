---
id: S-024
type: story
nature: feature
title: Prime agent sessions from CLAUDE.md
status: backlog
parent: E-005
owner: alex
created: 2026-09-15T18:26:32Z
updated: 2026-09-15T18:26:32Z
transitions: []
tags: [conventions, template]
---

# S-024 Prime agent sessions from CLAUDE.md

## Goal
The template's CLAUDE.md instructs every agent session to load the conventions first, so norms are in context before any work starts, and the root CLAUDE.md of this repository does the same.

## Acceptance criteria
- [ ] template/root/CLAUDE.md.tmpl opens with a "Prime your session" section: read design/conventions/README.md and every file it lists in order, then wip/agents/index.md, then the board, before any change
- [ ] The section says what to do when a convention conflicts with a user instruction (the user wins, the conflict is logged in the narrative and raised as an open question) and when a convention seems wrong (propose an edit, do not silently deviate)
- [ ] Baseline CLAUDE.md shrinks: layout table, priming section, pointers; the norms themselves are in design/conventions/
- [ ] This repository's CLAUDE.md is re-rendered from the template with its project section preserved
- [ ] Rendering the template with flai new produces a project whose CLAUDE.md and conventions pass flai check

## Tasks

## Notes
