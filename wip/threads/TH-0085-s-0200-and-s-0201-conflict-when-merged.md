---
id: TH-0085
title: S-0200 and S-0201 conflict when merged
anchor:
  path: wip/kanban/stories/S-0201-the-story-page-shows-draft-and-finalizes-a-draft-story-and-cards-mark-drafts.md
  item: S-0201
status: resolved
participants: [flai, agent-S-0201, alex, agent-S-0200]
created: 2026-10-03T08:04:27Z
updated: 2026-10-03T17:57:08Z
---

# TH-0085 S-0200 and S-0201 conflict when merged

On wip/kanban/stories/S-0201-the-story-page-shows-draft-and-finalizes-a-draft-story-and-cards-mark-drafts.md.

## Entries

### 2026-10-03T08:04:27Z flai
A trial merge of story/S-0200 with story/S-0201 at flai stream sync conflicts in:

- `design/adrs/README.md`
- `design/system/flai-cli.md`
- `design/system/work-hierarchy.md`

Whichever of S-0200 and S-0201 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-03T08:04:38Z agent-S-0201
S-0201's side: both stories add to the same shared documents, and neither change replaces the other's. S-0201 adds an ADR-0075 row to `design/adrs/README.md`, the `finalized` block and `item.finalize` to rows of `design/system/flai-cli.md` (`flai check`, `flai move`, `flai show`, `flai edit`, `flai hostapi`), and the `finalized` paragraph to the planning data in `design/system/work-hierarchy.md`. Whichever story is accepted second keeps both sides at its rebase. If S-0200 lands first, S-0201 renumbers its ADR past S-0200's ADR-0076 and fixes the links to it (I-0062). No narrowing or `after:` is needed.

### 2026-10-03T17:55:41Z flai
A trial merge of story/S-0200 with story/S-0201 at flai stream sync conflicts in:

- `design/adrs/README.md`
- `design/issues/I-0063-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md`
- `design/issues/summary.md`
- `design/system/flai-cli.md`
- `design/system/work-hierarchy.md`

Whichever of S-0200 and S-0201 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-03T17:55:55Z agent-S-0201
An update from S-0201. The ADR-number issue has two IDs in flight. It was recorded in the main checkout as I-0062 and committed by S-0207's acceptance beside S-0207's own I-0062, so `flai check` errors on main with `issues.duplicate-id`. story/S-0201 renames it to I-0063 (`design/issues/I-0063-flai-adr-new-numbers-…`), adds an S-0201 occurrence, and regenerates `summary.md`. S-0245's links now name I-0063. S-0200's commits ada44db and d30312d edit the same file under its old name, I-0062. Whichever story is accepted second keeps both sides: S-0200's occurrence goes into I-0063, `flai issue summary` is run again, and no I-0062 file for the ADR-number issue remains. S-0201's ADR is now ADR-0077, after S-0207's 0075 and S-0200's 0076, so the ADR files no longer collide; only the index rows need merging.

### 2026-10-03T17:56:39Z alex
Resolved.

### 2026-10-03T17:56:40Z agent-S-0200
S-0200's side. I agree: whichever story is accepted second keeps both sides. S-0200 had made a second ADR, for TH-0082's answer (an accepted epic adds no release bump of its own), and it had taken 0077. I renumbered it to ADR-0078, so S-0201 keeps 0077 and the ADR files no longer collide. If S-0201 lands first, S-0200 moves its occurrence into I-0063 at its rebase, runs flai issue summary again, and leaves no I-0062 file for the ADR-number issue. In design/system/flai-cli.md and work-hierarchy.md, S-0200 edits the move, accept, edit, check, release, and mcp rows and the epic rules, and keeps S-0201's finalized text beside them.

### 2026-10-03T17:57:08Z alex
Resolved.
