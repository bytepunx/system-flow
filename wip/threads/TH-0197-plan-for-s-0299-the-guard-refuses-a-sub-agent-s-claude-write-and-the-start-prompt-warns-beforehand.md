---
id: TH-0197
title: "Plan for S-0299: the guard refuses a sub-agent's .claude/ write, and the start prompt warns beforehand"
anchor:
  path: wip/kanban/stories/S-0299-a-story-s-sub-agent-that-writes-under-claude-blocks-its-layer-for-thirty-minutes-on-a-permission-thread-nobody-answers-and-the-start-prompt-does-not-warn-the-agent-beforehand.md
  item: S-0299
status: resolved
participants: [planner-S-0299, alex]
created: 2026-10-06T21:09:14Z
updated: 2026-10-06T21:17:41Z
---

# TH-0197 Plan for S-0299: the guard refuses a sub-agent's .claude/ write, and the start prompt warns beforehand

On wip/kanban/stories/S-0299-a-story-s-sub-agent-that-writes-under-claude-blocks-its-layer-for-thirty-minutes-on-a-permission-thread-nobody-answers-and-the-start-prompt-does-not-warn-the-agent-beforehand.md.

## Entries

### 2026-10-06T21:09:14Z planner-S-0299
The planner's plan for S-0299. It needs no answer to go on, but please redirect it before you finalize the draft if you want a different fix.

Proposed fix, from I-0093's one instance (S-0223's T-0957 writing under template/root/.claude/ while you were away): a sub-agent never makes a .claude/ write. The story's agent makes it itself after the layer, or hands it over as cp commands.

Tasks and layers:

- Layer 1, run together, no path in common:
  - T-1035: flai guard refuses at once a sub-agent's Edit, Write, MultiEdit, or NotebookEdit of a file in a .claude/ folder while auto-approve is off. The refusal tells the sub-agent to return the file's content in its final message. It touches the guard, flai/cmd/guard.go, the auto-approve description in hostapi/writes.go, and the generated reference. Its test reproduces I-0093.
  - T-1036: delegation() in flai/internal/harness/harness.go says the same in every story agent's start prompt, before it launches anything. It tells the agent to say so in each sub-agent's prompt, to make the writes itself once the layer is back, and, when you may be away and nothing else is left, to write the files into .flai-cache/ and open one thread with the cp commands, then end.
- Layer 2:
  - T-1037, after T-1035 and T-1036: a new ADR refining ADR-0060 and ADR-0086; design/system/flai-cli.md, docs/users/flai.md, and docs/operators/settings.md; a rewrite of the project addition in design/conventions/delegation.md, which today tells agents to Edit or Write .claude/ files like any other; and flai issue close I-0093.

Assumptions:

1. The guard is the right place for the fix. Its settings matcher already covers Edit|Write|NotebookEdit and it knows a sub-agent by agent_id, so no .claude/settings.json change is needed. MultiEdit is not in the matcher; I assumed that is fine, since agents do not use it here.
2. The story agent's own .claude/ write keeps going through permission_prompt as ADR-0086 has it. I did not plan a shorter hold for it, and the prompt steers it to the cp route when you may be away.
3. With auto-approve on, sub-agents may still write under .claude/, since nothing would hold them.
4. The template's baseline delegation.md is not changed. The rule is flai's own (guard and start prompt), and only this repository's project addition says the opposite.
5. Forecast 40m against flai's 18m, from the four closest stories (S-0257, S-0283, S-0285, S-0266). Cost of delay 15 USD a week, as flai cod gives it. I think the 6m input understates the cost: the hold could have run thirty minutes, and TH-0194 says S-0222 hit it too. Raising it is yours.

Nothing to split, merge, or drop.

### 2026-10-06T21:17:41Z alex
Resolved.
