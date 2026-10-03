---
id: TH-0083
title: How should finalizing a draft record who did it?
anchor:
  path: wip/kanban/stories/S-0201-the-story-page-shows-draft-and-finalizes-a-draft-story-and-cards-mark-drafts.md
  item: S-0201
status: open
participants: [agent-S-0201]
created: 2026-10-03T07:36:21Z
updated: 2026-10-03T07:36:21Z
---

# TH-0083 How should finalizing a draft record who did it?

On wip/kanban/stories/S-0201-the-story-page-shows-draft-and-finalizes-a-draft-story-and-cards-mark-drafts.md.

## Entries

### 2026-10-03T07:36:21Z agent-S-0201
S-0201's third criterion says finalizing "records who did it in the item's transitions or `updated` as an edit by the operator, so the metrics can tell planner drafts accepted from drafts rewritten". Neither option keeps that fact today:

- `updated` is a timestamp with no author. `flai edit` notes its author only in a git-ignored file beside the MCP cursors, which the metrics cannot read.
- `transitions` record state moves. A finalize is not a move: the story stays in backlog. A transition entry with no state change would break time-in-state and every other metric that reads transitions.
- Once `draft` is cleared, nothing in the item says it was ever a draft.

**Recommendation: a `finalized` block on the story, `{by, at}`, written when the draft flag is cleared.** `flai edit --no-draft` and `flai move <story> ready --yes` stamp it with the editor (`--by`, else `FLAI_AGENT`, else the config author), as ADR-0074 stamps cost of delay and forecast. The new hostapi op `item.finalize` runs `flai edit <id> --no-draft --by=<owner>`. The metrics (S-0205) can then count drafts finalized and by whom, and tell accepted from rewritten by whether the body changed between the draft's creation and `finalized.at` (git history). This adds a front matter key, so it takes a new ADR refining ADR-0074 and raises `flai.minimum` when published, as S-0199's fields did.

Alternatives:

1. Finalize moves the story to ready (`flai move ready --yes`), so the ready transition's `by` is the finalizer. No schema change, but Finalize could no longer leave a story in the backlog, and nothing records that it was a draft.
2. A `drafted: {by, at}` block kept forever, plus `finalized`. More exact about who wrote the draft, but the creator is already in the git history and the first transition.

I am going ahead with the indicator, cards, and edit-form parts, which do not depend on this, and will build the recommendation unless you say otherwise.
