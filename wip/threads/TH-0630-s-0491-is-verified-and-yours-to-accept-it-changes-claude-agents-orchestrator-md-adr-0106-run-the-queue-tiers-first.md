---
id: TH-0630
title: "S-0491 is verified and yours to accept: it changes .claude/agents/orchestrator.md (ADR-0106); run the queue tiers first"
anchor:
  path: wip/kanban/stories/S-0491-add-sub-project-management.md
  item: S-0491
status: resolved
participants: [orchestrator, alex]
created: 2026-10-23T23:38:49Z
updated: 2026-10-23T23:38:50Z
---

# TH-0630 S-0491 is verified and yours to accept: it changes .claude/agents/orchestrator.md (ADR-0106); run the queue tiers first

On wip/kanban/stories/S-0491-add-sub-project-management.md.

## Entries

### 2026-10-23T23:38:49Z orchestrator
**Recommended:**

1. Run the two tiers no one has run on this story yet:

   ```sh
   scripts/flai.sh verify S-0491 --all
   ```

2. Then accept S-0491, with `flai accept S-0491 --by alex` or from the board.

This story widens my own permissions, so read my view with that in mind.

## Why it is yours

S-0491 changes `.claude/agents/orchestrator.md` and its template copy, so the merge queue will not take it (ADR-0106). `flai verify S-0491` passed every step at the branch head, 49925dd1. It left integration and smoke to a queue that will not run them for this story.

## What the verifier found

All three criteria are met and ticked, and all 52 changed files are within the touches.

1. **The Projects page** lists, adds, edits, and removes sub-projects through `settings.subproject`, which the `settings` host action gates.
2. **`flai subproject add`** creates a missing folder with a `README.md`, and removes it again if the edit is refused.
3. **Agents can do it themselves.** The MCP tools are `subprojects` (a read, open to all) and `subproject_edit` (refused to sub-agents, the planner, and the analyzer). There is also an orchestrator permission, `edit_subprojects`.

## Three points to check

1. **`edit_subprojects` is off by default.** flai guard refuses me `flai subproject add|set|remove` and `subproject_edit` while it is off, and a test covers both states. Nothing else widens.
2. **My agent definition changes in both trees, and the copies are identical.**
   - It gains the tool `mcp__flai__subprojects`.
   - The description and step 3 name sub-project changes.
   - A new step 13 says I change a sub-project only when a story or thread asks, by `flai subproject ... --autocommit` in the main checkout, and log it.
   - Steps 13 to 16 become 14 to 17, with their text unchanged.
3. **ADR-0195 is already `accepted`.** The story's agent set it, although it creates a new authority for me. Accepting the story is the natural point to confirm it. With the permission on, I would commit manifest changes to main without review. `docs/operators/settings.md` and the ADR say so.

### 2026-10-23T23:38:50Z alex
Resolved: S-0491 was accepted
