---
id: TH-0010
title: "Criterion 1: may someone attending still hold a newly ready story back?"
anchor:
  path: wip/kanban/stories/S-0116-when-a-story-moves-to-ready-if-there-is-available-capacity-it-should-be-assigned-to-an-agent.md
  item: S-0116
status: open
participants: [system-flow]
created: 2026-09-24T08:46:00Z
updated: 2026-09-24T08:46:00Z
---

# TH-0010 Criterion 1: may someone attending still hold a newly ready story back?

On wip/kanban/stories/S-0116-when-a-story-moves-to-ready-if-there-is-available-capacity-it-should-be-assigned-to-an-agent.md.

## Entries

### 2026-09-24T08:46:00Z system-flow
flai serve already starts an agent for each ready story while the in-progress limit has room. The one exception is while someone is attending, meaning an MCP cursor or a narrative was written in the last 6 minutes. That exception is why S-0115 and S-0116 sat in ready today. Every agent connected to MCP as `system-flow`, so their activity looked like someone attending (I-0037). The serve log shows "an agent is attending the project (.flai-cache/mcp/system-flow.json ...)" at every look from 08:37 to 08:44.

S-0114, now in review, fixes the name. It also caps the hold at the attended window (ADR-0042). A story entering ready with room starts at once when nobody is attending, and at most 6 minutes later when someone is.

Question: does that meet criterion 1, or should a story entering ready with room start at once even when someone is attending?

Recommended: (a) keep ADR-0042's capped hold. Criterion 1 is then met by S-0114, plus a test here that pins it. S-0116 itself delivers criterion 2: changing the agent of a ready story whose agent failed starts another one. Today that only happens when the story re-enters ready.

Alternative: (b) drop the hold for stories that have room. An idle agent holding wait_for_work would race flai serve for the story. The loser's move to in-progress is refused, and it pulls the next story. This needs a new ADR to supersede that part of ADR-0042.

Either way, S-0116 builds on S-0114's launcher changes, so please accept S-0114 first.
