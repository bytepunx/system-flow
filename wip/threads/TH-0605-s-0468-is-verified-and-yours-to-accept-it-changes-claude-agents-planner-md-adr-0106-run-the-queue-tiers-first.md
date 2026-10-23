---
id: TH-0605
title: "S-0468 is verified and yours to accept: it changes .claude/agents/planner.md (ADR-0106); run the queue tiers first"
anchor:
  path: wip/kanban/stories/S-0468-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md
  item: S-0468
status: resolved
participants: [orchestrator, alex]
created: 2026-10-23T05:20:59Z
updated: 2026-10-23T17:53:06Z
---

# TH-0605 S-0468 is verified and yours to accept: it changes .claude/agents/planner.md (ADR-0106); run the queue tiers first

On wip/kanban/stories/S-0468-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md.

## Entries

### 2026-10-23T05:20:59Z orchestrator
**Recommended:** first run the two tiers that no one has run yet:

```sh
scripts/flai.sh verify S-0468 --all
```

Then accept S-0468 with `flai accept S-0468 --by alex`, or from the board.

## Why it is yours

S-0468 changes `.claude/agents/planner.md` and its template copy, so the merge queue will not take it (ADR-0106). The change appends two sentences to the planner's step 3, and nothing else:

- name files, not folders, in touches, keeping a folder only where the story's Planning note says why;
- name in `after` only a story whose accepted change is needed on main.

It widens no permission.

## Why `--all` first

`flai verify S-0468` passed every step at the branch head, 02851c72. It left integration and smoke to the merge queue, which will not run them for this story, so neither has run on it yet.

## What the verifier found

Both criteria are met and ticked, within the touches:

1. **A task's folder touch is refused** through `flai task new`, `flai edit`, `flai touches`, `item_new`, and `item_edit`. A planner's folder touch on a story is refused unless its Planning note names it. Your touches and the story's agent's are not refused. A folder inside `claims.shared` and a path not yet in the tree pass. `flai check` reports `touches.folder` and `wip.after-hub`. Thirteen new tests cover this, and ADR-0188 is accepted.
2. **I-0152 is closed** with its reason.

Both convention copies agree above the baseline marker. One nit: an unrelated merge of two bullets in `strategic-agents.md`.

### 2026-10-23T17:53:06Z alex
Resolved: S-0468 was accepted
