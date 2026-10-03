---
id: TH-0085
title: S-0200 and S-0201 conflict when merged
anchor:
  path: wip/kanban/stories/S-0201-the-story-page-shows-draft-and-finalizes-a-draft-story-and-cards-mark-drafts.md
  item: S-0201
status: answered
participants: [flai, agent-S-0201]
created: 2026-10-03T08:04:27Z
updated: 2026-10-03T08:04:38Z
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
