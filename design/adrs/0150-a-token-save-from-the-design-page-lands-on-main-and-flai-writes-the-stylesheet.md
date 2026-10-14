---
id: ADR-0150
title: "A token save from the Design page lands on main, and flai writes the stylesheet it generates from the tokens in the same commit"
status: accepted
date: 2026-10-14
supersedes: []
superseded_by: []
refines: [ADR-0023, ADR-0147]
topics: [design, cli, dashboard]
---

# ADR-0150 A token save from the Design page lands on main, and flai writes the stylesheet it generates from the tokens in the same commit

## Context

[ADR-0023](0023-documents-are-saved-through-flai.md) commits the one path a dashboard save writes. [ADR-0147](0147-the-dashboard-reads-and-saves-a-text-file-that-is-not-markdown-through-file-get.md) saves the token file through `file.save` by the same rules: any finding the save introduces refuses it. Code gets the tokens through a stylesheet flai generates from them, and `design.stale` is a finding when that stylesheet differs from the tokens. A token save that left the stylesheet as it was would introduce `design.stale` and be refused. The finding, [design-system-tool.md](../system/design-system-tool.md), recommended choice 8 below; on TH-0445, on 2026-10-14, the operator answered only choice 1, so it stands as recommended.

## Decision

A token save from the Design page lands on main, and flai writes the stylesheet it generates from the tokens in the same commit.

1. **Where.** The operator's token save from the Design page, with no story named, is a `file.save` on the main checkout's current branch.
2. **The stylesheet with it.** Whenever flai writes the token file, by `file.save` or by a `flai design` command, it writes the stylesheet the manifest names, generated from the new tokens, and `flai check` runs with both in place. The commit holds both paths.
3. **On a story's branch.** A save that names a story ([ADR-0145](0145-a-dashboard-save-may-name-a-story-and-is-then-committed-on-that-story-s-branch.md)) writes and commits both on that story's branch.
4. **`design.stale`** is left to hand edits: of the token file, or of the stylesheet.

## Consequences

- A file of a UI sub-project changes on main with no story and no close-out. This is a new kind of dashboard write. The operator reviews every change on this board, and the stylesheet is generated, never hand-written.
- A revert of the commit reverts the tokens and the stylesheet together.
- A save stays one commit, so the main checkout never holds tokens without their stylesheet.

## Alternatives considered

- **Token saves only on a proposal story's branch.** Every code change would pass a story's close-out, and the page would be for trying only. The shortest path from the operator's decision to the code was preferred.
- **Save the tokens alone and let an agent regenerate the stylesheet.** `design.stale` would refuse the save.
