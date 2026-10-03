---
id: TH-0086
title: S-0201 and S-0207 conflict when merged
anchor:
  path: wip/kanban/stories/S-0201-the-story-page-shows-draft-and-finalizes-a-draft-story-and-cards-mark-drafts.md
  item: S-0201
status: resolved
participants: [flai, agent-S-0201]
created: 2026-10-03T08:04:27Z
updated: 2026-10-03T17:55:41Z
---

# TH-0086 S-0201 and S-0207 conflict when merged

On wip/kanban/stories/S-0201-the-story-page-shows-draft-and-finalizes-a-draft-story-and-cards-mark-drafts.md.

## Entries

### 2026-10-03T08:04:27Z flai
A trial merge of story/S-0201 with story/S-0207 at flai stream sync conflicts in:

- `design/adrs/README.md`

Whichever of S-0201 and S-0207 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-03T08:04:39Z agent-S-0201
S-0201's side: the only conflict is the ADR index row. story/S-0207 adds ADR-0075 and story/S-0201 does too, because `flai adr new` numbers from the worktree (I-0062). S-0207 is in review and will probably be accepted first. S-0201 then renumbers its ADR to the next free number at its rebase and fixes every link to it. No other path conflicts, and no `after:` is needed.

### 2026-10-03T17:55:41Z flai
Resolved: S-0207 is done, no longer open, at the sync of S-0201
