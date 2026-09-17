---
id: S-0040
type: story
nature: feature
title: Document editing in the dashboard with commit through flai
status: backlog
parent: E-0006
owner: alex
created: 2026-09-17T19:46:29Z
updated: 2026-09-17T19:46:29Z
transitions: []
tags: [dashboard, cli]
---

# S-0040 Document editing in the dashboard with commit through flai

## Goal
The designer edits design, docs, and work item bodies in the dashboard with a markdown editor that protects flai-owned front matter, previews with the same renderer as the explorer, and saves through flai so the change is validated and committed on main with the designer as author.

## Acceptance criteria
- [ ] An edit view per document with the body editable and front matter shown read-only where flai owns it (items, narratives, board); design and docs front matter editable with schema validation
- [ ] Save runs `flai check` on the changed file and refuses with the finding; a saved change is committed on main as `docs:` or `chore:` with the designer as author and the dashboard as co-author trailer, or left uncommitted when the project sets `dashboard.autocommit: false`
- [ ] Concurrent edits: the save carries the file's content hash; a stale hash returns a conflict with a diff and the current version
- [ ] Threads (S-0038) can be opened from a selected heading in the editor; the `touches` badge (S-0037) warns before editing a document an in-progress story touches
- [ ] Tests for save, refuse, conflict; docs/users and flaiover-dashboard.md updated

## Tasks

## Notes
Requires S-0036 (authentication) and S-0037 (branches, so agent work is not on main). Autocommit on main keeps the designer's intent visible to `flai stream sync` immediately.
