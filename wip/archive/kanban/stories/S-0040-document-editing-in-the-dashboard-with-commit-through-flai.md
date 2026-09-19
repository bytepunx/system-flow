---
id: S-0040
type: story
nature: feature
title: Document editing in the dashboard with commit through flai
status: done
parent: E-0006
owner: alex
created: 2026-09-17T19:46:29Z
updated: 2026-09-19T02:43:23Z
transitions:
  - to: ready
    at: 2026-09-19T01:56:10Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:10:12Z
    by: alex
  - to: review
    at: 2026-09-19T02:32:53Z
    by: alex
  - to: done
    at: 2026-09-19T02:43:23Z
    by: alex
tags: [dashboard, cli]
---

# S-0040 Document editing in the dashboard with commit through flai

## Goal
The designer edits design, docs, and work item bodies in the dashboard with a markdown editor that protects flai-owned front matter, previews with the same renderer as the explorer, and saves through flai so the change is validated and committed on main with the designer as author.

## Acceptance criteria
- [x] An edit view per document with the body editable and front matter shown read-only where flai owns it (items, narratives, board); design and docs front matter editable with schema validation
- [x] Save runs `flai check` on the changed file and refuses with the finding; a saved change is committed on main as `docs:` or `chore:` with the designer as author and the dashboard as co-author trailer, or left uncommitted when the project sets `dashboard.autocommit: false`
- [x] Concurrent edits: the save carries the file's content hash; a stale hash returns a conflict with a diff and the current version
- [x] Threads (S-0038) can be opened from a selected heading in the editor; the `touches` badge (S-0037) warns before editing a document an in-progress story touches
- [x] Tests for save, refuse, conflict; docs/users and flaiover-dashboard.md updated

## Tasks
- T-0172 ADR-0023 and living design: documents are saved through flai, with ownership, a content hash, a check, and an autocommit
- T-0173 flai doc show and flai doc save, and the dashboard.autocommit manifest key
- T-0174 Dashboard API: read a document for editing and save it, with conflict and refusal responses
- T-0175 Editor view: body and front matter, preview, save, findings, and conflict handling
- T-0176 Editor: open a thread from the heading under the cursor, and warn when a story touches the document
- T-0177 User, design, and operator documentation for editing
- T-0178 Verify: all tiers, and edit, refusal, conflict, and commit through a real container

## Notes
Requires S-0036 (authentication) and S-0037 (branches, so agent work is not on main). Autocommit on main keeps the designer's intent visible to `flai stream sync` immediately.

Decided when pulled, 2026-09-19 (ADR-0023): saving is `flai doc save`, the editor learns what it may edit from `flai doc show`, and the dashboard's server holds none of the rules. "Schema validation" of design and docs front matter is the check's document rules. Because every one of those rules is a warning and this repository's gate is the strict check, a save is refused for any finding the edit introduces, at any level, not only for errors. Creating, renaming, and deleting documents are out of scope. Whether the editor should refuse edits to accepted ADRs is an open question in the narrative.

Verification, 2026-09-19. `make flai-test` passed (golangci-lint 0 issues; behavior, integration, smoke; markdown lint); flaiover lint, svelte-check 0 errors, 19 files and 130 tests, and a production build passed. In a browser, through the dev server against a scratch git project with a throwaway token: the touches warning showed and Save waited for the acknowledgement; the heading under the caret was offered for a thread; the preview rendered a Mermaid diagram live; a save committed as the scratch identity with a `docs:` subject from the message, the dashboard trailer, one path, and `updated` set to today; a change made on the host produced the conflict panel with its diff, and "Save mine over it" committed; a story body was saved at 390 pixels with read-only front matter. Through a real container built from the branch, started by the branch's `flai dashboard` on its own name and port, by API: load, save with the commit authored by the identity passed into the container, 422 with `doc.title` and the file restored and the tree clean, 409 with the current content and a diff naming `(current)` and `(yours)`, a front matter change to a story refused and its body edit committed as `chore: [S-0001]`, a generated file reported as not editable, and with `dashboard.autocommit: false` the save left uncommitted with HEAD unchanged. The operator's dashboard and repository were not touched. Not done: the editor was not used in a browser against the container build, only against the dev server; the container was exercised through the API the editor calls.
