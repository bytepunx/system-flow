---
id: TH-1023
title: "S-0699: may a pack over its budget list its link-step ADR briefs in the catalog, reversing TH-0032's \"briefs are never cut\"?"
anchor:
  path: wip/kanban/stories/S-0699-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md
  item: S-0699
status: open
participants: [planner-S-0699, orchestrator]
created: 2026-11-10T10:30:08Z
updated: 2026-11-10T11:02:24Z
---

# TH-1023 S-0699: may a pack over its budget list its link-step ADR briefs in the catalog, reversing TH-0032's "briefs are never cut"?

On wip/kanban/stories/S-0699-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md.

## Entries

### 2026-11-10T10:30:08Z planner-S-0699
**Recommendation: yes to both questions below.** Answering no to either drops that task. The story then fixes only the topic sources.

I measured today's pack (`flai prime --story --json`, 80 KB budget):

| Story | Pack | Conventions | Named | Briefs | Of the briefs, ADRs one link step away | `also` lists | Topic sources |
|-------|------|-------------|-------|--------|------|------|------|
| S-0425 | 840 KB | 106 KB | 362 KB | 342 KB, 288 briefs | 167 KB, 155 briefs | 134 KB | 12 KB, 136 sources |
| S-0586 | 552 KB | 106 KB | 74 KB | 342 KB, 295 briefs | 164 KB, 160 briefs | 132 KB | 13 KB, 130 sources |

One finding corrects the issue. Later parts do not repeat the topic sources: since T-0987, each part after the first carries an empty `topics` list (`flai/internal/context/parts.go`, `open`). The list costs part 1 alone, which is why part 1 held two conventions. The packs are large because:

1. **Briefs are never cut** (TH-0032, ADR-0241). On a pack over its budget, as every story's pack is from its conventions, every brief still loads. About half of them are ADRs reached only by one link step from a design document.
2. **Each brief's JSON carries every one of its other reasons** in `also`, which is 39% of the brief bytes.
3. **The topic sources hold one entry per claimed path per topic.**

Questions:

- **Q1. On a pack over its budget, may the ADR briefs reached only by a link step go to the catalog (path and title) instead of loading?** The briefs that topics select would still load. This reverses TH-0032 for those briefs and needs an ADR that refines ADR-0241. It saves about 150 KB a pack, about four parts. Alternative: keep every brief, and the story then barely changes the part count.
- **Q2. May a brief's JSON carry the count of its other reasons in place of the `also` list?** Loaded items would keep their list. This changes the line "the JSON keeps every reason" in `design/system/agent-context.md`. It saves about 130 KB a pack. Alternative: keep the list, capped at three reasons.

Fix 3 needs no answer: I am planning it now. One claim source per item and sub-project, with a count of its paths, takes the topics from 12 KB to under 2 KB.

### 2026-11-10T11:02:24Z orchestrator (recommendation)
**Recommended:** Q1 no, Q2 yes.

**Q1, link-step ADR briefs to the catalog: no.** You chose "briefs are never cut" twice: on TH-0032, and again on TH-0805, recorded in ADR-0241 as option 1. Nothing in I-0261 shows a cost that you did not have before you then. Keep every brief. S-0699 then has no ADR task for it.

**Q2, a count of a brief's other reasons in place of `also`: yes.** It saves about 130 KB a pack, which is 39% of the brief bytes. It cuts no brief and no document. A loaded item keeps its full list. It changes one line of `design/system/agent-context.md`, not a decision you recorded. The alternative, keeping the list capped at three reasons, saves less and still repeats.

With Q1 no and Q2 yes, S-0699 has three tasks: T-2992, the `also` count task, and T-2996.

Source: design/adrs/0241-a-pack-whose-conventions-alone-exceed-its-budget-still-loads-what-is-named-and.md
