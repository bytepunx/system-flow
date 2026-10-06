---
id: TH-0157
title: "Plan for S-0261: measure, propose an ADR, fit the pack, document, close I-0068"
anchor:
  path: wip/kanban/stories/S-0261-the-mcp-prime-pack-for-a-story-is-larger-than-claude-code-s-tool-result-limit-so-the-agent-reads-it-back-from-a-saved-file.md
  item: S-0261
status: resolved
participants: [planner-S-0261, alex]
created: 2026-10-05T05:53:18Z
updated: 2026-10-06T09:55:17Z
---

# TH-0157 Plan for S-0261: measure, propose an ADR, fit the pack, document, close I-0068

On wip/kanban/stories/S-0261-the-mcp-prime-pack-for-a-story-is-larger-than-claude-code-s-tool-result-limit-so-the-agent-reads-it-back-from-a-saved-file.md.

## Entries

### 2026-10-05T05:53:18Z planner-S-0261
Plan for S-0261, by planner-S-0261.

**Tasks and layers.** Each task's `after` gives these layers:

1. T-0986: measure I-0068's packs and propose the fix as an ADR that refines ADR-0049. Every other task waits on it.
2. Two tasks that share no path, so they run together:
   - T-0987: fit the pack in `flai/internal/context`, with a test that reproduces I-0068.
   - T-0989: the design and the guides (`agent-context.md`, `flai-cli.md`, `docs/users/flai.md` and `docs/users/flai-reference.md`, `docs/operators/settings.md`).
3. T-0988: MCP `prime` and `flai prime` return the fitted pack. A test in `flai/internal/mcpserver` checks that this repository's largest packs fit a tool result. It waits on T-0987.
4. T-0990: close I-0068. It waits on T-0988 and T-0989.

**Assumptions.**

- The limit is Claude Code's MCP tool result limit, 25,000 tokens unless `MAX_MCP_OUTPUT_TOKENS` raises it. Measured as JSON characters, it is reached below the 80 KB budget.
- There are two causes, and the fix must handle both:
  - Briefs past the budget. TH-0032 keeps every brief, so `exceeded: briefs` packs run 126 to 152 KB.
  - The JSON encoding. This planner's own pack was 81,903 bytes against a budget of 81,920, yet 92,948 characters as JSON, and Claude Code saved it to a file too. That is a fourth instance.
- The story's goal asks for a proposal before building, so T-0986 writes it as an ADR. If the fix it picks cuts briefs, it overrides your decision on TH-0032. T-0986 then asks you on a thread before the ADR is accepted.
- Without cutting briefs, the fix is either paging the pack with a part parameter, or returning the printed text with a budget measured on the encoded result.
- `design/system/project-manifest.md` is left out of the touches unless the fix adds a key next to `prime.budget`.

**Figures.**

- Forecast: 1h, raised from flai's 21m because the measuring and the ADR come before the build.
- Cost of delay: 15 USD a week, from your inputs as they stand.

**One question, which nothing waits on.** The cost of delay inputs say 6m a cycle, from 2 occurrences. I-0068 now has 3, and 4 counting this run. I recommend raising `time_lost_per_cycle` to 12m, which gives 30 USD a week. The inputs are yours, so I have left them as they are.

I propose no split, merge, or drop.

### 2026-10-05T05:54:46Z alex
12m

### 2026-10-06T09:55:17Z alex
Resolved.
