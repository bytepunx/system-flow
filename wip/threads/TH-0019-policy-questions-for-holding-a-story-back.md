---
id: TH-0019
title: Policy questions for holding a story back
anchor:
  path: wip/kanban/stories/S-0124-survey-state-of-the-art-coordination-mechanisms.md
  item: S-0124
status: resolved
participants: [agent-S-0124, alex]
created: 2026-09-26T07:24:42Z
updated: 2026-09-26T07:57:27Z
---

# TH-0019 Policy questions for holding a story back

On wip/kanban/stories/S-0124-survey-state-of-the-art-coordination-mechanisms.md.

## Entries

### 2026-09-26T07:24:42Z agent-S-0124
Whatever design you choose, these choices are yours. Each one has my recommendation first. Answer by number; "all as recommended" is fine.

1. **Skip ahead.** When the first ready story is held, may flai start a later ready story that overlaps nothing?
   - Recommended: yes. The held story keeps its place and starts first once it is clear.
   - Alternative: no. The pull order is strict, and nothing behind a held story starts.

2. **No touches declared.** What should a story with no `touches` do?
   - Recommended: treat it as overlapping everything. It is held while anything else is in progress, with the reason "declares no touches". When it is in progress, it holds everything else.
   - Alternative: it overlaps nothing and runs in parallel, flagged amber.

3. **Who declares touches.** Who writes a story's `touches`?
   - Recommended: whoever creates the story, which is you or an agent, as today. The pulling agent may widen them when it writes the tasks. Widening while in progress can hold nothing already running. It can only hold ready stories.
   - Alternative: touches are required to move a story to ready.

4. **Hand moves.** Should `flai move S-nnnn in-progress` by hand, and item_move, refuse a held story, or only warn as they do for the WIP limit?
   - Recommended: warn only. flai serve and wait_for_work respect the hold, and a person can always choose to proceed.

5. **Component names.** Should `touches` entries that are component names map to paths?
   - Recommended: yes. Entries such as `flaiover` or `cli` map to the sub-project paths in system-flow.yaml, so `cli` and `flai/cmd` overlap. Unknown names are compared as they are.

6. **Tasks.** Should tasks' `touches` count toward their story's claim?
   - Recommended: yes. The claim is the union of the story's touches and its open tasks' touches.

### 2026-09-26T07:57:22Z alex
continue with recommendations

### 2026-09-26T07:57:27Z agent-S-0124
Resolved: Designer took every recommendation: skip ahead yes; no touches overlaps everything; touches declared by the creator, widened by the pulling agent; hand moves warn only; component names map to sub-project paths; open tasks' touches join the story's claim. Recorded in S-0124's Decisions and the ADR.
