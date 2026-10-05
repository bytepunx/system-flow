---
id: T-0864
type: task
nature: remediation
title: work-hierarchy.md states the touch rule, and I-0071 is closed
status: cancelled
parent: S-0260
owner: alex
created: 2026-10-05T03:13:57Z
updated: 2026-10-05T03:21:10Z
transitions:
  - to: cancelled
    at: 2026-10-05T03:21:10Z
    by: agent-S-0260
stream: S-0260
tags: [flai]
touches: [design/system/work-hierarchy.md, design/issues]
after: [T-0862, T-0863]
---
# T-0864 work-hierarchy.md states the touch rule, and I-0071 is closed

## Work

- In `design/system/work-hierarchy.md`, in the paragraph on `touches`, say what a touch may be. Every way of setting touches checks it: `flai touches`, `--touches` on `new` and `flai edit`, MCP `item_new` and `item_edit`, and the dashboard. A touch may start with a dot, as `.claude/agents/planner.md` does. It may not start with a dash, be absolute, or have a `..` segment. Cite S-0260.
- If a command's help text changed in T-0862 or T-0863, regenerate `docs/users/flai-reference.md` and add it to this task's touches.
- Close I-0071 from the story's worktree with `flai issue close I-0071 --reason`. Say that touches now go through one rule, `workitem.CleanTouches`, which accepts a leading dot, and that every way of setting touches uses it (S-0260). Commit the issue and `design/issues/summary.md` on the story branch.

It waits for T-0862 and T-0863: the design describes the rule as they build it, and the issue is closed once the fix is in.

## Done when

- `work-hierarchy.md` states the touch rule, and the markdown lint passes on it.
- I-0071 is closed with its reason, and `design/issues/summary.md` lists it as closed.
- `flai check --strict` scoped to S-0260 reports nothing.

## Notes
- 2026-10-05T03:21:10Z: moved to cancelled: duplicate of T-0865 and T-0866, which agent-S-0260 wrote on pulling S-0260 (TH-0124)
