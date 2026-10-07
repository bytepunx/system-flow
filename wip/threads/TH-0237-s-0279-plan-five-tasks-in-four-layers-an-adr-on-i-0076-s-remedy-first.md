---
id: TH-0237
title: "S-0279 plan: five tasks in four layers, an ADR on I-0076's remedy first"
anchor:
  path: wip/kanban/stories/S-0279-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md
  item: S-0279
status: resolved
participants: [planner-S-0279, alex]
created: 2026-10-07T01:22:32Z
updated: 2026-10-07T01:38:02Z
---

# TH-0237 S-0279 plan: five tasks in four layers, an ADR on I-0076's remedy first

On wip/kanban/stories/S-0279-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md.

## Entries

### 2026-10-07T01:22:32Z planner-S-0279
S-0279 is planned: five tasks in four layers, 18 file touches plus one folder touch, a 45m forecast, and a cost of delay of 112.50 USD a week.

I-0076's 12 instances show three causes:

1. Every close-out records every overlap. ADR-0085 marks every `wip.overlap` outside the story, so each close-out records the overlaps of stories it is not part of, such as S-0251 recording S-0228 against S-0273.
2. Duplicates. `overlap()` pairs a story with another story's task as well, such as S-0262 against T-0877.
3. Raw touches. `overlap()` reads raw touches, not the claim the pull hold reads, so it reports folders that the tasks narrow, such as `flai/internal/serve` and `flai/internal/hostapi`.

Recommended remedy: a close-out records no `wip.overlap` and lists only the overlaps naming the story, and `wip.overlap` compares story claims, once per pair. The alternative T-1141 weighs is to keep recording, but only the overlaps that name the closing story.

The tasks and their layers:

| Layer | Task | After |
|-------|------|-------|
| 1 | T-1141 ADR refining ADR-0085 and ADR-0096 | none |
| 2 | T-1143 `wip.overlap` over claims, story IDs on each finding | T-1141 |
| 3 | T-1145 close-out records no overlap, notes only those naming the story | T-1143 |
| 4 | T-1149 close I-0076 | T-1145 |
| 4 | T-1150 convention (both copies), design, guides, reference | T-1145 |

T-1145 waits for T-1143 rather than running beside it: it scopes the findings by the story IDs T-1143 adds to `Finding` in `check.go`.

The assumptions I made:

- The remedy above is the one T-1141's ADR takes. The story's agent may choose the alternative, and T-1143 and T-1145 follow the ADR.
- The ADR goes to you for acceptance on a thread, and the code tasks go on meanwhile.
- `flai/internal/hostapi/people.go` (the designer's inbox) and `flai/internal/itemedit/itemedit.go` read `wip.overlap` findings but need no change. T-1143 runs their tests and widens its touches if one does.
- S-0277, which remediates I-0073 and is ready, is likely to change `scope.go`, `flai/cmd/check.go`, and the close-out rule in `work-management.md` too. The pull hold keeps the two apart, and whichever goes second rebases. I set no `after` between them.
- I added the tags `flai` and `template` and the topics `cli`, `conventions`, and `template` to the story. The template convention changes, and `flai check` asked for a component tag.

I changed neither the story's goal nor its criteria. There is nothing to split, merge, or drop.

### 2026-10-07T01:38:02Z alex
Resolved.
