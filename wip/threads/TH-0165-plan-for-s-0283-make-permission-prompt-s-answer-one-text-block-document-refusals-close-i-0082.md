---
id: TH-0165
title: "Plan for S-0283: make permission_prompt's answer one text block, document refusals, close I-0082"
anchor:
  path: wip/kanban/stories/S-0283-flai-s-permission-prompt-returns-a-result-claude-code-calls-invalid-so-a-claude-write-and-a-tmp-write-are-refused-at-once-with-no-thread.md
  item: S-0283
status: resolved
participants: [planner-S-0283, alex]
created: 2026-10-06T06:27:52Z
updated: 2026-10-06T06:48:04Z
---

# TH-0165 Plan for S-0283: make permission_prompt's answer one text block, document refusals, close I-0082

On wip/kanban/stories/S-0283-flai-s-permission-prompt-returns-a-result-claude-code-calls-invalid-so-a-claude-write-and-a-tmp-write-are-refused-at-once-with-no-thread.md.

## Entries

### 2026-10-06T06:27:52Z planner-S-0283
Tasks, in two layers:

1. Layer 1: T-1001, permission_prompt answers Claude Code with one text block and no structured content, with a test reproducing I-0082. This confirms the cause over the wire, proposes the fix in the narrative's Decisions, builds it in permission.go and folder.go, and traces why S-0219's .claude/ write opened no thread.
2. Layer 2, both after T-1001 and run together (no shared path): T-1002, the users' guide says a refused write reaches the agent with flai's reason; T-1003, close I-0082 saying what fixed it.

Figures: forecast 25m (flai's 16m plus 9m for the diagnosis), delivery 11:09Z. Cost of delay 25 USD a week, as flai cod gives from your 10m input.

Assumptions:

- The likely cause is the go-sdk's typed tool. mcp.AddTool with ToolHandlerFor[PermissionIn, PermissionOut] gives permission_prompt an output schema and adds structuredContent to its result. Claude Code wants exactly one text block. The tests read only Content[0], so they pass. T-1001 confirms this before building on it.
- A /tmp write being refused is what ADR-0086 intends: it is outside every story worktree. The defect is that the refusal arrives as an invalid result, not with its reason. So the plan does not widen the handler to allow /tmp. It makes the deny readable and documents it. No new ADR, since the decision stands and only the wire shape changes.
- If S-0219's .claude/ refusal was a deny ADR-0086 does not intend, such as a path check, T-1001 records it as an issue rather than widen this story.
- The fix reaches running agents only once a release with it is installed, because .mcp.json and flai serve run the installed flai.
- S-0284 changes the same files. flai's overlap hold runs the two stories one after the other, so I set no after.

No splits, merges, or drops proposed.

### 2026-10-06T06:48:04Z alex
Resolved.
