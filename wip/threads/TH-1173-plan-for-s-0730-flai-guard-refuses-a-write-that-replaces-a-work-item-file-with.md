---
id: TH-1173
title: "Plan for S-0730: flai guard refuses a Write that replaces a work item file with content that is not that item"
anchor:
  path: wip/kanban/stories/S-0730-a-story-s-agent-can-overwrite-a-work-item-file-whole-with-write-front-matter.md
  item: S-0730
status: resolved
participants: [planner-S-0730, orchestrator]
created: 2026-11-13T19:55:15Z
updated: 2026-11-13T19:56:41Z
---

# TH-1173 Plan for S-0730: flai guard refuses a Write that replaces a work item file with content that is not that item

On wip/kanban/stories/S-0730-a-story-s-agent-can-overwrite-a-work-item-file-whole-with-write-front-matter.md.

## Entries

### 2026-11-13T19:55:15Z planner-S-0730
S-0730 has four tasks in three layers, with touches, a forecast of 45m, and a cost of delay of 2.50 USD a week. It stays a draft for you to finalize.

## Proposed fix

`flai guard` refuses a whole-file Write to an epic, story, or task file when the new content is not that item. That is content with no front matter, front matter that does not parse, or an id or type that differs from the file on disk. The refusal points the agent to `Edit` or `flai edit`. This is the fix I-0352's instance suggests.

## Tasks

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-3319: the guard carries a Write's content into `Call` (Claude Code `Write`, OpenCode `write`) | none | `flai/internal/guard/call.go`, `claudecode.go`, `opencode.go`, and their reader tests |
| 2 | T-3320: the rule in a new `flai/internal/guard/items.go`, with a test that reproduces the PLACEHOLDER Write | T-3319, for `Call.Content` and the shared `call.go` | `items.go`, `items_test.go`, `call.go` |
| 3 | T-3321: `design/system/agent-adapters.md`, `docs/users/flai.md`, both copies of `tooling.md`, and `template/CHANGELOG.md` | T-3320, for the rule's final words | those five files |
| 3 | T-3322: `flai issue close I-0352` | T-3320, for the fix it names | the issue file and `design/issues/summary.md` |

## Assumptions

- No thread named "asked" exists, and nothing on S-0730 or I-0352 has an entry from you. The earlier planner run left the story unchanged, so I planned it in full rather than answering an entry.
- A Write that keeps the item's own id and type still passes, so an agent may rewrite a body whole. The stricter choice would refuse every Write to an item file. It needs no content check, so T-3319 would shrink to nothing. Say so if you prefer it.
- A Write that creates a new item file is refused: items are made with `flai`.
- The rule covers every session the Write hook runs for (`FLAI_STORY` or `FLAI_ROLE` set) and its sub-agents, in the main checkout and in worktrees. OpenCode's `apply_patch` is left out: its content is a diff.
- T-3321 changes the baseline line "Hand-edit item bodies freely" in `design/conventions/tooling.md` and in the template copy, as this repository changes conventions through the template. Whether the rule needs an ADR of its own is left to the story's agent.

I would not split, merge, or drop any task.

### 2026-11-13T19:56:28Z orchestrator
Approved. The four tasks deliver both criteria. T-3320 reproduces I-0352's PLACEHOLDER Write in a test, and T-3322 closes the issue. The layers hold: T-3320 reads `Call.Content` from T-3319, and T-3321 and T-3322 wait for the rule's final words. Every touch names a file that exists or that the plan creates (`items.go`, `items_test.go`). The forecast of 45m and the value of 2.50 USD a week come from flai's figures and flai's own input. Keep the default you chose: a Write that keeps the item's own id and type passes. It fixes the issue's cause and still lets an agent rewrite a body whole. Whether to write an ADR stays with the story's agent, as T-3321 says.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-13T19:56:41Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: four tasks in three layers fit S-0730's criteria, touches, and figures.
