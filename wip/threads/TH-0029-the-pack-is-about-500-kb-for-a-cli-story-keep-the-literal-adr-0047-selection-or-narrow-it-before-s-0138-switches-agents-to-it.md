---
id: TH-0029
title: The pack is about 500 KB for a cli story. Keep the literal ADR-0047 selection, or narrow it before S-0138 switches agents to it?
anchor:
  path: wip/kanban/stories/S-0137-flai-prime-story-adds-the-design-tech-and-adrs-a-story-s-topics-links-and-ranking-select-with-a-catalog-of-the-rest.md
  item: S-0137
status: open
participants: [agent-S-0137]
created: 2026-09-29T00:45:52Z
updated: 2026-09-29T00:45:52Z
---

# TH-0029 The pack is about 500 KB for a cli story. Keep the literal ADR-0047 selection, or narrow it before S-0138 switches agents to it?

On wip/kanban/stories/S-0137-flai-prime-story-adds-the-design-tech-and-adrs-a-story-s-topics-links-and-ranking-select-with-a-catalog-of-the-rest.md.

## Entries

### 2026-09-29T00:45:52Z agent-S-0137
`flai prime --story` now does everything ADR-0047 lists. On this repository a cli story's pack is about 500 KB, roughly 125k tokens. `flai prime --cat` is 56 KB.

| Story | Pack | Design by topics | Linked (ADRs, mostly one step) | Ranked |
|-------|------|------------------|--------------------------------|--------|
| S-0137 | 523 KB | 22 docs, ~260 KB | 34 | 10 |
| S-0125 | 490 KB | 22 docs | 33 | 10 |
| S-0116 (cli+dashboard) | 548 KB | 29 docs | 32 | 10 |

Why it is so large:
1. Topics on design files are whole-file and broad: seven `design/system` files are `[all]`, and eight more are `cli` (`flai-cli.md` alone is 80 KB).
2. I follow "the ADRs a selected section links" literally, from sections selected by topics too. That adds about 28 ADRs, about 90 KB, mostly from `flai-cli.md § Commands`.

Replay over the archive: 56 stories named 93 ADRs that existed when each was created. Without counting the story's own links:
- By topics alone the pack reaches 0 of them, because no ADR has topics yet.
- Adding one step from topic-selected sections reaches 91.
- Adding ranking reaches 93.
- Topics and ranking without that step reach 68.

So the step is what finds the ADRs, but only because topics select nearly all of `design/system`. Ranking alone reaches 47 of the 93.

I recommend (A): keep S-0137 as ADR-0047 says (no budget, step from every selected section) and narrow the topics before S-0138 switches agents. The single biggest win is heading topics on `flai-cli.md § Commands` and on the `[all]` files. (B) is to step only from what the story links (430 KB, and ranking-plus-links recall falls). (C) is to add a budget, which ADR-0047 declined.

I am going on with (A). Nothing here blocks S-0137's criteria. Say if you want B or C.
