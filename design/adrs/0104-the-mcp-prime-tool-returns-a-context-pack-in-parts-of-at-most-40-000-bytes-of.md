---
id: ADR-0104
title: "The MCP prime tool returns a context pack in parts of at most 40,000 bytes of JSON, and the agent reads every part"
status: accepted
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0049]
topics: [cli, conventions]
---

# ADR-0104 The MCP prime tool returns a context pack in parts of at most 40,000 bytes of JSON, and the agent reads every part

## Context

ADR-0049 fitted a story's context pack to a budget of 80 KB, "so that the pack fits an MCP tool result with room for the agent's own turn". It does not. I-0068 records eight stories, from S-0205 on 2026-10-03 to S-0296 on 2026-10-06, whose MCP `prime` result Claude Code refused to show and saved to a file, which the agent then read back with `jq`.

Claude Code 2.1.290 sets two limits on an MCP tool's result. A result over 25,000 tokens (`MAX_MCP_OUTPUT_TOKENS`) is refused and saved to a file. Before that, any tool result over 50,000 characters is saved to a file and only a preview is shown. The character limit is the one that binds: `inbox` at 49.5 KB was saved too.

Measured on this repository on 2026-10-06 with `flai prime --json`, in bytes of the JSON the MCP tool returns:

| Pack | Budgeted size | Exceeded | JSON | Conventions in it | Items in it | Largest element |
|------|---------------|----------|------|-------------------|-------------|-----------------|
| S-0205 | 140,271 | briefs | 182,265 | 71,213 | 101,508 | 13,778 |
| S-0210 | 138,953 | briefs | 181,878 | 71,213 | 100,511 | 13,778 |
| S-0211 | 149,319 | named | 197,656 | 71,213 | 114,703 | 13,778 |
| S-0261 | 152,249 | named | 190,212 | 71,213 | 110,103 | 16,975 |
| Planner's pack for S-0261 | 123,950 | briefs | 161,300 | 40,988 | 114,705 | 16,975 |
| Explorer's pack for S-0261 | 40,808 | | 45,324 | 29,488 | 12,650 | 10,786 |

Two causes add up, and neither alone explains the refusal:

- **The pack is over its own budget.** ADR-0049 never cuts conventions, what the story names, or briefs (TH-0032), so `exceeded` packs run 124 to 152 KB of text.
- **The encoding.** The JSON carries what the printed pack does not: the path, reason, `also`, and `whole` of every item, the conventions' `kept` and `left_out` pieces, and the escaping. S-0261's 110 briefs are 29,206 characters of text and 92,530 of JSON.

But the conventions alone encode to 71,213 bytes for a story's agent and 40,988 for the planner, and ADR-0047 and ADR-0049 never cut them. No pack a story's agent gets fits one result of 50,000 characters, however its briefs and encoding are trimmed.

## Decision

**The MCP `prime` tool returns the pack in parts, each encoding to at most 40,000 bytes of JSON, and the agent reads every part by calling it again with `part` until it has read `parts` of them.** This refines ADR-0049, whose pack, budget, and order stand: what changes is how the pack is delivered, not what is in it. Every brief is still printed (TH-0032).

- Part 1 carries the pack's header (story, title, role, goal, topics, budget, `exceeded`, size, README, open issues, and what was left out), then conventions and items in the pack's order while they fit. Each later part carries the story, its own number, and the next conventions and items. The catalog goes in the last part.
- Every part carries `part` and `parts`. A pack that fits one part is one part, `parts: 1`, and reads as the whole pack did.
- A part ends before the convention or item that would take it over 40,000 bytes; nothing is split that fits a part of its own. A convention or item too large for any part is split at line boundaries into chunks, each carrying `chunk` and `chunks`, and each chunk starts a part.
- `part` defaults to 1. A part past `parts` is refused with the number of parts.
- `flai prime --part N --json` prints the same part. `flai prime` without `--part` prints the whole pack, as before, for a terminal or a file.

40,000 bytes is Claude Code's 50,000 characters, less a fifth for the result's envelope and for a harness that counts differently. Bytes are never fewer than the characters Claude Code counts, and 40,000 bytes of this JSON is about 13k tokens, well under the 25,000-token cap. The limit is a constant in `flai/internal/context`, `PartLimit`, beside `DefaultBudget`.

## Consequences

- No pack a story's agent, a sub-agent, or a strategic agent gets is saved to a file. Each measured story pack, 161 to 198 KB of JSON, takes at least five calls where one failed.
- An agent that reads part 1 and stops has the header, the conventions that fit, and `parts` telling it what remains. The tool's description, the agent's start prompt, `session-start.md`, and the design say to read every part.
- The parts are cut fresh on each call. A pack that changes between two calls, because a document changed, can repeat or skip an element at the boundary; the next prime is right again.
- `flai/internal/context` gains `PartLimit` and a way to cut a pack into parts. `flai/internal/mcpserver` gains `part` on `prime`. `flai prime` gains `--part`.
- The budget and the briefs are unchanged, so a pack over its budget is still over it; paging only makes it readable. Narrowing the conventions' topics stays the remedy ADR-0049 names for a pack whose conventions fill the budget.

## Alternatives considered

| Fix | What it removes from the measured packs | Why not |
|-----|------------------------------------------|---------|
| Fit the budget to the encoded result | S-0261's JSON is 190 KB against 152 KB of text: about a fifth | The conventions alone encode to 71 KB. A budget measured on the JSON is still over 50,000 characters for every story's pack, and to get under it the conventions would have to be cut, which ADR-0047 and ADR-0049 forbid |
| Return the printed text instead of the JSON | 190 KB to 152 KB for S-0261, 182 KB to 140 KB for S-0205 | Still three times the limit. It also loses the structure `jq` and the dashboard read |
| Brief fewer ADRs once the briefs exceed the budget, cataloguing the rest | S-0261's 92 KB of encoded briefs, at most | Overrides TH-0032. The conventions and the named documents, 71 KB and 17 KB encoded, still exceed the limit alone |
| Lower the default budget | Ranked sections and nothing else: conventions, named documents, and briefs are never cut | Every measured pack is already over its budget; a lower budget changes none of them |
| Page the pack, as the inbox's cap did for I-0023 | Nothing from the pack; each result is under the limit | Chosen. It is the only fix under which every measured pack fits, and it keeps ADR-0049 and TH-0032 whole |
